package s3

import (
	"fmt"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awskms"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsssm"
	"github.com/aws/jsii-runtime-go"
)

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

func NewEventDrivenBucket(stack awscdk.Stack, name string, props BucketProps) awss3.Bucket {
	props.Stack = stack
	props.Name = name
	bucket := NewPrivateS3Bucket(props)
	bucket.EnableEventBridgeNotification()
	return bucket
}

func NewPublicS3Bucket(props PublicBucketProps) awss3.Bucket {
	bucket := newS3Bucket(props.BucketProps)

	for _, corsRule := range props.Cors {
		bucket.AddCorsRule(corsRule)
	}
	return bucket
}

type MultiRegionS3BucketProps struct {
	Name                             string
	OverrideBucketName               *string
	Cors                             *[]*awss3.CorsRule
	LogConfiguration                 BucketLogConfiguration
	LifecycleRules                   *[]*awss3.LifecycleRule
	Encryption                       awss3.BucketEncryption
	EncryptionKey                    awskms.IKey
	EventBridgeEnabled               *bool
	PrimaryRegion                    string
	SecondaryRegion                  string
	AccessControl                    awss3.BucketAccessControl
	IntelligentTieringConfigurations *[]*awss3.IntelligentTieringConfiguration
	RemovalPolicy                    awscdk.RemovalPolicy
	AutoDeleteObjects                *bool
}

func NewMultiRegionBucket(stack awscdk.Stack, props MultiRegionS3BucketProps) awss3.Bucket {
	if props.PrimaryRegion == "" {
		props.PrimaryRegion = "eu-west-1"
	}

	if props.SecondaryRegion == "" {
		props.SecondaryRegion = "eu-west-2"
	}
	// for secondary region, hardcoded bucket names need region suffix
	if *stack.Region() == props.SecondaryRegion {
		if props.OverrideBucketName != nil {
			if len(*props.OverrideBucketName)+len(props.SecondaryRegion) > 63 {
				name := *props.OverrideBucketName
				props.OverrideBucketName = jsii.String(name[:len(name)-len("-"+props.SecondaryRegion)])
			}

			props.OverrideBucketName = jsii.String(*props.OverrideBucketName + "-" + props.SecondaryRegion)
		}
	}

	// create bucket
	bucket := awss3.NewBucket(stack, jsii.String(props.Name), &awss3.BucketProps{
		Encryption:        props.Encryption,
		EncryptionKey:     props.EncryptionKey,
		LifecycleRules:    props.LifecycleRules,
		BlockPublicAccess: awss3.BlockPublicAccess_BLOCK_ALL(),
		AccessControl:     props.AccessControl,
		Versioned:         jsii.Bool(true),
		EnforceSSL:        jsii.Bool(true),
		ServerAccessLogsBucket: awss3.Bucket_FromBucketAttributes(
			stack, jsii.String(fmt.Sprintf("%s-s3AccessLogs", props.Name)), &awss3.BucketAttributes{
				BucketName: jsii.String(props.LogConfiguration.BucketName),
				Region:     jsii.String(props.LogConfiguration.Region),
			}),
		ServerAccessLogsPrefix:           jsii.String(props.LogConfiguration.LogPrefix),
		EventBridgeEnabled:               props.EventBridgeEnabled,
		BucketName:                       props.OverrideBucketName,
		Cors:                             props.Cors,
		IntelligentTieringConfigurations: props.IntelligentTieringConfigurations,
		RemovalPolicy:                    props.RemovalPolicy,
		AutoDeleteObjects:                props.AutoDeleteObjects,
	})

	//// publish the primary bucket name to SSM
	awsssm.NewStringParameter(stack, jsii.String(props.Name+"-SSMParam"),
		&awsssm.StringParameterProps{
			ParameterName: jsii.String(
				"/" + *stack.StackName() + "/" + *stack.Region() + "/bucket/" + props.Name + "/name"),
			StringValue: bucket.BucketName(),
		})

	return bucket
}
