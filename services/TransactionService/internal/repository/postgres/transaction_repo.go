package postgres

import (
	"log/slog"

	"github.com/ReilEgor/FinScale-backend/TransactionService/internal/domain"
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

func (TransactionRepository) SaveTransaction(transaction domain.Transaction) error {
	//TODO implement me
	return nil
}
