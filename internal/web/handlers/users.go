package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/untea/bottom_babruysk/internal/domain"
	"github.com/untea/bottom_babruysk/internal/repository"
)

func (h *Handler) MountUsers(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/", h.createUser)
		r.Get("/", h.getListUser)
		r.Get("/{id}", h.getUser)
		r.Patch("/{id}", h.updateUser)
		r.Delete("/{id}", h.deleteUser)
	})
}

func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	request, err := Decode[domain.CreateUserRequest](r)
	if err != nil {
		h.httpError(w, errors.New("cannot decode incoming message"), http.StatusInternalServerError)
		return
	}

	if request.Email == nil || request.PasswordHash == nil || request.DisplayName == nil {
		h.httpError(w, errors.New("email, password_hash, display_name — required"), http.StatusBadRequest)
		return
	}

	response, err := h.Repo.Users.CreateUser(r.Context(), request)
	if err != nil {
		h.Logger.Error("failed to createUser user", zap.Error(err))
		h.httpError(w, err, http.StatusInternalServerError)

		return
	}

	h.writeJson(w, response, http.StatusCreated)
}

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	request, err := Decode[domain.GetUserRequest](r)
	if err != nil {
		h.httpError(w, errors.New("cannot decode incoming message"), http.StatusInternalServerError)
		return
	}

	if request.ID == uuid.Nil {
		h.httpError(w, errors.New("ID — required"), http.StatusBadRequest)
		return
	}

	response, err := h.Repo.Users.GetUser(r.Context(), request)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.httpError(w, err, http.StatusNotFound)
			return
		}

		h.httpError(w, err, http.StatusNotFound)

		return
	}

	h.writeJson(w, response, http.StatusOK)
}

func (h *Handler) getListUser(w http.ResponseWriter, r *http.Request) {
	request, err := Decode[domain.GetListUserRequest](r)
	if err != nil {
		h.httpError(w, errors.New("cannot decode incoming message"), http.StatusInternalServerError)
		return
	}

	response, err := h.Repo.Users.GetListUser(r.Context(), request)
	if err != nil {
		h.httpError(w, err, http.StatusInternalServerError)
		return
	}

	h.writeJson(w, response, http.StatusOK)
}

func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request) {
	request, err := Decode[domain.UpdateUserRequest](r)
	if err != nil {
		h.httpError(w, errors.New("cannot decode incoming message"), http.StatusInternalServerError)
		return
	}

	if err = h.Repo.Users.UpdateUser(r.Context(), request); err != nil {
		h.httpError(w, err, http.StatusInternalServerError)
		return
	}

	h.writeJson(w, "user updated", http.StatusOK)
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	request, err := Decode[domain.DeleteUserRequest](r)
	if err != nil {
		h.httpError(w, errors.New("cannot decode incoming message"), http.StatusInternalServerError)
		return
	}

	if err = h.Repo.Users.DeleteUser(r.Context(), request); err != nil {
		h.httpError(w, err, http.StatusInternalServerError)
		return
	}

	h.writeJson(w, "user deleted", http.StatusOK)
}
