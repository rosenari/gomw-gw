package services

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"gomw-gw/app/internal/config"
	"gomw-gw/app/internal/models"
	"gomw-gw/app/pkg/logger"
)

type WebhookService struct {
	httpClient *http.Client
	config     *config.WebhookConfig
}

func NewWebhookService(cfg *config.WebhookConfig) *WebhookService {
	return &WebhookService{
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
		config: cfg,
	}
}

func (ws *WebhookService) NotifyConnection(session *models.Session) {
	if ws.config.OnConnectURL == "" {
		return
	}

	payload := &models.WebhookPayload{
		ConnectionID: session.ID,
		ClientIP:     session.ClientIP,
		QueryParams:  session.QueryParams,
		Timestamp:    session.ConnectedAt,
	}

	go ws.callWebhook(ws.config.OnConnectURL, payload, "connection")
}

func (ws *WebhookService) NotifyDisconnection(session *models.Session) {
	if ws.config.OnDisconnectURL == "" {
		return
	}

	payload := &models.WebhookPayload{
		ConnectionID: session.ID,
		ClientIP:     session.ClientIP,
		Timestamp:    time.Now(),
	}

	go ws.callWebhook(ws.config.OnDisconnectURL, payload, "disconnection")
}

func (ws *WebhookService) callWebhook(url string, payload interface{}, eventType string) {
	ctx, cancel := context.WithTimeout(context.Background(), ws.config.Timeout)
	defer cancel()

	jsonData, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Failed to marshal webhook payload", logger.Fields{
			"event_type": eventType,
			"error":      err.Error(),
		})
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonData))
	if err != nil {
		logger.Error("Failed to create webhook request", logger.Fields{
			"event_type": eventType,
			"url":        url,
			"error":      err.Error(),
		})
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "gomw-gw/1.0")

	resp, err := ws.httpClient.Do(req)
	if err != nil {
		logger.Warn("Webhook call failed", logger.Fields{
			"event_type": eventType,
			"url":        url,
			"error":      err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	io.Copy(io.Discard, resp.Body)

	if resp.StatusCode >= 400 {
		logger.Warn("Webhook returned error status", logger.Fields{
			"event_type":  eventType,
			"url":         url,
			"status_code": resp.StatusCode,
		})
		return
	}

	logger.Debug("Webhook call successful", logger.Fields{
		"event_type":  eventType,
		"url":         url,
		"status_code": resp.StatusCode,
	})
} 