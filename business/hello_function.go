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

func HelloFunction(logger *zapray.Logger, ctx context.Context, s3Manager s3_manager.S3Manager, client string, hit hits.Hits) string {
	logger.Info("HelloFunction", zap.String("client", client), zap.String("path", hit.Path))

	// time.Sleep(1 * time.Second)

	if client == "" {
		return "Hello Go world!"
	}

	message, err := s3Manager.GetFileContents(ctx, "hello.txt")

	if err != nil {
		panic("HelloFunction: failed to get file contents: " + err.Error())
	}

	if strings.Contains(hit.Path, "panic") {
		logger.Error("Panic!")
		panic(hit.Path)
	}

	return fmt.Sprintf("%s\nHello Go world at %s from %s hits: %d", message, hit.Path, client, hit.Count)
}
