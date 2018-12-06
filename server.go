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

func EncodeKey(metricType MetricType, metricId []byte, subtype int8, time int64) (result []byte) {
	sep := []byte("_")
	timeBuffer := make([]byte, 8)
	binary.BigEndian.PutUint64(timeBuffer, uint64(time))

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
func DecodeKey(key []byte) (metricType MetricType, metricId []byte, subtype int8, time int64) {
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
	timeLength := 8
	timeBuffer := key[length-timeLength:]
	time = int64(binary.BigEndian.Uint64(timeBuffer))
	subtype = int8(key[length-timeLength-2])
	metricId = key[3 : length-timeLength-3]
	return
}

func ApiServer(r *gin.Engine, rawKvClient *tikv.RawKVClient) {
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
		if metricTime < 1000000000000000 || metricTime > 9000000000000000 {
			errorResponse(c, "bad time range")
			return
		}

		// generate keys
		key := EncodeKey(ValueDataMetric, []byte(metricId), RawResolution, metricTime)

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
		keysInfoMetricKey := EncodeKey(KeysMetric, []byte(metricId), 0, 0)
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
		if limitStr != "" {
			limit64, err := strconv.ParseInt(limitStr, 10, 64)
			limit = int(limit64)
			if err != nil {
				errorResponse(c, "invalid limit")
				return
			}
		}

		sortStr := c.Query("sort")
		reverse := true
		if sortStr == "desc" || sortStr == "" {
			reverse = true
		} else if sortStr == "asc" {
			reverse = false
		} else {
			errorResponse(c, "invalid sort")
			return
		}

		var keys [][]byte
		var values [][]byte

		if reverse {
			startKey := EncodeKey(ValueDataMetric, targetId, RawResolution, upper)
			keys, values, _ = rawKvClient.ReverseScan(startKey, limit)
		} else {
			startKey := EncodeKey(ValueDataMetric, targetId, RawResolution, lower)
			keys, values, _ = rawKvClient.Scan(startKey, limit)
		}

		log.Printf("fond %v\n", len(keys))

		var responseRows []Row
		for i := range keys {
			metricType, metricId, resolution, time := DecodeKey(keys[i])
			if metricType != ValueDataMetric || resolution != RawResolution {
				break
			}
			if !bytes.Equal(metricId, targetId) {
				break
			}
			if time > upper { // out of range
				break
			}

			var unpacked interface{}
			_ = msgpack.Unmarshal(values[i], &unpacked)
			responseRows = append(responseRows, Row{Time: time, Value: unpacked})
		}

		res := MetricResponse{
			MetricId: string(targetId),
			Rows:     responseRows,
		}
		c.JSON(200, res)
	})

	r.GET("/metric/keys", func(c *gin.Context) {
		limitStr := c.Query("limit")
		limit := 1000
		if limitStr != "" {
			limit64, err := strconv.ParseInt(limitStr, 10, 64)
			limit = int(limit64)
			if err != nil {
				errorResponse(c, "invalid limit")
				return
			}
		}
		startKey := EncodeKey(KeysMetric, []byte{0}, 0, 0)
		keys, _, err := rawKvClient.Scan(startKey, limit)
		if err != nil {
			errorResponse(c, "cannot read db")
			return
		}
		var metricKeys []string
		for i := range keys {
			metricType, metricId, subtype, _ := DecodeKey(keys[i])
			if metricType != KeysMetric || subtype != 0 {
				break
			}
			metricKeys = append(metricKeys, string(metricId))
		}
		c.JSON(200, metricKeys)
	})
}

type MetricResponse struct {
	MetricId string `json:"metric_id"`
	Rows     []Row  `json:"rows"`
}
type Row struct {
	Time  int64       `json:"time"`
	Value interface{} `json:"value"`
}
