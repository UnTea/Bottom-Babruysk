package connect

import (
	"github.com/untea/bottom_babruysk/internal/domain"
	protov1 "github.com/untea/bottom_babruysk/internal/transport/gen/proto/v1"
)

func toProtoArtist(artist *domain.Artist) *protov1.Artist {
	if artist == nil {
		return nil
	}

	return &protov1.Artist{
		Id:        artist.ID.String(),
		Name:      ValueOrZero(artist.Name),
		Bio:       ValueOrZero(artist.Bio),
		CreatedAt: TimeToTimestamppb(artist.CreatedAt),
		UpdatedAt: TimeToTimestamppb(artist.UpdatedAt),
	}
}
