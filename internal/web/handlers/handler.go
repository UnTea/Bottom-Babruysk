package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/untea/bottom_babruysk/internal/application"
	"github.com/untea/bottom_babruysk/internal/domain"
)

type Handler struct {
	Logger   *zap.Logger
	Services *application.Services
}

func New(logger *zap.Logger, services *application.Services) *Handler {
	handler := &Handler{
		Logger:   logger,
		Services: services,
	}

	return handler
}

func (h *Handler) writeJson(w http.ResponseWriter, resp any, code int) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	if resp == nil {
		w.WriteHeader(code)
		return
	}

	response, err := json.Marshal(resp)
	if err != nil {
		h.Logger.Error("failed to marshal json response", zap.Error(err))

		h.httpError(w, errors.New("internal error"), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(code)

	_, err = w.Write(response)
	if err != nil {
		h.Logger.Error("failed to write response", zap.Error(err))
	}
}

func (h *Handler) httpError(w http.ResponseWriter, err error, code int) {
	errorDetail := ""
	if err != nil {
		errorDetail = err.Error()
	}

	h.writeJson(w, domain.ErrorResponse{Error: errorDetail}, code)
}
