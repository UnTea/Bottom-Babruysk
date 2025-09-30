package service

import (
	"context"

	"github.com/untea/bottom_babruysk/internal/domain"
)

type AlbumsService struct {
	repository Albums
}

func NewAlbumsService(repository Albums) *AlbumsService {
	return &AlbumsService{repository: repository}
}

func (s *AlbumsService) CreateAlbum(ctx context.Context, request domain.CreateAlbumRequest) (*domain.CreateAlbumResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, err
	}

	return s.repository.CreateAlbum(ctx, request)
}

func (s *AlbumsService) GetAlbum(ctx context.Context, request domain.GetAlbumRequest) (*domain.GetAlbumResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, err
	}

	return s.repository.GetAlbum(ctx, request)
}

func (s *AlbumsService) ListAlbums(ctx context.Context, request domain.ListAlbumsRequest) (*domain.ListAlbumsResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, err
	}

	return s.repository.ListAlbums(ctx, request)
}

func (s *AlbumsService) UpdateAlbum(ctx context.Context, request domain.UpdateAlbumRequest) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	return s.repository.UpdateAlbum(ctx, request)
}

func (s *AlbumsService) DeleteAlbum(ctx context.Context, request domain.DeleteAlbumRequest) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	return s.repository.DeleteAlbum(ctx, request)
}
