package agent

import (
	"SysTrace_Agent/internal/collector"
	"SysTrace_Agent/internal/data"
	"SysTrace_Agent/internal/handler"
	"SysTrace_Agent/internal/transport"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gen2brain/beeep"
)

type Agent struct {
	device          data.Device
	serverConnector transport.ServerConnector
}

func (a *Agent) StartAgent() {
	fmt.Println("Agent started...")
	a.CollectData()
	clientID := a.device.GetID()

	if clientID == "" {
		fmt.Println("Warning: ClientID is empty, using hostname as fallback")
		clientID = a.device.GetHostname()
	}

	fmt.Printf("Initializing agent with ClientID: %s\n", clientID)
	a.serverConnector = *transport.NewServerConnector(clientID)

	if err := a.serverConnector.Connect(context.Background()); err != nil {
		fmt.Println("Failed to connect to master server:", err)
		return
	}

	defer a.serverConnector.Close()

	// Commands from the socket loop are executed on this goroutine to keep UI calls thread-safe.
	commandQueue := make(chan string, 16)

	go a.serverConnector.ReadLoop(func(response data.WSResponse) {
		a.handleServerMessage(response, commandQueue)
	})

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case cmd := <-commandQueue:
			handler.HandleCommand(cmd)
		case <-ticker.C:
			a.CollectData()
			wsEvent := data.WSEvent{
				Type:   "update",
				Device: a.device,
			}
			payload, err := json.Marshal(wsEvent)
			if err != nil {
				fmt.Println("Error marshaling data:", err)
				continue
			}
			if err := a.serverConnector.Send(payload); err != nil {
				fmt.Println("Error sending data to master server:", err)
				fmt.Println("Trying to reconnect...")

				reconnected := false
				for i := 0; i < 15; i++ {
					a.serverConnector.Close()
					if err := a.serverConnector.Connect(context.Background()); err != nil {
						fmt.Printf("Reconnection attempt %d out of 15 failed: %v\n", i+1, err)

						delay := time.Duration(1<<uint(i)) * time.Second
						if delay > 30*time.Second {
							delay = 30 * time.Second
						}

						fmt.Printf("Waiting %v before next attempt...\n", delay)
						time.Sleep(delay)
						continue
					}

					reconnected = true
					break
				}

				if reconnected {
					fmt.Println("Reconnected successfully!")
				} else {
					fmt.Println("Failed to reconnect after 15 attempts. Giving up.")
					return
				}
			}
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

func (a *Agent) handleServerMessage(response data.WSResponse, commandQueue chan<- string) {
	fmt.Printf("Received message - Type: %s, Status: %s\n", response.Type, response.Status)

	switch response.Type {
	case "test":
		if err := beeep.Alert(response.Message, "Test", "Test"); err != nil {
			fmt.Println("Failed to show alert:", err)
		}
	case "command":
		select {
		case commandQueue <- response.Message:
		default:
			fmt.Println("Command queue is full, dropping command:", response.Message)
		}
	case "config":
		//TODO: Implement config handling logic
	case "error":
		//TODO: Implement error handling logic
	default:
		fmt.Printf("Unknown message type: %s, Message: %s\n", response.Type, response.Message)
	}
}
