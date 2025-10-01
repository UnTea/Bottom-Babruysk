package connect

import (
	"context"

	"connectrpc.com/connect"

	"github.com/untea/bottom_babruysk/internal/domain"
	"github.com/untea/bottom_babruysk/internal/service"
	protov1 "github.com/untea/bottom_babruysk/internal/transport/gen/proto/v1"
	"github.com/untea/bottom_babruysk/utils"
)

type UsersServer struct {
	usersService *service.UsersService
}

func NewUsersServer(usersService *service.UsersService) *UsersServer {
	return &UsersServer{usersService: usersService}
}

func (s *UsersServer) CreateUser(ctx context.Context, request *connect.Request[protov1.CreateUserRequest]) (*connect.Response[protov1.CreateUserResponse], error) {
	response, err := s.usersService.CreateUser(ctx, domain.CreateUserRequest{
		Email:        utils.Ptr(request.Msg.Email),
		PasswordHash: utils.Ptr(request.Msg.PasswordHash),
		DisplayName:  utils.Ptr(request.Msg.DisplayName),
		Role:         FromProtoRole(request.Msg.Role),
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
		ID: utils.StringToUUIDPtr(request.Msg.Id),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}

	return connect.NewResponse(&protov1.GetUserResponse{User: toProtoUser(response.User)}), nil
}

func (s *UsersServer) ListUsers(ctx context.Context, request *connect.Request[protov1.ListUsersRequest]) (*connect.Response[protov1.ListUsersResponse], error) {
	response, err := s.usersService.ListUsers(ctx, domain.ListUsersRequest{
		Limit:     utils.Ptr(int(request.Msg.Limit)),
		Offset:    utils.Ptr(int(request.Msg.Offset)),
		Role:      FromProtoRole(request.Msg.Role),
		SortField: utils.Ptr(request.Msg.SortField),
		SortOrder: utils.Ptr(request.Msg.SortOrder),
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
		ID:          utils.StringToUUIDPtr(request.Msg.Id),
		DisplayName: request.Msg.DisplayName,
		Role:        FromProtoRole(*request.Msg.Role),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}

	return connect.NewResponse(&protov1.UpdateUserResponse{}), nil
}

func (s *UsersServer) DeleteUser(ctx context.Context, request *connect.Request[protov1.DeleteUserRequest]) (*connect.Response[protov1.DeleteUserResponse], error) {
	err := s.usersService.DeleteUser(ctx, domain.DeleteUserRequest{
		ID: utils.StringToUUIDPtr(request.Msg.Id),
	})
	if err != nil {
		return nil, toConnectErr(err)
	}

	return connect.NewResponse(&protov1.DeleteUserResponse{}), nil
}
