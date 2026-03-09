package collector

import (
	"SysTrace_Agent/internal/data"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
)

type CPUCollector struct{}

func (C CPUCollector) Collect() data.Data {
	cpuPercent, _ := cpu.Percent(time.Second, false)
	cpuInfo, _ := cpu.Info()
	cpuCounts, _ := cpu.Counts(true)
	cpuCountsLogical, _ := cpu.Counts(false)

	cpuData := data.CPU{}
	if len(cpuPercent) > 0 {
		cpuData.SetUsage(cpuPercent[0])
	}
	if len(cpuInfo) > 0 {
		cpuData.SetModel(cpuInfo[0].ModelName)
	}
	cpuData.SetCores(cpuCounts)
	cpuData.SetThreads(cpuCountsLogical)

	return cpuData
}
