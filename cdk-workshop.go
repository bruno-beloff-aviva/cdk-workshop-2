// cdk deploy --profile bb

package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	awslambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

const Project = "CDK2"

type CdkWorkshopStackProps struct {
	awscdk.StackProps
}

type HelloCounterProps struct {
	// Downstream is the function for which we want to count hits
	Downstream awslambda.IFunction
}

func NewCdkTable(scope constructs.Construct, id string) awsdynamodb.Table {
	this := constructs.NewConstruct(scope, &id)

	table := awsdynamodb.NewTable(this, aws.String("Hits"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{Name: aws.String("path"), Type: awsdynamodb.AttributeType_STRING},
	})

	return table
}

func NewHelloHandler(stack awscdk.Stack, table awsdynamodb.Table) awslambdago.GoFunction {
	lambdaEnv := map[string]*string{
		"HITS_TABLE_NAME": table.TableName(),
	}

	helloHandler := awslambdago.NewGoFunction(stack, aws.String(Project+"HelloHandler"), &awslambdago.GoFunctionProps{
		Runtime:       awslambda.Runtime_PROVIDED_AL2(),
		Architecture:  awslambda.Architecture_ARM_64(),
		Entry:         aws.String("lambda/hello/"),
		Timeout:       awscdk.Duration_Seconds(aws.Float64(29)),
		LoggingFormat: awslambda.LoggingFormat_JSON,
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

	//	stack...
	stack := awscdk.NewStack(scope, &id, &sprops)

	// table...
	table := NewCdkTable(stack, Project+"HelloHitCounterTable")

	// hello lambda...
	helloHandler := NewHelloHandler(stack, table)
	table.GrantReadWriteData(helloHandler)

	// gateway...
	restApiProps := awsapigateway.LambdaRestApiProps{Handler: helloHandler}
	awsapigateway.NewLambdaRestApi(stack, aws.String(Project+"HelloEndpoint"), &restApiProps)

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewCdkWorkshopStack(app, Project+"WorkshopStack", &CdkWorkshopStackProps{})

	app.Synth(nil)
}
