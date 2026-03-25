package main

import (
	"context"
	_ "github.com/ReilEgor/FinScale-backend/CurrencyService/api/docs"
	config "github.com/ReilEgor/FinScale-backend/CurrencyService/internal/config"
	env "github.com/caarlos0/env/v11"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

// Swagger Metadata for API Documentation
// @title           FinScale Currency Service API
// @version         1.0
// @description     Microservice for real-time currency conversion and exchange rate tracking.
// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	var cfg config.Config
	err := env.Parse(&cfg)
	if err != nil {
		logger.Error("failed to load config",
			slog.Any("error", err),
		)
		os.Exit(1)
	}
	app, cleanup, err := InitializeApp(
		cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword, 0, cfg.CompareFinAPIURL, cfg.CompareFinAPIKey,
	)
	if err != nil {
		logger.Error("failed to initialize app", slog.Any("error", err))
		os.Exit(1)
	}
	defer cleanup()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Server.Run(ctx, ":"+cfg.HTTPPort)
	}()

	select {
	case <-ctx.Done():
		logger.Info("shutting down gracefully")
		if err := <-errCh; err != nil {
			logger.Error("server shutdown error", slog.Any("error", err))
		}
		logger.Info("server stopped")
	case err := <-errCh:
		logger.Error("server stopped unexpectedly", slog.Any("error", err))
	}
}
