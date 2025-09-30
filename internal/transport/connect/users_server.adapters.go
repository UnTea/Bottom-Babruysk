package connect

import (
	"github.com/untea/bottom_babruysk/internal/domain"
	protov1 "github.com/untea/bottom_babruysk/internal/transport/gen/proto/v1"
)

// TODO сделать автогенерацию типов для grpc уровня
func toProtoUser(user *domain.User) *protov1.User {
	if user == nil {
		return nil
	}

	return &protov1.User{
		Id:           user.ID.String(),
		Email:        ValueOrZero(user.Email),
		PasswordHash: ValueOrZero(user.PasswordHash),
		DisplayName:  ValueOrZero(user.DisplayName),
		Role:         ToProtoRole(user.Role),
		CreatedAt:    TimeToTimestamppb(user.CreatedAt),
		UpdatedAt:    TimeToTimestamppb(user.UpdatedAt),
	}
}

func ToProtoRole(role *domain.Role) protov1.Role {
	if role == nil {
		return protov1.Role_ROLE_UNSPECIFIED
	}

	switch *role {
	case domain.RoleUser:
		return protov1.Role_ROLE_USER
	case domain.RoleAdmin:
		return protov1.Role_ROLE_ADMIN
	default:
		return protov1.Role_ROLE_UNSPECIFIED
	}
}

func FromProtoRole(role protov1.Role) *domain.Role {
	switch role {
	case protov1.Role_ROLE_USER:
		x := domain.RoleUser
		return &x
	case protov1.Role_ROLE_ADMIN:
		x := domain.RoleAdmin
		return &x
	default:
		x := domain.RoleUser
		return &x
	}
}
