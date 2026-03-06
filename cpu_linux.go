//go:build !windows
package main

import (
	"github.com/shirou/gopsutil/v3/process"
)

func getCPUPercent(p *process.Process) (float64, error) {
	// Linux/Other: Use CPUPercent (average since process creation)
	return p.CPUPercent()
}
