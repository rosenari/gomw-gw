package handlers

import (
	"encoding/json"
	"net/http"

	"gomw-gw/app/internal/config"
	"gomw-gw/app/internal/models"
	"gomw-gw/app/internal/services"
	"gomw-gw/app/pkg/logger"
)

type InfoHandler struct {
	config         *config.Config
	sessionManager *services.SessionManager
}

func NewInfoHandler(cfg *config.Config, sessionManager *services.SessionManager) *InfoHandler {
	return &InfoHandler{
		config:         cfg,
		sessionManager: sessionManager,
	}
}

func (h *InfoHandler) HandleEnvironmentInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	envInfo := &models.EnvironmentInfo{
		ListenAddress:   h.config.Server.ListenAddress,
		OnConnectURL:    h.config.Webhook.OnConnectURL,
		OnDisconnectURL: h.config.Webhook.OnDisconnectURL,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(envInfo); err != nil {
		logger.Error("Failed to encode environment info", logger.Fields{
			"error": err.Error(),
		})
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Debug("Environment info requested", logger.Fields{
		"remote_addr": r.RemoteAddr,
	})
}

func (h *InfoHandler) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	activeConnections := h.sessionManager.GetSessionCount()

	healthInfo := map[string]interface{}{
		"status":             "healthy",
		"active_connections": activeConnections,
		"timestamp":          "2024-01-01T00:00:00Z",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(healthInfo); err != nil {
		logger.Error("Failed to encode health info", logger.Fields{
			"error": err.Error(),
		})
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *InfoHandler) HandleConnectionStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessions := h.sessionManager.GetAllSessions()
	
	connectionInfo := make([]map[string]interface{}, 0, len(sessions))
	for _, session := range sessions {
		info := map[string]interface{}{
			"connection_id": session.ID,
			"client_ip":     session.ClientIP,
			"connected_at":  session.ConnectedAt,
			"query_params":  session.QueryParams,
		}
		connectionInfo = append(connectionInfo, info)
	}

	statusInfo := map[string]interface{}{
		"total_connections": len(sessions),
		"connections":       connectionInfo,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(statusInfo); err != nil {
		logger.Error("Failed to encode connection status", logger.Fields{
			"error": err.Error(),
		})
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	logger.Debug("Connection status requested", logger.Fields{
		"remote_addr":        r.RemoteAddr,
		"active_connections": len(sessions),
	})
} 