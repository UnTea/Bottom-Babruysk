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
	"github.com/untea/bottom_babruysk/internal/repository"
	"github.com/untea/bottom_babruysk/internal/repository/postgres"
	"github.com/untea/bottom_babruysk/internal/server"
	"github.com/untea/bottom_babruysk/internal/service"
	"github.com/untea/bottom_babruysk/internal/web/handlers"
	"github.com/untea/bottom_babruysk/internal/web/router"
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

	databaseConfiguration := repository.Configuration{
		ConnectionString: config.DatabaseConnectionURL,
		Timeout:          30 * time.Second,
	}

	db, err := repository.New(context.Background(), databaseConfiguration)
	if err != nil {
		l.Info("failed to init db", zap.Error(err))
	}

	defer db.Close()

	usersRepo := postgres.NewUsersRepo(db)

	usersSvc := service.NewUsersSrv(usersRepo)

	h := handlers.New(l, struct{ Users service.Users }{Users: usersSvc})

	deps := router.Deps{
		Logger:     l,
		Users:      h,
		EnableCORS: true,
	}

	srv := server.New(config, l, deps)

	go func() {
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
