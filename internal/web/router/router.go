package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"

	"github.com/untea/bottom_babruysk/internal/repository"
	"github.com/untea/bottom_babruysk/internal/web/handlers"
	web "github.com/untea/bottom_babruysk/internal/web/middleware"
)

func New(db *repository.Client, log *zap.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(web.RequestLogger(log))

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		ctxzap.Extract(r.Context()).Debug("healthz ping")
		handlers.WriteJSON(w, http.StatusOK, map[string]any{
			"status":     "ok",
			"time_stamp": time.Now().UTC(),
		})
	})

	_ = db

	return r
}
