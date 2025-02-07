// https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html#golang-handler-signatures
// https://stackoverflow.com/questions/37365009/how-to-invoke-an-aws-lambda-function-asynchronously
// proper logging: https://github.com/awsdocs/aws-lambda-developer-guide/blob/main/sample-apps/blank-go/function/main.go

package main

import (
	"cdk-workshop-2/dynamomanager"
	"cdk-workshop-2/lambda/hello"
	"cdk-workshop-2/s3manager"
	"cdk-workshop-2/service"
	"context"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/joerdav/zapray"
)

func main() {
	logger, err1 := zapray.NewDevelopment() // log level is set using this: NewProduction(), NewDevelopment(), NewExample()

	if err1 != nil {
		panic("failed to create logger: " + err1.Error())
	}
	logger.Info("*** main")

	//	context...
	ctx := context.Background()
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
	dbManager := dynamomanager.NewDynamoManager(logger, cfg, tableName)
	tableIsAvailable := dbManager.TableIsAvailable(ctx)

	if !tableIsAvailable {
		panic("Table not available: " + tableName)
	}

	s3Manager := s3manager.NewS3Manager(logger, cfg, bucketName)
	bucket_is_available := s3Manager.BucketIsAvailable(ctx)

	if !bucket_is_available {
		panic("Bucket not available: " + bucketName)
	}

	hitManager := service.NewHitService(logger, dbManager)
	helloManager := service.NewHelloService(logger, s3Manager, objectName)

	handler := hello.NewHelloHandler(logger, helloManager, hitManager)

	lambda.StartWithOptions(handler.Handle, lambda.WithEnableSIGTERM(func() {
		logger.Info("Lambda container shutting down.")
	}))
}
