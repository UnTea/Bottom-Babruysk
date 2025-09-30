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
		Visibility:  ToProtoVisibility(track.Visibility),
		CreatedAt:   TimeToTimestamppb(track.CreatedAt),
		UpdatedAt:   TimeToTimestamppb(track.UpdatedAt),
		UploadedAt:  TimeToTimestamppb(track.UploadedAt),
	}
}

func ToProtoVisibility(visibility *domain.Visibility) protov1.Visibility {
	if visibility == nil {
		return protov1.Visibility_VISIBILITY_UNSPECIFIED
	}

	switch *visibility {
	case domain.VisibilityPrivate:
		return protov1.Visibility_VISIBILITY_PRIVATE
	case domain.VisibilityUnlisted:
		return protov1.Visibility_VISIBILITY_UNLISTED
	case domain.VisibilityPublic:
		return protov1.Visibility_VISIBILITY_PUBLIC
	default:
		return protov1.Visibility_VISIBILITY_UNSPECIFIED
	}
}

func FromProtoVisibility(visibility protov1.Visibility) *domain.Visibility {
	switch visibility {
	case protov1.Visibility_VISIBILITY_PRIVATE:
		x := domain.VisibilityPrivate
		return &x
	case protov1.Visibility_VISIBILITY_UNLISTED:
		x := domain.VisibilityUnlisted
		return &x
	case protov1.Visibility_VISIBILITY_PUBLIC:
		x := domain.VisibilityPublic
		return &x
	default:
		x := domain.VisibilityUnlisted
		return &x
	}
}
