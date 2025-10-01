package connect

import (
	"github.com/untea/bottom_babruysk/internal/domain"
	protov1 "github.com/untea/bottom_babruysk/internal/transport/gen/proto/v1"
	"github.com/untea/bottom_babruysk/utils"
)

// TODO сделать автогенерацию типов для grpc уровня
func toProtoAlbum(album *domain.Album) *protov1.Album {
	if album == nil {
		return nil
	}

	return &protov1.Album{
		Id:          album.ID.String(),
		OwnerId:     album.OwnerID.String(),
		Title:       utils.ValueOrZero(album.Title),
		Description: utils.ValueOrZero(album.Description),
		ReleaseDate: utils.TimeToTimestamppb(album.ReleaseDate),
		CreatedAt:   utils.TimeToTimestamppb(album.CreatedAt),
		UpdatedAt:   utils.TimeToTimestamppb(album.UpdatedAt),
	}
}
