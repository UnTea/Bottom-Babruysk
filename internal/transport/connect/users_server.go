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
		Email:        Ptr(request.Msg.Email),
		PasswordHash: Ptr(request.Msg.PasswordHash),
		DisplayName:  Ptr(request.Msg.DisplayName),
		Role:         ProtoRolePtrToDomain(Ptr(request.Msg.Role)),
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
		ID: StringToUUIDPtr(request.Msg.Id),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}

	return connect.NewResponse(&protov1.GetUserResponse{User: toProtoUser(response.User)}), nil
}

func (s *UsersServer) ListUsers(ctx context.Context, request *connect.Request[protov1.ListUsersRequest]) (*connect.Response[protov1.ListUsersResponse], error) {
	response, err := s.usersService.ListUsers(ctx, domain.ListUsersRequest{
		Limit:     Ptr(int(request.Msg.Limit)),
		Offset:    Ptr(int(request.Msg.Offset)),
		Role:      ProtoRolePtrToDomain(Ptr(request.Msg.Role)),
		SortField: Ptr(request.Msg.SortField),
		SortOrder: Ptr(request.Msg.SortOrder),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}

	result := &protov1.ListUsersResponse{Users: make([]*protov1.User, 0, len(response.Users))}

	for _, user := range response.Users {
		result.Users = append(result.Users, toProtoUser(user))
	}

	return connect.NewResponse(result), nil
}

func (s *UsersServer) UpdateUser(ctx context.Context, request *connect.Request[protov1.UpdateUserRequest]) (*connect.Response[protov1.UpdateUserResponse], error) {
	err := s.usersService.UpdateUser(ctx, domain.UpdateUserRequest{
		ID:          StringToUUIDPtr(request.Msg.Id),
		DisplayName: request.Msg.DisplayName,
		Role:        ProtoRolePtrToDomain(request.Msg.Role),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}

	return connect.NewResponse(&protov1.UpdateUserResponse{}), nil
}

func (s *UsersServer) DeleteUser(ctx context.Context, request *connect.Request[protov1.DeleteUserRequest]) (*connect.Response[protov1.DeleteUserResponse], error) {
	err := s.usersService.DeleteUser(ctx, domain.DeleteUserRequest{
		ID: StringToUUIDPtr(request.Msg.Id),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}

	return connect.NewResponse(&protov1.DeleteUserResponse{}), nil
}
