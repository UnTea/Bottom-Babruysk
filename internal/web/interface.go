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

type PlaylistsHTTP interface {
	MountPlaylists(r chi.Router)
}

type ArtistsHTTP interface {
	MountArtists(r chi.Router)
}

type TrackFilesHTTP interface {
	MountTrackFiles(r chi.Router)
}

type HandlerHTTP interface {
	UsersHTTP
	AlbumsHTTP
	TracksHTTP
	PlaylistsHTTP
	ArtistsHTTP
	TrackFilesHTTP
}
