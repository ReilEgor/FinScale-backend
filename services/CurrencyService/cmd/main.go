package main

import (
	"log/slog"
	"os"

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
}
