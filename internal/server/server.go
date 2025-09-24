package server

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/untea/bottom_babruysk/internal/configuration"
	"github.com/untea/bottom_babruysk/internal/repository"
	"github.com/untea/bottom_babruysk/internal/web/router"
)

type Server struct {
	configuration *configuration.Configuration
	db            *repository.Client
	http          *http.Server
}

func New(config *configuration.Configuration, log *zap.Logger) *Server {
	repositoryConfig := repository.Configuration{
		ConnectionString: config.DatabaseConnectionURL,
		Timeout:          time.Second * 30,
	}

	db, err := repository.New(context.Background(), repositoryConfig)
	if err != nil {
		panic(err)
	}

	handler := router.New(db, log)

	return &Server{
		configuration: config,
		db:            db,
		http: &http.Server{
			Addr:              config.HTTPAddress,
			Handler:           handler,
			ReadHeaderTimeout: time.Second * 10,
			ReadTimeout:       time.Second * 30,
			WriteTimeout:      time.Second * 30,
			IdleTimeout:       time.Second * 60,
		},
	}
}

func (s *Server) Start() error {
	return s.http.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	if s.db != nil {
		s.db.Close()
	}

	return s.http.Shutdown(ctx)
}
