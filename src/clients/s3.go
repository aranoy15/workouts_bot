package clients

import (
	"context"
	"errors"
	"fmt"
	"workouts_bot/src/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	client   *s3.Client
	bucket   string
	endpoint string
}

func NewS3Client(cfg *config.Config) (*S3Client, error) {
	if cfg.S3.AccessKeyID == "" || cfg.S3.SecretAccessKey == "" {
		return nil, errors.New("S3 credentials not configured")
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion(cfg.S3.Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.S3.AccessKeyID,
			cfg.S3.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.S3.Endpoint)
		o.UsePathStyle = true
	})

	return &S3Client{
		client:   client,
		bucket:   cfg.S3.BucketName,
		endpoint: cfg.S3.Endpoint,
	}, nil
}
