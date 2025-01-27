package business

import (
	"cdk-workshop-2/business/hits"
	"cdk-workshop-2/dynamo"
	"context"

	"github.com/joerdav/zapray"
	"go.uber.org/zap"
)

func Hit(logger *zapray.Logger, ctx context.Context, dbManager dynamo.DynamoManager, path string) hits.Hits {
	logger.Info("Hit", zap.String("path", path))

	hit := hits.NewHits(path)
	dbManager.Get(ctx, &hit)

	hit.Increment()
	dbManager.Insert(ctx, &hit)

	return hit
}
