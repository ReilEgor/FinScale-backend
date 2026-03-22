package aws

import (
	"context"
	domainConfig "github.com/ReilEgor/FinScale-backend/TransactionService/internal/config"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
)

func NewS3Client(
	region domainConfig.AWSRegionType,
	accessKeyID domainConfig.AWSAccessKeyIDType,
	secretAccessKey domainConfig.AWSSecretAccessKeyType,
) *s3.Client {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(string(region)),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				string(accessKeyID),
				string(secretAccessKey),
				"",
			),
		),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
	}

	return s3.NewFromConfig(cfg)
}
