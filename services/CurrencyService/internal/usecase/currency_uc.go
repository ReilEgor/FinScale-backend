package usecase

import (
	"context"
	"errors"
	"log/slog"

	"github.com/ReilEgor/FinScale-backend/CurrencyService/internal/domain"
)

type CurrencyUseCase struct {
	logger *slog.Logger
	repo   domain.CurrencyRepository
}

func NewCurrencyUseCase() *CurrencyUseCase {
	return &CurrencyUseCase{
		logger: slog.With(slog.String("useCase", "currencyUseCase")),
	}
}

func (u *CurrencyUseCase) ConvertCurrency(ctx context.Context, from string, to string, amount float64) (float64, error) {
	rate, err := u.repo.GetCurrency(ctx, from, to)
	if err != nil {
		if errors.Is(err, domain.ErrCurrencyRateNotFound) {
			u.logger.Warn("currency rate not found, fetching from external API", "from", from, "to", to)
		}
		return 0, err
	}

	convertedAmount := amount * rate
	u.logger.Info("currency converted successfully", "from", from, "to", to, "amount", amount, "convertedAmount", convertedAmount)
	return convertedAmount, nil
}

func (u *CurrencyUseCase) SaveCurrency(ctx context.Context, from string, to string, amount float64) error {
	return u.repo.SaveCurrency(ctx, from, to, amount)
}
