package awsutils

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewS3ClientUploader() (*s3.Client, *manager.Uploader) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		)),
	)
	if err != nil {
		panic("failed to load AWS config, " + err.Error())
	}

	s3Client := s3.NewFromConfig(cfg)
	uploader := manager.NewUploader(s3Client)

	return s3Client, uploader
}

// func NewS3ClientUploader() (*s3.Client, *manager.Uploader) {
// 	cfg, err := config.LoadDefaultConfig(context.TODO())
// 	if err != nil {
// 		panic("failed to load AWS config, " + err.Error())
// 	}

// 	s3Client := s3.NewFromConfig(cfg)
// 	uploader := manager.NewUploader(s3Client)

// 	return s3Client, uploader
// }

// LoadAWSConfig membaca konfigurasi AWS dari file shared credentials dan environment variables.
// func LoadAWSConfig() (*aws.Config, error) {
// 	cfg, err := config.LoadDefaultConfig(context.TODO())
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &cfg, nil
// }

func LoadAWSConfig() (*aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		)),
	)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
