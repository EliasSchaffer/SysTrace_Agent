package handler

import (
	"SysTrace_Agent/internal/transport"
	"context"
	"fmt"
)

func handleConfig(message string) {
	switch message {
	case "getConfig":
		//TODO: return config
	case "setMasterServer":
		//TODO: set master server
	case "setDeviceName":
		//TODO: set device name

	}

}

func setMasterServer(message string) {
	serverConnector := transport.NewServerConnector("-1")
	serverConnector.SetMasterServerURL(message)
	canConnect := serverConnector.TestConnection(context.Background())
	if !canConnect {
		fmt.Println("Connection test failed for new master server URL:", message)
		fmt.Println("Reverting to previous master server URL")
		return
	}

}
