package db

import (
	"context"
	"github.com/rag594/konfigStore/config"
)

type IConfigDbRepo[T config.TenantId, V any] interface {
	SaveConfig(ctx context.Context, config *config.Config[T, V]) (int64, error)
	GetConfigByKeyForEntity(ctx context.Context, entityId T) (*V, error)
}
