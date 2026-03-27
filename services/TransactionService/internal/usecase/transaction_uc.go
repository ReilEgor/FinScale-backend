package usecase

import (
	"context"
	"github.com/ReilEgor/FinScale-backend/TransactionService/internal/domain"
	sharedContextUtil "github.com/ReilEgor/FinScale-shared/pkg/contextutil"
	"github.com/google/uuid"
	"io"
	"log/slog"
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

func (uc *TransactionUseCase) RecordTransaction(ctx context.Context,
	transaction domain.Transaction,
) (uuid.UUID, error) {
	if err := validateTransaction(transaction); err != nil {
		uc.logger.Warn("transaction validation failed",
			slog.String("user_id", transaction.UserID),
			slog.Any("error", err),
		)
		return uuid.Nil, err
	}

	id, err := uc.repository.SaveTransaction(ctx, transaction)
	if err != nil {
		uc.logger.Error("failed to record transaction",
			slog.String("user_id", transaction.UserID),
			slog.Any("error", err),
		)
		return uuid.Nil, err
	}

	uc.logger.Info("transaction recorded successfully",
		slog.String("user_id", transaction.UserID),
		slog.String("transaction_id", id.String()),
	)
	return id, nil
}

func (uc *TransactionUseCase) RecordReceipt(
	ctx context.Context,
	openedFile io.Reader,
	fileName string,
) (string, error) {

	uid, ok := sharedContextUtil.GetUserID(ctx)

	logger := uc.logger.With(
		slog.String("operation", "TransactionUseCase.RecordReceipt"),
		slog.String("file_name", fileName),
	)

	if ok {
		logger = logger.With(slog.String("user_id", uid))
	} else {
		logger.WarnContext(ctx, "unauthorized receipt upload attempt")
		return "", domain.ErrUnauthorized
	}

	if openedFile == nil {
		logger.ErrorContext(ctx, "nil file provided")
		return "", domain.ErrNilFile
	}

	if fileName == "" {
		logger.ErrorContext(ctx, "empty filename")
		return "", domain.ErrInvalidFileName
	}

	url, err := uc.storage.UploadReceipt(ctx, fileName, openedFile)
	if err != nil {
		logger.ErrorContext(ctx, "failed to upload receipt",
			slog.Any("error", err),
		)
		return "", err
	}

	logger.InfoContext(ctx, "receipt uploaded successfully",
		slog.String("url", url),
	)

	return url, nil
}
func validateTransaction(t domain.Transaction) error {
	if t.UserID == "" {
		return domain.ErrInvalidUserID
	}
	if t.Amount <= 0 {
		return domain.ErrInvalidAmount
	}
	if t.Currency == "" {
		return domain.ErrInvalidCurrency
	}
	return nil
}
