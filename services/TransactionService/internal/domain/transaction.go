package domain

import (
	"context"
	"io"
	"time"
)

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
	RecordTransaction(ctx context.Context, transaction Transaction) error
	RecordReceipt(ctx context.Context, openedFile io.Reader, fileName string) (string, error)
}

type TransactionRepository interface {
	SaveTransaction(ctx context.Context, transaction Transaction) error
}

type FileStorage interface {
	UploadReceipt(ctx context.Context, fileName string, content io.Reader) (string, error)
	GetPresignedURL(ctx context.Context, fileKey string) (string, error)
}
