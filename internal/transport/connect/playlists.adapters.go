package connect

import (
	"github.com/untea/bottom_babruysk/internal/domain"
	protov1 "github.com/untea/bottom_babruysk/internal/transport/gen/proto/v1"
	"github.com/untea/bottom_babruysk/utils"
)

func toProtoPlaylist(p *domain.Playlist) *protov1.Playlist {
	if p == nil {
		return nil
	}
	return &protov1.Playlist{
		Id:          utils.UUIDPtrToString(p.ID),
		OwnerId:     utils.UUIDPtrToString(p.OwnerID),
		Title:       utils.ValueOrZero(p.Title),
		Description: utils.ValueOrZero(p.Description),
		Visibility:  ToProtoVisibility(p.Visibility),
		CreatedAt:   utils.TimeToTimestamppb(p.CreatedAt),
		UpdatedAt:   utils.TimeToTimestamppb(p.UpdatedAt),
	}
}
