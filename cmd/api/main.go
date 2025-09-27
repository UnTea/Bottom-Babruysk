package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/untea/bottom_babruysk/internal/application"
	"github.com/untea/bottom_babruysk/internal/configuration"
	"github.com/untea/bottom_babruysk/internal/logger"
	"github.com/untea/bottom_babruysk/internal/repository"
	"github.com/untea/bottom_babruysk/internal/server"
	"github.com/untea/bottom_babruysk/internal/web/handlers"
	"github.com/untea/bottom_babruysk/internal/web/router"
)

func main() {
	l, err := logger.New()
	if err != nil {
		panic(err)
	}

	cfg, err := configuration.Load()
	if err != nil {
		l.Fatal("failed to load cfg", zap.Error(err))
	}

	dbCfg := repository.Configuration{
		ConnectionString: cfg.DatabaseConnectionURL,
		Timeout:          30 * time.Second,
	}

	dbClient, err := repository.New(context.Background(), dbCfg)
	if err != nil {
		l.Fatal("failed to initialize dbClient", zap.Error(err))
	}

	defer dbClient.Close()

	container, err := application.BuildContainer(cfg, l, dbClient)
	if err != nil {
		panic(err)
	}

	h := handlers.New(l, &container.Services)

	dependencies := router.Dependencies{
		Logger:           l,
		Handlers:         h,
		Services:         container.Services,
		EnableCORS:       true,
		EnableReflection: true,
	}

	srv := server.New(cfg, l, dependencies)

	go func() {
		err := srv.Start()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			l.Fatal("failed to start HTTP server", zap.Error(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = srv.Stop(ctx)
	if err != nil {
		l.Info("failed to stop HTTP server", zap.Error(err))
	}
}
