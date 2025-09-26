package server

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/untea/bottom_babruysk/internal/configuration"
	"github.com/untea/bottom_babruysk/internal/web/router"
)

type Server struct {
	configuration *configuration.Configuration
	logger        *zap.Logger
	httpServer    *http.Server
}

func New(configuration *configuration.Configuration, logger *zap.Logger, dependencies router.Dependencies) *Server {
	r := router.New(dependencies)

	httpServer := &http.Server{
		Addr:              configuration.HTTPAddress,
		Handler:           r,
		ReadHeaderTimeout: 30 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       90 * time.Second,
	}

	return &Server{
		configuration: configuration,
		logger:        logger,
		httpServer:    httpServer,
	}
}

func (s *Server) Start() error {
	s.logger.Info("HTTP server start", zap.String("addr", s.configuration.HTTPAddress))

	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("HTTP server shutdown")

	return s.httpServer.Shutdown(ctx)
}
