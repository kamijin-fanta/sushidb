package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pingcap/tidb/store/tikv"
	"github.com/vmihailenco/msgpack"
	"log"
	"math"
	"strconv"
)

func errorResponse(c *gin.Context, message string) {
	c.JSON(500, gin.H{
		"error": message,
	})
}

const (
	RawResolution int8 = iota
	CompressResolution
	OneMinutesResolution
	OneHourResolution
	OneDayResolution
)

type MetricType int

const (
	ValueDataMetric MetricType = iota
	KeysMetric
)

func encodeKey(metricType MetricType, metricId []byte, subtype int8, time int64) (result []byte) {
	sep := []byte("_")
	timeBuffer := make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(timeBuffer, time)

	var prefix []byte
	switch metricType {
	case ValueDataMetric:
		prefix = []byte("v1")
	case KeysMetric:
		prefix = []byte("k1")
	default:
		panic("undefined metric Type")
	}

	// [prefix]_[metricId]_[subtype]_[time ns]
	result = append(result, prefix[:]...)          // 2 bytes
	result = append(result, sep...)                // 1 byte
	result = append(result, metricId...)           // n bytes
	result = append(result, sep...)                // 1 byte
	result = append(result, byte(subtype))         // 1 byte
	result = append(result, sep...)                // 1 byte
	result = append(result, []byte(timeBuffer)...) // 8 bytes
	log.Printf("encode key: %+v, %s", result, string(result))
	return
}
func decodeKey(key []byte) (metricType MetricType, metricId []byte, subtype int8, time int64) {
	length := len(key)
	prefix := key[:2]
	switch string(prefix) {
	case "v1":
		metricType = ValueDataMetric
	case "k1":
		metricType = KeysMetric
	default:
		log.Fatalf("undefined metric type: %s\n", string(prefix))
	}
	timeBuffer := key[length-binary.MaxVarintLen64:]
	time, _ = binary.Varint(timeBuffer)
	subtype = int8(key[length-binary.MaxVarintLen64-2])
	metricId = key[3 : length-binary.MaxVarintLen64-3]
	return
}

func Server(rawKvClient *tikv.RawKVClient) *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {

		var input = []byte(`{"hoge": 123, "fuga": true}`)
		var str interface{}
		err := json.Unmarshal(input, &str)
		if err != nil {
			errorResponse(c, "invalid json")
			return
		}
		b, _ := msgpack.Marshal(str)
		log.Printf("====> %v\n", str)
		log.Printf("====> %x\n", b)
		msgpack.Unmarshal(b, str)
		a, _ := json.Marshal(str)
		log.Printf("====> %s\n", string(a))

		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/cluster", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"cluster": rawKvClient.ClusterID(),
		})
	})
	r.POST("/metric/single/:id/:time", func(c *gin.Context) {
		metricId := c.Param("id")
		metricTimeStr := c.Param("time")

		metricTime, err := strconv.ParseInt(metricTimeStr, 10, 64)
		if err != nil {
			errorResponse(c, "can not parse nano second time")
			return
		}
		if metricTime < 1000000000000 || metricTime > 9000000000000 {
			errorResponse(c, "bad time range")
			return
		}

		// generate keys
		key := encodeKey(ValueDataMetric, []byte(metricId), RawResolution, metricTime)

		// receive body -> decode json -> encode msgPack
		buf := make([]byte, 4096)
		readLength, _ := c.Request.Body.Read(buf)

		var receiveJson interface{}
		err = json.Unmarshal(buf[:readLength], &receiveJson)
		if err != nil {
			errorResponse(c, "invalid json")
			return
		}
		packedValue, _ := msgpack.Marshal(receiveJson)

		// write value
		err = rawKvClient.Put(key, packedValue)

		// write keys info
		keysInfoMetricKey := encodeKey(KeysMetric, []byte(metricId), 0, 0)
		err2 := rawKvClient.Put(keysInfoMetricKey, []byte{0})
		if err != nil || err2 != nil {
			log.Printf("%+v\n", err)
			log.Printf("%+v\n", err2)
			errorResponse(c, "can not write storage")
			return
		}

		c.JSON(200, gin.H{
			"ok": 1,
		})
	})

	r.GET("/metric/single/:id", func(c *gin.Context) {
		targetIdStr := c.Param("id")
		if targetIdStr == "" {
			errorResponse(c, "invalid metric id")
			return
		}
		targetId := []byte(targetIdStr)

		var err error
		lowerStr := c.Query("lower")
		lower := int64(0)
		if lowerStr != "" {
			lower, err = strconv.ParseInt(lowerStr, 10, 64)
			if err != nil {
				errorResponse(c, "invalid lower")
				return
			}
		}

		upperStr := c.Query("upper")
		upper := int64(math.MaxInt64)
		if upperStr != "" {
			upper, err = strconv.ParseInt(upperStr, 10, 64)
			if err != nil {
				errorResponse(c, "invalid upper")
				return
			}
		}
		limitStr := c.Query("limit")
		limit := 1000
		if upperStr != "" {
			upper, err = strconv.ParseInt(limitStr, 10, 64)
			if err != nil {
				errorResponse(c, "invalid limit")
				return
			}
		}

		startKey := encodeKey(ValueDataMetric, targetId, RawResolution, lower)
		keys, values, _ := rawKvClient.Scan(startKey, limit)
		log.Printf("fond %v\n", len(keys))

		var responseRows []Row
		for i := range keys {
			_, metricId, _, time := decodeKey(keys[i])
			if !bytes.Equal(metricId, targetId) {
				break
			}
			if time > upper { // out of range
				break
			}

			var unpacked interface{}
			msgpack.Unmarshal(values[i], &unpacked)
			responseRows = append(responseRows, Row{Time: time, Value: unpacked})
		}

		res := MetricResponse{
			MetricId: string(targetId),
			Rows:     responseRows,
		}
		c.JSON(200, res)
	})
	return r
}

type MetricResponse struct {
	MetricId string `json:"metric_id"`
	Rows     []Row  `json:"rows"`
}
type Row struct {
	Time  int64       `json:"time"`
	Value interface{} `json:"value"`
}
