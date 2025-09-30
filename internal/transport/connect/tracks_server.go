package connect

import (
	"context"

	"connectrpc.com/connect"

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

func (s *TracksServer) GetTrack(ctx context.Context, request *connect.Request[protov1.GetTrackRequest]) (*connect.Response[protov1.GetTrackResponse], error) {
	response, err := s.tracksService.GetTrack(ctx, domain.GetTrackRequest{
		ID: StringToUUIDPtr(request.Msg.Id),
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	return connect.NewResponse(&protov1.GetTrackResponse{Track: toProtoTrack(response.Track)}), nil
}

func (s *TracksServer) ListTracks(ctx context.Context, request *connect.Request[protov1.ListTracksRequest]) (*connect.Response[protov1.ListTracksResponse], error) {
	response, err := s.tracksService.ListTracks(ctx, domain.ListTracksRequest{
		Limit:      Ptr(int(request.Msg.Limit)),
		Offset:     Ptr(int(request.Msg.Offset)),
		UploaderID: StringToUUIDPtr(request.Msg.UploaderId),
		Visibility: ProtoVisibilityPtrToDomain(Ptr(request.Msg.Visibility)),
		SortField:  Ptr(request.Msg.SortField),
		SortOrder:  Ptr(request.Msg.SortOrder),
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	result := &protov1.ListTracksResponse{Tracks: make([]*protov1.Track, 0, len(response.Tracks))}

	for _, track := range response.Tracks {
		result.Tracks = append(result.Tracks, toProtoTrack(track))
	}

	return connect.NewResponse(result), nil
}

func (s *TracksServer) UpdateTrack(ctx context.Context, request *connect.Request[protov1.UpdateTrackRequest]) (*connect.Response[protov1.UpdateTrackResponse], error) {
	err := s.tracksService.UpdateTrack(ctx, domain.UpdateTrackRequest{
		ID:          StringToUUIDPtr(request.Msg.Id),
		Title:       request.Msg.Title,
		Subtitle:    request.Msg.Subtitle,
		Description: request.Msg.Description,
		Visibility:  ProtoVisibilityPtrToDomain(request.Msg.Visibility),
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	return connect.NewResponse(&protov1.UpdateTrackResponse{}), nil
}

func (s *TracksServer) DeleteTrack(ctx context.Context, request *connect.Request[protov1.DeleteTrackRequest]) (*connect.Response[protov1.DeleteTrackResponse], error) {
	if err := s.tracksService.DeleteTrack(ctx, domain.DeleteTrackRequest{
		ID: StringToUUIDPtr(request.Msg.Id),
	}); err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	return connect.NewResponse(&protov1.DeleteTrackResponse{}), nil
}
