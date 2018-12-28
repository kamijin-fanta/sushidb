package main

import (
	"encoding/binary"
	"fmt"
	"github.com/bmizerany/perks/quantile"
	"github.com/joho/godotenv"
	"github.com/pingcap/tidb/config"
	"github.com/pingcap/tidb/store/tikv"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Config struct {
	Threads            int
	MaxCount           int
	Duration           time.Duration
	DataChannelBuffer  int
	RandomSourceLength int
	DataSize           int
}
type Status struct {
	ResponseNs int64
	Error      int
}
type BenchmarkInput struct {
	Value []byte
	Seq   int
}

func main() {
	conf := Config{
		Threads:            32,
		MaxCount:           100000,          // 0: infinity
		Duration:           0 * time.Second, // 0: infinity
		DataChannelBuffer:  100,
		RandomSourceLength: 10000,
		DataSize:           1000, // 1kb
	}

	// load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file\n")
	}

	// gen random source
	randomBody := make([]byte, conf.RandomSourceLength)
	rand.Read(randomBody)

	start := time.Now()

	doneCh := make(chan struct{})                               // notify when all benchmark are completed
	dataCh := make(chan BenchmarkInput, conf.DataChannelBuffer) // send random dataCh from generator to benchmark
	statsCh := make(chan Status, conf.Threads*2)
	statsDoneCh := make(chan struct{})

	go Stats(statsCh, statsDoneCh, start, conf)
	go Generator(dataCh, randomBody, conf.MaxCount, start, conf.Duration, conf.DataSize)
	for i := 0; i < conf.Threads; i++ {
		go WriteBenchmark(dataCh, statsCh, doneCh, i)
	}

	for i := 0; i < conf.Threads; i++ {
		// wait for all goroutine
		<-doneCh
	}
	close(statsCh)
	<-statsDoneCh
	fmt.Println("done")
}

func Generator(data chan BenchmarkInput, bytes []byte, count int, start time.Time, duration time.Duration, maxSize int) {
	l := len(bytes)
	length := 100
	for i := 0; i < count || count == 0; i++ {
		if duration != 0 && time.Now().Sub(start) > duration {
			break
		}
		pos := i % (l - length)
		data <- BenchmarkInput{
			Value: bytes[pos : pos+length],
			Seq:   i,
		}
	}
	close(data)
}

func Stats(stats chan Status, statsDoneCh chan struct{}, start time.Time, conf Config) {
	tick := time.NewTicker(1 * time.Second)
	q := quantile.NewTargeted(0.5, 0.90, 0.95, 0.99)

	var totalSamples quantile.Samples
	totalCount := 0
	totalError := 0
	count := 0
	errorCount := 0
	var lastTime time.Time

	writeStats := func() {
		ms := float64(1000000)
		currentTime := time.Now()
		sub := currentTime.Sub(start)
		perSeconds := float64(count) / (currentTime.Sub(lastTime).Seconds())

		percent := float64(totalCount) / float64(conf.MaxCount) * 100
		if conf.MaxCount == 0 {
			percent = 0
		}
		timePercent := sub.Seconds() / conf.Duration.Seconds() * 100
		if conf.Duration.Nanoseconds() == 0 {
			timePercent = 0
		}
		if percent < timePercent {
			percent = timePercent
		}

		fmt.Printf("Stats%5.0fs(%.1f%%) Count: %d(%.2f/s) Error: %d / 50%%:%7.3fms  90%%:%7.3fms  95%%:%7.3fms  99%%:%7.3fms\n",
			sub.Seconds(), percent, count, perSeconds, errorCount, q.Query(0.5)/ms, q.Query(0.90)/ms, q.Query(0.95)/ms, q.Query(0.99)/ms)
		totalSamples = append(totalSamples, q.Samples()...)

		q.Reset()
		totalCount += count
		count = 0
		totalError += errorCount
		errorCount = 0
		lastTime = currentTime
	}
loop:
	for {
		select {
		case stat, ok := <-stats:
			if !ok {
				break loop
			}
			q.Insert(float64(stat.ResponseNs))
			count += 1
			errorCount += stat.Error
		case <-tick.C:
			writeStats()
		}
	}
	tick.Stop()

	writeStats()

	fmt.Printf("\nSummury\n")
	count = totalCount
	errorCount = totalError
	lastTime = start
	q.Merge(totalSamples)
	writeStats()

	statsDoneCh <- struct{}{}
}

func WriteBenchmark(data chan BenchmarkInput, stats chan Status, done chan struct{}, threadNumber int) {
	// client init
	pdAddress := os.Getenv("PD_ADDRESS")
	ope := os.Getenv("OPERATION")

	addressList := strings.Split(pdAddress, ",")
	rawClient, err := tikv.NewRawKVClient(addressList, config.Security{})
	if err != nil {
		panic(err)
	}
	defer rawClient.Close()

	// main loop
	for row := range data {
		start := time.Now().UnixNano()
		//millisecond := int64(1000 * 1000)
		//r := 10*millisecond + rand.Int63n(50*millisecond)
		//time.Sleep(time.Duration(r))

		if false {
			fmt.Printf("  WriteBencharmk(%d): %v\n", threadNumber, row)
		}

		key := []byte("benc_")
		buf := make([]byte, 8)
		binary.BigEndian.PutUint32(buf, uint32(row.Seq))
		key = append(key, buf...)

		var err error
		switch ope {
		default: // put
			err = rawClient.Put(key, row.Value)
		case "get":
			_, err = rawClient.Get(key)
		case "del":
			err = rawClient.Delete(key)
		}

		errorCount := 0
		if err != nil {
			errorCount = 1
		}

		stats <- Status{
			ResponseNs: time.Now().UnixNano() - start,
			Error:      errorCount,
		}

	}
	done <- struct{}{}
}
