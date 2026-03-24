package redis

import (
	"context"
	"testing"

	"github.com/ReilEgor/FinScale-backend/CurrencyService/internal/domain"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestRepository(t *testing.T) (*CurrencyRepository, *miniredis.Miniredis) {
	t.Helper()
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	return NewCurrencyRepository(client), mr
}

func Test_CurrencyRepository_SaveCurrency(t *testing.T) {
	tests := []struct {
		name    string
		from    string
		to      string
		rate    float64
		setup   func(mr *miniredis.Miniredis)
		wantErr error
	}{
		{
			name:    "success: rate saved",
			from:    "USD",
			to:      "RUB",
			rate:    85.0,
			wantErr: nil,
		},
		{
			name: "error: redis unavailable",
			from: "USD",
			to:   "RUB",
			rate: 85.0,
			setup: func(mr *miniredis.Miniredis) {
				mr.Close()
			},
			wantErr: domain.ErrSaveCurrencyRate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, mr := newTestRepository(t)
			if tt.setup != nil {
				tt.setup(mr)
			}

			err := repo.SaveCurrency(context.Background(), tt.from, tt.to, tt.rate)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)

				val, redisErr := mr.Get(CacheKey(tt.from, tt.to))
				require.NoError(t, redisErr)
				assert.Equal(t, "85", val)
			}
		})
	}
}
func Test_CurrencyRepository_GetCurrency(t *testing.T) {
	tests := []struct {
		name     string
		from     string
		to       string
		setup    func(mr *miniredis.Miniredis)
		wantRate float64
		wantErr  error
	}{
		{
			name: "success: rate found in cache",
			from: "USD",
			to:   "RUB",
			setup: func(mr *miniredis.Miniredis) {
				mr.Set(CacheKey("USD", "RUB"), "85")
			},
			wantRate: 85.0,
		},
		{
			name:    "error: cache miss",
			from:    "USD",
			to:      "RUB",
			wantErr: domain.ErrCurrencyRateNotFound,
		},
		{
			name: "error: invalid value in cache",
			from: "USD",
			to:   "RUB",
			setup: func(mr *miniredis.Miniredis) {
				mr.Set(CacheKey("USD", "RUB"), "not-a-float")
			},
			wantErr: domain.ErrCurrencyRateParse,
		},
		{
			name: "error: redis unavailable",
			from: "USD",
			to:   "RUB",
			setup: func(mr *miniredis.Miniredis) {
				mr.Close()
			},
			wantErr: domain.ErrCurrencyRateRetrieve,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, mr := newTestRepository(t)
			if tt.setup != nil {
				tt.setup(mr)
			}

			rate, err := repo.GetCurrency(context.Background(), tt.from, tt.to)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, 0.0, rate)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantRate, rate)
			}
		})
	}
}
