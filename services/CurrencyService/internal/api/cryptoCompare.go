package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ReilEgor/FinScale-backend/CurrencyService/internal/config"
	"github.com/ReilEgor/FinScale-backend/CurrencyService/internal/domain"
	"github.com/sony/gobreaker"
	"net/http"
	"strings"
	"time"
)

type CryptoCompare struct {
	httpClient *http.Client
	apiURL     config.CompareFinAPIURLType
	apiKey     config.CompareFinAPIKeyType
	cb         *gobreaker.CircuitBreaker
}

func NewCryptoCompare(apiURL config.CompareFinAPIURLType, apiKey config.CompareFinAPIKeyType) *CryptoCompare {
	settings := gobreaker.Settings{
		Name:        "CryptoCompareAPI",
		MaxRequests: 3,
		Interval:    5 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 3
		},
	}

	return &CryptoCompare{
		httpClient: &http.Client{Timeout: time.Second * 10},
		apiURL:     apiURL,
		apiKey:     apiKey,
		cb:         gobreaker.NewCircuitBreaker(settings),
	}
}

func (c *CryptoCompare) GetRateFromCryptoCompare(ctx context.Context, from, to string) (float64, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	result, err := c.cb.Execute(func() (interface{}, error) {
		url := fmt.Sprintf("%s/data/price?fsym=%s&tsyms=%s",
			c.apiURL,
			strings.ToUpper(from),
			strings.ToUpper(to),
		)

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return 0, err
		}

		req.Header.Set("Authorization", "Apikey "+string(c.apiKey))

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return 0, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return 0, fmt.Errorf("%w: status %d", domain.ErrFetchFromExternalAPI, resp.StatusCode)
		}

		var result map[string]float64
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return 0, err
		}

		rate, ok := result[strings.ToUpper(to)]
		if !ok {
			return 0, fmt.Errorf("%w: %s", domain.ErrRateNotFound, to)
		}

		return rate, nil
	})
	if err != nil {
		return 0, fmt.Errorf("%w: %w", domain.ErrFetchFromExternalAPI, err)
	}
	rate, ok := result.(float64)
	if !ok {
		return 0, fmt.Errorf("unexpected result type from circuit breaker")
	}
	return rate, nil
}
