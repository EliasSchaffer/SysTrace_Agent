package collector

import (
	"SysTrace_Agent/internal/data/static"

	"github.com/shirou/gopsutil/v3/mem"
)

type MemoryCollector struct {
}

func (m MemoryCollector) Collect() static.Data {
	memData := static.Memory{}
	memStats, _ := mem.VirtualMemory()
	memData.SetTotal(memStats.Total)
	memData.SetUsed(memStats.Used)
	memData.SetAvailable(memStats.Available)
	memData.SetUsedPercent(memStats.UsedPercent)
	return memData
}
