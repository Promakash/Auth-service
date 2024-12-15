package main

import (
	pkgConfig "auth_service/pkg/config"
	"auth_service/service"
	"context"
	"errors"
	"golang.org/x/sync/errgroup"
	"log"
	"log/slog"
	"time"

	"auth_service/api/http"
	_ "auth_service/docs"
	pkghttp "auth_service/pkg/http"
	"auth_service/pkg/infra"
	pkglog "auth_service/pkg/log"
	"auth_service/pkg/shutdown"
	"auth_service/repository/postgres"
)

// @title           Auth service API
// @version         1.0
// @description     REST сервис аутентификации
// @termsOfService  http://swagger.io/terms/

// @host      localhost:8080
// @BasePath  /
func main() {
	cfg := pkgConfig.MustLoad()

	logger := pkglog.NewLogger("debug", "json")
	slog.SetDefault(&logger)
	logger.Info("Service started", "config", cfg)

	g, ctx := errgroup.WithContext(context.Background())

	defer shutdown.LogShutdownDuration(ctx, logger)()
	g.Go(func() error { return shutdown.ListenSignal(ctx, logger) })

	pg, err := infra.NewPostgres(cfg.PG)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func() { _ = pg.Close() }()

	authRepo := postgres.NewAuthRepository(pg)
	authService := service.NewAuth(authRepo, cfg.Secret, time.Duration(cfg.RefreshExp), time.Duration(cfg.AccessExp))

	userRepo := postgres.NewUserRepository(pg)
	userService := service.NewUser(authService, userRepo)
	userHandler := http.NewAuthHandler(logger, userService)

	publicHandler := pkghttp.NewHandler("/",
		pkghttp.WithLoggingMiddleware(logger),
		pkghttp.WithSwagger(),
		userHandler.WithAuthHandlers(),
		pkghttp.WithHealthHandler(),
	)

	g.Go(func() error {
		err = pkghttp.RunServer(ctx, cfg.Address, logger, publicHandler)
		if err != nil {
			logger.Error("Public server error", "error", err)
		}
		return err
	})

	err = g.Wait()
	if err != nil && !errors.Is(err, errors.New("operating system signal")) {
		logger.Error("Exit reason", "error", err)
	}
}
