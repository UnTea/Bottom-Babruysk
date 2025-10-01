package service

import (
	"context"

	"github.com/untea/bottom_babruysk/internal/domain"
)

type TrackFilesService struct {
	repository TrackFiles
}

func NewTrackFilesService(repository TrackFiles) *TrackFilesService {
	return &TrackFilesService{repository: repository}
}

func (s *TrackFilesService) CreateTrackFile(ctx context.Context, request domain.CreateTrackFileRequest) (*domain.CreateTrackFileResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, err
	}

	return s.repository.CreateTrackFile(ctx, request)
}

func (s *TrackFilesService) GetTrackFile(ctx context.Context, request domain.GetTrackFileRequest) (*domain.GetTrackFileResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, err
	}

	return s.repository.GetTrackFile(ctx, request)
}

func (s *TrackFilesService) ListTrackFiles(ctx context.Context, request domain.ListTrackFilesRequest) (*domain.ListTrackFilesResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, err
	}

	return s.repository.ListTrackFiles(ctx, request)
}

func (s *TrackFilesService) UpdateTrackFile(ctx context.Context, request domain.UpdateTrackFileRequest) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	return s.repository.UpdateTrackFile(ctx, request)
}

func (s *TrackFilesService) DeleteTrackFile(ctx context.Context, request domain.DeleteTrackFileRequest) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	return s.repository.DeleteTrackFile(ctx, request)
}
