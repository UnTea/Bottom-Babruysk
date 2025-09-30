package connect

import (
	"github.com/untea/bottom_babruysk/internal/domain"
	protov1 "github.com/untea/bottom_babruysk/internal/transport/gen/proto/v1"
)

func toProtoPlaylist(p *domain.Playlist) *protov1.Playlist {
	if p == nil {
		return nil
	}
	return &protov1.Playlist{
		Id:          UUIDPtrToString(p.ID),
		OwnerId:     UUIDPtrToString(p.OwnerID),
		Title:       ValueOrZero(p.Title),
		Description: ValueOrZero(p.Description),
		Visibility:  DomainVisibilityToProto(p.Visibility),
		CreatedAt:   TimeToTimestamppb(p.CreatedAt),
		UpdatedAt:   TimeToTimestamppb(p.UpdatedAt),
	}
}
