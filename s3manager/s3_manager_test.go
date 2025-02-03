package s3_manager

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/joerdav/zapray"
	"go.uber.org/zap"
)

const bucketName = "cdk2-hello-bucket"
const objectName = "hello.txt"

var logger *zapray.Logger

func init() {
	var err error
	logger, err = zapray.NewProduction() //	.NewDevelopment()

	if err != nil {
		panic("failed to create logger: " + err.Error())
	}
}

func TestBucketIsAvailable(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		panic("err: " + err.Error())
	}

	s3Manager := NewS3Manager(logger, cfg, bucketName)

	is_available := s3Manager.BucketIsAvailable(ctx)
	logger.Info("TestBucketIsAvailable: ", zap.Any("is_available", is_available))
}

func TestGetFileContents(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		panic("err: " + err.Error())
	}

	s3Manager := NewS3Manager(logger, cfg, bucketName)

	body, err := s3Manager.GetFileContents(ctx, objectName)
	logger.Info("TestGetFileContents: ", zap.Any("body", body), zap.Any("err", err))
}
