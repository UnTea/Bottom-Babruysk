package connect

import (
	"context"

	"connectrpc.com/connect"
	"github.com/google/uuid"

	"github.com/untea/bottom_babruysk/internal/domain"
	"github.com/untea/bottom_babruysk/internal/service"
	protov1 "github.com/untea/bottom_babruysk/internal/transport/gen/proto/v1"
)

type TracksServer struct {
	tracksService *service.TracksService
}

func NewTracksServer(tracksService *service.TracksService) *TracksServer {
	return &TracksServer{tracksService: tracksService}
}

func (s *TracksServer) GetTrack(ctx context.Context, req *connect.Request[protov1.GetTrackRequest]) (*connect.Response[protov1.GetTrackResponse], error) {
	out, err := s.tracksService.GetTrack(ctx, domain.GetTrackRequest{
		ID: StringToUUIDPtr(req.Msg.GetId()),
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	return connect.NewResponse(&protov1.GetTrackResponse{Track: toProtoTrack(out.Track)}), nil
}

func (s *TracksServer) ListTracks(ctx context.Context, request *connect.Request[protov1.ListTracksRequest]) (*connect.Response[protov1.ListTracksResponse], error) {
	out, err := s.tracksService.ListTracks(ctx, domain.ListTracksRequest{
		Limit:      Ptr(int(request.Msg.GetLimit())),
		Offset:     Ptr(int(request.Msg.GetOffset())),
		UploaderID: StringToUUIDPtr(request.Msg.GetUploaderId()),
		Visibility: (*domain.Visibility)(Ptr(request.Msg.GetVisibility())),
		SortField:  Ptr(request.Msg.GetSortField()),
		SortOrder:  Ptr(request.Msg.GetSortOrder()),
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	response := &protov1.ListTracksResponse{Tracks: make([]*protov1.Track, 0, len(out.Tracks))}

	for _, track := range out.Tracks {
		response.Tracks = append(response.Tracks, toProtoTrack(track))
	}

	return connect.NewResponse(response), nil
}

func (s *TracksServer) UpdateTrack(ctx context.Context, request *connect.Request[protov1.UpdateTrackRequest]) (*connect.Response[protov1.UpdateTrackResponse], error) {
	err := s.tracksService.UpdateTrack(ctx, domain.UpdateTrackRequest{
		ID:          StringToUUIDPtr(request.Msg.GetId()),
		Title:       Ptr(request.Msg.GetTitle()),
		Subtitle:    Ptr(request.Msg.GetSubtitle()),
		Description: Ptr(request.Msg.GetDescription()),
		Visibility:  (*domain.Visibility)(Ptr(request.Msg.GetVisibility())),
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	return connect.NewResponse(&protov1.UpdateTrackResponse{}), nil
}

func (s *TracksServer) DeleteTrack(ctx context.Context, req *connect.Request[protov1.DeleteTrackRequest]) (*connect.Response[protov1.DeleteTrackResponse], error) {
	id, _ := uuid.Parse(req.Msg.GetId())
	if err := s.tracksService.DeleteTrack(ctx, domain.DeleteTrackRequest{ID: &id}); err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}
	return connect.NewResponse(&protov1.DeleteTrackResponse{}), nil
}
