package handlers

import (
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
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/untea/bottom_babruysk/internal/domain"
	"github.com/untea/bottom_babruysk/internal/service"
)

type Handler struct {
	Logger *zap.Logger
	Repo   struct {
		Users service.Users
	}
}

func New(logger *zap.Logger, repos struct{ Users service.Users }) *Handler {
	h := &Handler{
		Logger: logger,
		Repo:   repos,
	}

	return h
}

func (h *Handler) writeJson(w http.ResponseWriter, resp any, code int) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	if resp == nil {
		w.WriteHeader(code)
		return
	}

	response, err := json.Marshal(resp)
	if err != nil {
		h.Logger.Error("failed to marshal json response", zap.Error(err))

		h.httpError(w, errors.New("internal error"), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(code)

	_, err = w.Write(response)
	if err != nil {
		h.Logger.Error("failed to write response", zap.Error(err))
	}
}

func (h *Handler) httpError(w http.ResponseWriter, err error, code int) {
	errorDetail := ""
	if err != nil {
		errorDetail = err.Error()
	}

	h.writeJson(w, domain.ErrorResponse{Error: errorDetail}, code)
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

	if err := tryFillStruct(reflect.ValueOf(&out).Elem(), r.URL.Query(), queryName, "query"); err != nil {
		return out, fmt.Errorf("decode query: %w", err)
	}

	if routeCtx := chi.RouteContext(r.Context()); routeCtx != nil {
		uv := url.Values{}

		for i := range routeCtx.URLParams.Keys {
			uv.Set(routeCtx.URLParams.Keys[i], routeCtx.URLParams.Values[i])
		}

		if err := tryFillStruct(reflect.ValueOf(&out).Elem(), uv, pathName, "path"); err != nil {
			return out, fmt.Errorf("decode path: %w", err)
		}
	}

	return out, nil
}

func tryFillStruct(reflectValue reflect.Value, urlValues url.Values, nameOf func(reflect.StructField) string, label string) error {
	switch reflectValue.Kind() {
	case reflect.Pointer:
		if reflectValue.IsNil() {
			reflectValue.Set(reflect.New(reflectValue.Type().Elem()))
		}

		return tryFillStruct(reflectValue.Elem(), urlValues, nameOf, label)
	case reflect.Struct:
	default:
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

		name := nameOf(structField)
		if name == "-" {
			continue
		}

		if name == "" && (fieldType.Kind() == reflect.Struct || (fieldType.Kind() == reflect.Pointer && fieldType.Elem().Kind() == reflect.Struct)) {
			if fieldValue.Kind() == reflect.Pointer && fieldValue.IsNil() {
				fieldValue.Set(reflect.New(fieldType.Elem()))
			}

			if err := tryFillStruct(indirectValue(fieldValue), urlValues, nameOf, label); err != nil {
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

		if err := assignValue(fieldValue, raw); err != nil {
			return fmt.Errorf("%s param %q: %w", label, name, err)
		}
	}

	return nil
}

func queryName(structField reflect.StructField) string {
	if tag := structField.Tag.Get("query"); tag != "" && tag != "-" {
		return tag
	}

	if tag := structField.Tag.Get("json"); tag != "-" {
		return "-"
	}

	return ""
}

func pathName(structField reflect.StructField) string {
	if tag := structField.Tag.Get("path"); tag != "" && tag != "-" {
		return tag
	}

	if tag := structField.Tag.Get("json"); tag == "-" {
		return "-"
	}

	return ""
}

func assignValue(dst reflect.Value, s string) error {
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

	case reflect.Struct:
		if dst.Type().PkgPath() == "time" && dst.Type().Name() == "Time" {
			tm, err := time.Parse(time.RFC3339, s)
			if err != nil {
				return err
			}

			dst.Set(reflect.ValueOf(tm))

			return nil
		}

		return fmt.Errorf("unsupported struct type %s", dst.Type())

	case reflect.Array:
		if dst.Type().PkgPath() == "github.com/google/uuid" && dst.Type().Name() == "UUID" {
			if err := uuid.Validate(s); err != nil {
				return err
			}

			u, _ := uuid.Parse(s)
			dst.Set(reflect.ValueOf(u))

			return nil
		}

		return fmt.Errorf("unsupported array type %s", dst.Type())

	default:
		under := dst.Type()
		if under.Kind() == reflect.String {
			dst.SetString(s)
			return nil
		}

		return fmt.Errorf("unsupported kind %s", dst.Kind())
	}
}

func indirectValue(value reflect.Value) reflect.Value {
	if value.Kind() == reflect.Pointer {
		return value.Elem()
	}

	return value
}

func toSnakeCase(str string) string {
	var b strings.Builder

	for i, r := range str {
		if unicode.IsUpper(r) {
			if i > 0 {
				prev := rune(str[i-1])
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
