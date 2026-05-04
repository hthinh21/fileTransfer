package storage

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var R2 *s3.Client

func InitR2() error {
	endpoint := os.Getenv("R2_ENDPOINT")
	accessKey := os.Getenv("R2_ACCESS_KEY")
	secretKey := os.Getenv("R2_SECRET_KEY")
	if endpoint == "" || accessKey == "" || secretKey == "" {
		return fmt.Errorf("missing R2 configuration")
	}

	R2 = s3.New(s3.Options{
		Region:       "auto",
		BaseEndpoint: aws.String(endpoint),
		Credentials: aws.NewCredentialsCache(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		),
	})

	return nil
}

func bucketName() (string, error) {
	bucket := os.Getenv("R2_BUCKET")
	if bucket == "" {
		return "", fmt.Errorf("missing R2_BUCKET")
	}

	return bucket, nil
}

func UploadFile(ctx context.Context, objectKey string, body io.Reader, contentType string) error {
	if R2 == nil {
		if err := InitR2(); err != nil {
			return err
		}
	}

	bucket, err := bucketName()
	if err != nil {
		return err
	}

	_, err = R2.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(objectKey),
		Body:        body,
		ContentType: aws.String(contentType),
	})
	return err
}

func DownloadFile(ctx context.Context, objectKey string) (io.ReadCloser, string, int64, error) {
	if R2 == nil {
		if err := InitR2(); err != nil {
			return nil, "", 0, err
		}
	}

	bucket, err := bucketName()
	if err != nil {
		return nil, "", 0, err
	}

	out, err := R2.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return nil, "", 0, err
	}

	contentType := "application/octet-stream"
	if out.ContentType != nil && *out.ContentType != "" {
		contentType = *out.ContentType
	}

	var contentLength int64
	if out.ContentLength != nil {
		contentLength = *out.ContentLength
	}

	return out.Body, contentType, contentLength, nil
}
