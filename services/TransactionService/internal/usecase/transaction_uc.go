package usecase

import (
	"log/slog"

	"github.com/ReilEgor/FinScale-backend/TransactionService/internal/domain"
)

type TransactionUseCase struct {
	repository domain.TransactionRepository
	logger     *slog.Logger
}

func NewTransactionUseCase(repository domain.TransactionRepository) *TransactionUseCase {
	return &TransactionUseCase{
		repository: repository,
		logger:     slog.With(slog.String("useCase", "NewTransactionUseCase")),
	}
}

func (uc *TransactionUseCase) RecordTransaction(transaction domain.Transaction) error {
	//TODO implement me
	panic("implement me")
}
