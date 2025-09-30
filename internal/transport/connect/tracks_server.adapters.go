package connect

import (
	"github.com/untea/bottom_babruysk/internal/domain"
	protov1 "github.com/untea/bottom_babruysk/internal/transport/gen/proto/v1"
)

// TODO сделать автогенерацию типов для grpc уровня
func toProtoTrack(track *domain.Track) *protov1.Track {
	if track == nil {
		return nil
	}

	return &protov1.Track{
		Id:          track.ID.String(),
		UploaderId:  track.UploaderID.String(),
		Title:       ValueOrZero(track.Title),
		Subtitle:    ValueOrZero(track.Subtitle),
		Description: ValueOrZero(track.Description),
		Duration:    DurationToDurationpb(track.Duration),
		Visibility:  DomainVisibilityToProto(track.Visibility),
		CreatedAt:   TimeToTimestamppb(track.CreatedAt),
		UpdatedAt:   TimeToTimestamppb(track.UpdatedAt),
		UploadedAt:  TimeToTimestamppb(track.UploadedAt),
	}
}

func ProtoVisibilityPtrToDomain(visibility *protov1.Visibility) *domain.Visibility {
	if visibility == nil {
		return nil
	}

	v := domain.Visibility(*visibility)

	return &v
}

func DomainVisibilityToProto(visibility *domain.Visibility) protov1.Visibility {
	if visibility == nil {
		return protov1.Visibility_VISIBILITY_UNSPECIFIED
	}

	return protov1.Visibility(*visibility)
}
