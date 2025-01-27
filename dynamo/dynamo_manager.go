package dynamo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/jsii-runtime-go"
	"github.com/joerdav/zapray"
	"go.uber.org/zap"
)

type DynamoAble interface {
	GetKeys() map[string]interface{}
}

type DynamoManager struct {
	Log            *zapray.Logger
	DynamoDbClient *dynamodb.Client
	TableName      string
}

func getDBKey(object DynamoAble) map[string]types.AttributeValue {
	var dbKey = make(map[string]types.AttributeValue)

	for key, value := range object.GetKeys() {
		attrValue, err := attributevalue.Marshal(value)

		if err != nil {
			panic(err)
		}

		dbKey[key] = attrValue
	}

	return dbKey
}

func (m DynamoManager) TableExists(ctx context.Context) bool {
	_, err := m.DynamoDbClient.DescribeTable(ctx, &dynamodb.DescribeTableInput{TableName: jsii.String(m.TableName)})

	return err == nil
}

func (m DynamoManager) Get(ctx context.Context, object DynamoAble) error {
	m.Log.Debug("Get: ", zap.Any("key", getDBKey(object)))

	response, err := m.DynamoDbClient.GetItem(ctx, &dynamodb.GetItemInput{
		Key: getDBKey(object), TableName: jsii.String(m.TableName),
	})
	m.Log.Debug("Get: ", zap.Any("response", response))

	if err != nil {
		m.Log.Error("GetItem: ", zap.Any("key", object.GetKeys()), zap.Error(err))
	} else {
		err = attributevalue.UnmarshalMap(response.Item, &object)
		if err != nil {
			m.Log.Error("GetItem UnmarshalMap: ", zap.Error(err))
		}
	}
	m.Log.Debug("Get: ", zap.Any("object", object))

	return err
}

func (m DynamoManager) Insert(ctx context.Context, object DynamoAble) error {
	m.Log.Debug("Insert: ", zap.Any("object", object), zap.Any("key", getDBKey(object)))

	item, err := attributevalue.MarshalMap(object)
	if err != nil {
		panic(err)
	}
	_, err = m.DynamoDbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: jsii.String(m.TableName), Item: item,
	})
	if err != nil {
		m.Log.Error("PutItem: ", zap.Error(err))
	}
	return err
}
