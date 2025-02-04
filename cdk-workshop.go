// cdk deploy --profile bb

// https://github.com/aviva-verde/cdk-standards.git
// https://docs.aws.amazon.com/cdk/v2/guide/resources.html

package main

import (
	"fmt"

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
const bucketName = "cdk2-hello-bucket"
const objectName = "hello.txt"

type BucketLogConfiguration struct {
	BucketName string
	Region     string
	LogPrefix  string
}

type BucketProps struct {
	Stack awscdk.Stack
	Name  string
	// OverrideFunctionName [optional] override the bucket name used for the bucket instead of setting the name and ID the same
	OverrideBucketName *string
	Versioned          bool
	LogConfiguration   BucketLogConfiguration
	LifecycleRules     []*awss3.LifecycleRule
	EventBridgeEnabled bool
}

type PublicBucketProps struct {
	BucketProps
	Cors []*awss3.CorsRule
}

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
	logConfig := BucketLogConfiguration{
		BucketName: name,
		LogPrefix:  "HelloLogPrefix",
	}

	props := BucketProps{
		Stack:              stack,
		Name:               name + "-props",
		OverrideBucketName: aws.String(name),
		Versioned:          false,
		EventBridgeEnabled: false,
		LogConfiguration:   logConfig,
	}

	fmt.Printf("props: %#v\n", props)

	return NewPrivateS3Bucket(props)
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

	// table...
	table := NewCdkTable(stack, project+"HelloHitCounterTable")

	// bucket...
	bucket := NewHelloBucket(stack, bucketName)

	// lambda...
	lambdaEnv := map[string]*string{
		"HITS_TABLE_NAME":   table.TableName(),
		"HELLO_BUCKET_NAME": bucket.BucketName(),
		"HELLO_OBJECT_NAME": aws.String(objectName),
		"HELLO_VERSION":     aws.String("0.1.1"),
	}

	helloHandler := NewHelloHandler(stack, lambdaEnv)

	table.GrantReadWriteData(helloHandler)
	bucket.GrantRead(helloHandler, nil)

	// searcher.GrantReadWrite(helloHandler)

	// gateway...
	restApiProps := awsapigateway.LambdaRestApiProps{Handler: helloHandler}
	awsapigateway.NewLambdaRestApi(stack, aws.String(project+"HelloEndpoint"), &restApiProps)

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)
	NewCdkWorkshopStack(app, project+"WorkshopStack", &CdkWorkshopStackProps{})

	app.Synth(nil)
}

func newS3Bucket(props BucketProps) awss3.Bucket {
	return awss3.NewBucket(props.Stack, jsii.String(props.Name), &awss3.BucketProps{
		Encryption:        awss3.BucketEncryption_S3_MANAGED,
		LifecycleRules:    &props.LifecycleRules,
		BlockPublicAccess: awss3.BlockPublicAccess_BLOCK_ALL(),
		Versioned:         jsii.Bool(props.Versioned),
		EnforceSSL:        jsii.Bool(true),
		ServerAccessLogsBucket: awss3.Bucket_FromBucketAttributes(props.Stack, jsii.String(fmt.Sprintf("%s-s3AccessLogs", props.Name)), &awss3.BucketAttributes{
			BucketName: jsii.String(props.LogConfiguration.BucketName),
			Region:     jsii.String(props.LogConfiguration.Region),
		}),
		ServerAccessLogsPrefix: jsii.String(props.LogConfiguration.LogPrefix),
		EventBridgeEnabled:     jsii.Bool(props.EventBridgeEnabled),
		BucketName:             props.OverrideBucketName,
	})
}

func NewPrivateS3Bucket(props BucketProps) awss3.Bucket {
	return newS3Bucket(props)
}
