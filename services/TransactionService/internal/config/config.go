package config

type (
	DSN string

	//	AWSconfig
	AWSRegionType          string
	AWSAccessKeyIDType     string
	AWSSecretAccessKeyType string
	AWSBucketType          string
)
type AWSConfig struct {
	Region          AWSRegionType          `env:"AWS_REGION" envDefault:"us-east-1"`
	AccessKeyID     AWSAccessKeyIDType     `env:"AWS_ACCESS_KEY_ID" envDefault:"your_access_key_id"`
	SecretAccessKey AWSSecretAccessKeyType `env:"AWS_SECRET_ACCESS_KEY" envDefault:"your_secret_access_key"`
	Bucket          AWSBucketType          `env:"AWS_BUCKET" envDefault:"your_bucket_name"`
}

type Config struct {
	DSN DSN `env:"DB_SOURCE" envDefault:"postgres://user:password@localhost:5432/transactionservice?sslmode=disable"`
	AWS AWSConfig
}
