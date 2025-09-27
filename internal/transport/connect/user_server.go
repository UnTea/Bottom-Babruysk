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

func NewUsersServer(svc *service.UsersService) *UsersServer {
	return &UsersServer{usersService: svc}
}

func (s *UsersServer) CreateUser(ctx context.Context, request *connect.Request[protov1.CreateUserRequest]) (*connect.Response[protov1.CreateUserResponse], error) {
	role := (*domain.UserRole)(nil)
	if r := request.Msg.GetRole(); r != "" {
		rr := domain.UserRole(r)
		role = &rr
	}

	resp, err := s.usersService.CreateUser(ctx, domain.CreateUserRequest{
		Email:        strPtrOrNil(request.Msg.GetEmail()),
		PasswordHash: strPtrOrNil(request.Msg.GetPasswordHash()),
		DisplayName:  strPtrOrNil(request.Msg.GetDisplayName()),
		Role:         role,
	})

	if err != nil {
		return nil, toConnectErr(err)
	}

	return connect.NewResponse(&protov1.CreateUserResponse{Id: uuidStr(resp.ID)}), nil
}

func (s *UsersServer) GetUser(ctx context.Context, request *connect.Request[protov1.GetUserRequest]) (*connect.Response[protov1.GetUserResponse], error) {
	out, err := s.usersService.GetUser(ctx, domain.GetUserRequest{ID: uuidFromStr(request.Msg.GetId())})
	if err != nil {
		return nil, toConnectErr(err)
	}

	return connect.NewResponse(&protov1.GetUserResponse{User: toProtoUser(out.User)}), nil
}

func (s *UsersServer) ListUsers(ctx context.Context, request *connect.Request[protov1.ListUsersRequest]) (*connect.Response[protov1.ListUsersResponse], error) {
	var role *domain.UserRole

	if r := request.Msg.GetRole(); r != "" {
		rr := domain.UserRole(r)
		role = &rr
	}

	out, err := s.usersService.ListUsers(ctx, domain.ListUsersRequest{
		Limit:     intPtr(int(request.Msg.GetLimit())),
		Offset:    intPtr(int(request.Msg.GetOffset())),
		Role:      role,
		SortField: strPtrOrNil(request.Msg.GetSortField()),
		SortOrder: strPtrOrNil(request.Msg.GetSortOrder()),
	})

	if err != nil {
		return nil, toConnectErr(err)
	}

	response := &protov1.ListUsersResponse{Users: make([]*protov1.User, 0, len(out.Users))}

	for _, u := range out.Users {
		response.Users = append(response.Users, toProtoUser(u))
	}

	return connect.NewResponse(response), nil
}

func (s *UsersServer) UpdateUser(ctx context.Context, request *connect.Request[protov1.UpdateUserRequest]) (*connect.Response[protov1.UpdateUserResponse], error) {
	var role *domain.UserRole

	if request.Msg.Role != nil {
		rr := domain.UserRole(request.Msg.GetRole())
		role = &rr
	}

	err := s.usersService.UpdateUser(ctx, domain.UpdateUserRequest{
		ID:          uuidPtrFromStr(request.Msg.GetId()),
		DisplayName: request.Msg.DisplayName,
		Role:        role,
	})

	if err != nil {
		return nil, toConnectErr(err)
	}

	return connect.NewResponse(&protov1.UpdateUserResponse{}), nil
}

func (s *UsersServer) DeleteUser(ctx context.Context, request *connect.Request[protov1.DeleteUserRequest]) (*connect.Response[protov1.DeleteUserResponse], error) {
	err := s.usersService.DeleteUser(ctx, domain.DeleteUserRequest{ID: uuidFromStr(request.Msg.GetId())})
	if err != nil {
		return nil, toConnectErr(err)
	}

	return connect.NewResponse(&protov1.DeleteUserResponse{}), nil
}

// TODO снести нахуй
func toProtoUser(user *domain.User) *protov1.User {
	if user == nil {
		return nil
	}

	return &protov1.User{
		Id:           uuidStr(user.ID),
		Email:        strOrEmpty(user.Email),
		PasswordHash: strOrEmpty(user.PasswordHash),
		DisplayName:  strOrEmpty(user.DisplayName),
		Role:         string(derefUserRole(user.Role)),
		CreatedAt:    tsPtr(user.CreatedAt),
		UpdatedAt:    tsPtr(user.UpdatedAt),
	}
}
