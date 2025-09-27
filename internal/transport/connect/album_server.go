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

func NewAlbumsServer(svc *service.AlbumsService) *AlbumsServer {
	return &AlbumsServer{albumsService: svc}
}

func (s *AlbumsServer) CreateAlbum(ctx context.Context, request *connect.Request[protov1.CreateAlbumRequest]) (*connect.Response[protov1.CreateAlbumResponse], error) {
	out, err := s.albumsService.CreateAlbum(ctx, domain.CreateAlbumRequest{
		OwnerID:     uuidPtrFromStr(request.Msg.GetOwnerId()),
		Title:       strPtrOrNil(request.Msg.GetTitle()),
		Description: strPtrOrNil(request.Msg.GetDescription()),
		ReleaseDate: timePtrFromTS(request.Msg.GetReleaseDate()),
	})

	if err != nil {
		return nil, toConnectErr(err)
	}

	return connect.NewResponse(&protov1.CreateAlbumResponse{Id: uuidStr(out.ID)}), nil
}

func (s *AlbumsServer) GetAlbum(ctx context.Context, request *connect.Request[protov1.GetAlbumRequest]) (*connect.Response[protov1.GetAlbumResponse], error) {
	out, err := s.albumsService.GetAlbum(ctx, domain.GetAlbumRequest{ID: uuidFromStr(request.Msg.GetId())})
	if err != nil {
		return nil, toConnectErr(err)
	}

	return connect.NewResponse(&protov1.GetAlbumResponse{Album: toProtoAlbum(out.Album)}), nil
}

func (s *AlbumsServer) ListAlbums(ctx context.Context, request *connect.Request[protov1.ListAlbumsRequest]) (*connect.Response[protov1.ListAlbumsResponse], error) {
	out, err := s.albumsService.ListAlbums(ctx, domain.ListAlbumsRequest{
		Limit:     intPtr(int(request.Msg.GetLimit())),
		Offset:    intPtr(int(request.Msg.GetOffset())),
		SortField: strPtrOrNil(request.Msg.GetSortField()),
		SortOrder: strPtrOrNil(request.Msg.GetSortOrder()),
	})

	if err != nil {
		return nil, toConnectErr(err)
	}

	response := &protov1.ListAlbumsResponse{Albums: make([]*protov1.Album, 0, len(out.Albums))}

	for _, a := range out.Albums {
		response.Albums = append(response.Albums, toProtoAlbum(a))
	}

	return connect.NewResponse(response), nil
}

func (s *AlbumsServer) UpdateAlbum(ctx context.Context, request *connect.Request[protov1.UpdateAlbumRequest]) (*connect.Response[protov1.UpdateAlbumResponse], error) {
	err := s.albumsService.UpdateAlbum(ctx, domain.UpdateAlbumRequest{
		ID:          uuidPtrFromStr(request.Msg.Id),
		Title:       request.Msg.Title,
		Description: request.Msg.Description,
		ReleaseDate: timePtrFromTS(request.Msg.ReleaseDate),
	})

	if err != nil {
		return nil, toConnectErr(err)
	}

	return connect.NewResponse(&protov1.UpdateAlbumResponse{}), nil
}

func (s *AlbumsServer) DeleteAlbum(ctx context.Context, req *connect.Request[protov1.DeleteAlbumRequest]) (*connect.Response[protov1.DeleteAlbumResponse], error) {
	err := s.albumsService.DeleteAlbum(ctx, domain.DeleteAlbumRequest{ID: uuidPtrFromStr(req.Msg.GetId())})
	if err != nil {
		return nil, toConnectErr(err)
	}

	return connect.NewResponse(&protov1.DeleteAlbumResponse{}), nil
}

// TODO снести нахуй
func toProtoAlbum(a *domain.Album) *protov1.Album {
	if a == nil {
		return nil
	}

	return &protov1.Album{
		Id:          uuidStr(a.ID),
		OwnerId:     uuidStr(a.OwnerID),
		Title:       strOrEmpty(a.Title),
		Description: strOrEmpty(a.Description),
		ReleaseDate: tsPtr(a.ReleaseDate),
		CreatedAt:   tsPtr(a.CreatedAt),
		UpdatedAt:   tsPtr(a.UpdatedAt),
	}
}
