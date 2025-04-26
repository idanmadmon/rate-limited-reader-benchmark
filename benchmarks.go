package main

import (
	"fmt"
	"io"
	"time"
)

type BenchmarkTest func(ReaderFactory)

func RateLimitBasicFunctionalityTest(readerFactory ReaderFactory) {
	dataSize := 102400 // 100 KB of data
	partsAmount := 4
	limit := int64(dataSize / partsAmount) // dataSize/partsAmount bytes per second
	var elapsed time.Duration

	readFunc := func(connReader io.ReadCloser) (int, error) {
		ratelimitedReader := readerFactory(connReader, limit)
		start := time.Now()
		buffer := make([]byte, dataSize)
		n, err := ratelimitedReader.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Printf("Unexpected error while reading: %v\n", err)
		}

		if n != dataSize {
			fmt.Printf("Read incomplete data, read: %d expected: %d\n", n, dataSize)
		}

		elapsed = time.Since(start)
		return n, err
	}

	err := receiveOnceTCPServer(dataSize, readFunc)
	if err != nil {
		fmt.Printf("Unexpected error from server: %v\n", err)
	}

	fmt.Printf("Took %v\n", elapsed)
	minTimeInSeconds := partsAmount
	maxTimeInSeconds := partsAmount + 1
	minTime := time.Duration(minTimeInSeconds) * time.Second
	maxTime := time.Duration(maxTimeInSeconds) * time.Second
	if elapsed.Abs().Round(time.Second) < minTime { // round to second - has a deviation of up to half a second
		fmt.Printf("Read completed too quickly, elapsed time: %v < min time: %v\n", elapsed, minTime)
	} else if elapsed.Abs().Round(time.Second) > maxTime { // round to second - has a deviation of up to half a second
		fmt.Printf("Read completed too slow, elapsed time: %v > max time: %v\n", elapsed, maxTime)
	}
}

func TestReaderBehavior2(readerFactory ReaderFactory) {
	fmt.Println("Running TestReaderBehavior2")
}

func TestReaderBehavior3(readerFactory ReaderFactory) {
	fmt.Println("Running TestReaderBehavior3")
}
