package handler

import (
	"SysTrace_Agent/internal/data/ws"
	"SysTrace_Agent/internal/transport"
	"context"
	"fmt"
	"strings"
)

func HandleConfig(resp ws.WSRequest, connector *transport.ServerConnector) {
	action := strings.TrimSpace(resp.Payload)
	message := strings.TrimSpace(resp.Message)

	fmt.Printf("HandleConfig called - Payload: %q, Message: %q\n", action, message)
	switch action {
	case "getConfig":
		//TODO: return config
	case "setMasterServer":
		setMasterServer(message, connector)
	case "setDeviceName":
		//TODO: set device name
	default:
		// Fallback: wenn kein Action-Feld gesendet wird, Message als URL interpretieren.
		if message != "" {
			setMasterServer(message, connector)
			return
		}
		fmt.Println("Config message ignored: no action and no message")
	}
}

func setMasterServer(newURL string, connector *transport.ServerConnector) {
	newURL = strings.TrimSpace(newURL)
	if newURL == "" {
		fmt.Println("setMasterServer ignored: empty URL")
		return
	}

	tempConnector := transport.NewServerConnectorWithID(connector.ClientID(), newURL)
	if !tempConnector.TestConnection(context.Background()) {
		fmt.Println("Connection test failed for new master server URL:", newURL)
		fmt.Println("Reverting to previous master server URL")
		return
	}

	oldURL := connector.MasterServerURL()
	_ = connector.Close()
	connector.SetMasterServerURL(newURL)
	if err := connector.Connect(context.Background()); err != nil {
		fmt.Printf("Failed to connect to new master server URL %s: %v\n", newURL, err)
		fmt.Println("Trying to revert to old master server URL:", oldURL)
		connector.SetMasterServerURL(oldURL)
		if err := connector.Connect(context.Background()); err != nil {
			fmt.Println("Failed to revert to old master server URL:", err)
		}
		return
	}

	fmt.Println("Master server URL changed to:", newURL)
}
