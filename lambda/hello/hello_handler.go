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
	helloService service.HelloService
	hitService   service.HitService
}

func NewHelloHandler(logger *zapray.Logger, helloManager service.HelloService, hitManager service.HitService) HelloHandler {
	return HelloHandler{
		logger:       logger,
		helloService: helloManager,
		hitService:   hitManager,
	}
}

func (h HelloHandler) Handle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	h.logger.Info("Handle: ", zap.String("request", fmt.Sprintf("%v", request)))

	sourceIP := request.RequestContext.Identity.SourceIP

	hit := h.hitService.Tally(ctx, request.Path)
	message, err := h.helloService.SayHello(ctx, sourceIP, hit)

	if err != nil {
		return response.New(500, err.Error()), err
	}

	return response.New(200, message), nil
}
