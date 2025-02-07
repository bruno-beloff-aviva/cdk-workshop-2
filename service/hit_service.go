package service

import (
	"cdk-workshop-2/dynamomanager"
	"cdk-workshop-2/service/hits"
	"context"

	"github.com/joerdav/zapray"
	"go.uber.org/zap"
)

type HitService struct {
	logger    *zapray.Logger
	dbManager dynamomanager.DynamoManager
}

func NewHitService(logger *zapray.Logger, dbManager dynamomanager.DynamoManager) HitService {
	return HitService{logger: logger, dbManager: dbManager}
}

func (m *HitService) HitFunction(ctx context.Context, path string) hits.Hits {
	m.logger.Info("HitFunction", zap.String("path", path))

	hit := hits.NewHits(path)
	// m.dbManager.Get(ctx, &hit)

	// hit.Increment()
	// m.dbManager.Put(ctx, &hit) // TODO: make this atomic

	err := m.dbManager.Increment(ctx, &hit, "count")
	m.logger.Info("HitFunction: got err from Increment: ", zap.Any("err", err))

	m.dbManager.Get(ctx, &hit)

	return hit
}
