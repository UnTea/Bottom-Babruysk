package connect

import (
	"time"

	"golang.org/x/exp/constraints"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TimeToTimestamppb(value *time.Time) *timestamppb.Timestamp {
	if value == nil {
		return nil
	}

	return timestamppb.New(*value)
}

func TimestamppbToTime(value *timestamppb.Timestamp) *time.Time {
	if value == nil {
		return nil
	}

	t := value.AsTime()

	return &t
}

func DurationToDurationpb(value *time.Duration) *durationpb.Duration {
	if value == nil {
		return nil
	}

	return durationpb.New(*value)
}

func DurationpbToDuration(value *durationpb.Duration) *time.Duration {
	if value == nil {
		return nil
	}

	t := value.AsDuration()

	return &t
}

func Int16ToInt32(value *int16) *int32 {
	if value == nil {
		return nil
	}

	num := int32(*value)

	return &num
}

func Int32ToInt16(value *int32) *int16 {
	if value == nil {
		return nil
	}

	num := int16(*value)

	return &num
}

func Integer[T2, T1 constraints.Integer](value T1) T2 {
	return T2(value)
}

func String[T2, T1 ~string](value T1) T2 {
	return T2(value)
}

func BoolToInt(value bool) int {
	if value {
		return 1
	}

	return 0
}

func IntToBool(value int) bool {
	return value != 0
}
