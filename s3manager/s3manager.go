// https://github.com/awsdocs/aws-doc-sdk-examples/blob/main/gov2/s3/actions/bucket_basics.go
// https://github.com/aws/aws-sdk-go/issues/1837

// https://stackoverflow.com/questions/45405434/dynamodb-dynamic-atomic-update-of-mapped-values-with-aws-lambda-nodejs-runtime
// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/WorkingWithItems.html#WorkingWithItems.ConditionalUpdate

package s3manager

import (
	"context"
	"errors"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
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

func (m S3Manager) BucketIsAvailable(ctx context.Context) bool {
	_, err := m.s3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(m.bucketName),
	})

	if err != nil {
		m.logger.Error("BucketIsAvailable: ", zap.Any("bucketName", m.bucketName), zap.Any("err", err))
	}

	return err == nil
}

func (m S3Manager) GetFileContents(ctx context.Context, objectKey string) (string, error) {
	result, err := m.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(m.bucketName),
		Key:    aws.String(objectKey),
	})

	if err != nil {
		var noKey *types.NoSuchKey
		if errors.As(err, &noKey) {
			err = noKey
		}
		m.logger.Error("GetFileContents: ", zap.Any("objectKey", objectKey), zap.Any("err", err))
		return "", err
	}

	body, err := io.ReadAll(result.Body)

	return string(body), err
}
