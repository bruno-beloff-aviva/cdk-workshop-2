// https://stackoverflow.com/questions/24030059/skip-some-tests-with-go-test

package s3manager

import (
	"cdk-workshop-2/test"
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/joerdav/zapray"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

const bucketName = "cdk2-hello-bucket"
const objectName = "hello.txt"

var logger *zapray.Logger

func init() {
	var err error
	logger, err = zapray.NewDevelopment()

	if err != nil {
		panic("failed to create logger: " + err.Error())
	}
}

func TestBucketIsAvailable(t *testing.T) {
	test.SkipCI(t)

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		panic("err: " + err.Error())
	}

	s3Manager := NewS3Manager(logger, cfg, bucketName)

	isAvailable := s3Manager.BucketIsAvailable(ctx)
	logger.Info("TestBucketIsAvailable: ", zap.Any("isAvailable", isAvailable))

	assert.Equal(t, isAvailable, true)
}

func TestGetFileContents(t *testing.T) {
	test.SkipCI(t)

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		panic("err: " + err.Error())
	}

	s3Manager := NewS3Manager(logger, cfg, bucketName)

	body, err := s3Manager.GetFileContents(ctx, objectName)
	logger.Info("TestGetFileContents: ", zap.Any("body", body), zap.Any("err", err))

	assert.Equal(t, body, "April is the cruellest month, breeding\nLilacs out of the dead land, mixing\nMemory and desire, stirring\nDull roots with spring rain.\nWinter kept us warm, covering\nEarth in forgetful snow, feeding\nA little life with dried tubers.\n")
}
