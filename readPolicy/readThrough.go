package readPolicy

import (
	"context"
	"github.com/rag594/konfigStore/cache"
	"github.com/rag594/konfigStore/model"
)

type ReadThroughPolicy[T model.TenantId, V any] struct {
	ConfigCacheOps cache.IConfigCacheRepo[T, V]
}

func NewReadThroughPolicy[T model.TenantId, V any](configCacheOps cache.IConfigCacheRepo[T, V]) *ReadThroughPolicy[T, V] {
	return &ReadThroughPolicy[T, V]{
		ConfigCacheOps: configCacheOps,
	}
}

func (r *ReadThroughPolicy[T, V]) GetConfig(ctx context.Context, key string, entityId T) (*V, error) {
	// https://docs.hazelcast.org/docs/latest/javadoc/com/hazelcast/map/MapLoader.html#loadAllKeys()
	// TODO: think on how to synchronise or load the cache from DB in a distributed mode
	return r.ConfigCacheOps.GetConfig(ctx, key, entityId)
}
