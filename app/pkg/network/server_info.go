package network

import (
	"net"
	"strings"
)

type ServerInfo struct {
	IP   string
	Port string
}

func GetServerInfo(listenAddress string) *ServerInfo {
	serverIP := getOutboundIP()
	serverPort := extractPort(listenAddress)
	
	return &ServerInfo{
		IP:   serverIP,
		Port: serverPort,
	}
}

func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "unknown"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func extractPort(listenAddress string) string {
	if strings.HasPrefix(listenAddress, ":") {
		return listenAddress[1:]
	}
	
	_, port, err := net.SplitHostPort(listenAddress)
	if err != nil {
		return "unknown"
	}
	
	return port
} 