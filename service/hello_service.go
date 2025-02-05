package service

import (
	"cdk-workshop-2/s3manager"
	"cdk-workshop-2/service/hits"
	"context"
	"errors"

	"fmt"
	"strings"

	"github.com/joerdav/zapray"
	"go.uber.org/zap"
)

type HelloService struct {
	logger     *zapray.Logger
	s3Manager  s3manager.S3Manager
	objectName string
}

func NewHelloService(logger *zapray.Logger, s3Manager s3manager.S3Manager, objectName string) HelloService {
	return HelloService{logger: logger, s3Manager: s3Manager, objectName: objectName}
}

func (m HelloService) HelloFunction(ctx context.Context, client string, hits hits.Hits) (string, error) {
	m.logger.Info("HelloFunction", zap.String("client", client), zap.String("path", hits.Path))

	// time.Sleep(1 * time.Second)

	if client == "" {
		return "Hello Go world!", nil
	}

	message, err := m.s3Manager.GetFileContents(ctx, m.objectName)

	if err != nil {
		panic("HelloFunction: failed to get file contents: " + err.Error())
	}

	if strings.Contains(hits.Path, "panic") {
		m.logger.Error("Panic!")
		panic(hits.Path)
	}

	if strings.Contains(hits.Path, "error") {
		return fmt.Sprintf("%s\nError at %s from %s hits: %d", message, hits.Path, client, hits.Count), errors.New("Error at: " + hits.Path)
	}

	return fmt.Sprintf("%s\nHello Go world at %s from %s hits: %d", message, hits.Path, client, hits.Count), nil
}
