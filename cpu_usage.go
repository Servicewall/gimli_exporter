package main

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

var (
	gimliProcCpuUsageMap = sync.Map{}
)

func init() {
	go func() {
		updateGimliCpuUsage()
		ticker := time.NewTicker(10 * time.Second)
		for range ticker.C {
			updateGimliCpuUsage()
		}
	}()
}

func updateGimliCpuUsage() {
	processes, err := process.Processes()
	if err != nil {
		log.Printf("Error fetching processes: %v\n", err)
		return
	}
	for _, p := range processes {
		name, err := p.Name()
		if err != nil {
			continue
		}
		switch name {
		case "gimli":
			// linux 探针
			go func() {
				percent, err := p.Percent(1 * time.Second)
				if err != nil {
					log.Printf("Error get gimli(%d) process percent: %v\n", p.Pid, err)
					return
				}
				if percent > 0 {
					gimliProcCpuUsageMap.Store(p.Pid, percent)
				}
			}()
		case "gimli.exe":
			// windows 探针，cpu 使用率需要除以cpu核数
			go func() {
				percent, err := p.Percent(1 * time.Second)
				if err != nil {
					log.Printf("Error get gimli(%d) process percent: %v\n", p.Pid, err)
					return
				}
				if percent > 0 {
					gimliProcCpuUsageMap.Store(p.Pid, (percent / float64(runtime.NumCPU())))
				}
			}()
		}
	}
}

func getCPUPercent(pid int32) (float64, error) {
	cpuUsageVal, ok := gimliProcCpuUsageMap.Load(pid)
	if !ok {
		return 0, fmt.Errorf("cpu usage not exists")
	}
	cpuUsage, ok := cpuUsageVal.(float64)
	if !ok {
		return 0, fmt.Errorf("cpu usage type is not float64")
	}
	return cpuUsage, nil
}
