package agent

import (
	"SysTrace_Agent/internal/collector"
	"SysTrace_Agent/internal/data"
	"SysTrace_Agent/internal/transport"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type Agent struct {
	device          data.Device
	serverConnector transport.ServerConnector
}

func (a *Agent) StartAgent() {
	fmt.Println("Agent started...")

	a.serverConnector = *transport.NewServerConnector()

	if err := a.serverConnector.Connect(context.Background()); err != nil {
		fmt.Println("Failed to connect to master server:", err)
		return
	}

	defer a.serverConnector.Close()

	go a.serverConnector.ReadLoop(func(messageType int, msg []byte) {
		fmt.Printf("Received message from master server: %s\n", string(msg))
	})

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		a.CollectData()
		payload, err := json.Marshal(a.device)
		if err != nil {
			fmt.Println("Error marshaling data:", err)
			continue
		}
		if err := a.serverConnector.Send(payload); err != nil {
			fmt.Println("Error sending data to master server:", err)
		}

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
	device, ok := collector.HardwareCollector{}.Collect().(data.Device)
	if !ok {
		panic("Failed to collect hardware data")
	}
	a.device = device
	gpsdata, ok := collector.GPSCollector{}.Collect().(data.GPS)
	if !ok {
		panic("Failed to collect GPS data")
	}
	a.device.SetGPS(gpsdata)
}
