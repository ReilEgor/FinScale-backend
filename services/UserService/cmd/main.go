package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ReilEgor/FinScale-backend/UserService/internal/config"
	"github.com/caarlos0/env/v11"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)
	logger = slog.With(slog.String("service", "main"))

	var cfg config.Config
	if err := env.Parse(&cfg); err != nil {
		logger.Error("failed to load config",
			slog.Any("error", err),
		)
		os.Exit(1)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	app, cleanup, err := InitializeApp(ctx,
		string(cfg.DSN),
	)
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
