package db

import (
	"context"
	"github.com/rag594/konfigStore/model"
)

type IConfigDbRepo[T model.TenantId, V any] interface {
	SaveConfig(ctx context.Context, config *model.Config[T, V]) (int64, error)
	GetConfigByKeyForEntity(ctx context.Context, entityId T) (*V, error)
}
