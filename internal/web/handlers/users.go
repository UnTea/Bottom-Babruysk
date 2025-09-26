package handlers

import (
	"github.com/go-chi/chi/v5"
)

func (h *Handler) MountUsers(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/", Handle(h, h.Services.UsersServices.CreateUser))
		r.Get("/", Handle(h, h.Services.UsersServices.ListUsers))
		r.Get("/{id}", Handle(h, h.Services.UsersServices.GetUser))
		r.Patch("/{id}", Handle(h, Lift(h.Services.UsersServices.UpdateUser)))
		r.Delete("/{id}", Handle(h, Lift(h.Services.UsersServices.DeleteUser)))
	})
}
