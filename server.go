package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/kamijin-fanta/sushidb/fetcher"
	"github.com/kamijin-fanta/sushidb/kvstore"
	"github.com/kamijin-fanta/sushidb/querying"
	"io"
	"log"
	"math"
	"strconv"
	"time"
)

func errorResponse(c *gin.Context, message string) {
	c.JSON(500, gin.H{
		"error": message,
	})
}

func ApiServer(r *gin.Engine, store *kvstore.Store) {
	/********** PING **********/
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	/********** Cluster Info **********/
	r.GET("/cluster", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"cluster": store.ClusterID(),
		})
	})

	/********** PostMetrics **********/
	r.POST("/metric/:type/:id/:time", func(c *gin.Context) {
		metricId := c.Param("id")
		metricIdBytes := []byte(metricId)
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

		var writeError error

		// write value
		switch metricType {
		case MetricSingle:
			floatValue, success := receiveJson.(float64)
			if !success {
				errorResponse(c, "invalid body. You can post a numerical value.")
			}
			writeError = store.PutSingleMetric(metricIdBytes, metricTime, kvstore.SubRawResolution, floatValue)
			break
		case MetricMessage:
			writeError = store.PutMessageMetric(metricIdBytes, metricTime, kvstore.SubRawResolution, receiveJson)
			break
		}

		// display errors
		if writeError != nil {
			log.Printf("%+v\n", writeError)
			errorResponse(c, "can not write storage")
			return
		}

		c.JSON(200, gin.H{
			"ok": 1,
		})
	})

	/********** Query Metrics **********/
	r.GET("/metric/:type/:id", func(c *gin.Context) {
		c.Set("req", time.Now().UnixNano())

		var err error
		metricType, err := parseMetricType(c)
		if err != nil {
			errorResponse(c, "bad metric type")
			return
		}

		targetIdStr := c.Param("id")
		if targetIdStr == "" {
			errorResponse(c, "invalid metric id")
			return
		}
		targetId := []byte(targetIdStr)

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

		var rows []kvstore.SingleMetricResponseRow
		var fetchErr error

		switch metricType {
		case MetricSingle:
			rows, fetchErr = store.FetchSingleMetric(targetId, lower, upper, limit, kvstore.SubRawResolution, reverse, false)
		case MetricMessage:
			rows, fetchErr = store.FetchMessageMetric(targetId, lower, upper, limit, kvstore.SubRawResolution, reverse, false)
		}
		if fetchErr != nil {
			errorResponse(c, "fetch error")
			return
		}

		res := MetricResponse{
			MetricId:    string(targetId),
			Rows:        rows,
			QueryTimeNs: time.Now().UnixNano() - c.GetInt64("req"),
		}
		c.JSON(200, res)
	})

	/********** Advanced Query Metrics **********/
	r.POST("/query/:type/:id", func(c *gin.Context) {
		c.Set("req", time.Now().UnixNano())

		metricType, err := parseMetricType(c)
		if err != nil {
			errorResponse(c, "bad metric type")
			return
		}

		targetIdStr := c.Param("id")
		if targetIdStr == "" {
			errorResponse(c, "invalid metric id")
			return
		}
		targetId := []byte(targetIdStr)

		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, c.Request.Body)
		if err != nil {
			errorResponse(c, "cannot request body")
			return
		}
		postData := buf.Bytes()

		query, err := querying.New(postData)
		if err != nil {
			errorResponse(c, "invalid query jsondata")
			return
		}

		reverse := true
		switch query.Query.Sort {
		case "desc", "":
			reverse = true
		case "asc":
			reverse = false
		default:
			errorResponse(c, "invalid sort")
			return
		}

		var prefixTypes kvstore.PrefixTypes
		switch metricType {
		case MetricSingle:
			prefixTypes = kvstore.PrefixSingleValueMetric
		case MetricMessage:
			prefixTypes = kvstore.PrefixMessageDataMetric
		}

		cursorTimestamp, _, cursorErr := query.Query.ParseCursor()
		if cursorErr != nil {
			errorResponse(c, "cannot parse cursor")
			return
		}

		resource := kvstore.StoreResourceImpl{
			Store:       store,
			Limit:       100, // todo batch size
			PrefixTypes: prefixTypes,
		}
		var storeFetcher fetcher.Fetcher
		var keys [][]byte
		// keys = append(keys, targetId) // todo multi
		for i := range query.Query.MetricKeys {
			keys = append(keys, []byte(query.Query.MetricKeys[i]))
		}
		switch query.Query.Sort {
		default: // desc
			storeFetcher = fetcher.NewFetcher(keys, query.Query.Upper, query.Query.Lower, false, &resource)
		case "asc":
			storeFetcher = fetcher.NewFetcher(keys, query.Query.Upper, query.Query.Lower, true, &resource)
		}

		var rows []fetcher.Row
		var fetchErr error
		filteredRes := make([]kvstore.SingleMetricResponseRow, 0)
		var lastTimestamp int64 = 0
		skipCount := 0

		fetchErr = storeFetcher.PreFetch()
		if fetchErr != nil {
			errorResponse(c, "fetch error")
			return
		}
		for len(filteredRes) < query.Query.Limit && skipCount < query.Query.MaxSkip {
			includeLastBorder := false
			lower := query.Query.Lower
			upper := query.Query.Upper
			if cursorTimestamp != 0 {
				if reverse && upper >= cursorTimestamp { //desc
					upper = cursorTimestamp
				} else if !reverse && lower <= cursorTimestamp { // asc
					lower = cursorTimestamp
					includeLastBorder = true // skip first row
				}
			}
			limit := query.Query.Limit - len(filteredRes) + query.Query.MaxSkip/2

			resource.IncludeLastBorder = includeLastBorder // todo
			rows, fetchErr = storeFetcher.Next(limit)
			if fetchErr != nil {
				errorResponse(c, "fetch error")
				return
			}

			for _, row := range rows {
				condition, err := query.FilterRow(row.Value)
				if err != nil {
					errorResponse(c, "query error"+err.Error())
					return
				}
				if condition {
					filteredRes = append(filteredRes, kvstore.SingleMetricResponseRow{
						Value:     row.Value,
						Time:      row.TimeStamp,
						MetricKey: string(row.MetricKey),
					})
				} else {
					skipCount += 1
				}
				lastTimestamp = row.TimeStamp

				if len(filteredRes) >= query.Query.Limit || skipCount >= query.Query.MaxSkip {
					break
				}
			}
			cursorTimestamp = lastTimestamp

			if len(rows) < limit { // If len does not reach limit, there is no next row
				break
			}
		}

		res := MetricResponse{
			MetricId:    string(targetId),
			Rows:        filteredRes,
			QueryTimeNs: time.Now().UnixNano() - c.GetInt64("req"),
			Cursor:      strconv.FormatInt(lastTimestamp, 10) + "," + strconv.Itoa(0),
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

		metricKeys, err := store.FetchKeys([]byte{0}, limit)
		if err != nil {
			errorResponse(c, "invalid sort")
			return
		}
		c.JSON(200, metricKeys)
	})

	/********** PD List **********/
	r.GET("/pd/", func(c *gin.Context) {
		res := store.GetPdList()
		c.JSON(200, res)
	})

	/********** PD Infos **********/
	r.GET("/pd/api/*any", func(c *gin.Context) {
		res, err := store.PdRequest(c.Request.URL.Path)
		if err != nil {
			errorResponse(c, err.Error())
			return
		}

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Content-Length", strconv.Itoa(len(res)))
		c.Writer.WriteHeader(200)
		c.Writer.Write(res)
	})
}

type MetricResponse struct {
	MetricId    string                            `json:"metric_id"`
	Rows        []kvstore.SingleMetricResponseRow `json:"rows"`
	QueryTimeNs int64                             `json:"query_time_ns"`
	Cursor      string                            `json:"cursor"`
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
