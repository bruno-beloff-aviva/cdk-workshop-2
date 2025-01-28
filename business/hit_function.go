package business

import (
	"cdk-workshop-2/business/hits"
	"cdk-workshop-2/dynamo_manager"
	"context"

	"github.com/joerdav/zapray"
	"go.uber.org/zap"
)

func HitFunction(logger *zapray.Logger, ctx context.Context, dbManager dynamo_manager.DynamoManager, path string) hits.Hits {
	logger.Info("HitFunction", zap.String("path", path))

	hit := hits.NewHits(path)
	dbManager.Get(ctx, &hit)

	hit.Increment()
	dbManager.Insert(ctx, &hit)

	return hit
}
