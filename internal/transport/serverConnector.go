package transport

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type ServerConnector struct {
	masterServerURL string
	wsURL           string

	conn    *websocket.Conn
	writeMu sync.Mutex
}

func NewServerConnector() *ServerConnector {
	envLoader := ENVLoader{}
	masterServerURL := envLoader.GetMasterServerURL()

	return &ServerConnector{
		masterServerURL: masterServerURL,
		wsURL:           toWSURL(masterServerURL),
	}
}

func (s *ServerConnector) TestConnection(ctx context.Context) bool {
	if err := s.Connect(ctx); err != nil {
		fmt.Println("Connection test failed:", err)
		return false
	}
	_ = s.Close()
	return true
}

func toWSURL(raw string) string {
	raw = strings.TrimRight(raw, "/")
	u, err := url.Parse(raw)
	if err != nil || u.Scheme == "" {
		// Fallback, falls nur host:port in .env steht
		if strings.HasPrefix(raw, "ws://") || strings.HasPrefix(raw, "wss://") {
			return raw
		}
		return "ws://" + raw
	}

	switch u.Scheme {
	case "http":
		u.Scheme = "ws"
	case "https":
		u.Scheme = "wss"
	}
	return u.String()
}

func (s *ServerConnector) Connect(ctx context.Context) error {
	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	conn, resp, err := dialer.DialContext(ctx, s.wsURL, nil)
	if err != nil {
		if resp != nil {
			return fmt.Errorf("ws connection failed with status %d: %v", resp.StatusCode, err)
		}
		return fmt.Errorf("ws connection failed: %v", err)
	}

	s.conn = conn
	fmt.Println("WebSocket connection established with master server")
	return nil
}

func (s *ServerConnector) Send(data []byte) error {

	if s.conn == nil {
		return errors.New("WebSocket connection not established")
	}

	s.writeMu.Lock()
	defer s.writeMu.Unlock()

	if wErr := s.conn.WriteMessage(websocket.TextMessage, data); wErr != nil {
		return wErr
	}
	return nil
}

func (s *ServerConnector) ReadLoop(onMessage func(messageType int, payload []byte)) {
	if s.conn == nil {
		return
	}
	for {
		mt, msg, err := s.conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			return
		}
		if onMessage != nil {
			onMessage(mt, msg)
		}

	}
}

func (s *ServerConnector) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}

	err := s.conn.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		time.Now().Add(2*time.Second),
	)
	_ = err
	cErr := s.conn.Close()
	s.conn = nil
	return cErr
}
