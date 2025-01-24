package main

import (
	"cdk-workshop-2/business"
	"cdk-workshop-2/dynamo"
	"cdk-workshop-2/logging"

	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var TableName string

func init() {
	fmt.Println("init!!")

	level := os.Getenv("LOG_LEVEL") // may also be set by ApplicationLogLevelV2: awslambda.ApplicationLogLevel_INFO
	fmt.Println("log level: ", level)

	logging.Init("hello2", level)

	TableName = os.Getenv("HITS_TABLE_NAME")
	logging.Info("TableName", TableName)
}

func handleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logging.Info(request)
	logging.Info("handleRequest", request.Path)

	var message string
	sourceIP := request.RequestContext.Identity.SourceIP

	message = business.Hello(sourceIP, request.Path)

	return new200Response(message), nil
}

func new200Response(message string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       message + "\n",
	}
}

func main() {
	fmt.Println("main!!")

	cfg, err1 := config.LoadDefaultConfig(context.TODO()) //	config.WithSharedConfigProfile("bb")

	fmt.Printf("cfg: %v\n", cfg)
	fmt.Printf("err1: %v\n", err1)

	basics := &dynamo.TableBasics{TableName: TableName, DynamoDbClient: dynamodb.NewFromConfig(cfg)}
	fmt.Printf("basics: %v\n", basics)

	exists := basics.TableExists(context.Background())
	fmt.Printf("exists: %v\n", exists)

	lambda.Start(handleRequest)
}
