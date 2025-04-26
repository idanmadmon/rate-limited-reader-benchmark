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

	fmt.Printf("RateLimitBasicFunctionalityTest Took %v\n", elapsed)
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

func MaxReadOverTimeSyntheticTest(readerFactory ReaderFactory) {
	const durationInSeconds = 10
	const bufferSize = 32 * 1024 // 32KB buffer
	fmt.Printf("Duration set: %d seconds\n", durationInSeconds)

	buffer := make([]byte, bufferSize)
	var totalBytes int64

	reader := infiniteReader{}
	ratelimitedReader := readerFactory(reader, 0) // no limit
	deadline := time.Now().Add(durationInSeconds * time.Second)

	for time.Now().Before(deadline) {
		n, err := ratelimitedReader.Read(buffer)
		if n > 0 {
			totalBytes += int64(n)
		}
		if err != nil {
			fmt.Printf("Read error: %v\n", err)
			break
		}
	}

	mb := float64(totalBytes) / 1024.0 / 1024.0
	fmt.Printf("MaxReadOverTimeSyntheticTest: Read %.3f MB in 10 seconds\n", mb)
}

type infiniteReader struct{}

func (infiniteReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = 'A'
	}
	return len(p), nil
}

func (infiniteReader) Close() error {
	return nil
}

func LargeReadFromNetTest(readerFactory ReaderFactory) {
	dataSize := 1 * 1024 * 1024 * 1024 // 1 GB of data
	var elapsed time.Duration

	readFunc := func(connReader io.ReadCloser) (int, error) {
		ratelimitedReader := readerFactory(connReader, 0) // no limit
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

	fmt.Printf("LargeReadFromNetTest Took %v\n", elapsed)
}
