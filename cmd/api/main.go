package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bottom_babruysk/internal/configuration"
	"bottom_babruysk/internal/server"
)

func main() {
	databaseConnectionURL := "postgres://admin:admin@localhost:5432/bottom_babruysk?sslmode=disable"
	httpAddress := ":8080"

	config := configuration.New(databaseConnectionURL, httpAddress)
	srv := server.New(config)

	go func() {
		log.Printf("HTTP listening on %s", config.HTTPAddress)

		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server start: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := srv.Stop(ctx); err != nil {
		log.Printf("server stop: %v", err)
	}
}
