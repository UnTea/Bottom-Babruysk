package service

import (
	"context"

	"github.com/untea/bottom_babruysk/internal/domain"
	"github.com/untea/bottom_babruysk/internal/repository"
)

type TracksService struct {
	repository repository.Tracks
}

func NewTracksService(repository repository.Tracks) *TracksService {
	return &TracksService{repository: repository}
}

func (s *TracksService) CreateTrack(ctx context.Context, request domain.CreateTrackRequest) (*domain.CreateTrackResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, err
	}

	return s.repository.CreateTrack(ctx, request)
}

func (s *TracksService) GetTrack(ctx context.Context, request domain.GetTrackRequest) (*domain.GetTrackResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, err
	}

	return s.repository.GetTrack(ctx, request)
}

func (s *TracksService) ListTracks(ctx context.Context, request domain.ListTracksRequest) (*domain.ListTracksResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, err
	}

	return s.repository.ListTracks(ctx, request)
}

func (s *TracksService) UpdateTrack(ctx context.Context, request domain.UpdateTrackRequest) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	return s.repository.UpdateTrack(ctx, request)
}

func (s *TracksService) DeleteTrack(ctx context.Context, request domain.DeleteTrackRequest) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	return s.repository.DeleteTrack(ctx, request)
}
