// https://github.com/awsdocs/aws-doc-sdk-examples/blob/main/gov2/s3/actions/bucket_basics.go#L200
// https://github.com/aws/aws-sdk-go/issues/1837

package s3_manager

import (
	"context"
	"errors"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/joerdav/zapray"
	"go.uber.org/zap"
)

type S3Manager struct {
	logger     *zapray.Logger
	s3Client   *s3.Client
	bucketName string
}

func NewS3Manager(logger *zapray.Logger, cfg aws.Config, bucketName string) S3Manager {
	cfg.ResponseChecksumValidation = aws.ResponseChecksumValidationUnset

	s3Client := s3.NewFromConfig(cfg)
	return S3Manager{logger: logger, s3Client: s3Client, bucketName: bucketName}
}

func (m S3Manager) BucketExists(ctx context.Context) (bool, error) {
	_, err := m.s3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(m.bucketName),
	})
	exists := true
	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) {
			switch apiError.(type) {
			case *types.NotFound:
				exists = false
				err = nil
			default:
				m.logger.Error("BucketExists: ", zap.Any("bucketName", m.bucketName), zap.Any("err", err))
			}
		}
	}

	return exists, err
}

func (m S3Manager) GetFileContents(ctx context.Context, objectKey string) (string, error) {
	result, err := m.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(m.bucketName),
		Key:    aws.String(objectKey),
	})

	if err != nil {
		var noKey *types.NoSuchKey
		if errors.As(err, &noKey) {
			log.Printf("Can't get object %s from bucket %s. No such key exists.\n", objectKey, m.bucketName)
			err = noKey
		} else {
			log.Printf("Couldn't get object %v:%v. Here's why: %v\n", m.bucketName, objectKey, err)
		}
		return "", err
	}

	body, err := io.ReadAll(result.Body)

	return string(body), err
}
