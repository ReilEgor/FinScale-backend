package domain

import "errors"

var (
	ErrFailedToSaveTransaction = errors.New("failed to save transaction")
	ErrFailedToUploadObject    = errors.New("failed to upload object to s3")
	ErrFailedToPresignURL      = errors.New("failed to presign s3 url")

	ErrInvalidUserID   = errors.New("user id is empty")
	ErrInvalidAmount   = errors.New("amount must be greater than zero")
	ErrInvalidCurrency = errors.New("currency is empty")
	ErrInvalidFileName = errors.New("file name is empty")
	ErrNilFile         = errors.New("file is nil")

	ErrUnauthorized = errors.New("unauthorized")
)
