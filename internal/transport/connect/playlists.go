package connect

import (
	"context"

	"connectrpc.com/connect"

	"github.com/untea/bottom_babruysk/internal/domain"
	"github.com/untea/bottom_babruysk/internal/service"
	protov1 "github.com/untea/bottom_babruysk/internal/transport/gen/proto/v1"
	"github.com/untea/bottom_babruysk/utils"
)

type PlaylistsServer struct {
	playlistsService *service.PlaylistsService
}

func NewPlaylistsServer(playlistsService *service.PlaylistsService) *PlaylistsServer {
	return &PlaylistsServer{playlistsService: playlistsService}
}

func (s *PlaylistsServer) CreatePlaylist(ctx context.Context, request *connect.Request[protov1.CreatePlaylistRequest]) (*connect.Response[protov1.CreatePlaylistResponse], error) {
	response, err := s.playlistsService.CreatePlaylist(ctx, domain.CreatePlaylistRequest{
		OwnerID:     utils.StringToUUIDPtr(request.Msg.OwnerId),
		Title:       utils.Ptr(request.Msg.Title),
		Description: utils.Ptr(request.Msg.Description),
		Visibility:  FromProtoVisibility(request.Msg.Visibility),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}

	return connect.NewResponse(&protov1.CreatePlaylistResponse{Id: response.ID.String()}), nil
}

func (s *PlaylistsServer) GetPlaylist(ctx context.Context, request *connect.Request[protov1.GetPlaylistRequest]) (*connect.Response[protov1.GetPlaylistResponse], error) {
	out, err := s.playlistsService.GetPlaylist(ctx, domain.GetPlaylistRequest{
		ID: utils.StringToUUIDPtr(request.Msg.Id),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}

	return connect.NewResponse(&protov1.GetPlaylistResponse{Playlist: toProtoPlaylist(out.Playlist)}), nil
}

func (s *PlaylistsServer) ListPlaylists(ctx context.Context, request *connect.Request[protov1.ListPlaylistsRequest]) (*connect.Response[protov1.ListPlaylistsResponse], error) {
	response, err := s.playlistsService.ListPlaylists(ctx, domain.ListPlaylistsRequest{
		Limit:      utils.Ptr(int(request.Msg.Limit)),
		Offset:     utils.Ptr(int(request.Msg.Offset)),
		OwnerID:    utils.StringToUUIDPtr(request.Msg.OwnerId),
		Visibility: FromProtoVisibility(request.Msg.Visibility),
		SortField:  utils.Ptr(request.Msg.SortField),
		SortOrder:  utils.Ptr(request.Msg.SortOrder),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}

	result := &protov1.ListPlaylistsResponse{Playlists: make([]*protov1.Playlist, 0, len(response.Playlists))}

	for _, playlist := range response.Playlists {
		result.Playlists = append(result.Playlists, toProtoPlaylist(playlist))
	}

	return connect.NewResponse(result), nil
}

func (s *PlaylistsServer) UpdatePlaylist(ctx context.Context, request *connect.Request[protov1.UpdatePlaylistRequest]) (*connect.Response[protov1.UpdatePlaylistResponse], error) {
	err := s.playlistsService.UpdatePlaylist(ctx, domain.UpdatePlaylistRequest{
		ID:          utils.StringToUUIDPtr(request.Msg.Id),
		Title:       request.Msg.Title,
		Description: request.Msg.Description,
		Visibility:  FromProtoVisibility(*request.Msg.Visibility),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}

	return connect.NewResponse(&protov1.UpdatePlaylistResponse{}), nil
}

func (s *PlaylistsServer) DeletePlaylist(ctx context.Context, request *connect.Request[protov1.DeletePlaylistRequest]) (*connect.Response[protov1.DeletePlaylistResponse], error) {
	err := s.playlistsService.DeletePlaylist(ctx, domain.DeletePlaylistRequest{
		ID: utils.StringToUUIDPtr(request.Msg.Id),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}

	return connect.NewResponse(&protov1.DeletePlaylistResponse{}), nil
}
