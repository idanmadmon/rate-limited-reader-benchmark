package main

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/juju/ratelimit"
)

func main() {
	const dataSize = 32 * 1024
	reader := bytes.NewBuffer(make([]byte, dataSize))

	// limit := dataSize / 4
	limit := dataSize / 5
	limiter := ratelimit.NewBucketWithRate(float64(limit), int64(limit))

	var total int
	buffer := make([]byte, 1024)
	start := time.Now()
	for {
		n, err := reader.Read(buffer)
		total += n
		limiter.Wait(int64(n))
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
