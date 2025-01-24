package main

import (
	"cdk-workshop-2/loglevel"
	// "cdk-workshop/go_dynamo"

	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/labstack/gommon/log"
)

var Level log.Lvl
var Logger *log.Logger

func init() {
	fmt.Println("init!!")

	Logger = log.New("gohello")
	Level = loglevel.Level(os.Getenv("LOG_LEVEL"))
	fmt.Printf("gommon log level: %v\n", Level)

	Logger.SetLevel(Level) // may also be set by ApplicationLogLevelV2: awslambda.ApplicationLogLevel_INFO,
}

func handleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	Logger.Info(request)
	Logger.Info("handleRequest", request.Path)

	var message string
	sourceIP := request.RequestContext.Identity.SourceIP

	message = businessFunction(sourceIP, request.Path)

	return new200Response(message), nil
}

func new200Response(message string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       message + "\n",
	}
}

func businessFunction(sourceIP string, path string) string {
	Logger.Infof("businessFunction %s %s", sourceIP, path)

	time.Sleep(1 * time.Second)

	if sourceIP == "" {
		return "Hello Go world!"
	}

	if strings.Contains(path, "panic") {
		Logger.Panic("Panic!")
		panic(path)
	}

	return fmt.Sprintf("Hello Go world at %s from %s", path, sourceIP)
}

func main() {
	fmt.Println("main!!")

	cfg, err1 := config.LoadDefaultConfig(context.TODO()) //	config.WithSharedConfigProfile("bb")

	fmt.Printf("cfg: %v\n", cfg)
	fmt.Printf("err1: %v\n", err1)

	// basics := &go_dynamo.TableBasics{TableName: "CdkWorkshopStack-HelloHitCounterHits7AAEBF80-1UGY3S88AOYH0", DynamoDbClient: dynamodb.NewFromConfig(cfg)}
	// fmt.Printf("basics: %v\n", basics)

	// all, err2 := basics.TableExists(context.Background())

	// fmt.Printf("all: %v\n", all)
	// fmt.Printf("err2: %v\n", err2)

	lambda.Start(handleRequest)
}
