package server

import (
	"net/http"

	"gomw-gw/app/internal/handlers"
	"gomw-gw/app/pkg/logger"
)

type Router struct {
	mux               *http.ServeMux
	websocketHandler  *handlers.WebSocketHandler
	messageHandler    *handlers.MessageHandler
	infoHandler       *handlers.InfoHandler
}

func NewRouter(
	wsHandler *handlers.WebSocketHandler,
	msgHandler *handlers.MessageHandler,
	infoHandler *handlers.InfoHandler,
) *Router {
	return &Router{
		mux:              http.NewServeMux(),
		websocketHandler: wsHandler,
		messageHandler:   msgHandler,
		infoHandler:      infoHandler,
	}
}

func (r *Router) SetupRoutes() {
	r.mux.HandleFunc("/ws", r.websocketHandler.HandleConnection)
	r.mux.HandleFunc("/send", r.messageHandler.HandleSendMessage)
	r.mux.HandleFunc("/env", r.infoHandler.HandleEnvironmentInfo)
	r.mux.HandleFunc("/health", r.infoHandler.HandleHealthCheck)
	r.mux.HandleFunc("/status", r.infoHandler.HandleConnectionStatus)

	logger.Info("Routes configured", logger.Fields{
		"routes": []string{"/ws", "/send", "/env", "/health", "/status"},
	})
}

func (r *Router) GetHandler() http.Handler {
	return r.mux
} 