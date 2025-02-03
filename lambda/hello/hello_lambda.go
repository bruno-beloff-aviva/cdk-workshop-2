// https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html#golang-handler-signatures
// https://stackoverflow.com/questions/37365009/how-to-invoke-an-aws-lambda-function-asynchronously
// proper logging: https://github.com/awsdocs/aws-lambda-developer-guide/blob/main/sample-apps/blank-go/function/main.go

package main

import (
	"cdk-workshop-2/business"
	"cdk-workshop-2/dynamomanager"
	"cdk-workshop-2/lambda/response"
	"cdk-workshop-2/s3manager"
	"fmt"

	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/joerdav/zapray"
	"go.uber.org/zap"
)

var logger *zapray.Logger

var hitManager business.HitManager
var helloManager business.HelloManager

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
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logger.Info("handler: ", zap.String("request", fmt.Sprintf("%v", request)))

	sourceIP := request.RequestContext.Identity.SourceIP

	hit := hitManager.HitFunction(ctx, request.Path)
	message := helloManager.HelloFunction(ctx, sourceIP, hit)

	return response.New200(message), nil
}

func main() {
	logger.Info("*** main")

	//	context...
	ctx := context.Background() //	context.TODO(), config.WithSharedConfigProfile("bb")
	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		panic("err: " + err.Error())
	}

	//	environment...
	tableName := os.Getenv("HITS_TABLE_NAME")
	logger.Info("tableName: " + tableName)

	bucketName := os.Getenv("HELLO_BUCKET_NAME")
	logger.Info("bucketName: " + bucketName)

	objectName := os.Getenv("HELLO_OBJECT_NAME")
	logger.Info("objectName: " + objectName)

	//	managers...
	dbManager := dynamo_manager.NewDynamoManager(logger, cfg, tableName)
	table_is_available := dbManager.TableIsAvailable(ctx)

	if !table_is_available {
		panic("Table not available: " + tableName)
	}

	s3Manager := s3_manager.NewS3Manager(logger, cfg, bucketName)
	bucket_is_available := s3Manager.BucketIsAvailable(ctx)

	if !bucket_is_available {
		panic("Bucket not available: " + bucketName)
	}

	hitManager = business.NewHitManager(logger, dbManager)
	helloManager = business.NewHelloManager(logger, s3Manager, objectName)

	lambda.Start(handler)
}
