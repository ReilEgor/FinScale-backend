package domain

import "context"

type ConvertCurrencyRequest struct {
	From   string  `json:"from" binding:"required"`
	To     string  `json:"to" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
}

type CurrencyUseCase interface {
	ConvertCurrency(ctx context.Context, from string, to string, amount float64) (float64, error)
	SaveCurrency(ctx context.Context, from string, to string, rate float64) error
}

type CurrencyRepository interface {
	GetCurrency(ctx context.Context, from string, to string) (float64, error)
	SaveCurrency(ctx context.Context, from string, to string, rate float64) error
}
