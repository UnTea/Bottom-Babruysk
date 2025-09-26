package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"

	"github.com/untea/bottom_babruysk/internal/repository"
)

type Empty struct{}

func (h *Handler) toHTTPStatus(err error) int {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return http.StatusNotFound
	case isValidationErr(err):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func isValidationErr(err error) bool {
	if err == nil {
		return false
	}

	var ve validation.Errors

	if errors.As(err, &ve) {
		return true
	}

	var e *validation.Error

	return errors.As(err, e)
}

func Lift[R any](f func(ctx context.Context, request R) error) func(ctx context.Context, request R) (Empty, error) {
	return func(ctx context.Context, request R) (Empty, error) {
		return Empty{}, f(ctx, request)
	}
}

func Handle[R any, T any](h *Handler, action func(ctx context.Context, request R) (T, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := Decode[R](r)
		if err != nil {
			h.httpError(w, err, http.StatusBadRequest)
			return
		}

		response, err := action(r.Context(), request)
		if err != nil {
			h.httpError(w, err, h.toHTTPStatus(err))
			return
		}

		h.writeJson(w, response, http.StatusOK)
	}
}

func Decode[T any](r *http.Request) (T, error) {
	var out T

	if r.Body != nil {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			return out, fmt.Errorf("read body: %w", err)
		}

		_ = r.Body.Close()

		if len(data) > 0 {
			if err := json.Unmarshal(data, &out); err != nil {
				return out, fmt.Errorf("unmarshal body: %w", err)
			}
		}
	}

	err := fillFrom(reflect.ValueOf(&out).Elem(), r.URL.Query(), queryName, "query")
	if err != nil {
		return out, fmt.Errorf("decode query: %w", err)
	}

	rc := chi.RouteContext(r.Context())
	if rc != nil {
		uv := url.Values{}
		for i := range rc.URLParams.Keys {
			uv.Set(rc.URLParams.Keys[i], rc.URLParams.Values[i])
		}

		err = fillFrom(reflect.ValueOf(&out).Elem(), uv, pathName, "path")
		if err != nil {
			return out, fmt.Errorf("decode path: %w", err)
		}
	}

	return out, nil
}

type nameFunc = func(reflect.StructField) (name string, skip bool)

func fillFrom(reflectValue reflect.Value, urlValues url.Values, nameFunc nameFunc, label string) error {
	if reflectValue.Kind() == reflect.Pointer {
		if reflectValue.IsNil() {
			reflectValue.Set(reflect.New(reflectValue.Type().Elem()))
		}

		return fillFrom(reflectValue.Elem(), urlValues, nameFunc, label)
	}

	if reflectValue.Kind() != reflect.Struct {
		return nil
	}

	t := reflectValue.Type()
	for i := 0; i < t.NumField(); i++ {
		structField := t.Field(i)
		if structField.PkgPath != "" {
			continue
		}

		fieldValue := reflectValue.Field(i)
		fieldType := structField.Type

		name, skip := nameFunc(structField)
		if skip {
			continue
		}

		isStruct := fieldType.Kind() == reflect.Struct
		isStructBehindPointer := fieldType.Kind() == reflect.Pointer && fieldType.Elem().Kind() == reflect.Struct

		if name == "" && (isStruct || isStructBehindPointer) {
			if fieldValue.Kind() == reflect.Pointer && fieldValue.IsNil() {
				fieldValue.Set(reflect.New(fieldType.Elem()))
			}

			err := fillFrom(indirect(fieldValue), urlValues, nameFunc, label)
			if err != nil {
				return err
			}

			continue
		}

		if name == "" {
			name = toSnakeCase(structField.Name)
		}

		raw := urlValues.Get(name)
		if raw == "" {
			continue
		}

		err := setValue(fieldValue, raw)
		if err != nil {
			return fmt.Errorf("%s param %q: %w", label, name, err)
		}
	}

	return nil
}

func queryName(sf reflect.StructField) (string, bool) {
	if tag, ok, skip := pickTag(sf, "query"); ok {
		return tag, skip
	}

	if isJSONOmitted(sf) {
		return "", true
	}

	return "", false
}

func pathName(sf reflect.StructField) (string, bool) {
	if tag, ok, skip := pickTag(sf, "path"); ok {
		return tag, skip
	}
	if isJSONOmitted(sf) {
		return "", true
	}
	return "", false
}

func pickTag(sf reflect.StructField, key string) (name string, has bool, skip bool) {
	tag := sf.Tag.Get(key)
	if tag == "" {
		return "", false, false
	}
	if tag == "-" {
		return "", true, true
	}
	if idx := strings.IndexByte(tag, ','); idx >= 0 {
		tag = tag[:idx]
	}
	return tag, true, false
}

func isJSONOmitted(sf reflect.StructField) bool {
	j := sf.Tag.Get("json")

	return j == "-"
}

func indirect(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Pointer {
		return v.Elem()
	}

	return v
}

func setValue(dst reflect.Value, s string) error {
	if dst.Kind() == reflect.Pointer {
		if dst.IsNil() {
			dst.Set(reflect.New(dst.Type().Elem()))
		}

		dst = dst.Elem()
	}

	switch dst.Kind() {
	case reflect.String:
		dst.SetString(s)
		return nil

	case reflect.Bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}

		dst.SetBool(b)

		return nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if dst.Type().PkgPath() == "time" && dst.Type().Name() == "Duration" {
			d, err := time.ParseDuration(s)
			if err != nil {
				return err
			}

			dst.SetInt(int64(d))

			return nil
		}

		n, err := strconv.ParseInt(s, 10, dst.Type().Bits())
		if err != nil {
			return err
		}

		dst.SetInt(n)

		return nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(s, 10, dst.Type().Bits())
		if err != nil {
			return err
		}

		dst.SetUint(n)

		return nil

	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(s, dst.Type().Bits())
		if err != nil {
			return err
		}

		dst.SetFloat(f)

		return nil

	case reflect.Struct:
		if dst.Type().PkgPath() == "time" && dst.Type().Name() == "Time" {
			tm, err := time.Parse(time.RFC3339, s)
			if err != nil {
				return err
			}

			dst.Set(reflect.ValueOf(tm))

			return nil
		}

		if dst.Type().PkgPath() == "github.com/google/uuid" && dst.Type().Name() == "UUID" {
			if err := uuid.Validate(s); err != nil {
				return err
			}

			u := uuid.MustParse(s)
			dst.Set(reflect.ValueOf(u))

			return nil
		}

		under := dst.Type()
		if under.Kind() == reflect.String {
			dst.SetString(s)

			return nil
		}

		return fmt.Errorf("unsupported struct type %s", dst.Type())

	default:
		under := dst.Type()

		switch under.Kind() {
		case reflect.String:
			dst.SetString(s)

			return nil

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n, err := strconv.ParseInt(s, 10, under.Bits())
			if err != nil {
				return err
			}

			dst.SetInt(n)

			return nil
		}

		return fmt.Errorf("unsupported kind %s", dst.Kind())
	}
}

func toSnakeCase(s string) string {
	var b strings.Builder

	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				prev := rune(s[i-1])
				if prev != '_' && (unicode.IsLower(prev) || unicode.IsDigit(prev)) {
					b.WriteByte('_')
				}
			}

			b.WriteRune(unicode.ToLower(r))
		} else {
			b.WriteRune(r)
		}
	}

	return b.String()
}
