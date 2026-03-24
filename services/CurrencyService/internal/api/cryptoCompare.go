package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ReilEgor/FinScale-backend/CurrencyService/internal/config"
	"github.com/ReilEgor/FinScale-backend/CurrencyService/internal/domain"
)

type CryptoCompare struct {
	httpClient      *http.Client
	CoinGeckoAPIURL config.CompareFinAPIURLType
	CoinGeckoAPIKey config.CompareFinAPIKeyType
}

func NewCryptoCompare(apiURL config.CompareFinAPIURLType, apiKey config.CompareFinAPIKeyType) *CryptoCompare {
	return &CryptoCompare{
		httpClient:      &http.Client{Timeout: time.Second * 10},
		CoinGeckoAPIURL: apiURL,
		CoinGeckoAPIKey: apiKey,
	}
}

func (c *CryptoCompare) GetRateFromCryptoCompare(ctx context.Context, from, to string) (float64, error) {
	url := fmt.Sprintf("%s/data/price?fsym=%s&tsyms=%s",
		c.CoinGeckoAPIURL,
		strings.ToUpper(from),
		strings.ToUpper(to),
	)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Authorization", "Apikey "+string(c.CoinGeckoAPIKey))

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
}
