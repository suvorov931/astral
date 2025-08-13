// @title Authentication-Service API
// @version 1.0
// @description Документация API Authentication-Service
// @host localhost:8082
// @BasePath /
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @scheme bearer
// @bearerFormat JWT
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	//_ "authentication-service/docs"

	"astral/internal/api/handler"
	mmiddleware "astral/internal/api/middleware"
	"astral/internal/auth"
	rredisClient "astral/internal/cache/redisCLient"
	cconfig "astral/internal/config"
	llogger "astral/internal/logger"
	ppostgresClient "astral/internal/storage/postgresClient"
)

const (
	pathToConfig     = "./config/config.env"
	pathToMigrations = "file://./database/migrations"
	shoutdownTime    = 15 * time.Second
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancel()

	config, err := cconfig.New(pathToConfig)
	if err != nil {
		log.Fatalf("failed to initialize config: %v", err)
	}

	logger, err := llogger.New(&config.Logger)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	postgresClient, err := ppostgresClient.New(ctx, &config.Postgres, logger, pathToMigrations)
	if err != nil {
		logger.Fatal("failed to initialize postgres client", zap.Error(err))
	}

	redisClient, err := rredisClient.New(ctx, &config.Redis, logger)
	if err != nil {
		logger.Fatal("failed to initialize redis client", zap.Error(err))
	}

	router := chi.NewRouter()

	router.Use(middleware.RealIP)
	router.Use(middleware.RequestID)
	router.Use(llogger.MiddlewareLogger(logger, &config.Logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	authService := auth.New(&config.Auth, logger)

	router.With(mmiddleware.RequireAdminToken(authService, logger)).
		Post("/api/register", handler.Register(postgresClient, authService, logger))

	router.Post("/api/auth", handler.Auth(postgresClient, redisClient, authService, logger))

	server := http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.HttpServer.Host, config.HttpServer.Port),
		Handler: router,
	}

	go func() {
		logger.Info("starting http server", zap.String("addr", server.Addr))
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("failed to start server", zap.Error(err))
		}
	}()

	<-ctx.Done()

	logger.Info("received shutdown signal")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shoutdownTime)
	defer shutdownCancel()

	if err = server.Shutdown(shutdownCtx); err != nil {
		logger.Error("cannot shutdown http server", zap.Error(err))
		return
	}

	//postgresClient.Close()

	logger.Info("stopping http server", zap.String("addr", server.Addr))

	logger.Info("application shutdown completed successfully")
}

// TODO: documentation
// TODO: tests
