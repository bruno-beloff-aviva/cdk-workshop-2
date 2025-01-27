// https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html#golang-handler-signatures
// https://stackoverflow.com/questions/37365009/how-to-invoke-an-aws-lambda-function-asynchronously
// proper logging: https://github.com/awsdocs/aws-lambda-developer-guide/blob/main/sample-apps/blank-go/function/main.go

package main

import (
	"cdk-workshop-2/business"
	"cdk-workshop-2/business/hits"
	"cdk-workshop-2/dynamo"
	"cdk-workshop-2/lambda/response"
	"encoding/json"
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

var TableName string
var Log *zapray.Logger
var Err error

func init() {
	Log, Err = zapray.NewProduction()
	if Err != nil {
		panic("failed to create logger: " + Err.Error())
	}
	Log.Info("hello_lambda init!!")

	Log.Info("log level: " + Log.Level().String())

	// level := os.Getenv("LOG_LEVEL") // may also be set by ApplicationLogLevelV2: awslambda.ApplicationLogLevel_INFO
	// fmt.Println("log level: ", level)

	TableName = os.Getenv("HITS_TABLE_NAME")
	Log.Info("TableName: " + TableName)
}

func handleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	Log.Info("handler log level: " + Log.Level().String())

	requestJSON, err := json.Marshal(request)
	if err != nil {
		Log.Error("Failed to marshal request: " + err.Error())
	} else {
		Log.Info(string(requestJSON))
	}
	Log.Info("handleRequest" + request.Path)

	var message string
	sourceIP := request.RequestContext.Identity.SourceIP

	message = business.Hello(Log, sourceIP, request.Path)

	return response.New200(message), nil
}

func main() {
	Log.Info("hello_lambda main!!")

	ctx := context.Background() //	context.TODO(), config.WithSharedConfigProfile("bb")
	cfg, err1 := config.LoadDefaultConfig(ctx)

	Log.Info("main: ", zap.String("cfg", fmt.Sprintf("%v", cfg)))
	if err1 != nil {
		Log.Info("err1: " + err1.Error())
	}

	manager := dynamo.DynamoManager{Log: Log, TableName: TableName, DynamoDbClient: dynamodb.NewFromConfig(cfg)}
	Log.Info("main: ", zap.Any("manager", manager))

	exists := manager.TableExists(context.Background())
	Log.Info("main: ", zap.Any("exists", exists))

	hit := hits.NewHits("/test")
	Log.Info("main: ", zap.Any("hit", hit))

	manager.Get(ctx, &hit)
	Log.Info("main got: ", zap.Any("hit", hit))

	hit.Increment()
	Log.Info("main incremented: ", zap.Any("hit", hit))

	manager.Insert(ctx, &hit)
	Log.Info("main inserted: ", zap.Any("hit", hit))

	lambda.Start(handleRequest)
}
