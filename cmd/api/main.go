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

	"github.com/untea/bottom_babruysk/internal/configuration"
	"github.com/untea/bottom_babruysk/internal/logger"
	"github.com/untea/bottom_babruysk/internal/server"
)

func main() {
	l, err := logger.New()
	if err != nil {
		panic(err)
	}

	config, err := configuration.Load()
	if err != nil {
		l.Info("failed to load config", zap.Error(err))
	}

	srv := server.New(config, l)

	go func() {
		l.Info("HTTP server listening", zap.Any("address", config.HTTPAddress))

		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			l.Info("failed to start HTTP server", zap.Error(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err = srv.Stop(ctx); err != nil {
		l.Info("failed to stop HTTP server", zap.Error(err))
	}
}
