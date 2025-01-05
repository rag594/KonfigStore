package readPolicy

import (
	"context"
	"errors"
	"github.com/rag594/konfigStore/cache"
	"github.com/rag594/konfigStore/db"
	"github.com/rag594/konfigStore/model"
	"github.com/redis/go-redis/v9"
)

type ReadThroughPolicy[T model.TenantId, V any] struct {
	ConfigDbOps    db.IConfigDbRepo[T, V]
	ConfigCacheOps cache.IConfigCacheRepo[T, V]
}

func NewReadThroughPolicy[T model.TenantId, V any](configDbOps db.IConfigDbRepo[T, V], configCacheOps cache.IConfigCacheRepo[T, V]) *ReadThroughPolicy[T, V] {
	return &ReadThroughPolicy[T, V]{
		ConfigCacheOps: configCacheOps,
		ConfigDbOps:    configDbOps,
	}
}

func (r *ReadThroughPolicy[T, V]) GetConfig(ctx context.Context, key string, entityId T) (*V, error) {
	v, err := r.ConfigCacheOps.GetConfigByKeyForEntity(ctx, key)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	// config present in cache
	if v != nil {
		return v, nil
	}

	// fetch the config from primary data source
	vDb, err := r.ConfigDbOps.GetConfigByKeyForEntity(ctx, key, entityId)

	if err != nil {
		return nil, err
	}

	// set the config in cache
	err = r.ConfigCacheOps.SaveConfig(ctx, key, vDb)

	if err != nil {
		return nil, err
	}

	return vDb, nil
}
