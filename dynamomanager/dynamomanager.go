// https://stackoverflow.com/questions/45405434/dynamodb-dynamic-atomic-update-of-mapped-values-with-aws-lambda-nodejs-runtime
// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/WorkingWithItems.html#WorkingWithItems.ConditionalUpdate
// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Expressions.ConditionExpressions.html

// To test with a live DynamoDB table, use:
// assume "bb"

// https://www.youtube.com/watch?v=bLY7-kTsQBM

package dynamomanager

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/jsii-runtime-go"
	"github.com/joerdav/zapray"
	"go.uber.org/zap"
)

type DynamoAble interface {
	GetKey() map[string]any
}

type DynamoManager struct {
	logger    *zapray.Logger
	dBClient  *dynamodb.Client
	tableName string
}

func NewDynamoManager(logger *zapray.Logger, cfg aws.Config, tableName string) DynamoManager {
	dBClient := dynamodb.NewFromConfig(cfg)
	return DynamoManager{logger: logger, dBClient: dBClient, tableName: tableName}
}

// --------------------------------------------------------------------------------------------------------------------

func (m DynamoManager) TableIsAvailable(ctx context.Context) bool {
	_, err := m.dBClient.DescribeTable(ctx, &dynamodb.DescribeTableInput{TableName: jsii.String(m.tableName)})

	if err != nil {
		m.logger.Error("TableIsAvailable: ", zap.Any("tableName", m.tableName), zap.Any("err", err))
	}

	return err == nil
}

func (m DynamoManager) Get(ctx context.Context, object DynamoAble) error {
	m.logger.Debug("Get: ", zap.Any("key", object.GetKey()))

	params := dynamodb.GetItemInput{
		Key:       getDBKey(object),
		TableName: jsii.String(m.tableName),
	}

	response, err := m.dBClient.GetItem(ctx, &params)

	if err != nil {
		m.logger.Error("GetItem: ", zap.Any("key", object.GetKey()), zap.Error(err))
	} else {
		err = attributevalue.UnmarshalMap(response.Item, &object)
		if err != nil {
			panic(err)
		}
	}

	return err
}

func (m DynamoManager) Put(ctx context.Context, object DynamoAble) error {
	m.logger.Debug("Put: ", zap.Any("object", object), zap.Any("key", object.GetKey()))

	item, err := attributevalue.MarshalMap(object)
	if err != nil {
		panic(err)
	}

	params := dynamodb.PutItemInput{
		TableName: jsii.String(m.tableName),
		Item:      item,
	}

	_, err = m.dBClient.PutItem(ctx, &params)
	if err != nil {
		m.logger.Error("PutItem: ", zap.Error(err))
	}

	return err
}

func (m DynamoManager) Increment(ctx context.Context, object DynamoAble, field string) (err error) {
	m.logger.Debug("Increment: ", zap.Any("object", object), zap.Any("key", object.GetKey()))

	defer func() {
		if err != nil && strings.Contains(err.Error(), "does not exist") {
			err = m.Put(ctx, object)
		}
	}()

	// increment
	update_params := dynamodb.UpdateItemInput{
		Key:                       getDBKey(object),
		TableName:                 jsii.String(m.tableName),
		ExpressionAttributeNames:  map[string]string{"#field": field},
		ExpressionAttributeValues: map[string]types.AttributeValue{":inc": &types.AttributeValueMemberN{Value: "1"}},
		UpdateExpression:          jsii.String("SET #field = #field + :inc"),
		ReturnValues:              types.ReturnValueAllNew,
	}

	_, err = m.dBClient.UpdateItem(ctx, &update_params)

	return err
}

// --------------------------------------------------------------------------------------------------------------------

func keyMap(objectKey map[string]any, marshal func(any) types.AttributeValue) map[string]types.AttributeValue {
	dBKey := make(map[string]types.AttributeValue, len(objectKey))

	for key, value := range objectKey {
		dBKey[key] = marshal(value)
	}

	return dBKey
}

func keyMarshal(objectValue any) types.AttributeValue {
	dBValue, err := attributevalue.Marshal(objectValue)

	if err != nil {
		panic(err)
	}

	return dBValue
}

func getDBKey(object DynamoAble) map[string]types.AttributeValue {
	return keyMap(object.GetKey(), keyMarshal)
}
