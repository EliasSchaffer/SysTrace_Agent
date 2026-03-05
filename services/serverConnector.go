package services

import (
	"bytes"
	"fmt"
	"net/http"
)

type ServerConnector struct {
	masterServerURL string
}

func NewServerConnector() *ServerConnector {
	envLoader := ENVLoader{}
	masterServerURL := envLoader.GetMasterServerURL()

	return &ServerConnector{
		masterServerURL: masterServerURL,
	}
}

func (s *ServerConnector) TestConnection() bool {
	resp, err := http.Get(s.masterServerURL + "/status")
	if err != nil {
		fmt.Println("Connection error:", err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func (s *ServerConnector) SendDataToMasterServer(data []byte, err error) {

	if err != nil {
		fmt.Println("Error preparing data:", err)
		return
	}

	resp, err := http.Post(
		s.masterServerURL+"/metrics",
		"application/json",
		bytes.NewBuffer(data),
	)

	if err != nil {
		fmt.Println("Send error:", err)
		return
	}

	defer resp.Body.Close()

}
