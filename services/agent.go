package services

import (
	"SysTrace_Agent/data"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

type Agent struct {
	device          data.Device
	envLoader       ENVLoader
	serverConnector ServerConnector
}

func (a *Agent) StartAgent() {
	fmt.Println("Agent gestartet...")

	a.serverConnector = *NewServerConnector()
	if !a.serverConnector.TestConnection() {
		fmt.Println("No connection to master server. Please check the URL and try again.")
		return
	}

	a.envLoader = ENVLoader{}
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		a.CollectData()
		a.printStats()
		go a.serverConnector.SendDataToMasterServer(json.Marshal(a.device))
	}
}

func (a *Agent) collectDeviceInfo() {
	hostInfo, _ := host.Info()

	a.device.SetHostname(hostInfo.Hostname)
	a.device.SetOS(hostInfo.OS + " " + hostInfo.Platform + " " + hostInfo.PlatformVersion)
	a.device.SetID(hostInfo.HostID)
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

func (a *Agent) printStats() {
	fmt.Println("=====================================")
	fmt.Printf("Hostname: %s\n", a.device.GetHostname())
	fmt.Printf("OS: %s\n", a.device.GetOS())
	fmt.Printf("CPU: %.2f%%", a.device.GetHardware().GetCPU().GetUsage())
	fmt.Printf("RAM: %.2f%% (Used: %d MB / Total: %d MB)\n",
		a.device.GetHardware().GetMemory().GetUsedPercent(),
		a.device.GetHardware().GetMemory().GetUsed()/1024/1024,
		a.device.GetHardware().GetMemory().GetTotal()/1024/1024)
	fmt.Printf("GPS - City: %s, Region: %s, Country: %s\n",
		a.device.GetGPS().GetCity(),
		a.device.GetGPS().GetRegion(),
		a.device.GetGPS().GetCountry())
	fmt.Printf("GPS - Latitude: %.4f, Longitude: %.4f\n",
		a.device.GetGPS().GetLatitude(),
		a.device.GetGPS().GetLongitude())
	fmt.Println("=====================================")
}

func (a *Agent) CollectData() {
	a.collectDeviceInfo()
	a.collectHardwareData()
	a.CollectGPSData()
}

func (a *Agent) CollectGPSData() {
	apiKey := a.envLoader.GetGeoLocationAPIKey()
	url := fmt.Sprintf("https://api.ipgeolocation.io/ipgeo?apiKey=%s", apiKey)

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	gps := &data.GPS{}
	gps.SetCity(fmt.Sprintf("%v", result["city"]))
	gps.SetRegion(fmt.Sprintf("%v", result["state_prov"]))
	gps.SetCountry(fmt.Sprintf("%v", result["country_name"]))

	lat, _ := strconv.ParseFloat(fmt.Sprintf("%v", result["latitude"]), 64)
	gps.SetLatitude(lat)
	lon, _ := strconv.ParseFloat(fmt.Sprintf("%v", result["longitude"]), 64)
	gps.SetLongitude(lon)

	a.device.SetGPS(*gps)
}
