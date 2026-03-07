package services

import (
	"SysTrace_Agent/data"
	"net"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

func (a *Agent) collectDeviceInfo() {
	hostInfo, _ := host.Info()

	a.device.SetHostname(hostInfo.Hostname)
	a.device.SetOS(hostInfo.OS + " " + hostInfo.Platform + " " + hostInfo.PlatformVersion)
	a.device.SetID(hostInfo.HostID)
	a.device.SetIP(a.collectIPAddress())
	a.collectHardwareData()
}

func (a *Agent) collectIPAddress() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.String()

}

func (a *Agent) collectHardwareData() {
	cpuData := a.getCPUInfo()
	memData := a.getMemoryInfo()

	hardware := data.Hardware{}
	hardware.SetCPU(cpuData)
	hardware.SetMemory(memData)

	a.device.SetHardware(hardware)
}

func (a *Agent) getCPUInfo() data.CPU {
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

func (a *Agent) getMemoryInfo() data.Memory {
	memData := data.Memory{}
	memStats, _ := mem.VirtualMemory()
	memData.SetTotal(memStats.Total)
	memData.SetUsed(memStats.Used)
	memData.SetAvailable(memStats.Available)
	memData.SetUsedPercent(memStats.UsedPercent)
	return memData
}
