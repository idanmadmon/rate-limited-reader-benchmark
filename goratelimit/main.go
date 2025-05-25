package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"golang.org/x/time/rate"
)

func main() {
	const dataSize = 32 * 1024 * 1024
	reader := bytes.NewBuffer(make([]byte, dataSize))

	// 8 events per second, burst of 1 (to see the limitation)
	// should take 32 / 8 = 4 seconds
	// put 33 burst to see the burst in action
	limiter := rate.NewLimiter(rate.Every(time.Second/8), 1)

	var total int
	buffer := make([]byte, 1024*1024)
	start := time.Now()
	for {
		err := limiter.Wait(context.Background())
		if err != nil {
			fmt.Printf("limiter.Wait err: %v\n", err)
		}

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
