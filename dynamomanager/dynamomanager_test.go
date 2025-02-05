// https://stackoverflow.com/questions/24030059/skip-some-tests-with-go-test

package dynamomanager

import (
	"cdk-workshop-2/service/hits"
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/joerdav/zapray"
	"github.com/stretchr/testify/assert"
)

const tableName = "CDK2WorkshopStack-CDK2HelloHitCounterTableHits06BD259E-18FEY6SM4USEL"

var logger *zapray.Logger

func init() {
	var err error
	logger, err = zapray.NewDevelopment()

	if err != nil {
		panic("failed to create logger: " + err.Error())
	}
}

func skipCI(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")
	}
}

func TestGetDBKey(t *testing.T) {
	path := "/test"

	hit := hits.NewHits(path)
	fmt.Printf("hit: %#v\n", hit)

	var keyEntry types.AttributeValue

	key := getDBKey(&hit)
	keyEntry = key["path"]
	fmt.Printf("keyEntry: %#v\n", keyEntry.(*types.AttributeValueMemberS).Value)

	assert.Equal(t, keyEntry.(*types.AttributeValueMemberS).Value, path)
}

func TestTableIsAvailable(t *testing.T) {
	skipCI(t)

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
