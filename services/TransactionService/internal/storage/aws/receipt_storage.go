package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/ReilEgor/FinScale-backend/TransactionService/internal/config"
	"github.com/ReilEgor/FinScale-backend/TransactionService/internal/domain"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"io"
	"log/slog"
)

type ReceiptStorage struct {
	client *s3.Client
	bucket string
	logger *slog.Logger
}

func NewReceiptStorage(client *s3.Client, bucket config.AWSBucketType) domain.FileStorage {
	return &ReceiptStorage{client: client, bucket: string(bucket), logger: slog.With(slog.String("component", "ReceiptStorage"))}
}

func (s *ReceiptStorage) UploadReceipt(ctx context.Context, fileName string, content io.Reader) (string, error) {
	fileKey := fmt.Sprintf("receipts/%s", fileName)
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(fileKey),
		Body:        content,
		ACL:         types.ObjectCannedACLPrivate,
		ContentType: aws.String("image/jpeg"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload object to s3: %w", err)
	}
	return fileKey, nil
}

func (s *ReceiptStorage) GetPresignedURL(ctx context.Context, fileKey string) (string, error) {
	if fileKey == "" {
		return "", nil
	}
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fileKey),
	}
	psClient := s3.NewPresignClient(s.client)
	resp, err := psClient.PresignGetObject(ctx, input, func(o *s3.PresignOptions) {
		o.Expires = time.Duration(15 * time.Minute)
	})

	if err != nil {
		return "", fmt.Errorf("failed to sign s3 url: %w", err)
	}

	return resp.URL, nil
}
