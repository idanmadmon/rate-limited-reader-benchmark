package main

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"go.uber.org/ratelimit"
)

func main() {
	const dataSize = 32 * 1024 * 1024
	reader := bytes.NewBuffer(make([]byte, dataSize))

	// 8 events per second, 125ms between calls
	// should take 32 / 8 = 32 * 125ms = 4 seconds
	limiter := ratelimit.New(8 * 1024)

	var total int
	buffer := make([]byte, 1024)
	start := time.Now()
	for {
		limiter.Take()

		n, err := reader.Read(buffer)
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
