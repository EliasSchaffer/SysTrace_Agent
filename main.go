package main

import (
	"SysTrace_Agent/internal/agent"
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

func main() {
	cpuPercent, _ := cpu.Percent(time.Second, false)

	memStats, _ := mem.VirtualMemory()

	fmt.Println("CPU:", cpuPercent[0], "%")
	fmt.Println("RAM:", memStats.UsedPercent, "%")

	a := agent.Agent{}
	a.StartAgent()
}
