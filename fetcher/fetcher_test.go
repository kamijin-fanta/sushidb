package fetcher

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockDefined struct {
	MetricKey string
	TimeStart int64
	TimeEnd   int64
	TimeStep  int64
}

var defs = []MockDefined{
	{"aaa", 1000, 1500, 20},
	{"bbb", 1200, 1700, 20},
	{"ccc", 1110, 1500, 20},
}

type MockResourceImpl struct{}

func (r MockResourceImpl) fetch(key []byte, timestamp int64, asc bool) (rows []Row, error error) {
	limit := 5
	for _, def := range defs {
		if bytes.Compare([]byte(def.MetricKey), key) == 0 {
			var offset int64
			if asc && def.TimeStart < timestamp {
				diff := timestamp - def.TimeStart
				offset = diff / def.TimeStep
				if diff%def.TimeStep != 0 {
					offset++
				}
			}
			if !asc && def.TimeEnd > timestamp {
				diff := timestamp - def.TimeStart
				if timestamp == 0 {
					diff = def.TimeEnd - def.TimeStart
				}
				offset = diff / def.TimeStep
				if timestamp != 0 && diff%def.TimeStep == 0 {
					offset--
				}
			}
			for i := 0; i < limit; i++ {
				var time int64
				if asc {
					time = def.TimeStart + def.TimeStep*(offset+int64(i))
				} else {
					time = def.TimeStart + def.TimeStep*(offset-int64(i))
				}
				if def.TimeEnd >= time && def.TimeStart <= time {
					rows = append(rows, Row{
						MetricKey: []byte(def.MetricKey),
						TimeStamp: time,
					})
				}
			}
			break
		}
	}
	return rows, nil
}

func PrintRows(rows []Row) {
	for _, row := range rows {
		fmt.Printf("Key: %v  Time: %d\n", string(row.MetricKey), row.TimeStamp)
	}
}

func TestMockResource(t *testing.T) {
	var mockResource Resource = MockResourceImpl{}
	rows, err := mockResource.Fetch([]byte("aaa"), 0, true)
	assert.Nil(t, err)
	assert.Equal(t, rows[0].TimeStamp, int64(1000))

	rows, err = mockResource.Fetch([]byte("aaa"), 1120, true)
	assert.Nil(t, err)
	assert.Equal(t, rows[0].TimeStamp, int64(1120))

	rows, err = mockResource.Fetch([]byte("aaa"), 1130, true)
	assert.Nil(t, err)
	assert.Equal(t, rows[0].TimeStamp, int64(1140))

	rows, err = mockResource.Fetch([]byte("aaa"), 0, false)
	assert.Nil(t, err)
	assert.Equal(t, rows[0].TimeStamp, int64(1500))

	rows, err = mockResource.Fetch([]byte("aaa"), 1380, false)
	assert.Nil(t, err)
	assert.Equal(t, rows[0].TimeStamp, int64(1360))

	rows, err = mockResource.Fetch([]byte("aaa"), 1370, false)
	assert.Nil(t, err)
	assert.Equal(t, rows[0].TimeStamp, int64(1360))
}

func ExampleFetchSingleAsc() {
	var mockResource Resource = MockResourceImpl{}
	requestKeys := [][]byte{
		[]byte("aaa"),
	}
	fetcher := NewFetcher(requestKeys, 1180, 0, true, mockResource)
	rows, _ := fetcher.Next(5)
	PrintRows(rows)

	// Output:
	// Key: aaa  Time: 1180
	// Key: aaa  Time: 1200
	// Key: aaa  Time: 1220
	// Key: aaa  Time: 1240
	// Key: aaa  Time: 1260
}

func ExampleFetchMultiAsc() {
	var mockResource Resource = MockResourceImpl{}
	requestKeys := [][]byte{
		[]byte("aaa"),
		[]byte("bbb"),
		[]byte("ccc"),
	}
	fetcher := NewFetcher(requestKeys, 1180, 0, true, mockResource)
	rows, _ := fetcher.Next(5)
	PrintRows(rows)

	// Output:
	// Key: aaa  Time: 1180
	// Key: ccc  Time: 1190
	// Key: aaa  Time: 1200
	// Key: bbb  Time: 1200
	// Key: ccc  Time: 1210
}

func ExampleFetchAscLimited() {
	var mockResource Resource = MockResourceImpl{}
	requestKeys := [][]byte{
		[]byte("aaa"),
		[]byte("bbb"),
		[]byte("ccc"),
	}
	fetcher := NewFetcher(requestKeys, 1190, 1230, true, mockResource)
	rows, _ := fetcher.Next(100)
	PrintRows(rows)

	// Output:
	// Key: ccc  Time: 1190
	// Key: aaa  Time: 1200
	// Key: bbb  Time: 1200
	// Key: ccc  Time: 1210
	// Key: aaa  Time: 1220
	// Key: bbb  Time: 1220
}

func ExampleFetchSingleDesc() {
	var mockResource Resource = MockResourceImpl{}
	requestKeys := [][]byte{
		[]byte("aaa"),
	}
	fetcher := NewFetcher(requestKeys, 1270, 0, false, mockResource)
	rows, _ := fetcher.Next(5)
	PrintRows(rows)

	// Output:
	// Key: aaa  Time: 1260
	// Key: aaa  Time: 1240
	// Key: aaa  Time: 1220
	// Key: aaa  Time: 1200
	// Key: aaa  Time: 1180
}

func ExampleFetchMultiDesc() {
	var mockResource Resource = MockResourceImpl{}
	requestKeys := [][]byte{
		[]byte("ccc"),
		[]byte("aaa"),
		[]byte("bbb"),
	}
	fetcher := NewFetcher(requestKeys, 1220, 1140, false, mockResource)
	rows, _ := fetcher.Next(5)
	PrintRows(rows)
	fmt.Println(fetcher.MaybeHasNext)

	rows, _ = fetcher.Next(5)
	PrintRows(rows)
	fmt.Println(fetcher.MaybeHasNext)

	// Output:
	// Key: ccc  Time: 1210
	// Key: aaa  Time: 1200
	// Key: bbb  Time: 1200
	// Key: ccc  Time: 1190
	// Key: aaa  Time: 1180
	// true
	// Key: ccc  Time: 1170
	// Key: aaa  Time: 1160
	// Key: ccc  Time: 1150
	// Key: aaa  Time: 1140
	// false
}
