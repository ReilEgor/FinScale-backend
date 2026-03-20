package redis

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/ReilEgor/FinScale-backend/CurrencyService/internal/domain"
	"github.com/redis/go-redis/v9"
)

type CurrencyRepository struct {
	client *redis.Client
	logger *slog.Logger
}

func NewCurrencyRepository(client *redis.Client) *CurrencyRepository {
	return &CurrencyRepository{
		client: client,
		logger: slog.With(slog.String("component", "currency_repository")),
	}
}

func (r *CurrencyRepository) SaveCurrency(ctx context.Context, from string, to string, rate float64) error {
	cacheKey := fmt.Sprintf("rate:%s:%s", strings.ToUpper(from), strings.ToUpper(to))
	return r.client.Set(ctx, cacheKey, rate, 5*time.Minute).Err()
}
func (r *CurrencyRepository) GetCurrency(ctx context.Context, from string, to string) (float64, error) {
	cacheKey := fmt.Sprintf("rate:%s:%s", strings.ToUpper(from), strings.ToUpper(to))
	val, err := r.client.Get(ctx, cacheKey).Result()
	if errors.Is(err, redis.Nil) {
		return 0, domain.ErrCurrencyRateNotFound
	} else if err != nil {
		return 0, fmt.Errorf("%w: %w", domain.ErrCurrencyRateRetrieve, err)
	}

	rate, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: %w", domain.ErrCurrencyRateParse, err)
	}

	return rate, nil
}
