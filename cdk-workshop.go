package main

import (
	// hitcounter "cdk-workshop/cdk-hitcounter"

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

func NewHelloHandler(stack awscdk.Stack, table awsdynamodb.Table, props HelloCounterProps) awslambdago.GoFunction {
	lambdaEnv := map[string]*string{
		"HITS_TABLE_NAME":          table.TableName(),
		"DOWNSTREAM_FUNCTION_NAME": props.Downstream.FunctionName(),
	}

	helloHandler := awslambdago.NewGoFunction(stack, aws.String(Project+"HelloHandler"), &awslambdago.GoFunctionProps{
		Runtime:               awslambda.Runtime_PROVIDED_AL2(),
		Architecture:          awslambda.Architecture_ARM_64(),
		Entry:                 aws.String("lambda/hello/"),
		Timeout:               awscdk.Duration_Seconds(aws.Float64(29)),
		ApplicationLogLevelV2: awslambda.ApplicationLogLevel_DEBUG,
		LoggingFormat:         awslambda.LoggingFormat_JSON,
		Environment:           &lambdaEnv,
	})

	table.GrantReadWriteData(helloHandler)
	props.Downstream.GrantInvoke(helloHandler)

	return helloHandler
}

func NewHitHandler(stack awscdk.Stack, table awsdynamodb.Table) awslambdago.GoFunction {
	lambdaEnv := map[string]*string{
		"HITS_TABLE_NAME": table.TableName(),
	}

	hitHandler := awslambdago.NewGoFunction(stack, aws.String(Project+"HitHandler"), &awslambdago.GoFunctionProps{
		Runtime:               awslambda.Runtime_PROVIDED_AL2(),
		Architecture:          awslambda.Architecture_ARM_64(),
		Entry:                 aws.String("lambda/hitcounter/"),
		Timeout:               awscdk.Duration_Seconds(aws.Float64(29)),
		ApplicationLogLevelV2: awslambda.ApplicationLogLevel_DEBUG,
		LoggingFormat:         awslambda.LoggingFormat_JSON,
		Environment:           &lambdaEnv,
	})

	table.GrantReadWriteData(hitHandler)

	return hitHandler
}

// func NewHitCounter(scope constructs.Construct, id string, props *HitCounterProps) (HitCounter, awsdynamodb.Table) {
// 	this := constructs.NewConstruct(scope, &id)

// 	table := awsdynamodb.NewTable(this, jsii.String("Hits"), &awsdynamodb.TableProps{
// 		PartitionKey: &awsdynamodb.Attribute{Name: jsii.String("path"), Type: awsdynamodb.AttributeType_STRING},
// 	})

// 	handler := awslambda.NewFunction(this, jsii.String("HitCounterHandler"), &awslambda.FunctionProps{
// 		Runtime: awslambda.Runtime_NODEJS_16_X(),
// 		Handler: jsii.String("hitcounter.handler"),
// 		Code:    awslambda.Code_FromAsset(jsii.String("js_lambda"), nil),
// 		Timeout: awscdk.Duration_Seconds(jsii.Number(30)),
// 		Environment: &map[string]*string{
// 			"DOWNSTREAM_FUNCTION_NAME": props.Downstream.FunctionName(),
// 			"HITS_TABLE_NAME":          table.TableName(),
// 		},
// 	})

// 	table.GrantReadWriteData(handler)
// 	props.Downstream.GrantInvoke(handler)

// 	return &hitCounter{this, handler}, table
// }

func NewCdkWorkshopStack(scope constructs.Construct, id string, props *CdkWorkshopStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}

	//	stack...
	stack := awscdk.NewStack(scope, &id, &sprops)

	// table...
	table := NewCdkTable(stack, Project+"HelloHitCounterTable")

	// hitcounter lambda...
	hitHandler := NewHitHandler(stack, table)

	// hello lambda...
	helloProps := HelloCounterProps{Downstream: hitHandler}
	helloHandler := NewHelloHandler(stack, table, helloProps)

	// props.Downstream.GrantInvoke(hitHandler)

	// hitcounter lambda...
	// hitcounter, table := hitcounter.NewHitCounter(stack, "HelloHitCounter", &hitcounter.HitCounterProps{
	// 	Downstream: helloHandler,
	// })

	// gateway...
	awsapigateway.NewLambdaRestApi(stack, aws.String(Project+"HelloEndpoint"), &awsapigateway.LambdaRestApiProps{
		// Handler: hitcounter.Handler(),
		// EndpointConfiguration: &awsapigateway.EndpointConfiguration{
		// 	Types: &[]awsapigateway.EndpointType{awsapigateway.EndpointType_REGIONAL},
		// },
		Handler: helloHandler,
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	// NewCdkWorkshopStack(app, "CdkWorkshopStack", &CdkWorkshopStackProps{
	// 	awscdk.StackProps{
	// 		Env: env(),
	// 	},
	// })

	NewCdkWorkshopStack(app, Project+"WorkshopStack", &CdkWorkshopStackProps{})

	app.Synth(nil)
}
