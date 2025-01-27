package business

import (
	"fmt"
	"strings"
	"time"

	"github.com/joerdav/zapray"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
)

func Hello(logger *zapray.Logger, client string, request string) string {
	logger.Info("businessFunction", zap.String("client", client), zap.String("request", request))

	time.Sleep(1 * time.Second)

	if client == "" {
		return "Hello Go world!"
	}

	if strings.Contains(request, "panic") {
		log.Error("Panic!")
		panic(request)
	}

	return fmt.Sprintf("Hello Go world at %s from %s", request, client)
}
