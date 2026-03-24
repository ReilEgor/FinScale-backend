package domain

import "context"

//go:generate mockery --name CurrencyUseCase --output ../mocks/domain --outpkg domain --case=underscore
type CurrencyUseCase interface {
	ConvertCurrency(ctx context.Context, from string, to string, amount float64) (float64, error)
	SaveCurrency(ctx context.Context, from string, to string, rate float64) error
}

//go:generate mockery --name CurrencyRepository --output ../mocks/domain --outpkg domain --case=underscore
type CurrencyRepository interface {
	GetCurrency(ctx context.Context, from string, to string) (float64, error)
	SaveCurrency(ctx context.Context, from string, to string, rate float64) error
}

//go:generate mockery --name CurrencyFetcher --output ../mocks/domain --outpkg domain --case=underscore
type CurrencyFetcher interface {
	GetRateFromCryptoCompare(ctx context.Context, from, to string) (float64, error)
}
