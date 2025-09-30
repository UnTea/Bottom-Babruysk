package service

import (
	"context"

	"github.com/untea/bottom_babruysk/internal/domain"
)

type ArtistsService struct {
	repository Artists
}

func NewArtistsService(repo Artists) *ArtistsService {
	return &ArtistsService{repository: repo}
}

func (s *ArtistsService) CreateArtist(ctx context.Context, request domain.CreateArtistRequest) (*domain.CreateArtistResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, err
	}

	return s.repository.CreateArtist(ctx, request)
}

func (s *ArtistsService) GetArtist(ctx context.Context, request domain.GetArtistRequest) (*domain.GetArtistResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, err
	}

	return s.repository.GetArtist(ctx, request)
}

func (s *ArtistsService) ListArtists(ctx context.Context, request domain.ListArtistsRequest) (*domain.ListArtistsResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, err
	}

	return s.repository.ListArtists(ctx, request)
}

func (s *ArtistsService) UpdateArtist(ctx context.Context, request domain.UpdateArtistRequest) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	return s.repository.UpdateArtist(ctx, request)
}

func (s *ArtistsService) DeleteArtist(ctx context.Context, request domain.DeleteArtistRequest) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	return s.repository.DeleteArtist(ctx, request)
}
