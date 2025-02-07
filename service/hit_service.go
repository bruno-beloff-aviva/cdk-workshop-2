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

	err := m.dbManager.Increment(ctx, &hit, "count")
	if err != nil {
		m.logger.Error("HitFunction: ", zap.Error(err))
	}

	m.dbManager.Get(ctx, &hit)

	return hit
}
