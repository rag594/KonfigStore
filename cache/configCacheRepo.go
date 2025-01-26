package cache

import (
	"context"
	"github.com/rag594/konfigStore/model"
)

type IConfigCacheRepo[T model.TenantId, V any] interface {
	SaveConfig(ctx context.Context, key string, config *V) error
	GetConfigByKeyForEntity(ctx context.Context, key string) (*V, error)
	GetConfig(ctx context.Context, key string, entityId T) (*V, error)
	IsConfigCacheKeyPresent(ctx context.Context, key string) (bool, error)
}
