package main

import (
	// hitcounter "cdk-workshop/cdk-hitcounter"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	awslambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

const Project = "CDK2"

type CdkWorkshopStackProps struct {
	awscdk.StackProps
}

func NewCdkTable(scope constructs.Construct, id string) awsdynamodb.Table {
	this := constructs.NewConstruct(scope, &id)

	table := awsdynamodb.NewTable(this, jsii.String("Hits"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{Name: jsii.String("path"), Type: awsdynamodb.AttributeType_STRING},
	})

	return table
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

	// lambda...
	lambdaEnv := map[string]*string{
		"LOG_LEVEL":       jsii.String("INFO"),
		"HITS_TABLE_NAME": table.TableName(),
	}

	helloHandler := awslambdago.NewGoFunction(stack, jsii.String(Project+"HelloHandler"), &awslambdago.GoFunctionProps{
		Runtime:       awslambda.Runtime_PROVIDED_AL2(),
		Architecture:  awslambda.Architecture_ARM_64(),
		Entry:         jsii.String("lambda/"),
		Timeout:       awscdk.Duration_Seconds(jsii.Number(29)),
		LoggingFormat: awslambda.LoggingFormat_JSON,
		Environment:   &lambdaEnv,
	})

	// hitcounter...
	// hitcounter, table := hitcounter.NewHitCounter(stack, "HelloHitCounter", &hitcounter.HitCounterProps{
	// 	Downstream: helloHandler,
	// })

	table.GrantReadWriteData(helloHandler)

	// gateway...
	awsapigateway.NewLambdaRestApi(stack, jsii.String(Project+"HelloEndpoint"), &awsapigateway.LambdaRestApiProps{
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
