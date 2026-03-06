//go:build windows

package main

import (
	"log"
	"math"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

var (
	// 使用原子 uint64 存储 float64 的位模式
	atomicCachedCPU atomic.Uint64
)

func init() {
	go func() {
		atomicCachedCPU.Store(math.Float64bits(collectTotalGimliCPU()))
		// 1. init() 函数改为 10s Ticker 循环
		ticker := time.NewTicker(10 * time.Second)
		for range ticker.C {
			val := collectTotalGimliCPU()
			// 如果 cpu 的计算值为 0, 则不刷新变量
			if val > 0 {
				atomicCachedCPU.Store(math.Float64bits(val))
			}
		}
	}()
}

// collectTotalGimliCPU 实现获取 PID 数组、遍历并单独计算每个 PID 使用率、最后累加逻辑
func collectTotalGimliCPU() float64 {
	processes, err := process.Processes()
	if err != nil {
		log.Printf("Error fetching processes: %v\n", err)
		return 0
	}
	var totalUsage float64
	for _, p := range processes {
		name, err := p.Name()
		if err != nil {
			continue
		}
		if name == "gimli.exe" {
			percent, err := p.Percent(1 * time.Second)
			if err != nil {
				log.Printf("Error get process percent: %v\n", err)
				continue
			}
			totalUsage += (percent / float64(runtime.NumCPU()))
		}
	}
	// 4. 函数返回值是所有进程id 的cpu使用率的累加值
	return totalUsage
}

func getCPUPercent(p *process.Process) (float64, error) {
	return math.Float64frombits(atomicCachedCPU.Load()), nil
}
