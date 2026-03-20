package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/ReilEgor/FinScale-backend/CurrencyService/internal/config"
	"github.com/ReilEgor/FinScale-backend/CurrencyService/internal/domain"
)

type CurrencyUseCase struct {
	logger          *slog.Logger
	repo            domain.CurrencyRepository
	CoinGeckoAPIURL config.CompareFinAPIURLType
	CoinGeckoAPIKey config.CompareFinAPIKeyType
}

func NewCurrencyUseCase(repo domain.CurrencyRepository, compareFinAPIURL config.CompareFinAPIURLType, compareFinAPIKey config.CompareFinAPIKeyType) *CurrencyUseCase {
	return &CurrencyUseCase{
		logger:          slog.With(slog.String("useCase", "currencyUseCase")),
		repo:            repo,
		CoinGeckoAPIURL: compareFinAPIURL,
		CoinGeckoAPIKey: compareFinAPIKey,
	}
}

func (u *CurrencyUseCase) ConvertCurrency(ctx context.Context, from string, to string, amount float64) (float64, error) {
	rate, err := u.repo.GetCurrency(ctx, from, to)
	if err != nil {
		if errors.Is(err, domain.ErrCurrencyRateNotFound) {
			u.logger.Warn("Cache miss, fetching from external API", "from", from, "to", to)

			rate, err = u.getRateFromCryptoCompare(ctx, from, to)
			if err != nil {
				return 0, fmt.Errorf("failed to fetch from external API: %w", err)
			}
			_ = u.repo.SaveCurrency(ctx, from, to, rate)
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
	return u.repo.SaveCurrency(ctx, from, to, rate)
}

func (u *CurrencyUseCase) getRateFromCryptoCompare(ctx context.Context, from, to string) (float64, error) {
	url := fmt.Sprintf("https://min-api.cryptocompare.com/data/price?fsym=%s&tsyms=%s",
		strings.ToUpper(from), strings.ToUpper(to))

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Authorization", "Apikey "+string(u.CoinGeckoAPIKey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	rate, ok := result[strings.ToUpper(to)]
	if !ok {
		return 0, fmt.Errorf("rate not found for %s", to)
	}

	return rate, nil
}
