package main

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

var SyntheticRXBytes atomic.Uint64

func monitorLoop(ctx context.Context) {
	SyntheticRXBytes.Store(0) // reset for monitor
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	var prevRx uint64
	var prevSyntheticRx uint64

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

			currSyntheticRx := SyntheticRXBytes.Load()
			syntheticRxDelta := currSyntheticRx - prevSyntheticRx
			prevSyntheticRx = currSyntheticRx

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

			fmt.Printf("RX: %d bytes |CPU: %.2f%% | RAM: %.2fMB | SyntheticRX: %d bytes | TotalSyntheticRX: %d bytes\n",
				rxDelta, cpuPercent[0], float64(vmStat.Used)/1024.0/1024.0, syntheticRxDelta, prevSyntheticRx)
		}
	}
}
