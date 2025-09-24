package web

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

type statusRW struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (w *statusRW) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusRW) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}

	n, err := w.ResponseWriter.Write(b)
	w.bytes += n

	return n, err
}

// RequestLogger кладёт request-scoped zap.Logger в контекст + логирует запрос/ответ.
func RequestLogger(base *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			requestID := middleware.GetReqID(r.Context())

			l := base.With(
				zap.String("req_id", requestID),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote", r.RemoteAddr),
				zap.String("ua", r.UserAgent()),
			)

			ctx := ctxzap.ToContext(r.Context(), l)

			srw := &statusRW{ResponseWriter: w}
			next.ServeHTTP(srw, r.WithContext(ctx))

			l.Info("http_request",
				zap.Int("status", srw.status),
				zap.Int("bytes", srw.bytes),
				zap.Duration("duration", time.Since(start)),
			)
		})
	}
}
