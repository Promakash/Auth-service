package shutdown

import (
	"log/slog"

	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func LogShutdownDuration(ctx context.Context, logger slog.Logger) func() {
	var shutdownTime time.Time
	go func() {
		<-ctx.Done()
		shutdownTime = time.Now()
	}()
	return func() {
		logger.Info("Shutdown duration: ", elapsedMs(shutdownTime))
	}
}

func ListenSignal(ctx context.Context, logger slog.Logger) error {
	sigquit := make(chan os.Signal, 1)
	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(sigquit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		return nil
	case sig := <-sigquit:
		logger.Info("Captured signal: ", sig)
		logger.Info("Gracefully shutting down server...")
		return errors.New("operating system signal")
	}
}

func DurationToMs(d time.Duration) float64 {
	return float64(d.Nanoseconds()) / float64(time.Millisecond)
}

func elapsedMs(since time.Time) float64 {
	return DurationToMs(time.Since(since))
}
