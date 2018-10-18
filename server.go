package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pingcap/tidb/store/tikv"
	"github.com/vmihailenco/msgpack"
	"log"
	"strconv"
)

func errorResponse(c *gin.Context, message string) {
	c.JSON(500, gin.H{
		"error": message,
	})
}

func encodeKey(metricId []byte, time int64) (result []byte) {
	prefix := []byte("v1_") // prefix
	sep := []byte("__")
	timeBuffer := make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(timeBuffer, time)

	result = append(result, prefix...)
	result = append(result, metricId...)
	result = append(result, sep...)
	result = append(result, []byte(timeBuffer)...)
	return
}
func decodeKey(key []byte) (metricId []byte, time int64) {
	length := len(key)
	timeBuffer := key[length-binary.MaxVarintLen64:]
	time, _ = binary.Varint(timeBuffer)
	metricId = key[3 : length-binary.MaxVarintLen64-2]
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
		// v1_metric-id__time-ns
		key := encodeKey([]byte(metricId), metricTime)

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

		err = rawKvClient.Put(key, packedValue)
		if err != nil {
			errorResponse(c, "can not write storage")
			return
		}

		c.JSON(200, gin.H{
			"ok": 1,
		})
	})

	r.GET("/metric/single/:id", func(c *gin.Context) {
		targetId := []byte("example-device")
		keys, values, _ := rawKvClient.Scan(encodeKey(targetId, 0), 100)
		log.Printf("fond %v\n", len(keys))

		responseRows := []Row{}
		for i := range keys {
			metricId, time := decodeKey(keys[i])
			if !bytes.Equal(metricId, targetId) {
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
