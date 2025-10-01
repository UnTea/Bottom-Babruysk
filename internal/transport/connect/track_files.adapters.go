package connect

import (
	"github.com/untea/bottom_babruysk/internal/domain"
	protov1 "github.com/untea/bottom_babruysk/internal/transport/gen/proto/v1"
	"github.com/untea/bottom_babruysk/utils"
)

func toProtoTrackFile(trackFile *domain.TrackFile) *protov1.TrackFile {
	if trackFile == nil {
		return nil
	}

	return &protov1.TrackFile{
		Id:         trackFile.ID.String(),
		TrackId:    trackFile.TrackID.String(),
		Filename:   utils.ValueOrZero(trackFile.Filename),
		S3Key:      trackFile.S3Key,
		Mime:       utils.ValueOrZero(trackFile.Mime),
		Format:     ToProtoFormat(trackFile.Format),
		Codec:      ToProtoCodec(trackFile.Codec),
		Bitrate:    utils.ValueOrZero(utils.IntToInt32(trackFile.Bitrate)),
		SampleRate: utils.ValueOrZero(utils.IntToInt32(trackFile.SampleRate)),
		Channels:   utils.ValueOrZero(utils.IntToInt32(trackFile.Channels)),
		Size:       utils.ValueOrZero(trackFile.Size),
		Duration:   utils.DurationToDurationpb(trackFile.Duration),
		Checksum:   utils.ValueOrZero(trackFile.Checksum),
		CreatedAt:  utils.TimeToTimestamppb(trackFile.CreatedAt),
		UpdatedAt:  utils.TimeToTimestamppb(trackFile.UpdatedAt),
		UploadedAt: utils.TimeToTimestamppb(trackFile.UploadedAt),
	}
}

func ToProtoFormat(format *domain.Format) protov1.Format {
	if format == nil {
		return protov1.Format_FORMAT_UNSPECIFIED
	}

	switch *format {
	case domain.FormatMP3:
		return protov1.Format_FORMAT_MP3
	case domain.FormatMP4:
		return protov1.Format_FORMAT_MP4
	case domain.FormatM4A:
		return protov1.Format_FORMAT_M4A
	case domain.FormatOGG:
		return protov1.Format_FORMAT_OGG
	case domain.FormatFLAC:
		return protov1.Format_FORMAT_FLAC
	case domain.FormatWAV:
		return protov1.Format_FORMAT_WAV
	case domain.FormatWEBM:
		return protov1.Format_FORMAT_WEBM
	case domain.FormatAAC:
		return protov1.Format_FORMAT_AAC
	default:
		return protov1.Format_FORMAT_UNSPECIFIED
	}
}

func FromProtoFormat(format protov1.Format) *domain.Format {
	switch format {
	case protov1.Format_FORMAT_MP3:
		x := domain.FormatMP3
		return &x
	case protov1.Format_FORMAT_MP4:
		x := domain.FormatMP4
		return &x
	case protov1.Format_FORMAT_M4A:
		x := domain.FormatM4A
		return &x
	case protov1.Format_FORMAT_OGG:
		x := domain.FormatOGG
		return &x
	case protov1.Format_FORMAT_FLAC:
		x := domain.FormatFLAC
		return &x
	case protov1.Format_FORMAT_WAV:
		x := domain.FormatWAV
		return &x
	case protov1.Format_FORMAT_WEBM:
		x := domain.FormatWEBM
		return &x
	case protov1.Format_FORMAT_AAC:
		x := domain.FormatAAC
		return &x
	default:
		x := domain.FormatUnspecified
		return &x
	}
}

func ToProtoCodec(codec *domain.Codec) protov1.Codec {
	if codec == nil {
		return protov1.Codec_CODEC_UNSPECIFIED
	}

	switch *codec {
	case domain.CodecWAV:
		return protov1.Codec_CODEC_WMA
	case domain.CodecFLAC:
		return protov1.Codec_CODEC_FLAC
	case domain.CodecALAC:
		return protov1.Codec_CODEC_ALAC
	case domain.CodecAPE:
		return protov1.Codec_CODEC_APE
	case domain.CodecSHN:
		return protov1.Codec_CODEC_SHN
	case domain.CodecMP3:
		return protov1.Codec_CODEC_MP3
	case domain.CodecAAC:
		return protov1.Codec_CODEC_AAC
	case domain.CodecOPUS:
		return protov1.Codec_CODEC_OPUS
	case domain.CodecVORBIS:
		return protov1.Codec_CODEC_VORBIS
	case domain.CodecWMA:
		return protov1.Codec_CODEC_WMA
	default:
		return protov1.Codec_CODEC_UNSPECIFIED
	}
}

func FromProtoCodec(codec protov1.Codec) *domain.Codec {
	switch codec {
	case protov1.Codec_CODEC_WAV:
		x := domain.CodecWAV
		return &x
	case protov1.Codec_CODEC_FLAC:
		x := domain.CodecFLAC
		return &x
	case protov1.Codec_CODEC_ALAC:
		x := domain.CodecALAC
		return &x
	case protov1.Codec_CODEC_APE:
		x := domain.CodecAPE
		return &x
	case protov1.Codec_CODEC_SHN:
		x := domain.CodecSHN
		return &x
	case protov1.Codec_CODEC_MP3:
		x := domain.CodecMP3
		return &x
	case protov1.Codec_CODEC_AAC:
		x := domain.CodecAAC
		return &x
	case protov1.Codec_CODEC_OPUS:
		x := domain.CodecOPUS
		return &x
	case protov1.Codec_CODEC_VORBIS:
		x := domain.CodecVORBIS
		return &x
	case protov1.Codec_CODEC_WMA:
		x := domain.CodecWMA
		return &x
	default:
		x := domain.CodecUnspecified
		return &x
	}
}
