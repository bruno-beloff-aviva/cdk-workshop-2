package dynamo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
)

type TableBasics struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
}

func (basics TableBasics) TableExists(ctx context.Context) bool {
	_, err := basics.DynamoDbClient.DescribeTable(ctx, &dynamodb.DescribeTableInput{TableName: aws.String(basics.TableName)})

	return err == nil
}
