package http

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

func LoggingMiddleware(logger slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			scheme := "http"
			if r.TLS != nil {
				scheme = "https"
			}

			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			defer func() {
				if rec := recover(); rec != nil {
					logger.Error("Request panic",
						"method", r.Method,
						"url", r.URL.String(),
						"remote_addr", r.RemoteAddr,
						"user_agent", r.UserAgent(),
						"scheme", scheme,
						"panic", rec,
						"stacktrace", string(debug.Stack()),
					)
					http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(rw, r)

			duration := time.Since(start)
			logger.Info("Request completed",
				"method", r.Method,
				"url", r.URL.String(),
				"status", rw.statusCode,
				"duration_ms", duration.Milliseconds(),
				"bytes_written", rw.bytesWritten,
				"user_agent", r.UserAgent(),
				"scheme", scheme,
			)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(data)
	rw.bytesWritten += size
	return size, err
}
