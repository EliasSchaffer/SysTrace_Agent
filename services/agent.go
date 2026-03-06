package services

import (
	"SysTrace_Agent/data"
	"encoding/json"
	"fmt"
	"time"
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
	a.CollectGPSData()
}
