package connect

import (
	"context"

	"connectrpc.com/connect"

	"github.com/untea/bottom_babruysk/internal/domain"
	"github.com/untea/bottom_babruysk/internal/service"
	protov1 "github.com/untea/bottom_babruysk/internal/transport/gen/proto/v1"
)

type UsersServer struct {
	usersService *service.UsersService
}

func NewUsersServer(usersService *service.UsersService) *UsersServer {
	return &UsersServer{usersService: usersService}
}

func (s *UsersServer) CreateUser(ctx context.Context, request *connect.Request[protov1.CreateUserRequest]) (*connect.Response[protov1.CreateUserResponse], error) {
	response, err := s.usersService.CreateUser(ctx, domain.CreateUserRequest{
		Email:        Ptr(request.Msg.GetEmail()),
		PasswordHash: Ptr(request.Msg.GetPasswordHash()),
		DisplayName:  Ptr(request.Msg.GetDisplayName()),
		Role:         (*domain.UserRole)(Ptr(request.Msg.GetRole())),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}

	return connect.NewResponse(&protov1.CreateUserResponse{
		Id: response.ID.String(),
	}), nil
}

func (s *UsersServer) GetUser(ctx context.Context, request *connect.Request[protov1.GetUserRequest]) (*connect.Response[protov1.GetUserResponse], error) {
	response, err := s.usersService.GetUser(ctx, domain.GetUserRequest{
		ID: StringToUUIDPtr(request.Msg.GetId()),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}

	return connect.NewResponse(&protov1.GetUserResponse{User: toProtoUser(response.User)}), nil
}

func (s *UsersServer) ListUsers(ctx context.Context, request *connect.Request[protov1.ListUsersRequest]) (*connect.Response[protov1.ListUsersResponse], error) {
	out, err := s.usersService.ListUsers(ctx, domain.ListUsersRequest{
		Limit:     Ptr(int(request.Msg.GetLimit())),
		Offset:    Ptr(int(request.Msg.GetOffset())),
		Role:      (*domain.UserRole)(Ptr(request.Msg.GetRole())),
		SortField: Ptr(request.Msg.GetSortField()),
		SortOrder: Ptr(request.Msg.GetSortOrder()),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}

	response := &protov1.ListUsersResponse{Users: make([]*protov1.User, 0, len(out.Users))}

	for _, user := range out.Users {
		response.Users = append(response.Users, toProtoUser(user))
	}

	return connect.NewResponse(response), nil
}

func (s *UsersServer) UpdateUser(ctx context.Context, request *connect.Request[protov1.UpdateUserRequest]) (*connect.Response[protov1.UpdateUserResponse], error) {
	err := s.usersService.UpdateUser(ctx, domain.UpdateUserRequest{
		ID:          StringToUUIDPtr(request.Msg.GetId()),
		DisplayName: Ptr(request.Msg.GetDisplayName()),
		Role:        (*domain.UserRole)(Ptr(request.Msg.GetRole())),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}

	return connect.NewResponse(&protov1.UpdateUserResponse{}), nil
}

func (s *UsersServer) DeleteUser(ctx context.Context, request *connect.Request[protov1.DeleteUserRequest]) (*connect.Response[protov1.DeleteUserResponse], error) {
	err := s.usersService.DeleteUser(ctx, domain.DeleteUserRequest{
		ID: StringToUUIDPtr(request.Msg.GetId()),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}

	return connect.NewResponse(&protov1.DeleteUserResponse{}), nil
}
