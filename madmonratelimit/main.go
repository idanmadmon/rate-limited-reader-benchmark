package main

import (
	"bytes"
	"fmt"
	"io"
	"time"

	ratelimitedreader "github.com/idanmadmon/rate-limited-reader"
)

func main() {
	dataSize := 32 * 1024 * 1024
	reader := bytes.NewBuffer(make([]byte, dataSize))

	// allow 1/4 of the data size per second,
	// should take 32 / (32/4) = 4 seconds
	// reads interval divided evenly by ratelimitedreader.ReadIntervalMilliseconds
	limitedReader := ratelimitedreader.NewRateLimitedReader(reader, int64(dataSize/4))

	var total int
	buffer := make([]byte, 1024)
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
	fmt.Printf("Total: %d, Elapsed: %s\n", total, elapsed)
}
