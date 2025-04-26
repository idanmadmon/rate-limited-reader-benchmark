package main

import (
	"context"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go monitorLoop(ctx)
	time.Sleep(500 * time.Millisecond)

	// readerFactory := IdanMadmonRateLimitReaderFactory
	readerFactory := NoLimitReaderFactory

	runTest(RateLimitBasicFunctionalityTest, readerFactory)
	runTest(TestReaderBehavior2, readerFactory)
	runTest(TestReaderBehavior3, readerFactory)
}

func monitorLoop(ctx context.Context) {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	var prevRx uint64

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Monitor loop stopped.")
			return
		case <-ticker.C:
			ioCounters, err := net.IOCounters(false)
			if err != nil || len(ioCounters) == 0 {
				fmt.Printf("Error reading RX bytes: %v\n", err)
				continue
			}
			currRx := ioCounters[0].BytesRecv
			rxDelta := currRx - prevRx
			prevRx = currRx

			cpuPercent, err := cpu.Percent(0, false)
			if err != nil || len(cpuPercent) == 0 {
				fmt.Printf("Error reading CPU usage: %v\n", err)
				continue
			}

			vmStat, err := mem.VirtualMemory()
			if err != nil {
				fmt.Printf("Error reading memory usage: %v\n", err)
				continue
			}

			fmt.Printf("RX: %d bytes | CPU: %.2f%% | RAM: %.2fMB\n",
				rxDelta, cpuPercent[0], float64(vmStat.Used)/1024.0/1024.0)
		}
	}
}

func runTest(testFn BenchmarkTest, factory ReaderFactory) {
	testName := strings.TrimPrefix(filepath.Ext(runtime.FuncForPC(reflect.ValueOf(testFn).Pointer()).Name()), ".")
	fmt.Printf("Starting %s...\n", testName)
	testFn(factory)
	fmt.Printf("Finished %s\n", testName)
	time.Sleep(250 * time.Millisecond)
}
