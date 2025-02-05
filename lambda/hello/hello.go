package hello

import (
	"fmt"

	"cdk-workshop-2/lambda/response"
	"cdk-workshop-2/service"
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/joerdav/zapray"
	"go.uber.org/zap"
)

type HelloHandler struct {
	logger       *zapray.Logger
	helloManager service.HelloService
	hitManager   service.HitService
}

func NewHelloHandler(logger *zapray.Logger, helloManager service.HelloService, hitManager service.HitService) HelloHandler {
	return HelloHandler{
		logger:       logger,
		helloManager: helloManager,
		hitManager:   hitManager,
	}
}

func (h HelloHandler) Handle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	h.logger.Info("Handle: ", zap.String("request", fmt.Sprintf("%v", request)))

	sourceIP := request.RequestContext.Identity.SourceIP

	hit := h.hitManager.HitFunction(ctx, request.Path)
	message := h.helloManager.HelloFunction(ctx, sourceIP, hit)

	return response.New200(message), nil
}
