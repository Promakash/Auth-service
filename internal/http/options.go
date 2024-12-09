package http

import (
	"net/http"

	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
)

func RouterOptions(options ...RouterOption) func(chi.Router) {
	return func(r chi.Router) {
		for _, option := range options {
			option(r)
		}
	}
}

type RouterOption func(chi.Router)

func WithHealthHandler() RouterOption {
	return func(r chi.Router) {
		r.Mount("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))
		}))
	}
}

func WithSwagger() RouterOption {
	return func(r chi.Router) {
		r.Get("/docs/*", httpSwagger.Handler(
			httpSwagger.URL("docs/doc.json"),
		))
	}
}
