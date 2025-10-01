package connect

import (
	"context"

	"connectrpc.com/connect"

	"github.com/untea/bottom_babruysk/internal/domain"
	"github.com/untea/bottom_babruysk/internal/service"
	protov1 "github.com/untea/bottom_babruysk/internal/transport/gen/proto/v1"
	"github.com/untea/bottom_babruysk/utils"
)

type TrackFilesServer struct {
	trackService *service.TrackFilesService
}

func NewTrackFilesServer(trackFilesService *service.TrackFilesService) *TrackFilesServer {
	return &TrackFilesServer{trackService: trackFilesService}
}

func (s *TrackFilesServer) GetTrackFile(ctx context.Context, request *connect.Request[protov1.GetTrackFileRequest]) (*connect.Response[protov1.GetTrackFileResponse], error) {
	response, err := s.trackService.GetTrackFile(ctx, domain.GetTrackFileRequest{
		ID:      utils.StringToUUIDPtr(request.Msg.Id),
		TrackID: utils.StringToUUIDPtr(request.Msg.TrackId),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}

	return connect.NewResponse(&protov1.GetTrackFileResponse{File: toProtoTrackFile(response.TrackFile)}), nil
}

func (s *TrackFilesServer) ListTrackFiles(ctx context.Context, request *connect.Request[protov1.ListTrackFilesRequest]) (*connect.Response[protov1.ListTrackFilesResponse], error) {
	response, err := s.trackService.ListTrackFiles(ctx, domain.ListTrackFilesRequest{
		TrackID:   utils.StringToUUIDPtr(request.Msg.TrackId),
		Limit:     utils.Ptr(int(request.Msg.Limit)),
		Offset:    utils.Ptr(int(request.Msg.Offset)),
		SortField: utils.Ptr(request.Msg.SortField),
		SortOrder: utils.Ptr(request.Msg.SortOrder),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}

	result := &protov1.ListTrackFilesResponse{Files: make([]*protov1.TrackFile, 0, len(response.TrackFiles))}

	for _, trackFile := range response.TrackFiles {
		result.Files = append(result.Files, toProtoTrackFile(trackFile))
	}

	return connect.NewResponse(result), nil
}

func (s *TrackFilesServer) UpdateTrackFile(ctx context.Context, request *connect.Request[protov1.UpdateTrackFileRequest]) (*connect.Response[protov1.UpdateTrackFileResponse], error) {
	err := s.trackService.UpdateTrackFile(ctx, domain.UpdateTrackFileRequest{
		ID:         utils.StringToUUIDPtr(request.Msg.Id),
		TrackID:    utils.StringToUUIDPtr(request.Msg.TrackId),
		Filename:   request.Msg.Filename,
		S3Key:      request.Msg.S3Key,
		Mime:       request.Msg.Mime,
		Format:     FromProtoFormat(*request.Msg.Format),
		Codec:      FromProtoCodec(*request.Msg.Codec),
		Bitrate:    utils.Int32ToInt(request.Msg.Bitrate),
		SampleRate: utils.Int32ToInt(request.Msg.SampleRate),
		Channels:   utils.Int32ToInt(request.Msg.Channels),
		Size:       request.Msg.Size,
		Duration:   utils.DurationpbToDuration(request.Msg.Duration),
		Checksum:   request.Msg.Checksum,
	})
	if err != nil {
		return nil, toConnectErr(err)
	}

	return connect.NewResponse(&protov1.UpdateTrackFileResponse{}), nil
}

func (s *TrackFilesServer) DeleteTrackFile(ctx context.Context, request *connect.Request[protov1.DeleteTrackFileRequest]) (*connect.Response[protov1.DeleteTrackFileResponse], error) {
	err := s.trackService.DeleteTrackFile(ctx, domain.DeleteTrackFileRequest{
		ID:      utils.StringToUUIDPtr(request.Msg.Id),
		TrackID: utils.StringToUUIDPtr(request.Msg.TrackId),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}

	return connect.NewResponse(&protov1.DeleteTrackFileResponse{}), nil
}
