// cdk deploy --profile bb

// https://github.com/aviva-verde/cdk-standards.git
// https://docs.aws.amazon.com/cdk/v2/guide/resources.html

package main

import (
	s3 "cdk-workshop-2/s3aviva"
	"fmt"
	"strings"

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
const version = "0.1.2"

var bucketName = strings.ToLower(project) + "-hello-bucket"

const objectName = "hello.txt"
const tableName = project + "HelloHitCounterTable"
const endpointName = project + "HelloEndpoint"

type CdkWorkshopStackProps struct {
	awscdk.StackProps
}

func NewCdkTable(scope constructs.Construct, id string) awsdynamodb.Table {
	this := constructs.NewConstruct(scope, &id)

	table := awsdynamodb.NewTable(this, aws.String("Hits"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{Name: aws.String("path"), Type: awsdynamodb.AttributeType_STRING},
	})

	return table
}

func NewHelloBucket(stack awscdk.Stack, name string) awss3.Bucket { // using cdk-standards
	logConfig := s3.BucketLogConfiguration{
		BucketName: name,
		LogPrefix:  "HelloLogPrefix",
	}

	props := s3.BucketProps{
		Stack:              stack,
		Name:               name + "-props",
		OverrideBucketName: aws.String(name),
		Versioned:          false,
		EventBridgeEnabled: false,
		LogConfiguration:   logConfig,
	}

	fmt.Printf("props: %#v\n", props)

	return s3.NewPrivateS3Bucket(props)
}

func NewHelloHandler(stack awscdk.Stack, lambdaEnv map[string]*string) awslambdago.GoFunction {
	helloHandler := awslambdago.NewGoFunction(stack, aws.String(project+"HelloHandler"), &awslambdago.GoFunctionProps{
		Runtime:       awslambda.Runtime_PROVIDED_AL2(),
		Architecture:  awslambda.Architecture_ARM_64(),
		Entry:         aws.String("lambda/hello/"),
		Timeout:       awscdk.Duration_Seconds(aws.Float64(29)),
		LoggingFormat: awslambda.LoggingFormat_JSON,
		LogRetention:  awslogs.RetentionDays_FIVE_DAYS,
		Environment:   &lambdaEnv,
	})

	// ApplicationLogLevelV2: awslambda.ApplicationLogLevel_INFO,
	// SystemLogLevelV2:      awslambda.SystemLogLevel_INFO,

	return helloHandler
}

func NewCdkWorkshopStack(scope constructs.Construct, id string, props *CdkWorkshopStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}

	// searcher := awsopensearchservice.NewDomain(scope, aws.String("HelloSearchDomain"), &awsopensearchservice.DomainProps{
	// 	DomainName: aws.String("hello-search"),
	// })

	//	stack...
	stack := awscdk.NewStack(scope, &id, &sprops)

	// lambda...
	lambdaEnv := map[string]*string{
		"HELLO_OBJECT_NAME": aws.String(objectName),
		"HELLO_VERSION":     aws.String(version),
	}

	helloHandler := NewHelloHandler(stack, lambdaEnv)

	// bucket...
	var bucket awss3.IBucket

	bucket = awss3.Bucket_FromBucketName(stack, aws.String(bucketName), aws.String(bucketName))
	fmt.Printf("*** existing bucket: %#v\n", bucket.BucketName())

	if bucket == nil {
		bucket = NewHelloBucket(stack, bucketName)
	}

	bucket.GrantRead(helloHandler, nil)
	lambdaEnv["HELLO_BUCKET_NAME"] = bucket.BucketName()

	// table...
	var table awsdynamodb.ITable

	table = awsdynamodb.Table_FromTableName(stack, aws.String(tableName), aws.String(tableName))
	fmt.Printf("*** existing table: %#v\n", table.TableName())

	if table == nil {
		table = NewCdkTable(stack, tableName)
	}

	table.GrantReadWriteData(helloHandler)
	lambdaEnv["HITS_TABLE_NAME"] = table.TableName()

	// gateway...
	restApiProps := awsapigateway.LambdaRestApiProps{Handler: helloHandler}
	awsapigateway.NewLambdaRestApi(stack, aws.String(endpointName), &restApiProps)

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)
	NewCdkWorkshopStack(app, project+"WorkshopStack", &CdkWorkshopStackProps{})

	app.Synth(nil)
}
