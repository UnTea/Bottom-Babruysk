package router

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

func HandleStaticReflect(r chi.Router, basePath string, h http.Handler) {
	p := strings.TrimRight(basePath, "/")

	r.Handle(p+"*", h)
}
