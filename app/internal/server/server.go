package server

import (
	"context"
	"net/http"

	"gomw-gw/app/internal/config"
	"gomw-gw/app/pkg/logger"
)

type Server struct {
	httpServer *http.Server
	config     *config.ServerConfig
}

func NewServer(cfg *config.ServerConfig, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         cfg.ListenAddress,
			Handler:      handler,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		},
		config: cfg,
	}
}

func (s *Server) Start() error {
	logger.Info("Starting HTTP server", logger.Fields{
		"listen_address": s.config.ListenAddress,
		"read_timeout":   s.config.ReadTimeout,
		"write_timeout":  s.config.WriteTimeout,
	})

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("Server failed to start", logger.Fields{
			"error": err.Error(),
		})
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	logger.Info("Shutting down HTTP server", logger.Fields{})

	if err := s.httpServer.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed", logger.Fields{
			"error": err.Error(),
		})
		return err
	}

	logger.Info("Server shutdown completed", logger.Fields{})
	return nil
} 