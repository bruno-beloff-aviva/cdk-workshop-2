package business

import (
	"cdk-workshop-2/business/hits"
	"cdk-workshop-2/dynamo_manager"
	"context"

	"github.com/joerdav/zapray"
	"go.uber.org/zap"
)

type HitManager struct {
	logger    *zapray.Logger
	dbManager dynamo_manager.DynamoManager
}

func NewHitManager(logger *zapray.Logger, dbManager dynamo_manager.DynamoManager) HitManager {
	return HitManager{logger: logger, dbManager: dbManager}
}

func (m *HitManager) HitFunction(ctx context.Context, path string) hits.Hits {
	m.logger.Info("HitFunction", zap.String("path", path))

	hit := hits.NewHits(path) // TODO: make this atomic
	m.dbManager.Get(ctx, &hit)

	hit.Increment()
	m.dbManager.Put(ctx, &hit)

	return hit
}
