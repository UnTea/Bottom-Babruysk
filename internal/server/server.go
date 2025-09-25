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
	log           *zap.Logger
	http          *http.Server
}

func New(config *configuration.Configuration, log *zap.Logger, deps router.Deps) *Server {
	r := router.New(deps)

	httpServer := &http.Server{
		Addr:              config.HTTPAddress,
		Handler:           r,
		ReadHeaderTimeout: time.Second * 30,
		ReadTimeout:       time.Second * 30,
		WriteTimeout:      time.Second * 60,
		IdleTimeout:       time.Second * 90,
	}

	return &Server{
		configuration: config,
		log:           log,
		http:          httpServer,
	}
}

func (s *Server) Start() error {
	s.log.Info("HTTP server start", zap.String("addr", s.configuration.HTTPAddress))

	return s.http.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("HTTP server shutdown")

	return s.http.Shutdown(ctx)
}
