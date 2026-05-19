package clients

import (
	"workouts_bot/src/config"

	"github.com/aranoy15/go-s3"
)

func NewS3Client(cfg *config.S3Config) (*s3.Client, error) {
	sc := &s3.Config{
		Endpoint:        cfg.Endpoint,
		AccessKeyID:     cfg.AccessKeyID,
		SecretAccessKey: cfg.SecretAccessKey,
		BucketName:      cfg.BucketName,
		Region:          cfg.Region,
	}

	return s3.New(sc)
}
