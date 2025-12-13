package middleware

import (
	"net/http"
	"time"

	"github.com/nomdb/backend/internal/logger"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	bytes      int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.bytes += n
	return n, err
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the response writer to capture status code
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Log incoming request
		logger.Debug("→ %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		// Call the next handler
		next.ServeHTTP(rw, r)

		// Calculate duration
		duration := time.Since(start)

		// Log response with appropriate level based on status code
		statusEmoji := getStatusEmoji(rw.statusCode)
		if rw.statusCode >= 500 {
			logger.Error("← %s %s %s %d (%s) %d bytes",
				statusEmoji, r.Method, r.URL.Path, rw.statusCode, duration, rw.bytes)
		} else if rw.statusCode >= 400 {
			logger.Warn("← %s %s %s %d (%s) %d bytes",
				statusEmoji, r.Method, r.URL.Path, rw.statusCode, duration, rw.bytes)
		} else {
			logger.Info("← %s %s %s %d (%s) %d bytes",
				statusEmoji, r.Method, r.URL.Path, rw.statusCode, duration, rw.bytes)
		}
	})
}

func getStatusEmoji(statusCode int) string {
	switch {
	case statusCode >= 500:
		return "💥"
	case statusCode >= 400:
		return "⚠️ "
	case statusCode >= 300:
		return "↪️ "
	case statusCode >= 200:
		return "✅"
	default:
		return "ℹ️ "
	}
}
