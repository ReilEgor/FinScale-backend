package postgres

import (
	"context"
	"log/slog"

	"github.com/ReilEgor/FinScale-backend/TransactionService/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{
		db:     db,
		logger: slog.With(slog.String("repository", "TransactionRepository")),
	}
}

const saveTransactionQuery = `
    INSERT INTO transactions (user_id, amount, currency, categories, account, receipt_url, transaction_time)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING id;
`

func (r *TransactionRepository) SaveTransaction(ctx context.Context, transaction domain.Transaction) error {
	var lastID uuid.UUID
	err := r.db.QueryRow(ctx, saveTransactionQuery,
		transaction.UserID,          // $1
		transaction.Amount,          // $2
		transaction.Currency,        // $3
		transaction.Categories,      // $4
		transaction.Account,         // $5
		transaction.ReceiptURL,      // $6
		transaction.TransactionTime, // $7
	).Scan(&lastID)
	if err != nil {
		r.logger.Error("failed to save transaction",
			slog.Any("error", err),
			slog.Any("transaction", transaction),
			slog.Any("last_id", lastID),
		)
		return err
	}
	return nil
}
