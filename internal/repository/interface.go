package repository

import (
	"context"

	"github.com/untea/bottom_babruysk/internal/domain"
)

type Users interface {
	CreateUser(context.Context, domain.CreateUserRequest) (*domain.CreateUserResponse, error)
	GetUser(context.Context, domain.GetUserRequest) (*domain.GetUserResponse, error)
	ListUsers(context.Context, domain.ListUsersRequest) (*domain.ListUsersResponse, error)
	UpdateUser(context.Context, domain.UpdateUserRequest) error
	DeleteUser(context.Context, domain.DeleteUserRequest) error
}

type Albums interface {
	CreateAlbum(context.Context, domain.CreateAlbumRequest) (*domain.CreateAlbumResponse, error)
	GetAlbum(context.Context, domain.GetAlbumRequest) (*domain.GetAlbumResponse, error)
	ListAlbums(context.Context, domain.ListAlbumsRequest) (*domain.ListAlbumsResponse, error)
	UpdateAlbum(context.Context, domain.UpdateAlbumRequest) error
	DeleteAlbum(context.Context, domain.DeleteAlbumRequest) error
}

type Tracks interface {
	CreateTrack(context.Context, domain.CreateTrackRequest) (*domain.CreateTrackResponse, error)
	GetTrack(context.Context, domain.GetTrackRequest) (*domain.GetTrackResponse, error)
	ListTracks(context.Context, domain.ListTracksRequest) (*domain.ListTracksResponse, error)
	UpdateTrack(context.Context, domain.UpdateTrackRequest) error
	DeleteTrack(context.Context, domain.DeleteTrackRequest) error
}
