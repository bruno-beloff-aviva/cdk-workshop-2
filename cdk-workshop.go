// cdk deploy --profile bb

// https://github.com/aviva-verde/cdk-standards.git
// https://docs.aws.amazon.com/cdk/v2/guide/resources.html

package main

import (
	s3 "cdk-workshop-2/s3aviva"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	awslambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

const project = "CDK2"
const version = "0.1.4"
const region = "eu-west-2"

const bucketName = "cdk2-hello-bucket"

const objectName = "hello.txt"
const tableName = project + "HelloHitCounterTable"
const handlerId = project + "HelloHandler"
const endpointId = project + "HelloEndpoint"
const stackId = project + "WorkshopStack"

const logPrefix = project + "Hello" // not used by zap logger

type CdkWorkshopStackProps struct {
	awscdk.StackProps
}

func existingHitsTable(scope constructs.Construct, id string) awsdynamodb.ITable {
	return awsdynamodb.TableV2_FromTableName(scope, aws.String(id), aws.String(id))
}

func NewHitsTable(scope constructs.Construct, id string) awsdynamodb.ITable {
	defer existingHitsTable(scope, id) // after panic on "aldready exists"

	this := constructs.NewConstruct(scope, &id)

	table := awsdynamodb.NewTable(this, aws.String(tableName), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{Name: aws.String("path"), Type: awsdynamodb.AttributeType_STRING},
		TableName:    aws.String(tableName),
	})

	return table
}

func existingHelloBucket(stack awscdk.Stack, name string) awss3.IBucket {
	return awss3.Bucket_FromBucketName(stack, aws.String(name), aws.String(name))
}

func NewHelloBucket(stack awscdk.Stack, name string) awss3.IBucket {
	defer existingHelloBucket(stack, name) // after panic on "aldready exists"

	logConfig := s3.BucketLogConfiguration{
		BucketName: name,
		Region:     region,
		LogPrefix:  logPrefix,
	}

	props := s3.BucketProps{
		Stack:              stack,
		Name:               name + "-props",
		OverrideBucketName: aws.String(name),
		Versioned:          false,
		EventBridgeEnabled: false,
		LogConfiguration:   logConfig,
	}

	return s3.NewPrivateS3Bucket(props)
}

func NewHelloHandler(stack awscdk.Stack, lambdaEnv map[string]*string) awslambdago.GoFunction {
	helloHandler := awslambdago.NewGoFunction(stack, aws.String(handlerId), &awslambdago.GoFunctionProps{
		Runtime:       awslambda.Runtime_PROVIDED_AL2(),
		Architecture:  awslambda.Architecture_ARM_64(),
		Entry:         aws.String("lambda/"),
		Timeout:       awscdk.Duration_Seconds(aws.Float64(29)),
		LoggingFormat: awslambda.LoggingFormat_JSON,
		LogRetention:  awslogs.RetentionDays_FIVE_DAYS,
		Environment:   &lambdaEnv,
	})

	return helloHandler
}

func NewCdkWorkshopStack(scope constructs.Construct, id string, props *CdkWorkshopStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}

	//	stack...
	stack := awscdk.NewStack(scope, &id, &sprops)

	// lambda...
	lambdaEnv := map[string]*string{
		"HELLO_VERSION":     aws.String(version),
		"HELLO_BUCKET_NAME": aws.String(bucketName),
		"HELLO_OBJECT_NAME": aws.String(objectName),
		"HITS_TABLE_NAME":   aws.String(tableName),
	}

	helloHandler := NewHelloHandler(stack, lambdaEnv)

	// bucket...
	bucket := NewHelloBucket(stack, bucketName)
	bucket.GrantRead(helloHandler, nil)

	// table...
	table := NewHitsTable(stack, tableName)
	table.GrantReadWriteData(helloHandler)

	// gateway...
	restApiProps := awsapigateway.LambdaRestApiProps{Handler: helloHandler}
	awsapigateway.NewLambdaRestApi(stack, aws.String(endpointId), &restApiProps)

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)
	NewCdkWorkshopStack(app, stackId, &CdkWorkshopStackProps{})

	app.Synth(nil)
}
