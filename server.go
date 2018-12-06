package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
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

// Metric subtypes
const (
	SubRawResolution int8 = iota
	SubCompressResolution
	SubOneMinutesResolution
	SubOneHourResolution
	SubOneDayResolution
)

// Keys subtype
const (
	SubSingleKeys int8 = iota
	SubMessageKeys
)

// Prefix Types
type MetricType int

const (
	PrefixSingleValueMetric MetricType = iota
	PrefixMessageDataMetric
	PrefixKeysMetric
)

func EncodeKey(metricType MetricType, metricId []byte, subtype int8, time int64) (result []byte) {
	sep := []byte("_")
	timeBuffer := make([]byte, 8)
	binary.BigEndian.PutUint64(timeBuffer, uint64(time))

	var prefix []byte
	switch metricType {
	case PrefixSingleValueMetric:
		prefix = []byte("s1")
	case PrefixMessageDataMetric:
		prefix = []byte("m1")
	case PrefixKeysMetric:
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
	case "s1":
		metricType = PrefixSingleValueMetric
	case "m1":
		metricType = PrefixMessageDataMetric
	case "k1":
		metricType = PrefixKeysMetric
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
	/********** PING **********/
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

	/********** Cluster Info **********/
	r.GET("/cluster", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"cluster": rawKvClient.ClusterID(),
		})
	})

	/********** PostMetrics **********/
	r.POST("/metric/:type/:id/:time", func(c *gin.Context) {
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

		metricType, err := parseMetricType(c)
		if err != nil {
			errorResponse(c, "bad metric type")
			return
		}

		// receive body -> decode json
		buf := make([]byte, 4096)
		readLength, _ := c.Request.Body.Read(buf)
		var receiveJson interface{}
		err = json.Unmarshal(buf[:readLength], &receiveJson)
		if err != nil {
			errorResponse(c, "invalid json")
			return
		}

		var writeValueError error
		var writeKeyInfoError error

		switch metricType {
		case MetricSingle:
			floatValue, success := receiveJson.(float64)
			if !success {
				errorResponse(c, "invalid body. You can post a numerical value.")
			}

			// encode msgPack
			packedValue, _ := msgpack.Marshal(floatValue)

			// write value
			key := EncodeKey(PrefixSingleValueMetric, []byte(metricId), SubRawResolution, metricTime)
			writeValueError = rawKvClient.Put(key, packedValue)

			// write keys info
			keysInfoMetricKey := EncodeKey(PrefixKeysMetric, []byte(metricId), SubSingleKeys, 0)
			writeKeyInfoError = rawKvClient.Put(keysInfoMetricKey, []byte{0})

			break
		case MetricMessage:
			// encode msgPack
			packedValue, _ := msgpack.Marshal(receiveJson)

			// write value
			key := EncodeKey(PrefixMessageDataMetric, []byte(metricId), SubRawResolution, metricTime)
			writeValueError = rawKvClient.Put(key, packedValue)

			// write keys info
			keysInfoMetricKey := EncodeKey(PrefixKeysMetric, []byte(metricId), SubMessageKeys, 0)
			writeKeyInfoError = rawKvClient.Put(keysInfoMetricKey, []byte{0})
			break
		}

		// display errors
		if writeValueError != nil || writeKeyInfoError != nil {
			log.Printf("%+v\n", writeValueError)
			log.Printf("%+v\n", writeKeyInfoError)
			errorResponse(c, "can not write storage")
			return
		}

		c.JSON(200, gin.H{
			"ok": 1,
		})
	})

	/********** Query Metrics **********/
	r.GET("/metric/:type/:id", func(c *gin.Context) {
		targetIdStr := c.Param("id")
		if targetIdStr == "" {
			errorResponse(c, "invalid metric id")
			return
		}
		targetId := []byte(targetIdStr)

		var err error
		metricType, err := parseMetricType(c)
		if err != nil {
			errorResponse(c, "bad metric type")
			return
		}
		var prefix MetricType
		switch metricType {
		case MetricSingle:
			prefix = PrefixSingleValueMetric
			break
		case MetricMessage:
			prefix = PrefixMessageDataMetric
			break
		}

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
			startKey := EncodeKey(prefix, targetId, SubRawResolution, upper)
			keys, values, _ = rawKvClient.ReverseScan(startKey, limit)
		} else {
			startKey := EncodeKey(prefix, targetId, SubRawResolution, lower)
			keys, values, _ = rawKvClient.Scan(startKey, limit)
		}

		log.Printf("fond %v\n", len(keys))

		var responseRows []Row
		for i := range keys {
			metricType, metricId, resolution, time := DecodeKey(keys[i])
			if metricType != prefix || !bytes.Equal(metricId, targetId) || resolution != SubRawResolution {
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

	/********** Query Keys **********/
	r.GET("/keys", func(c *gin.Context) {
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

		startKey := EncodeKey(PrefixKeysMetric, []byte{0}, 0, 0)

		keys, _, err := rawKvClient.Scan(startKey, limit)
		if err != nil {
			errorResponse(c, "cannot read db")
			return
		}
		metricKeys := make([]Key, 0)
		for i := range keys {
			metricType, metricId, subtypeId, _ := DecodeKey(keys[i])
			log.Printf("%+v / %+v / %+v\n", metricType, string(metricId), subtypeId)
			if metricType != PrefixKeysMetric {
				break
			}
			var subtype string
			switch subtypeId {
			case SubMessageKeys:
				subtype = "message"
				break
			case SubSingleKeys:
				subtype = "single"
				break
			}
			metricKeys = append(metricKeys, Key{
				MetricId: string(metricId),
				Type:     subtype,
			})
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
type Key struct {
	MetricId string `json:"metric_id"`
	Type     string `json:"type"`
}

const (
	MetricSingle = iota
	MetricMessage
)

func parseMetricType(c *gin.Context) (int, error) {
	str := c.Param("type")
	switch str {
	case "single":
		return MetricSingle, nil
	case "message":
		return MetricMessage, nil
	default:
		return 0, errors.New("parse error")
	}
}
