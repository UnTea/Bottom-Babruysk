package connect

import (
	"github.com/untea/bottom_babruysk/internal/domain"
	protov1 "github.com/untea/bottom_babruysk/internal/transport/gen/proto/v1"
	"github.com/untea/bottom_babruysk/utils"
)

func toProtoArtist(artist *domain.Artist) *protov1.Artist {
	if artist == nil {
		return nil
	}

	return &protov1.Artist{
		Id:        artist.ID.String(),
		Name:      utils.ValueOrZero(artist.Name),
		Bio:       utils.ValueOrZero(artist.Bio),
		CreatedAt: utils.TimeToTimestamppb(artist.CreatedAt),
		UpdatedAt: utils.TimeToTimestamppb(artist.UpdatedAt),
	}
}
