package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ReilEgor/FinScale-backend/CurrencyService/internal/config"
	"github.com/ReilEgor/FinScale-backend/CurrencyService/internal/domain"
)

type CryptoCompare struct {
	CoinGeckoAPIURL config.CompareFinAPIURLType
	CoinGeckoAPIKey config.CompareFinAPIKeyType
}

func NewCryptoCompare(apiURL config.CompareFinAPIURLType, apiKey config.CompareFinAPIKeyType) *CryptoCompare {
	return &CryptoCompare{
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
		return 0, fmt.Errorf("%w: %s", domain.ErrRateNotFound, to)
	}

	return rate, nil
}
