package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ReilEgor/FinScale-backend/TransactionService/internal/config"
	"github.com/caarlos0/env/v11"
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
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	app, cleanup, err := InitializeApp(ctx, cfg.DSN, cfg.AWS.Region, cfg.AWS.AccessKeyID, cfg.AWS.SecretAccessKey, cfg.AWS.Bucket)
	if err != nil {
		logger.Error("failed to initialize app",
			slog.Any("error", err),
		)
		os.Exit(1)
	}
	defer cleanup()

	go func() {
		port := os.Getenv("HTTP_PORT")
		if port == "" {
			port = "8080"
		}
		if err := app.Server.Run(":" + port); err != nil {
			logger.Error("failed to start server", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down gracefully")

}
