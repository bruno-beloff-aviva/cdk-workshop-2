// https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html#golang-handler-signatures
// https://stackoverflow.com/questions/37365009/how-to-invoke-an-aws-lambda-function-asynchronously
// proper logging: https://github.com/awsdocs/aws-lambda-developer-guide/blob/main/sample-apps/blank-go/function/main.go

package main

import (
	"cdk-workshop-2/business"
	"cdk-workshop-2/dynamo"
	"cdk-workshop-2/lambda/response"
	"fmt"

	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/joerdav/zapray"
	"go.uber.org/zap"
)

var tableName string
var logger *zapray.Logger
var dbManager dynamo.DynamoManager

func init() {
	var err error
	logger, err = zapray.NewProduction() //	.NewDevelopment()

	if err != nil {
		panic("failed to create logger: " + err.Error())
	}
	logger.Info("*** init")

	// logger.Info("logger level: " + logger.Level().String())
	// level := os.Getenv("LOG_LEVEL") // may also be set by ApplicationLogLevelV2: awslambda.ApplicationLogLevel_INFO
	// fmt.Println("log level: ", level)

	tableName = os.Getenv("HITS_TABLE_NAME")
	logger.Info("TableName: " + tableName)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logger.Info("handler: ", zap.String("request", fmt.Sprintf("%v", request)))

	sourceIP := request.RequestContext.Identity.SourceIP

	hit := business.Hit(logger, ctx, dbManager, request.Path)
	message := business.Hello(logger, sourceIP, hit)

	return response.New200(message), nil
}

func main() {
	logger.Info("*** main")

	ctx := context.Background() //	context.TODO(), config.WithSharedConfigProfile("bb")
	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		logger.Info("err: " + err.Error())
	}

	dbManager = dynamo.DynamoManager{Log: logger, TableName: tableName, DynamoDbClient: dynamodb.NewFromConfig(cfg)}
	is_available := dbManager.TableIsAvailable(ctx)

	if !is_available {
		panic("Table is not available")
	}

	lambda.Start(handler)
}
