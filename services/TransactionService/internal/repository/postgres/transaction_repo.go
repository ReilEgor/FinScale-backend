package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ReilEgor/FinScale-backend/TransactionService/internal/domain"
	sharedContextUtil "github.com/ReilEgor/FinScale-shared/pkg/contextutil"
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

func (r *TransactionRepository) SaveTransaction(ctx context.Context, transaction domain.Transaction) (uuid.UUID, error) {
	uid, ok := sharedContextUtil.GetUserID(ctx)

	logger := r.logger.With(
		slog.String("operation", "TransactionRepository.SaveTransaction"),
		slog.Any("transaction", transaction),
	)
	if ok {
		logger = logger.With(slog.String("user_id", uid))
	} else {
		logger.WarnContext(ctx, "unauthorized receipt upload attempt")
		return uuid.Nil, domain.ErrUnauthorized
	}

	var id uuid.UUID
	err := r.db.QueryRow(ctx, saveTransactionQuery,
		transaction.UserID,
		transaction.Amount,
		transaction.Currency,
		transaction.Categories,
		transaction.Account,
		transaction.ReceiptURL,
		transaction.TransactionTime,
	).Scan(&id)
	if err != nil {
		r.logger.Error(domain.ErrFailedToSaveTransaction.Error(),
			slog.Any("error", err),
		)
		return uuid.Nil, fmt.Errorf("%w: %w", domain.ErrFailedToSaveTransaction, err)
	}
	return id, nil
}
