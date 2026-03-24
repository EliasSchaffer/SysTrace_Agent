package agent

import (
	"SysTrace_Agent/internal/collector"
	"SysTrace_Agent/internal/data/static"
	"SysTrace_Agent/internal/data/ws"
	"SysTrace_Agent/internal/handler"
	"SysTrace_Agent/internal/transport"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gen2brain/beeep"
)

type Agent struct {
	device             static.Device
	serverConnector    transport.ServerConnector
	lastSettingsChange time.Time
}

func (a *Agent) StartAgent() {
	settingsHandler := handler.SettingsHandler{}
	a.writeLog("Agent started")
	a.CollectData()
	clientID := a.device.GetID()

	if clientID == "" {
		a.writeWarn("Device ID is empty, falling back to hostname as ClientID")
		clientID = a.device.GetHostname()
	}

	a.writeLog(fmt.Sprintf("Initializing agent with ClientID: %s", clientID))
	a.serverConnector = *transport.NewServerConnector(clientID)

	if err := a.serverConnector.Connect(context.Background()); err != nil {
		a.writeError(fmt.Sprintf("Failed to connect to master server: %v", err))
		return
	}

	a.writeLog("Successfully connected to master server")
	defer a.serverConnector.Close()
	queue := make(chan ws.WSRequest, 16)

	go a.serverConnector.ReadLoop(func(request ws.WSRequest) {
		a.handleServerMessage(request, queue)
	})

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		fmt.Println("Check for settings change")
		info, err := os.Stat(".env")
		if err != nil {
			fmt.Println("Fehler:", err)
			return
		}
		if info.ModTime().After(a.lastSettingsChange) {
			fmt.Println("Settings file changed, reloading...")
			a.lastSettingsChange = info.ModTime()
			env := transport.ENVLoader{}
			settings := env.GetSettings()
			go settingsHandler.HandleSettingsChange(settings, a.SetNewMasterServer, a.changeStaticGPSData)

		} else {
			fmt.Println("No Changes found")
		}
		select {
		case cmd := <-queue:
			a.writeLog(fmt.Sprintf("Dispatching message type: %s\n", cmd.Type))
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
				a.writeWarn(fmt.Sprintf("Received unknown message type: %s", cmd.Type))
			}
		case <-ticker.C:
			a.CollectData()
			wsEvent := ws.WSEvent{
				Type:   "update",
				Device: a.device,
			}
			payload, err := json.Marshal(wsEvent)
			if err != nil {
				a.writeError(fmt.Sprintf("Error marshaling data: %v", err))
				continue
			}
			if err := a.serverConnector.Send(payload); err != nil {
				a.writeError(fmt.Sprintf("Error sending data to master server: %v", err))
				a.writeLog("Attempting to reconnect to master server...")

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
					a.writeLog("Reconnected to master server successfully.")
					go a.serverConnector.ReadLoop(func(request ws.WSRequest) {
						a.handleServerMessage(request, queue)
					})
				} else {
					a.writeLog("Failed to reconnect to master server after multiple attempts, exiting agent.")
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

func (a *Agent) changeStaticGPSData(gps static.GPS) {
	collector.GPSCollector{}.SetStaticGPSData(gps)
}

func (a *Agent) changeSendStaticGPS(send bool) {
	collector.GPSCollector{}.SetSendStaticGPS(send)
}

func (a *Agent) changeSendGPS(send bool) {
	collector.GPSCollector{}.SetSendGPS(send)
}

func (a *Agent) CollectData() {
	device, ok := collector.HardwareCollector{}.Collect().(static.Device)
	if !ok {
		panic("Failed to collect hardware data")
	}
	a.device = device
	gpsdata, ok := collector.GPSCollector{}.Collect().(static.GPS)
	if !ok {
		panic("Failed to collect GPS data")
	}
	a.device.SetGPS(gpsdata)
}

func (a *Agent) handleServerMessage(request ws.WSRequest, queue chan<- ws.WSRequest) {
	fmt.Printf("Received message - Type: %s, Message: %s\n", request.Type, request.Message)
	select {
	case queue <- request:
	default:
		fmt.Printf("Incoming queue full, dropping message type=%s message=%s\n", request.Type, request.Message)
		response := ws.WSResponse{
			"response",
			request.RequestID,
			503,
		}
		go a.serverConnector.SendResponse(response)
	}
}
