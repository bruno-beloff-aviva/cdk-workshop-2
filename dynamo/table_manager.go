package dynamo

import (
	"cdk-workshop-2/log"
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/jsii-runtime-go"
)

type DynamoAble interface {
	GetKeys() map[string]interface{}
}

type TableManager struct {
	//	TODO: add context?
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

func (manager TableManager) TableExists(ctx context.Context) bool {
	_, err := manager.DynamoDbClient.DescribeTable(ctx, &dynamodb.DescribeTableInput{TableName: jsii.String(manager.TableName)})

	return err == nil
}

func (manager TableManager) Get(ctx context.Context, object DynamoAble) error {
	log.Logger.Debugf("Get key:%#v", getDBKey(object))

	response, err := manager.DynamoDbClient.GetItem(ctx, &dynamodb.GetItemInput{
		Key: getDBKey(object), TableName: jsii.String(manager.TableName),
	})
	log.Logger.Debugf("Get response:%#v", response)

	if err != nil {
		log.Logger.Errorf("GetItem for %v: %v\n", object.GetKeys(), err)
	} else {
		err = attributevalue.UnmarshalMap(response.Item, &object)
		if err != nil {
			log.Logger.Errorf("GetItem UnmarshalMap: %v\n", err)
		}
	}
	log.Logger.Debugf("Get object:%#v\n", object)

	return err
}

func (manager TableManager) Insert(ctx context.Context, object DynamoAble) error {
	log.Logger.Debugf("Insert object:%#v\n", object)
	log.Logger.Debugf("Insert key:%#v\n", getDBKey(object))

	item, err := attributevalue.MarshalMap(object)
	if err != nil {
		panic(err)
	}
	_, err = manager.DynamoDbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: jsii.String(manager.TableName), Item: item,
	})
	if err != nil {
		log.Logger.Errorf("PutItem: %v\n", err)
	}
	return err
}
