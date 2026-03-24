package domain

import "errors"

var (
	ErrCurrencyRateNotFound = errors.New("currency rate not found in cache")
	ErrCurrencyRateRetrieve = errors.New("error retrieving currency rate from cache")
	ErrCurrencyRateParse    = errors.New("failed to parse rate from cache")

	ErrFetchFromExternalAPI = errors.New("failed to fetch from external API")
	ErrRateNotFound         = errors.New("rate not found for currency")

	ErrSaveCurrencyRate = errors.New("save currency rate")
	ErrInvalidResponse  = errors.New("invalid response")
)
