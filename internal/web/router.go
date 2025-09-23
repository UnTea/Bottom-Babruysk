package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"bottom_babruysk/internal/handlers"
	"bottom_babruysk/internal/repository"
	"bottom_babruysk/internal/repository/postgres"
)

func New(db *repository.Client) http.Handler {
	r := chi.NewRouter()

	var (
		userRepo = postgres.NewUsersRepo(db)
	)

	usersHandler := handlers.NewUsersHandler(userRepo)
	usersHandler.Mount(r)

	// ping
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

	return r
}
