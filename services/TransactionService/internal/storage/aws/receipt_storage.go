package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/ReilEgor/FinScale-backend/TransactionService/internal/config"
	"github.com/ReilEgor/FinScale-backend/TransactionService/internal/domain"
	sharedContextUtil "github.com/ReilEgor/FinScale-shared/pkg/contextutil"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"io"
	"log/slog"
)

const (
	receiptsPrefix     = "receipts/"
	presignedURLExpiry = 15 * time.Minute
	defaultContentType = "image/jpeg"
)

type ReceiptStorage struct {
	client    *s3.Client
	presigner *s3.PresignClient
	bucket    string
	logger    *slog.Logger
}

func NewReceiptStorage(client *s3.Client, bucket config.AWSBucketType) domain.FileStorage {
	return &ReceiptStorage{client: client, presigner: s3.NewPresignClient(client), bucket: string(bucket), logger: slog.With(slog.String("component", "ReceiptStorage"))}
}

func (s *ReceiptStorage) UploadReceipt(ctx context.Context, fileName string, content io.Reader) (string, error) {
	uid, ok := sharedContextUtil.GetUserID(ctx)

	logger := s.logger.With(
		slog.String("operation", "ReceiptStorage.UploadReceipt"),
	)

	if ok {
		logger = logger.With(slog.String("user_id", uid))
	} else {
		logger.WarnContext(ctx, "unauthorized receipt upload attempt")
		return "", domain.ErrUnauthorized
	}

	fileKey := receiptsPrefix + fileName
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(fileKey),
		Body:        content,
		ACL:         types.ObjectCannedACLPrivate,
		ContentType: aws.String(defaultContentType),
	})
	if err != nil {
		s.logger.Error(domain.ErrFailedToUploadObject.Error(),
			slog.String("file_key", fileKey),
			slog.Any("error", err),
		)
		return "", fmt.Errorf("%w: %w", domain.ErrFailedToUploadObject, err)
	}
	return fileKey, nil
}

func (s *ReceiptStorage) GetPresignedURL(ctx context.Context, fileKey string) (string, error) {
	uid, ok := sharedContextUtil.GetUserID(ctx)

	logger := s.logger.With(
		slog.String("operation", "ReceiptStorage.GetPresignedURL"),
		slog.String("bucket", s.bucket),
		slog.String("file_key", fileKey),
	)

	if ok {
		logger = logger.With(slog.String("user_id", uid))
	} else {
		logger.WarnContext(ctx, "unauthorized receipt upload attempt")
		return "", domain.ErrUnauthorized
	}

	if fileKey == "" {
		logger.DebugContext(ctx, "empty file key provided, skipping presigning")
		return "", nil
	}
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fileKey),
	}
	psClient := s3.NewPresignClient(s.client)

	start := time.Now()
	resp, err := psClient.PresignGetObject(ctx, input, func(o *s3.PresignOptions) {
		o.Expires = time.Duration(15 * time.Minute)
	})

	if err != nil {
		logger.ErrorContext(ctx, "failed to generate presigned S3 URL",
			slog.Any("error", err),
			slog.Duration("duration", time.Since(start)),
		)
		return "", fmt.Errorf("failed to sign s3 url: %w", err)
	}

	logger.InfoContext(ctx, "presigned URL generated successfully",
		slog.Duration("duration", time.Since(start)),
		slog.String("expires_in", "15m"),
	)

	return resp.URL, nil
}
