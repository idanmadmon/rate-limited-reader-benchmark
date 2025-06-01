package main

import (
	"context"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go monitorLoop(ctx)
	time.Sleep(500 * time.Millisecond)

	readerFactory := IdanMadmonDeterministicRateLimitReaderFactory
	// readerFactory := GolangBurstsRateLimitReaderFactory
	// readerFactory := JujuBurstsRateLimitReaderFactory
	// readerFactory := UberDeterministicRateLimitReaderFactory
	// readerFactory := NoLimitReaderFactory

	// runTest(RateLimitingSyntheticTest, readerFactory)
	// runTest(RateLimitingRealWorldLocalTest, readerFactory)
	// runTest(MaxReadOverTimeSyntheticTest, readerFactory)
	runTest(SpikeRecoveryRealWorldLocalTest, readerFactory)
}

func runTest(testFn BenchmarkTest, factory ReaderFactory) {
	testName := strings.TrimPrefix(filepath.Ext(runtime.FuncForPC(reflect.ValueOf(testFn).Pointer()).Name()), ".")
	fmt.Printf("Starting %s...\n", testName)
	testFn(factory)
	time.Sleep(250 * time.Millisecond)
	fmt.Printf("Finished %s\n", testName)
	time.Sleep(250 * time.Millisecond)
}
