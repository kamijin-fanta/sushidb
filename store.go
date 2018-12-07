package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/pingcap/pd/client"
	"github.com/pingcap/tidb/store/tikv"
	"github.com/vmihailenco/msgpack"
	"io/ioutil"
	"net/http"
)

type Store struct {
	rawKvClient tikv.RawKVClient
	pbClient    pd.Client
}

type KeyResponseRow struct {
	MetricId string `json:"metric_id"`
	Type     string `json:"type"`
}

func (s *Store) FetchKeys(start []byte, limit int) ([]KeyResponseRow, error) {
	responseKeys := make([]KeyResponseRow, 0)
	startKey := EncodeKey(PrefixKeysMetric, start, 0, 0)

	keys, _, err := s.rawKvClient.Scan(startKey, limit)
	if err != nil {
		return responseKeys, err
	}

	for i := range keys {
		metricType, metricId, subtypeId, _ := DecodeKey(keys[i])
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
		responseKeys = append(responseKeys, KeyResponseRow{
			MetricId: string(metricId),
			Type:     subtype,
		})
	}
	return responseKeys, nil
}

type SingleMetricResponseRow struct {
	Time  int64       `json:"time"`
	Value interface{} `json:"value"`
}

func (s *Store) FetchSingleMetric(metricId []byte, lower int64, upper int64, limit int, resolution int8, reverse bool) ([]SingleMetricResponseRow, error) {
	return s.FetchMetric(PrefixSingleValueMetric, metricId, lower, upper, limit, resolution, reverse)
}
func (s *Store) FetchMessageMetric(metricId []byte, lower int64, upper int64, limit int, resolution int8, reverse bool) ([]SingleMetricResponseRow, error) {
	return s.FetchMetric(PrefixMessageDataMetric, metricId, lower, upper, limit, resolution, reverse)
}

func (s *Store) FetchMetric(prefix PrefixTypes, metricId []byte, lower int64, upper int64, limit int, resolution int8, reverse bool) ([]SingleMetricResponseRow, error) {
	var keys [][]byte
	var values [][]byte
	var err error

	var responseRows []SingleMetricResponseRow

	if reverse {
		startKey := EncodeKey(prefix, metricId, resolution, upper)
		keys, values, err = s.rawKvClient.ReverseScan(startKey, limit)
	} else {
		startKey := EncodeKey(prefix, metricId, resolution, lower)
		keys, values, err = s.rawKvClient.Scan(startKey, limit)
	}
	if err != nil {
		return responseRows, err
	}

	for i := range keys {
		metricType, respondMetricId, respondResolution, time := DecodeKey(keys[i])
		if metricType != prefix || !bytes.Equal(respondMetricId, metricId) || resolution != respondResolution {
			break
		}
		if (!reverse && time > upper) || (reverse && time < lower) { // out of range
			break
		}

		var unpacked interface{}
		err = msgpack.Unmarshal(values[i], &unpacked)
		if err != nil {
			break
		}
		responseRows = append(responseRows, SingleMetricResponseRow{
			Time:  time,
			Value: unpacked,
		})
	}

	return responseRows, err
}

func (s *Store) PutSingleMetric(metricId []byte, time int64, resolution int8, value float64) error {
	packedValue, _ := msgpack.Marshal(value)
	return s.PutMetric(PrefixSingleValueMetric, metricId, time, resolution, packedValue)
}
func (s *Store) PutMessageMetric(metricId []byte, time int64, resolution int8, object interface{}) error {
	packedValue, _ := msgpack.Marshal(object)
	return s.PutMetric(PrefixMessageDataMetric, metricId, time, resolution, packedValue)
}

func (s *Store) PutMetric(prefix PrefixTypes, metricId []byte, time int64, resolution int8, body []byte) error {
	var subType int8
	switch prefix {
	case PrefixSingleValueMetric:
		subType = SubSingleKeys
		break
	case PrefixMessageDataMetric:
		subType = SubMessageKeys
		break
	default:
		panic("undefined prefix type")
	}

	// write value
	key := EncodeKey(prefix, []byte(metricId), resolution, time)
	writeValueError := s.rawKvClient.Put(key, body)
	if writeValueError != nil {
		return writeValueError
	}

	// write keys info
	keysInfoMetricKey := EncodeKey(prefix, []byte(metricId), subType, 0)
	writeKeyInfoError := s.rawKvClient.Put(keysInfoMetricKey, []byte{0})
	if writeKeyInfoError != nil {
		return writeKeyInfoError
	}

	return nil
}

type client interface {
	GetURLs() []string
}

func (s *Store) PdRequest(path string) ([]byte, error) {
	urls := s.pbClient.(client).GetURLs()
	maxInitClusterRetries := 10
	for i := 0; i < maxInitClusterRetries; i++ {
		for _, url := range urls {
			fmt.Println(url + path)
			res, err := http.Get(url + path)
			if err != nil {
				return nil, err
			}
			defer res.Body.Close()
			body, err := ioutil.ReadAll(res.Body)
			return body, err
		}
	}
	return nil, errors.New("cannot connect to pd")
}

func (s *Store) GetPdList() []string {
	return s.pbClient.(client).GetURLs()
}