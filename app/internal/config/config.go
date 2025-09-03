package config

import (
	"os"
	"time"
)

type Config struct {
	Server    ServerConfig    `json:"server"`
	Webhook   WebhookConfig   `json:"webhook"`
	WebSocket WebSocketConfig `json:"websocket"`
}

type ServerConfig struct {
	ListenAddress string        `json:"listen_address"`
	ReadTimeout   time.Duration `json:"read_timeout"`
	WriteTimeout  time.Duration `json:"write_timeout"`
}

type WebhookConfig struct {
	OnConnectURL    string        `json:"on_connect_url"`
	OnDisconnectURL string        `json:"on_disconnect_url"`
	Timeout         time.Duration `json:"timeout"`
}

type WebSocketConfig struct {
	ReadBufferSize  int `json:"read_buffer_size"`
	WriteBufferSize int `json:"write_buffer_size"`
	CheckOrigin     bool `json:"check_origin"`
}

func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			ListenAddress: getEnvOrDefault("LISTEN_ADDR", ":8080"),
			ReadTimeout:   5 * time.Second,
			WriteTimeout:  10 * time.Second,
		},
		Webhook: WebhookConfig{
			OnConnectURL:    os.Getenv("ONCONNECT_URL"),
			OnDisconnectURL: os.Getenv("DISCONNECT_URL"),
			Timeout:         5 * time.Second,
		},
		WebSocket: WebSocketConfig{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     true,
		},
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 