package services

import (
	"SysTrace_Agent/data"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
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
	gpsData := a.GetGPSDataByLocationAPI()
	if gpsData != nil {
		a.device.SetGPS(*gpsData)
	} else {
		fmt.Println("Failed to get GPS data from Location API, trying IP-based geolocation...")
		gpsData = a.GetGPSDataByIP()
		if gpsData != nil {
			fmt.Println("")
			a.device.SetGPS(*gpsData)
		} else {
			fmt.Println("Failed to get GPS data from both methods.")
		}
	}

}

func (a *Agent) GetGPSDataByIP() *data.GPS {
	apiKey := a.envLoader.GetGeoLocationAPIKey()
	url := fmt.Sprintf("https://api.ipgeolocation.io/ipgeo?apiKey=%s", apiKey)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching GPS data: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Printf("Error decoding GPS data: %v\n", err)
		return nil
	}

	gps := &data.GPS{}

	if city, ok := result["city"].(string); ok {
		gps.SetCity(city)
	}
	if region, ok := result["state_prov"].(string); ok {
		gps.SetRegion(region)
	}
	if country, ok := result["country_name"].(string); ok {
		gps.SetCountry(country)
	}

	if latStr, ok := result["latitude"].(string); ok {
		if lat, err := strconv.ParseFloat(latStr, 64); err == nil {
			gps.SetLatitude(lat)
		}
	} else if latFloat, ok := result["latitude"].(float64); ok {
		gps.SetLatitude(latFloat)
	}

	if lonStr, ok := result["longitude"].(string); ok {
		if lon, err := strconv.ParseFloat(lonStr, 64); err == nil {
			gps.SetLongitude(lon)
		}
	} else if lonFloat, ok := result["longitude"].(float64); ok {
		gps.SetLongitude(lonFloat)
	}

	return gps
}

// Installationspfad der MSIX-App verwenden
func getGpsHelperPath() string {
	// Direkt den WindowsApps Pfad verwenden
	cmd := exec.Command("powershell", "-Command",
		"(Get-AppxPackage -Name 'GpsHelper').InstallLocation")
	out, err := cmd.Output()
	if err != nil || strings.TrimSpace(string(out)) == "" {
		return ""
	}
	return filepath.Join(strings.TrimSpace(string(out)), "gpshelper.exe")
}

func (a *Agent) GetGPSDataByLocationAPI() *data.GPS {
	gpsHelperPath := getGpsHelperPath()
	if gpsHelperPath == "" {
		fmt.Println("Error: GpsHelper App nicht installiert!")
		return nil
	}

	cmd := exec.Command(gpsHelperPath)
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error executing GPS helper: %v\n", err)
		return nil
	}

	gps := &data.GPS{}
	if err = json.Unmarshal(output, gps); err != nil {
		fmt.Printf("Error unmarshaling GPS data: %v\n", err)
		return nil
	}

	//	if err = enrichGPSData(gps); err != nil {
	//		fmt.Printf("Warning: Could not enrich GPS data: %v\n", err)
	//	}

	return gps
}
