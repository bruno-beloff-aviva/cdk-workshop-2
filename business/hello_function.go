package business

import (
	"cdk-workshop-2/business/hits"

	"fmt"
	"strings"
	"time"

	"github.com/joerdav/zapray"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
)

func Hello(logger *zapray.Logger, client string, hit hits.Hits) string {
	logger.Info("Hello", zap.String("client", client), zap.String("path", hit.Path))

	time.Sleep(1 * time.Second)

	if client == "" {
		return "Hello Go world!"
	}

	if strings.Contains(hit.Path, "panic") {
		log.Error("Panic!")
		panic(hit.Path)
	}

	return fmt.Sprintf("Hello Go world at %s from %s hits: %d", hit.Path, client, hit.Count)
}
