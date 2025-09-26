package validation

import (
	"fmt"
	"reflect"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Set[T comparable] map[T]struct{}

func NewSet[T comparable](vals ...T) Set[T] {
	s := make(Set[T], len(vals))
	for _, v := range vals {
		s[v] = struct{}{}
	}

	return s
}

func (s Set[T]) Has(v T) bool {
	_, ok := s[v]

	return ok
}

func (s Set[string]) StringKeys() []string {
	out := make([]string, 0, len(s))
	for k := range s {
		out = append(out, k)
	}

	return out
}

func InSetPtr[T comparable](set Set[T]) validation.Rule {
	return validation.By(func(v any) error {
		if v == nil {
			return nil
		}

		switch vv := v.(type) {
		case *T:
			if vv == nil {
				return nil
			}

			if !set.Has(*vv) {
				return validation.NewError("validation", "invalid value")
			}

			return nil

		default:
			rv := reflect.ValueOf(v)
			if rv.Kind() == reflect.Ptr && !rv.IsNil() {
				rv = rv.Elem()
			}

			if rv.IsValid() && rv.Type() == reflect.TypeOf(*new(T)) {
				val := rv.Interface().(T)
				if !set.Has(val) {
					return validation.NewError("validation", "invalid value")
				}

				return nil
			}

			return nil
		}
	})
}

func InStringsPtr(set Set[string], fieldName string) validation.Rule {
	return validation.By(func(v any) error {
		if v == nil {
			return nil
		}

		var s string

		switch vv := v.(type) {
		case *string:
			if vv == nil {
				return nil
			}

			s = *vv

		case string:
			s = vv

		default:
			return nil
		}

		if !set.Has(strings.ToLower(s)) {
			allowed := strings.Join(set.StringKeys(), ", ")
			return validation.NewError("validation", fmt.Sprintf("%s must be one of: %s", fieldName, allowed))
		}

		return nil
	})
}
