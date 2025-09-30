package connect

import (
	"github.com/untea/bottom_babruysk/internal/domain"
	protov1 "github.com/untea/bottom_babruysk/internal/transport/gen/proto/v1"
)

// TODO сделать автогенерацию типов для grpc уровня
func toProtoAlbum(album *domain.Album) *protov1.Album {
	if album == nil {
		return nil
	}

	return &protov1.Album{
		Id:          album.ID.String(),
		OwnerId:     album.OwnerID.String(),
		Title:       ValueOrZero(album.Title),
		Description: ValueOrZero(album.Description),
		ReleaseDate: TimeToTimestamppb(album.ReleaseDate),
		CreatedAt:   TimeToTimestamppb(album.CreatedAt),
		UpdatedAt:   TimeToTimestamppb(album.UpdatedAt),
	}
}
