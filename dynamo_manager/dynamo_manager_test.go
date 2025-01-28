package dynamo_manager

import (
	"cdk-workshop-2/business/hits"
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/joerdav/zapray"
	// "github.com/stretchr/testify/assert"
)

const tableName = "CDK2WorkshopStack-CDK2HelloHitCounterTableHits06BD259E-18FEY6SM4USEL"

var logger *zapray.Logger

func init() {
	var err error
	logger, err = zapray.NewProduction() //	.NewDevelopment()

	if err != nil {
		panic("failed to create logger: " + err.Error())
	}
}

func TestGetDBKey(t *testing.T) {
	hit := hits.NewHits("/test")
	fmt.Printf("hit: %#v\n", hit)

	key := getDBKey(&hit)
	fmt.Printf("key: %#v\n", key)
}

func TestTableIsAvailable(t *testing.T) {
	ctx := context.Background() //	context.TODO(), config.WithSharedConfigProfile("bb")
	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		panic("err: " + err.Error())
	}

	dbManager := NewDynamoManager(logger, cfg, tableName)
	fmt.Printf("dbManager: %#v\n", dbManager)

	// https://github.com/gusaul/go-dynamock

	is_available := dbManager.TableIsAvailable(ctx)
	fmt.Printf("is_available: %#v\n", is_available)
}
