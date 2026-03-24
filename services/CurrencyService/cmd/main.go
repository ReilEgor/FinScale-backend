package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	config "github.com/ReilEgor/FinScale-backend/CurrencyService/internal/config"
	env "github.com/caarlos0/env/v11"
)

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
