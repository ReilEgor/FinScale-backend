package usecase

import (
	"context"
	"io"
	"log/slog"

	"github.com/ReilEgor/FinScale-backend/TransactionService/internal/domain"
)

type TransactionUseCase struct {
	repository domain.TransactionRepository
	logger     *slog.Logger
	storage    domain.FileStorage
}

func NewTransactionUseCase(repository domain.TransactionRepository, storage domain.FileStorage) *TransactionUseCase {
	return &TransactionUseCase{
		repository: repository,
		logger:     slog.With(slog.String("useCase", "NewTransactionUseCase")),
		storage:    storage,
	}
}

func (uc *TransactionUseCase) RecordTransaction(ctx context.Context, transaction domain.Transaction) error {
	return uc.repository.SaveTransaction(ctx, transaction)
}
func (uc *TransactionUseCase) RecordReceipt(ctx context.Context, openedFile io.Reader, fileName string) (string, error) {
	receipt, err := uc.storage.UploadReceipt(ctx, fileName, openedFile)
	if err != nil {
		return "", err
	}
	return receipt, nil
}
