package models

import (
	"encoding/json"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type ConnectionID string

type Session struct {
	ID         ConnectionID    `json:"id"`
	Connection *websocket.Conn `json:"-"`
	ClientIP   string          `json:"client_ip"`
	QueryParams url.Values     `json:"query_params"`
	ConnectedAt time.Time      `json:"connected_at"`
}

type SendMessageRequest struct {
	ConnectionID ConnectionID    `json:"connection_id"`
	Message      json.RawMessage `json:"message"`
}

type WebhookPayload struct {
	ConnectionID ConnectionID `json:"connection_id"`
	ClientIP     string       `json:"client_ip"`
	QueryParams  url.Values   `json:"query_params,omitempty"`
	Timestamp    time.Time    `json:"timestamp"`
	ServerIP     string       `json:"server_ip"`
	ServerPort   string       `json:"server_port"`
}

type EnvironmentInfo struct {
	ListenAddress   string `json:"listen_address"`
	OnConnectURL    string `json:"on_connect_url"`
	OnDisconnectURL string `json:"on_disconnect_url"`
}

func (s *Session) IsValid() bool {
	return s.Connection != nil && s.ID != ""
}

func (s *Session) Close() error {
	if s.Connection != nil {
		return s.Connection.Close()
	}
	return nil
} 