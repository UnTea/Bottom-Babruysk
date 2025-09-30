package handlers

import "github.com/go-chi/chi/v5"

func (h *Handler) MountUsers(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/", Handle(h, h.Services.UsersServices.CreateUser))
		r.Get("/", Handle(h, h.Services.UsersServices.ListUsers))
		r.Get("/{id}", Handle(h, h.Services.UsersServices.GetUser))
		r.Patch("/{id}", Handle(h, Lift(h.Services.UsersServices.UpdateUser)))
		r.Delete("/{id}", Handle(h, Lift(h.Services.UsersServices.DeleteUser)))
	})
}

func (h *Handler) MountAlbums(r chi.Router) {
	r.Route("/albums", func(r chi.Router) {
		r.Post("/", Handle(h, h.Services.AlbumServices.CreateAlbum))
		r.Get("/", Handle(h, h.Services.AlbumServices.ListAlbums))
		r.Get("/{id}", Handle(h, h.Services.AlbumServices.GetAlbum))
		r.Patch("/{id}", Handle(h, Lift(h.Services.AlbumServices.UpdateAlbum)))
		r.Delete("/{id}", Handle(h, Lift(h.Services.AlbumServices.DeleteAlbum)))
	})
}

func (h *Handler) MountTracks(r chi.Router) {
	r.Route("/tracks", func(r chi.Router) {
		r.Get("/", Handle(h, h.Services.TacksServices.ListTracks))
		r.Get("/{id}", Handle(h, h.Services.TacksServices.GetTrack))
		r.Patch("/{id}", Handle(h, Lift(h.Services.TacksServices.UpdateTrack)))
		r.Delete("/{id}", Handle(h, Lift(h.Services.TacksServices.DeleteTrack)))
	})
}

func (h *Handler) MountPlaylists(r chi.Router) {
	r.Route("/playlist", func(r chi.Router) {
		r.Post("/", Handle(h, h.Services.PlaylistsService.CreatePlaylist))
		r.Get("/", Handle(h, h.Services.PlaylistsService.ListPlaylists))
		r.Get("/{id}", Handle(h, h.Services.PlaylistsService.GetPlaylist))
		r.Patch("/{id}", Handle(h, Lift(h.Services.PlaylistsService.UpdatePlaylist)))
		r.Delete("/{id}", Handle(h, Lift(h.Services.PlaylistsService.DeletePlaylist)))
	})
}
