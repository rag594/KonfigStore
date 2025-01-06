package readPolicy

import (
	"context"
	"errors"
	"github.com/rag594/konfigStore/cache"
	"github.com/rag594/konfigStore/db"
	"github.com/rag594/konfigStore/model"
	"github.com/redis/go-redis/v9"
)

type CacheAsidePolicy[T model.TenantId, V any] struct {
	ConfigCacheOps cache.IConfigCacheRepo[T, V]
	ConfigDbRepo   db.IConfigDbRepo[T, V]
}

func NewCacheAsidePolicy[T model.TenantId, V any](configCacheOps cache.IConfigCacheRepo[T, V], configDbOps db.IConfigDbRepo[T, V]) *CacheAsidePolicy[T, V] {
	return &CacheAsidePolicy[T, V]{
		ConfigCacheOps: configCacheOps,
		ConfigDbRepo:   configDbOps,
	}
}

func (r *CacheAsidePolicy[T, V]) GetConfig(ctx context.Context, key string, entityId T) (*V, error) {
	// TODO: prevent the cache stampede(all hits lead to cache miss and request goes to DB)
	v, err := r.ConfigCacheOps.GetConfigByKeyForEntity(ctx, key)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	// config present in cache
	if v != nil {
		return v, nil
	}

	// fetch the config from primary data source
	vDb, err := r.ConfigDbRepo.GetConfigByKeyForEntity(ctx, entityId)

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
