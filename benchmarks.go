package main

import (
	"fmt"
	"io"
	"strings"
	"time"
)

type BenchmarkTest func(ReaderFactory)

func RateLimitingSyntheticTest(readerFactory ReaderFactory) {
	const dataSize = 100 * 1024 * 1024 // 100MB
	const bufferSize = 32 * 1024       // 32KB classic io.Copy
	const limit = dataSize / 4         // should take 4 seconds
	var total int

	reader := &syntheticReader{size: dataSize}
	limitedReader := readerFactory(reader, bufferSize, limit)
	buffer := make([]byte, bufferSize)

	start := time.Now()
	for {
		n, err := limitedReader.Read(buffer)
		total += n
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Error: %v\n", err)
			}
			break
		}
	}
	elapsed := time.Since(start)

	if total != dataSize {
		fmt.Printf("Read incomplete data, read: %d expected: %d\n", total, dataSize)
	}

	fmt.Printf("RateLimitingSyntheticTest Took %v\n", elapsed)
}

func MaxReadOverTimeSyntheticTest(readerFactory ReaderFactory) {
	const durationInSeconds = 10
	const bufferSize = 32 * 1024 // 32KB classic io.Copy
	const limit = 0              //bufferSize * 10000 // large limit
	fmt.Printf("Duration set: %d seconds\n", durationInSeconds)

	buffer := make([]byte, bufferSize)
	var totalBytes int64

	reader := &syntheticReader{}
	ratelimitedReader := readerFactory(reader, bufferSize, limit)

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

func RateLimitingRealWorldLocalTest(readerFactory ReaderFactory) {
	const dataSize = 100 * 1024 * 1024 // 100MB
	const bufferSize = 32 * 1024       // 32KB classic io.Copy
	const limit = dataSize / 4         // should take 4 seconds
	var elapsed time.Duration

	rf := func(connReader io.ReadCloser) (int, error) {
		ratelimitedReader := readerFactory(connReader, bufferSize, limit)

		var total, n int
		var err error
		buffer := make([]byte, bufferSize)
		start := time.Now()
		for {
			n, err = ratelimitedReader.Read(buffer)
			total += n
			if err != nil {
				if err != io.EOF {
					fmt.Printf("Unexpected error while reading: %v\n", err)
				}
				break
			}
		}

		elapsed = time.Since(start)
		if total != dataSize {
			fmt.Printf("Read incomplete data, read: %d expected: %d\n", n, dataSize)
		}

		return total, err
	}

	wf := func(connWriter io.Writer) (int, error) {
		message := strings.Repeat("A", dataSize)
		return connWriter.Write([]byte(message))
	}

	go func() {
		// give the server a sec to start
		time.Sleep(1 * time.Second)

		n, err := sendTCPMessage(wf)
		if err != nil {
			fmt.Println("Failed to send message:", err)
		}
		if n != dataSize {
			fmt.Printf("Failed to send message: sent insufficient size=%d expectedSize=%d\n", n, dataSize)
		}
	}()

	n, err := receiveOnceTCPServer(rf)
	if err != nil && err != io.EOF {
		fmt.Printf("Unexpected error from server: %v\n", err)
	}
	if n != dataSize {
		fmt.Printf("Failed to get message: got insufficient size=%d expectedSize=%d\n", n, dataSize)
	}

	fmt.Printf("RateLimitingRealWorldLocalTest Took %v\n", elapsed)
}

func RateLimitingRealWorldServerTest(readerFactory ReaderFactory) {
	const dataSize = 100 * 1024 * 1024 // 100MB
	const bufferSize = 32 * 1024       // 32KB classic io.Copy
	const limit = dataSize / 4         // should take 4 seconds
	var elapsed time.Duration

	rf := func(connReader io.ReadCloser) (int, error) {
		ratelimitedReader := readerFactory(connReader, bufferSize, limit)

		var total, n int
		var err error
		buffer := make([]byte, bufferSize)
		start := time.Now()
		for {
			n, err = ratelimitedReader.Read(buffer)
			total += n
			if err != nil {
				if err != io.EOF {
					fmt.Printf("Unexpected error while reading: %v\n", err)
				}
				break
			}
		}

		elapsed = time.Since(start)
		if total != dataSize {
			fmt.Printf("Read incomplete data, read: %d expected: %d\n", n, dataSize)
		}

		return total, err
	}

	n, err := receiveOnceTCPServer(rf)
	if err != nil && err != io.EOF {
		fmt.Printf("Unexpected error from server: %v\n", err)
	}
	if n != dataSize {
		fmt.Printf("Failed to get message: got insufficient size=%d expectedSize=%d\n", n, dataSize)
	}

	fmt.Printf("RateLimitingRealWorldServerTest Took %v\n", elapsed)
}

func SpikeRecoveryRealWorldLocalTest(readerFactory ReaderFactory) {
	const dataSize = 100 * 1024 * 1024 // 100MB
	const bufferSize = 32 * 1024       // 32KB classic io.Copy
	const limit = dataSize / 8         // should take 8 seconds
	var elapsed time.Duration

	rf := func(connReader io.ReadCloser) (int, error) {
		ratelimitedReader := readerFactory(connReader, bufferSize, limit)

		var total, n int
		var err error
		buffer := make([]byte, bufferSize)
		start := time.Now()
		for {
			n, err = ratelimitedReader.Read(buffer)
			total += n
			if err != nil {
				if err != io.EOF {
					fmt.Printf("Unexpected error while reading: %v\n", err)
				}
				break
			}
		}

		elapsed = time.Since(start)
		if total != dataSize {
			fmt.Printf("Read incomplete data, read: %d expected: %d\n", n, dataSize)
		}

		return total, err
	}

	wf := func(connWriter io.Writer) (int, error) {
		chunkSize := dataSize / 8
		var total int
		var n int
		var err error

		for i := 0; i < 1; i++ {
			n, err = connWriter.Write([]byte(strings.Repeat("A", chunkSize)))
			total += n
			if err != nil {
				return total, nil
			}
			time.Sleep(1 * time.Second)
		}

		n, err = connWriter.Write([]byte(strings.Repeat("A", chunkSize*5)))
		total += n
		if err != nil {
			return total, nil
		}
		time.Sleep(1 * time.Second)

		for i := 0; i < 2; i++ {
			n, err = connWriter.Write([]byte(strings.Repeat("A", chunkSize)))
			total += n
			if err != nil {
				return total, nil
			}
			time.Sleep(1 * time.Second)
		}

		return total, err
	}

	go func() {
		// give the server a sec to start
		time.Sleep(1 * time.Second)

		n, err := sendTCPMessage(wf)
		if err != nil {
			fmt.Println("Failed to send message:", err)
		}
		if n != dataSize {
			fmt.Printf("Failed to send message: sent insufficient size=%d expectedSize=%d\n", n, dataSize)
		}
	}()

	n, err := receiveOnceTCPServer(rf)
	if err != nil && err != io.EOF {
		fmt.Printf("Unexpected error from server: %v\n", err)
	}
	if n != dataSize {
		fmt.Printf("Failed to get message: got insufficient size=%d expectedSize=%d\n", n, dataSize)
	}

	fmt.Printf("SpikeRecoveryRealWorldLocalTest Took %v\n", elapsed)
}

func SpikeRecoveryRealWorldServerTest(readerFactory ReaderFactory) {
	const dataSize = 100 * 1024 * 1024 // 100MB
	const bufferSize = 32 * 1024       // 32KB classic io.Copy
	const limit = dataSize / 8         // should take 8 seconds
	var elapsed time.Duration

	rf := func(connReader io.ReadCloser) (int, error) {
		ratelimitedReader := readerFactory(connReader, bufferSize, limit)

		var total, n int
		var err error
		buffer := make([]byte, bufferSize)
		start := time.Now()
		for {
			n, err = ratelimitedReader.Read(buffer)
			total += n
			if err != nil {
				if err != io.EOF {
					fmt.Printf("Unexpected error while reading: %v\n", err)
				}
				break
			}
		}

		elapsed = time.Since(start)
		if total != dataSize {
			fmt.Printf("Read incomplete data, read: %d expected: %d\n", n, dataSize)
		}

		return total, err
	}

	n, err := receiveOnceTCPServer(rf)
	if err != nil && err != io.EOF {
		fmt.Printf("Unexpected error from server: %v\n", err)
	}
	if n != dataSize {
		fmt.Printf("Failed to get message: got insufficient size=%d expectedSize=%d\n", n, dataSize)
	}

	fmt.Printf("SpikeRecoveryRealWorldLocalTest Took %v\n", elapsed)
}
