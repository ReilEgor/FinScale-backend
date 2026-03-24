package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/ReilEgor/FinScale-backend/CurrencyService/internal/config"
	"github.com/ReilEgor/FinScale-backend/CurrencyService/internal/domain"
)

type CurrencyUseCase struct {
	logger  *slog.Logger
	repo    domain.CurrencyRepository
	fetcher domain.CurrencyFetcher
}

func NewCurrencyUseCase(repo domain.CurrencyRepository, fetcher domain.CurrencyFetcher, compareFinAPIURL config.CompareFinAPIURLType, compareFinAPIKey config.CompareFinAPIKeyType) *CurrencyUseCase {
	return &CurrencyUseCase{
		logger:  slog.With(slog.String("useCase", "currencyUseCase")),
		repo:    repo,
		fetcher: fetcher,
	}
}

func (u *CurrencyUseCase) ConvertCurrency(ctx context.Context, from string, to string, amount float64) (float64, error) {
	rate, err := u.repo.GetCurrency(ctx, from, to)
	if err != nil {
		if errors.Is(err, domain.ErrCurrencyRateNotFound) {
			u.logger.Warn("Cache miss, fetching from external API", "from", from, "to", to)

			rate, err = u.fetcher.GetRateFromCryptoCompare(ctx, from, to)
			if err != nil {
				return 0, fmt.Errorf("%w: %w", domain.ErrFetchFromExternalAPI, err)
			}
			err = u.repo.SaveCurrency(ctx, from, to, rate)
			if err != nil {
				u.logger.Error("Failed to save currency rate to repository", "error", err)
			}
		} else {
			return 0, err
		}
	}

	convertedAmount := amount * rate

	u.logger.Info("currency converted successfully",
		"from", from,
		"to", to,
		"amount", amount,
		"rate", rate,
		"convertedAmount", convertedAmount)

	return convertedAmount, nil
}

func (u *CurrencyUseCase) SaveCurrency(ctx context.Context, from string, to string, rate float64) error {
	if err := u.repo.SaveCurrency(ctx, from, to, rate); err != nil {
		u.logger.Error("failed to save currency rate",
			"from", from,
			"to", to,
			"rate", rate,
			"error", err,
		)
		return fmt.Errorf("%w: %w", domain.ErrSaveCurrencyRate, err)
	}

	u.logger.Debug("currency rate saved successfully",
		"from", from,
		"to", to,
		"rate", rate,
	)

	return nil
}
