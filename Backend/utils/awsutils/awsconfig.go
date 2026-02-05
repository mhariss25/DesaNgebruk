package awsutils

import (
	"context"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func GeneratePublicURL(bucket, key string) (string, error) {
	encodedBucket := url.QueryEscape(bucket)
	encodedKey := url.QueryEscape(key)

	return "https://" + encodedBucket + ".s3.amazonaws.com/" + encodedKey, nil
}

func GeneratePresignURL(bucket, key string) (string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", err
	}

	client := s3.NewFromConfig(cfg)
	presigner := s3.NewPresignClient(client)

	input := &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}

	resp, err := presigner.PresignGetObject(context.TODO(), input,
		func(po *s3.PresignOptions) {
			po.Expires = 15 * time.Hour
		},
	)
	if err != nil {
		return "", err
	}

	return resp.URL, nil
}
