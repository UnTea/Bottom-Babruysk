package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/untea/bottom_babruysk/internal/repository"
)

type Handler struct {
	Logger *zap.Logger
	Repo   struct {
		Users repository.Users
	}
}

func New(logger *zap.Logger, repos struct{ Users repository.Users }) *Handler {
	h := &Handler{
		Logger: logger,
		Repo:   repos,
	}

	return h
}

func WriteJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_ = json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, err error, code int) {
	WriteJSON(w, code, map[string]error{"error": err})
}
