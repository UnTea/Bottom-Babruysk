package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"

	web2 "github.com/untea/bottom_babruysk/internal/web"
	web "github.com/untea/bottom_babruysk/internal/web/middleware"
)

type Deps struct {
	Logger     *zap.Logger
	Users      web2.UsersHTTP
	EnableCORS bool
}

func New(deps Deps) *chi.Mux {
	r := chi.NewRouter()

	// базовые middlewares
	if deps.EnableCORS {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}))
	}

	r.Use(
		chiMiddleware.RequestID,
		chiMiddleware.RealIP,
		chiMiddleware.Recoverer,
		chiMiddleware.Timeout(60*time.Second),
		web.RequestLogger(deps.Logger),
	)

	// health
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Route("/api/v1", func(api chi.Router) {
		// /api/v1/users/...
		if deps.Users != nil {
			deps.Users.MountUsers(api)
		}
	})

	return r
}
