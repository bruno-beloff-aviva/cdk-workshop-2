// https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html#golang-handler-signatures
// https://stackoverflow.com/questions/37365009/how-to-invoke-an-aws-lambda-function-asynchronously
// proper logging: https://github.com/awsdocs/aws-lambda-developer-guide/blob/main/sample-apps/blank-go/function/main.go

package main

import (
	"cdk-workshop-2/business"
	"cdk-workshop-2/business/hits"
	"cdk-workshop-2/dynamo"
	"cdk-workshop-2/lambda/response"
	"cdk-workshop-2/log"

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
	fmt.Println("hello_lambda init!!")

	// level := os.Getenv("LOG_LEVEL") // may also be set by ApplicationLogLevelV2: awslambda.ApplicationLogLevel_INFO
	// fmt.Println("log level: ", level)

	log.Init("hello")

	TableName = os.Getenv("HITS_TABLE_NAME")
	log.Logger.Info("TableName", TableName)
}

func handleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Logger.Info(request)
	log.Logger.Info("handleRequest", request.Path)

	var message string
	sourceIP := request.RequestContext.Identity.SourceIP

	message = business.Hello(sourceIP, request.Path)

	// https://github.com/awsdocs/aws-lambda-developer-guide/blob/main/sample-apps/blank-go/function/main.go

	return response.New200(message), nil
}

func main() {
	fmt.Println("hello_lambda main!!")

	cfg, err1 := config.LoadDefaultConfig(context.TODO()) //	config.WithSharedConfigProfile("bb")

	log.Logger.Infof("cfg: %v", cfg)
	log.Logger.Infof("err1: %v", err1)

	manager := &dynamo.TableManager{TableName: TableName, DynamoDbClient: dynamodb.NewFromConfig(cfg)}
	log.Logger.Infof("manager: %v", manager)

	exists := manager.TableExists(context.Background())
	log.Logger.Infof("exists: %v", exists)

	var hit hits.Hits
	var err error

	// hit = hits.NewHits("/test")
	// log.Logger.Infof("hit: %v", hit)
	// err = manager.Insert(context.Background(), &hit)
	// log.Logger.Infof("hit: %v", hit)
	// log.Logger.Infof("err: %v", err)

	hit = hits.NewHits("/test")
	log.Logger.Infof("hit: %#v", hit)

	err = manager.Get(context.Background(), &hit)
	if err != nil {
		log.Logger.Error("Failed to get hit", err)
		return
	}
	log.Logger.Infof("got hit: %#v", hit)

	hit.Increment()
	log.Logger.Infof("incremented hit: %#v", hit)

	manager.Insert(context.Background(), &hit)
	log.Logger.Infof("inserted hit: %#v", hit)

	lambda.Start(handleRequest)
}
