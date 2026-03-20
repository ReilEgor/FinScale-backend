package domain

import "errors"

var (
	ErrCurrencyRateNotFound = errors.New("currency rate not found in cache")
	ErrCurrencyRateRetrieve = errors.New("error retrieving currency rate from cache")
	ErrCurrencyRateParse    = errors.New("failed to parse rate from cache")
)
