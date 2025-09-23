package server

import (
	"context"
	"net/http"
	"time"

	"bottom_babruysk/internal/configuration"
	"bottom_babruysk/internal/repository"
	"bottom_babruysk/internal/web"
)

type Server struct {
	config *configuration.Config
	db     *repository.Client
	http   *http.Server
}

func New(config *configuration.Config) *Server {
	repositoryConfig := repository.Config{
		ConnectionString: config.DatabaseConnectionURL,
		Timeout:          30 * time.Second,
	}

	db, err := repository.New(context.Background(), repositoryConfig)
	if err != nil {
		panic(err)
	}

	handler := web.New(db)

	return &Server{
		config: config,
		db:     db,
		http: &http.Server{
			Addr:    config.HTTPAddress,
			Handler: handler,
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
