package kvstore

import (
	"encoding/binary"
	"log"
)

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
type PrefixTypes int

const (
	PrefixSingleValueMetric PrefixTypes = iota
	PrefixMessageDataMetric
	PrefixKeysMetric
	PrefixKnown = 1000000000
)

func EncodeKey(metricType PrefixTypes, metricKey []byte, subtype int8, time int64) (result []byte) {
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

	// [prefix]_[metricKey]_[subtype]_[time ns]
	result = append(result, prefix[:]...)          // 2 bytes
	result = append(result, sep...)                // 1 byte
	result = append(result, metricKey...)          // n bytes
	result = append(result, sep...)                // 1 byte
	result = append(result, byte(subtype))         // 1 byte
	result = append(result, sep...)                // 1 byte
	result = append(result, []byte(timeBuffer)...) // 8 bytes
	log.Printf("encode key: %+v, %s", result, string(result))
	return
}
func DecodeKey(key []byte) (metricType PrefixTypes, metricKey []byte, subtype int8, time int64) {
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
		return PrefixKnown, metricKey, subtype, time
	}
	timeLength := 8
	timeBuffer := key[length-timeLength:]
	time = int64(binary.BigEndian.Uint64(timeBuffer))
	subtype = int8(key[length-timeLength-2])
	metricKey = key[3 : length-timeLength-3]
	return
}
