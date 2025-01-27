package dynamo

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
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
	logger    *zapray.Logger
	dBClient  *dynamodb.Client
	tableName string
}

func NewDynamoManager(logger *zapray.Logger, cfg aws.Config, tableName string) DynamoManager {
	dBClient := dynamodb.NewFromConfig(cfg)
	return DynamoManager{logger: logger, dBClient: dBClient, tableName: tableName}
}

func (m DynamoManager) TableIsAvailable(ctx context.Context) bool {
	_, err := m.dBClient.DescribeTable(ctx, &dynamodb.DescribeTableInput{TableName: jsii.String(m.tableName)})
	return err == nil
}

func (m DynamoManager) Get(ctx context.Context, object DynamoAble) error {
	m.logger.Debug("Get: ", zap.Any("key", getDBKey(object)))

	response, err := m.dBClient.GetItem(ctx, &dynamodb.GetItemInput{
		Key: getDBKey(object), TableName: jsii.String(m.tableName),
	})
	m.logger.Debug("Get: ", zap.Any("response", response))

	if err != nil {
		m.logger.Error("GetItem: ", zap.Any("key", object.GetKeys()), zap.Error(err))
	} else {
		err = attributevalue.UnmarshalMap(response.Item, &object)
		if err != nil {
			m.logger.Error("GetItem UnmarshalMap: ", zap.Error(err))
		}
	}
	m.logger.Debug("Get: ", zap.Any("object", object))

	return err
}

func (m DynamoManager) Insert(ctx context.Context, object DynamoAble) error {
	m.logger.Debug("Insert: ", zap.Any("object", object), zap.Any("key", getDBKey(object)))

	item, err := attributevalue.MarshalMap(object)
	if err != nil {
		panic(err)
	}
	_, err = m.dBClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: jsii.String(m.tableName), Item: item,
	})
	if err != nil {
		m.logger.Error("PutItem: ", zap.Error(err))
	}
	return err
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
