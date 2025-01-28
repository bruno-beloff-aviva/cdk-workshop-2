package business

import (
	"cdk-workshop-2/business/hits"
	"cdk-workshop-2/s3_manager"
	"context"

	"fmt"
	"strings"

	"github.com/joerdav/zapray"
	"go.uber.org/zap"
)

type HelloManager struct {
	logger     *zapray.Logger
	s3Manager  s3_manager.S3Manager
	objectName string
}

func NewHelloManager(logger *zapray.Logger, s3Manager s3_manager.S3Manager, objectName string) HelloManager {
	return HelloManager{logger: logger, s3Manager: s3Manager, objectName: objectName}
}

func (m HelloManager) HelloFunction(ctx context.Context, client string, hits hits.Hits) string {
	m.logger.Info("HelloFunction", zap.String("client", client), zap.String("path", hits.Path))

	// time.Sleep(1 * time.Second)

	if client == "" {
		return "Hello Go world!"
	}

	message, err := m.s3Manager.GetFileContents(ctx, m.objectName)

	if err != nil {
		panic("HelloFunction: failed to get file contents: " + err.Error())
	}

	if strings.Contains(hits.Path, "panic") {
		m.logger.Error("Panic!")
		panic(hits.Path)
	}

	return fmt.Sprintf("%s\nHello Go world at %s from %s hits: %d", message, hits.Path, client, hits.Count)
}
