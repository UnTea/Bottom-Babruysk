package service

import (
	"context"

	"github.com/untea/bottom_babruysk/internal/domain"
)

type PlaylistsService struct {
	repository Playlists
}

func NewPlaylistsService(repository Playlists) *PlaylistsService {
	return &PlaylistsService{repository: repository}
}

func (s *PlaylistsService) CreatePlaylist(ctx context.Context, request domain.CreatePlaylistRequest) (*domain.CreatePlaylistResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, err
	}

	return s.repository.CreatePlaylist(ctx, request)
}

func (s *PlaylistsService) GetPlaylist(ctx context.Context, request domain.GetPlaylistRequest) (*domain.GetPlaylistResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, err
	}

	return s.repository.GetPlaylist(ctx, request)
}

func (s *PlaylistsService) ListPlaylists(ctx context.Context, request domain.ListPlaylistsRequest) (*domain.ListPlaylistsResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, err
	}

	return s.repository.ListPlaylists(ctx, request)
}

func (s *PlaylistsService) UpdatePlaylist(ctx context.Context, request domain.UpdatePlaylistRequest) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	return s.repository.UpdatePlaylist(ctx, request)
}

func (s *PlaylistsService) DeletePlaylist(ctx context.Context, request domain.DeletePlaylistRequest) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	return s.repository.DeletePlaylist(ctx, request)
}
