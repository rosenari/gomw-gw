package handlers

import (
	"encoding/json"
	"net/http"

	"gomw-gw/app/internal/models"
	"gomw-gw/app/internal/services"
	"gomw-gw/app/pkg/logger"

	"github.com/gorilla/websocket"
)

type MessageHandler struct {
	sessionManager *services.SessionManager
}

func NewMessageHandler(sessionManager *services.SessionManager) *MessageHandler {
	return &MessageHandler{
		sessionManager: sessionManager,
	}
}

func (h *MessageHandler) HandleSendMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request models.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		logger.Warn("Invalid JSON in send message request", logger.Fields{
			"error":       err.Error(),
			"remote_addr": r.RemoteAddr,
		})
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if request.ConnectionID == "" {
		http.Error(w, "connection_id is required", http.StatusBadRequest)
		return
	}

	if len(request.Message) == 0 {
		http.Error(w, "message is required", http.StatusBadRequest)
		return
	}

	session, exists := h.sessionManager.GetSession(request.ConnectionID)
	if !exists {
		logger.Warn("Connection not found for send request", logger.Fields{
			"connection_id": string(request.ConnectionID),
			"remote_addr":   r.RemoteAddr,
		})
		http.Error(w, "Connection not found", http.StatusNotFound)
		return
	}

	if !session.IsValid() {
		logger.Warn("Invalid session for send request", logger.Fields{
			"connection_id": string(request.ConnectionID),
		})
		http.Error(w, "Invalid connection", http.StatusGone)
		return
	}

	if err := session.Connection.WriteMessage(websocket.TextMessage, request.Message); err != nil {
		logger.Error("Failed to send message to WebSocket", logger.Fields{
			"connection_id": string(request.ConnectionID),
			"error":         err.Error(),
		})
		
		h.sessionManager.RemoveSession(request.ConnectionID)
		
		http.Error(w, "Failed to send message", http.StatusBadGateway)
		return
	}

	logger.Info("Message sent successfully", logger.Fields{
		"connection_id": string(request.ConnectionID),
		"message_size":  len(request.Message),
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":       true,
		"connection_id": request.ConnectionID,
		"message":       "Message sent successfully",
	})
} 