package handlers

import (
	"net"
	"net/http"
	"time"

	"gomw-gw/app/internal/config"
	"gomw-gw/app/internal/models"
	"gomw-gw/app/internal/services"
	"gomw-gw/app/pkg/logger"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	upgrader       websocket.Upgrader
	sessionManager *services.SessionManager
	webhookService *services.WebhookService
}

func NewWebSocketHandler(
	cfg *config.WebSocketConfig,
	sessionManager *services.SessionManager,
	webhookService *services.WebhookService,
) *WebSocketHandler {
	return &WebSocketHandler{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  cfg.ReadBufferSize,
			WriteBufferSize: cfg.WriteBufferSize,
			CheckOrigin: func(r *http.Request) bool {
				return cfg.CheckOrigin
			},
		},
		sessionManager: sessionManager,
		webhookService: webhookService,
	}
}

func (h *WebSocketHandler) HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Warn("WebSocket upgrade failed", logger.Fields{
			"error":      err.Error(),
			"remote_addr": r.RemoteAddr,
		})
		return
	}

	connectionID := models.ConnectionID(uuid.NewString())
	clientIP := h.extractClientIP(r)

	session := &models.Session{
		ID:          connectionID,
		Connection:  conn,
		ClientIP:    clientIP,
		QueryParams: r.URL.Query(),
		ConnectedAt: time.Now(),
	}

	h.sessionManager.AddSession(session)

	logger.Info("Client connected", logger.Fields{
		"connection_id": string(connectionID),
		"client_ip":     clientIP,
		"query_params":  r.URL.Query(),
	})

	h.webhookService.NotifyConnection(session)

	go h.handleConnectionLoop(session)
}

func (h *WebSocketHandler) handleConnectionLoop(session *models.Session) {
	defer func() {
		h.sessionManager.RemoveSession(session.ID)
		h.webhookService.NotifyDisconnection(session)
		
		logger.Info("Client disconnected", logger.Fields{
			"connection_id": string(session.ID),
			"client_ip":     session.ClientIP,
		})
	}()

	for {
		_, _, err := session.Connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
				websocket.CloseNormalClosure) {
				logger.Warn("Unexpected WebSocket close", logger.Fields{
					"connection_id": string(session.ID),
					"error":         err.Error(),
				})
			} else {
				logger.Debug("WebSocket closed normally", logger.Fields{
					"connection_id": string(session.ID),
					"error":         err.Error(),
				})
			}
			break
		}
	}
}

func (h *WebSocketHandler) extractClientIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return forwarded
	}
	
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}
	
	clientIP, _, _ := net.SplitHostPort(r.RemoteAddr)
	return clientIP
} 