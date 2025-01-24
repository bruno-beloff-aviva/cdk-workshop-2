package business

import (
	"cdk-workshop-2/logging"
	"fmt"
	"strings"
	"time"
)

func Hello(client string, request string) string {
	logging.Infof("businessFunction - client:%s request:%s", client, request)

	time.Sleep(1 * time.Second)

	if client == "" {
		return "Hello Go world!"
	}

	if strings.Contains(request, "panic") {
		logging.Error("Panic!")
		panic(request)
	}

	return fmt.Sprintf("Hello Go world at %s from %s", request, client)
}
