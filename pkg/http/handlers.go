package http

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

const (
	defaultReadHeaderTimeout = 10 * time.Second
	defaultWriteTimeout      = 10 * time.Second
	shutdownTimeout          = 5 * time.Second
)

func NewHandler(basePath string, opts ...RouterOption) http.Handler {
	baseRouter := chi.NewRouter()
	baseRouter.Route(basePath, func(r chi.Router) {
		for _, opt := range opts {
			opt(r)
		}
	})
	return baseRouter
}

func NewServer(addr string, logger *log.Logger, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         addr,
		Handler:      handler,
		ErrorLog:     logger,
		ReadTimeout:  defaultReadHeaderTimeout,
		WriteTimeout: defaultWriteTimeout,
	}
}

func RunServer(ctx context.Context, addr string, logger slog.Logger, handler http.Handler) error {
	errLog := slog.NewLogLogger(logger.Handler(), slog.LevelError)
	server := NewServer(addr, errLog, handler)
	errListen := make(chan error, 1)
	go func() {
		logger.Info("Starting the server...")
		errListen <- server.ListenAndServe()
	}()
	select {
	case <-ctx.Done():
		ctxShutdown, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		err := server.Shutdown(ctxShutdown)
		if err != nil {
			return fmt.Errorf("can't shutdown server: %w", err)
		}
		return nil
	case err := <-errListen:
		return fmt.Errorf("can't run server: %w", err)
	}
}
