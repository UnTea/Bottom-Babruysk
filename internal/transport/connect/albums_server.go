package connect

import (
	"context"

	"connectrpc.com/connect"

	"github.com/untea/bottom_babruysk/internal/domain"
	"github.com/untea/bottom_babruysk/internal/service"
	protov1 "github.com/untea/bottom_babruysk/internal/transport/gen/proto/v1"
	"github.com/untea/bottom_babruysk/utils"
)

type AlbumsServer struct {
	albumsService *service.AlbumsService
}

func NewAlbumsServer(albumsService *service.AlbumsService) *AlbumsServer {
	return &AlbumsServer{albumsService: albumsService}
}

func (s *AlbumsServer) CreateAlbum(ctx context.Context, request *connect.Request[protov1.CreateAlbumRequest]) (*connect.Response[protov1.CreateAlbumResponse], error) {
	response, err := s.albumsService.CreateAlbum(ctx, domain.CreateAlbumRequest{
		OwnerID:     utils.StringToUUIDPtr(request.Msg.OwnerId),
		Title:       utils.Ptr(request.Msg.Title),
		Description: utils.Ptr(request.Msg.Description),
		ReleaseDate: utils.TimestamppbToTime(request.Msg.ReleaseDate),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}

	return connect.NewResponse(&protov1.CreateAlbumResponse{Id: utils.UUIDPtrToString(response.ID)}), nil
}

func (s *AlbumsServer) GetAlbum(ctx context.Context, request *connect.Request[protov1.GetAlbumRequest]) (*connect.Response[protov1.GetAlbumResponse], error) {
	response, err := s.albumsService.GetAlbum(ctx, domain.GetAlbumRequest{
		ID: utils.StringToUUIDPtr(request.Msg.Id),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}

	return connect.NewResponse(&protov1.GetAlbumResponse{Album: toProtoAlbum(response.Album)}), nil
}

func (s *AlbumsServer) ListAlbums(ctx context.Context, request *connect.Request[protov1.ListAlbumsRequest]) (*connect.Response[protov1.ListAlbumsResponse], error) {
	response, err := s.albumsService.ListAlbums(ctx, domain.ListAlbumsRequest{
		Limit:     utils.Ptr(int(request.Msg.Limit)),
		Offset:    utils.Ptr(int(request.Msg.Offset)),
		SortField: utils.Ptr(request.Msg.SortField),
		SortOrder: utils.Ptr(request.Msg.SortOrder),
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
		ID:          utils.StringToUUIDPtr(request.Msg.Id),
		Title:       request.Msg.Title,
		Description: request.Msg.Description,
		ReleaseDate: utils.TimestamppbToTime(request.Msg.ReleaseDate),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}

	return connect.NewResponse(&protov1.UpdateAlbumResponse{}), nil
}

func (s *AlbumsServer) DeleteAlbum(ctx context.Context, request *connect.Request[protov1.DeleteAlbumRequest]) (*connect.Response[protov1.DeleteAlbumResponse], error) {
	err := s.albumsService.DeleteAlbum(ctx, domain.DeleteAlbumRequest{
		ID: utils.StringToUUIDPtr(request.Msg.Id),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}

	return connect.NewResponse(&protov1.DeleteAlbumResponse{}), nil
}
