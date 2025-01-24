package main

import (
	// hitcounter "cdk-workshop/cdk-hitcounter"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	awslambdago "github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type CdkWorkshopStackProps struct {
	awscdk.StackProps
}

func NewCdkWorkshopStack(scope constructs.Construct, id string, props *CdkWorkshopStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	helloHandler := awslambdago.NewGoFunction(stack, jsii.String("CDKHelloHandler2"), &awslambdago.GoFunctionProps{
		Runtime:       awslambda.Runtime_PROVIDED_AL2(),
		Architecture:  awslambda.Architecture_ARM_64(),
		Entry:         jsii.String("lambda/"),
		Timeout:       awscdk.Duration_Seconds(jsii.Number(30)),
		LoggingFormat: awslambda.LoggingFormat_JSON,
		Environment:   &map[string]*string{"LOG_LEVEL": jsii.String("INFO")},
	})

	// hitcounter, table := hitcounter.NewHitCounter(stack, "HelloHitCounter", &hitcounter.HitCounterProps{
	// 	Downstream: helloHandler,
	// })

	// table.GrantReadWriteData(helloHandler)

	awsapigateway.NewLambdaRestApi(stack, jsii.String("CDKHelloEndpoint2"), &awsapigateway.LambdaRestApiProps{
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

	NewCdkWorkshopStack(app, "CDKWorkshopStack2", &CdkWorkshopStackProps{})

	app.Synth(nil)
}
