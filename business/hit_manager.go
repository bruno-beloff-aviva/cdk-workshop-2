package business

import (
	"cdk-workshop-2/business/hits"
	"cdk-workshop-2/dynamomanager"
	"context"

	"github.com/joerdav/zapray"
	"go.uber.org/zap"
)

type HitManager struct {
	logger    *zapray.Logger
	dbManager dynamomanager.DynamoManager
}

func NewHitManager(logger *zapray.Logger, dbManager dynamomanager.DynamoManager) HitManager {
	return HitManager{logger: logger, dbManager: dbManager}
}

func (m *HitManager) HitFunction(ctx context.Context, path string) hits.Hits {
	m.logger.Info("HitFunction", zap.String("path", path))

	hit := hits.NewHits(path)
	m.dbManager.Get(ctx, &hit)

	hit.Increment()
	m.dbManager.Put(ctx, &hit) // TODO: make this atomic

	return hit
}
