package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// CompressionMiddleware adds gzip compression to responses when client supports it
func CompressionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if client supports gzip
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// Don't compress if content is already compressed or is an image
		contentType := w.Header().Get("Content-Type")
		if strings.Contains(contentType, "image/") ||
		   strings.Contains(contentType, "video/") ||
		   strings.Contains(contentType, "application/zip") ||
		   strings.Contains(contentType, "application/gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// Set gzip headers
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Del("Content-Length") // Content-Length will be different after compression

		// Create gzip writer
		gz := gzip.NewWriter(w)
		defer gz.Close()

		// Wrap response writer
		gzw := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		next.ServeHTTP(gzw, r)
	})
}
