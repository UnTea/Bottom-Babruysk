package connect

import (
	"context"

	"connectrpc.com/connect"

	"github.com/untea/bottom_babruysk/internal/domain"
	"github.com/untea/bottom_babruysk/internal/service"
	protov1 "github.com/untea/bottom_babruysk/internal/transport/gen/proto/v1"
)

type ArtistsServer struct {
	svc *service.ArtistsService
}

func NewArtistsServer(svc *service.ArtistsService) *ArtistsServer { return &ArtistsServer{svc: svc} }

func (s *ArtistsServer) CreateArtist(ctx context.Context, req *connect.Request[protov1.CreateArtistRequest]) (*connect.Response[protov1.CreateArtistResponse], error) {
	out, err := s.svc.CreateArtist(ctx, domain.CreateArtistRequest{
		Name: Ptr(req.Msg.Name),
		Bio:  Ptr(req.Msg.Bio),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}

	return connect.NewResponse(&protov1.CreateArtistResponse{Id: out.ID.String()}), nil
}

func (s *ArtistsServer) GetArtist(ctx context.Context, request *connect.Request[protov1.GetArtistRequest]) (*connect.Response[protov1.GetArtistResponse], error) {
	out, err := s.svc.GetArtist(ctx, domain.GetArtistRequest{
		ID: StringToUUIDPtr(request.Msg.Id),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}

	return connect.NewResponse(&protov1.GetArtistResponse{Artist: toProtoArtist(out.Artist)}), nil
}

func (s *ArtistsServer) ListArtists(ctx context.Context, request *connect.Request[protov1.ListArtistsRequest]) (*connect.Response[protov1.ListArtistsResponse], error) {
	response, err := s.svc.ListArtists(ctx, domain.ListArtistsRequest{
		Name:      Ptr(request.Msg.Name),
		Bio:       Ptr(request.Msg.Bio),
		Limit:     Ptr(int(request.Msg.Limit)),
		Offset:    Ptr(int(request.Msg.Offset)),
		SortField: Ptr(request.Msg.SortField),
		SortOrder: Ptr(request.Msg.SortOrder),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}

	result := &protov1.ListArtistsResponse{Artists: make([]*protov1.Artist, 0, len(response.Artists))}

	for _, artist := range response.Artists {
		result.Artists = append(result.Artists, toProtoArtist(artist))
	}

	return connect.NewResponse(result), nil
}

func (s *ArtistsServer) UpdateArtist(ctx context.Context, request *connect.Request[protov1.UpdateArtistRequest]) (*connect.Response[protov1.UpdateArtistResponse], error) {
	err := s.svc.UpdateArtist(ctx, domain.UpdateArtistRequest{
		ID:   StringToUUIDPtr(request.Msg.Id),
		Name: request.Msg.Name,
		Bio:  request.Msg.Bio,
	})
	if err != nil {
		return nil, toConnectErr(err)
	}

	return connect.NewResponse(&protov1.UpdateArtistResponse{}), nil
}

func (s *ArtistsServer) DeleteArtist(ctx context.Context, request *connect.Request[protov1.DeleteArtistRequest]) (*connect.Response[protov1.DeleteArtistResponse], error) {
	err := s.svc.DeleteArtist(ctx, domain.DeleteArtistRequest{
		ID: StringToUUIDPtr(request.Msg.Id),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}

	return connect.NewResponse(&protov1.DeleteArtistResponse{}), nil
}
