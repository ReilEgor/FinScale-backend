package domain

import "time"

type Transaction struct {
	UserID          string    `json:"user_id" db:"user_id"`
	Amount          float64   `json:"amount" db:"amount"`
	Currency        string    `json:"currency" db:"currency"`
	TransactionTime time.Time `json:"transaction_time" db:"transaction_time"`
	Categories      []string  `json:"categories" db:"categories"`
	Account         string    `json:"account" db:"account"`
	ReceiptURL      string    `json:"receipt_url" db:"receipt_url"`
}

type TransactionUseCase interface {
	RecordTransaction(transaction Transaction) error
}

type TransactionRepository interface {
	SaveTransaction(transaction Transaction) error
}
