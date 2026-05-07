package delivery

import (
	"net/http"
	"time"
	"github.com/yousefggg/common-lib/pkg/logger"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		logger.Info("HTTP Request processed",
			"method",      r.Method,
			"path",        r.URL.Path,
			"status",      wrapped.status,
			"duration",    time.Since(start).String(),
			"remote_addr", r.RemoteAddr,
			"user_agent",  r.UserAgent(),
		)
	})
}