package collector

import (
	"SysTrace_Agent/internal/data/static"
	"net"

	"github.com/shirou/gopsutil/v3/host"
)

type HardwareCollector struct {
}

func (h HardwareCollector) Collect() static.Data {
	hostInfo, _ := host.Info()
	device := static.Device{}
	device.SetHostname(hostInfo.Hostname)
	device.SetOS(hostInfo.OS + " " + hostInfo.Platform + " " + hostInfo.PlatformVersion)
	device.SetID(hostInfo.HostID)
	device.SetIP(collectIPAddress())
	collectHardwareData(&device)

	return device
}

func collectIPAddress() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.String()

}

func collectHardwareData(device *static.Device) *static.Device {
	cpuData := CPUCollector{}.Collect()
	cpu, ok := cpuData.(static.CPU)
	if !ok {
		panic("Failed to collect CPU data")
	}
	memData := MemoryCollector{}.Collect()
	memory, ok := memData.(static.Memory)
	if !ok {
		panic("Failed to collect Memory data")
	}

	hardware := static.Hardware{}
	hardware.SetCPU(cpu)
	hardware.SetMemory(memory)

	device.SetHardware(hardware)
	return device
}
