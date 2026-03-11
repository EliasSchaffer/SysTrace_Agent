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
	queue := make(chan data.WSResponse, 16)

	go a.serverConnector.ReadLoop(func(response data.WSResponse) {
		a.handleServerMessage(response, queue)
	})

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case cmd := <-queue:
			fmt.Printf("Dispatching message type: %s\n", cmd.Type)
			switch cmd.Type {
			case "test":
				if err := beeep.Alert(cmd.Message, "Test", "Test"); err != nil {
					fmt.Println("Failed to show alert:", err)
				}
			case "command":
				handler.HandleCommand(cmd.Message)
			case "config":
				handler.HandleConfig(cmd, &a.serverConnector)
			case "error":
				//TODO: Implement error handling logic
			default:
				fmt.Printf("Unknown message type: %s, Message: %s\n", cmd.Type, cmd.Message)
			}
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
					go a.serverConnector.ReadLoop(func(response data.WSResponse) {
						a.handleServerMessage(response, queue)
					})
				} else {
					fmt.Println("Failed to reconnect after 15 attempts. Giving up.")
					return
				}
			}
		}
	}
}

func (a *Agent) StopAgent() {
	fmt.Println("Agent stopped.")
	a.serverConnector.Close()

}

func (a *Agent) SetNewMasterServer(url string) {
	a.serverConnector.Close()
	a.serverConnector.SetMasterServerURL(url)
	err := a.serverConnector.Connect(context.Background())
	if err != nil {
		fmt.Println("Failed to connect to new master server:", err)
	} else {
		fmt.Println("Successfully connected to new master server:", url)
	}
	envLoader := transport.ENVLoader{}
	envLoader.SetMasterServerURL(url)
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

func (a *Agent) handleServerMessage(response data.WSResponse, queue chan<- data.WSResponse) {
	fmt.Printf("Received message - Type: %s, Message: %s\n", response.Type, response.Message)
	select {
	case queue <- response:
	default:
		fmt.Printf("Incoming queue full, dropping message type=%s message=%s\n", response.Type, response.Message)
	}
}
