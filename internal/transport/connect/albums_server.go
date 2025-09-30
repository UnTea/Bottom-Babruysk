package connect

import (
	"context"

	"connectrpc.com/connect"

	"github.com/untea/bottom_babruysk/internal/domain"
	"github.com/untea/bottom_babruysk/internal/service"
	protov1 "github.com/untea/bottom_babruysk/internal/transport/gen/proto/v1"
)

type AlbumsServer struct {
	albumsService *service.AlbumsService
}

func NewAlbumsServer(albumsService *service.AlbumsService) *AlbumsServer {
	return &AlbumsServer{albumsService: albumsService}
}

func (s *AlbumsServer) CreateAlbum(ctx context.Context, request *connect.Request[protov1.CreateAlbumRequest]) (*connect.Response[protov1.CreateAlbumResponse], error) {
	response, err := s.albumsService.CreateAlbum(ctx, domain.CreateAlbumRequest{
		OwnerID:     StringToUUIDPtr(request.Msg.GetOwnerId()),
		Title:       Ptr(request.Msg.GetTitle()),
		Description: Ptr(request.Msg.GetDescription()),
		ReleaseDate: TimestamppbToTime(request.Msg.GetReleaseDate()),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}

	return connect.NewResponse(&protov1.CreateAlbumResponse{
		Id: UUIDPtrToString(response.ID),
	}), nil
}

func (s *AlbumsServer) GetAlbum(ctx context.Context, request *connect.Request[protov1.GetAlbumRequest]) (*connect.Response[protov1.GetAlbumResponse], error) {
	response, err := s.albumsService.GetAlbum(ctx, domain.GetAlbumRequest{
		ID: StringToUUIDPtr(request.Msg.GetId()),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}

	return connect.NewResponse(&protov1.GetAlbumResponse{Album: toProtoAlbum(response.Album)}), nil
}

func (s *AlbumsServer) ListAlbums(ctx context.Context, request *connect.Request[protov1.ListAlbumsRequest]) (*connect.Response[protov1.ListAlbumsResponse], error) {
	response, err := s.albumsService.ListAlbums(ctx, domain.ListAlbumsRequest{
		Limit:     Ptr(int(request.Msg.GetLimit())),
		Offset:    Ptr(int(request.Msg.GetOffset())),
		SortField: Ptr(request.Msg.GetSortField()),
		SortOrder: Ptr(request.Msg.GetSortOrder()),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}

	result := &protov1.ListAlbumsResponse{
		Albums: make([]*protov1.Album, 0, len(response.Albums)),
	}

	for _, album := range response.Albums {
		result.Albums = append(result.Albums, toProtoAlbum(album))
	}

	return connect.NewResponse(result), nil
}

func (s *AlbumsServer) UpdateAlbum(ctx context.Context, request *connect.Request[protov1.UpdateAlbumRequest]) (*connect.Response[protov1.UpdateAlbumResponse], error) {
	err := s.albumsService.UpdateAlbum(ctx, domain.UpdateAlbumRequest{
		ID:          StringToUUIDPtr(request.Msg.GetId()),
		Title:       Ptr(request.Msg.GetTitle()),
		Description: Ptr(request.Msg.GetDescription()),
		ReleaseDate: TimestamppbToTime(request.Msg.GetReleaseDate()),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}

	return connect.NewResponse(&protov1.UpdateAlbumResponse{}), nil
}

func (s *AlbumsServer) DeleteAlbum(ctx context.Context, request *connect.Request[protov1.DeleteAlbumRequest]) (*connect.Response[protov1.DeleteAlbumResponse], error) {
	err := s.albumsService.DeleteAlbum(ctx, domain.DeleteAlbumRequest{
		ID: StringToUUIDPtr(request.Msg.GetId()),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}

	return connect.NewResponse(&protov1.DeleteAlbumResponse{}), nil
}
