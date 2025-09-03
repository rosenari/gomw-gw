package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gomw-gw/app/internal/config"
	"gomw-gw/app/internal/handlers"
	"gomw-gw/app/internal/server"
	"gomw-gw/app/internal/services"
	"gomw-gw/app/pkg/logger"
)

func main() {
	cfg := config.LoadConfig()

	logger.Info("Starting gomw-gw", logger.Fields{
		"listen_address":    cfg.Server.ListenAddress,
		"on_connect_url":    cfg.Webhook.OnConnectURL,
		"on_disconnect_url": cfg.Webhook.OnDisconnectURL,
	})

	sessionManager := services.NewSessionManager()
	webhookService := services.NewWebhookService(&cfg.Webhook, &cfg.Server)

	wsHandler := handlers.NewWebSocketHandler(&cfg.WebSocket, sessionManager, webhookService)
	msgHandler := handlers.NewMessageHandler(sessionManager)
	infoHandler := handlers.NewInfoHandler(cfg, sessionManager)

	router := server.NewRouter(wsHandler, msgHandler, infoHandler)
	router.SetupRoutes()

	srv := server.NewServer(&cfg.Server, router.GetHandler())

	go func() {
		if err := srv.Start(); err != nil {
			logger.Fatal("Server startup failed", logger.Fields{
				"error": err.Error(),
			})
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...", logger.Fields{})

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", logger.Fields{
			"error": err.Error(),
		})
	}

	logger.Info("Server exited", logger.Fields{})
} 