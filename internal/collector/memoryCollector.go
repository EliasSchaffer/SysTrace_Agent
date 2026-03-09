package collector

import (
	"SysTrace_Agent/internal/data"

	"github.com/shirou/gopsutil/v3/mem"
)

type MemoryCollector struct {
}

func (m MemoryCollector) Collect() data.Data {
	memData := data.Memory{}
	memStats, _ := mem.VirtualMemory()
	memData.SetTotal(memStats.Total)
	memData.SetUsed(memStats.Used)
	memData.SetAvailable(memStats.Available)
	memData.SetUsedPercent(memStats.UsedPercent)
	return memData
}
