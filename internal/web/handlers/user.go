package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/untea/bottom_babruysk/internal/domain"
	"github.com/untea/bottom_babruysk/internal/repository"
)

type UsersHandler struct {
	repo repository.Users
}

func NewUsersHandler(repo repository.Users) *UsersHandler {
	return &UsersHandler{
		repo: repo,
	}
}

func (h *UsersHandler) Mount(r chi.Router) {
	r.Route("/api/v1/users", func(r chi.Router) {
		r.Get("/", h.list)
		r.Post("/", h.create)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.get)
			r.Patch("/", h.update)
			r.Delete("/", h.delete)
		})
	})
}

func (h *UsersHandler) list(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))
	role := q.Get("role")
	search := q.Get("search")

	var rolePtr, searchPtr *string

	if role != "" {
		rolePtr = &role
	}

	if search != "" {
		searchPtr = &search
	}

	users, err := h.repo.List(r.Context(), domain.Page{Limit: limit, Offset: offset}, rolePtr, searchPtr)
	if err != nil {
		WriteError(w, err, http.StatusInternalServerError)
		return
	}

	WriteJSON(w, http.StatusOK, users)
}

func (h *UsersHandler) create(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, err, http.StatusBadRequest)
		return
	}

	user, err := h.repo.Create(r.Context(), req)
	if err != nil {
		WriteError(w, err, http.StatusInternalServerError)
		return
	}

	WriteJSON(w, http.StatusCreated, user)
}

func (h *UsersHandler) get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		WriteError(w, err, http.StatusBadRequest)
		return
	}

	user, err := h.repo.Get(r.Context(), id)
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, repository.ErrNotFound) {
			status = http.StatusNotFound
		}

		WriteError(w, err, status)

		return
	}

	WriteJSON(w, http.StatusOK, user)
}

func (h *UsersHandler) update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		WriteError(w, err, http.StatusBadRequest)
		return
	}

	var req domain.UpdateUserRequest

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, err, http.StatusBadRequest)
		return
	}

	if err = h.repo.Update(r.Context(), id, req); err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, repository.ErrNotFound) {
			status = http.StatusNotFound
		}

		WriteError(w, err, status)

		return
	}

	WriteJSON(w, http.StatusOK, "")
}

func (h *UsersHandler) delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		WriteError(w, err, http.StatusBadRequest)
		return
	}

	if err = h.repo.Delete(r.Context(), id); err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, repository.ErrNotFound) {
			status = http.StatusNotFound
		}

		WriteError(w, err, status)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
