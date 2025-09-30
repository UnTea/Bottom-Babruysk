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
		Role:         DomainRoleToProto(user.Role),
		CreatedAt:    TimeToTimestamppb(user.CreatedAt),
		UpdatedAt:    TimeToTimestamppb(user.UpdatedAt),
	}
}

func ProtoRolePtrToDomain(p *protov1.Role) *domain.Role {
	if p == nil {
		return nil
	}

	v := domain.Role(*p)

	return &v
}

func DomainRoleToProto(role *domain.Role) protov1.Role {
	if role == nil {
		return protov1.Role_ROLE_UNSPECIFIED
	}

	return protov1.Role(*role)
}
