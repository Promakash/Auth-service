package http

import (
	"bytes"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
)

const docsURI = "docs"
const healthURI = "health"

func LoggingMiddleware(logger slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			if strings.Contains(path, docsURI) || strings.Contains(path, healthURI) {
				next.ServeHTTP(w, r)
				return
			}

			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK, body: &bytes.Buffer{}}

			defer func() {
				if rec := recover(); rec != nil {
					logger.Error("Request panic",
						"endpoint", path,
						"method", r.Method,
						"remote_addr", r.RemoteAddr,
						"panic", rec,
						"stacktrace", string(debug.Stack()),
					)
					http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(rw, r)

			var responseBody string
			if rw.body.Len() > 0 {
				responseBody = rw.body.String()
			}
			logger.Info("Request completed",
				"endpoint", path,
				"method", r.Method,
				"status", rw.statusCode,
				"response_body", responseBody,
			)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
	body         *bytes.Buffer
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	rw.body.Write(data)
	size, err := rw.ResponseWriter.Write(data)
	rw.bytesWritten += size
	return size, err
}
