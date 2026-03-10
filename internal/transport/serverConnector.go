package transport

import (
	"SysTrace_Agent/internal/data"
	"context"
	"encoding/json"
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
	clientID        string

	conn    *websocket.Conn
	writeMu sync.Mutex
}

func NewServerConnector(clientID string) *ServerConnector {
	envLoader := ENVLoader{}
	masterServerURL := envLoader.GetMasterServerURL()

	return &ServerConnector{
		masterServerURL: masterServerURL,
		clientID:        clientID,
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

func toWSURL(raw, clientID string) string {
	raw = strings.TrimRight(raw, "/")
	u, err := url.Parse(raw)
	if err != nil || u.Scheme == "" {
		if !strings.HasPrefix(raw, "ws://") && !strings.HasPrefix(raw, "wss://") {
			raw = "ws://" + raw
		}
		u, _ = url.Parse(raw)
	}

	switch u.Scheme {
	case "http":
		u.Scheme = "ws"
	case "https":
		u.Scheme = "wss"
	}

	u.Path = "/ws"
	q := u.Query()
	q.Set("clientId", clientID)
	u.RawQuery = q.Encode()

	return u.String()
}

func (s *ServerConnector) Connect(ctx context.Context) error {
	wsURL := toWSURL(s.masterServerURL, s.clientID)

	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	conn, resp, err := dialer.DialContext(ctx, wsURL, nil)
	if err != nil {
		if resp != nil {
			return fmt.Errorf("ws connection failed with status %d: %v", resp.StatusCode, err)
		}
		return fmt.Errorf("ws connection failed: %v", err)
	}

	s.conn = conn
	fmt.Printf("WebSocket connection established with master server (ClientID: %s)\n", s.clientID)
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

func (s *ServerConnector) ReadLoop(onResponse func(resp data.WSResponse)) {
	if s.conn == nil {
		return
	}

	for {
		mt, msg, err := s.conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			return
		}
		if mt != websocket.TextMessage {
			continue
		}

		var resp data.WSResponse
		if err := json.Unmarshal(msg, &resp); err != nil {
			fmt.Println("Invalid WSResponse:", err, "raw:", string(msg))
			continue
		}

		if onResponse != nil {
			onResponse(resp)
		}
	}
}

func (s *ServerConnector) Close() error {
	if s.conn == nil {
		return nil
	}

	_ = s.conn.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		time.Now().Add(2*time.Second),
	)

	err := s.conn.Close()
	s.conn = nil
	return err
}
