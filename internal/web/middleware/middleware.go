package web

import (
	"net/http"
	"time"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
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
func RequestLogger(baseLogger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now().UTC()
			requestID := chiMiddleware.GetReqID(r.Context())

			if requestID != "" {
				w.Header().Set("X-Request-ID", requestID)
			}

			l := baseLogger.With(
				zap.Any("request id", requestID),
				zap.Any("method", r.Method),
				zap.Any("path", r.URL.Path),
				zap.Any("query", r.URL.RawQuery),
				zap.Any("remote", r.RemoteAddr),
				zap.Any("user agent", r.UserAgent()),
				zap.Any("referer", r.Referer()),
			)

			ctx := ctxzap.ToContext(r.Context(), l)

			srw := &statusRW{ResponseWriter: w}
			next.ServeHTTP(srw, r.WithContext(ctx))

			l.Info("http_request",
				zap.Any("status", srw.status),
				zap.Any("bytes", srw.bytes),
				zap.Duration("duration", time.Since(start)),
			)
		})
	}
}
