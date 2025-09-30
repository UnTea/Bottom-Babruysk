package web

import "github.com/go-chi/chi/v5"

type UsersHTTP interface {
	MountUsers(r chi.Router)
}

type AlbumsHTTP interface {
	MountAlbums(r chi.Router)
}

type TracksHTTP interface {
	MountTracks(r chi.Router)
}

type HandlerHTTP interface {
	UsersHTTP
	AlbumsHTTP
	TracksHTTP
}
