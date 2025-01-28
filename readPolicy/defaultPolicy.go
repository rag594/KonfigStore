package readPolicy

import (
	"context"
	"errors"
	"github.com/rag594/konfig-store/cache"
	"github.com/rag594/konfig-store/config"
	"github.com/rag594/konfig-store/db"
	"github.com/rag594/konfig-store/requestCoalescing"
	"github.com/redis/go-redis/v9"
)

type DefaultReadPolicy[T config.TenantId, V any] struct {
	ConfigCacheOps cache.IConfigCacheRepo[T, V]
	ConfigDbRepo   db.IConfigDbRepo[T, V]
}

func NewDefaultReadPolicy[T config.TenantId, V any](configCacheOps cache.IConfigCacheRepo[T, V], configDbOps db.IConfigDbRepo[T, V]) *DefaultReadPolicy[T, V] {
	return &DefaultReadPolicy[T, V]{
		ConfigCacheOps: configCacheOps,
		ConfigDbRepo:   configDbOps,
	}
}

func (r *DefaultReadPolicy[T, V]) GetConfig(ctx context.Context, key string, entityId T) (*V, error) {
	v, err := r.ConfigCacheOps.GetConfigByKeyForEntity(ctx, key)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	// config present in cache
	if v != nil {
		return v, nil
	}

	g := &requestCoalescing.Group[V]{}

	f := func() (*V, error) {
		// fetch the config from primary data source
		vDb, err := r.ConfigDbRepo.GetConfigByKeyForEntity(ctx, entityId)

		if err != nil {
			return nil, err
		}

		return vDb, nil
	}

	val, err := g.Do(key, f)

	if err != nil {
		return nil, err
	}

	// set the config in cache
	err = r.ConfigCacheOps.SaveConfig(ctx, key, val)

	if err != nil {
		return nil, err
	}

	return val, nil
}
