package writePolicy

import (
	"context"
	"github.com/rag594/konfig-store/cache"
	"github.com/rag594/konfig-store/config"
	"github.com/rag594/konfig-store/db"
	"github.com/rag594/konfig-store/locks"
)

type WriteThroughPolicy[T config.TenantId, V any] struct {
	ConfigCacheOps cache.IConfigCacheRepo[T, V]
	ConfigDbRepo   db.IConfigDbRepo[T, V]
	LockManager    locks.LockManager
}

func NewWriteThroughPolicy[T config.TenantId, V any](configCacheOps cache.IConfigCacheRepo[T, V], configDbOps db.IConfigDbRepo[T, V], locksManager locks.LockManager) *WriteThroughPolicy[T, V] {
	return &WriteThroughPolicy[T, V]{
		ConfigCacheOps: configCacheOps,
		ConfigDbRepo:   configDbOps,
		LockManager:    locksManager,
	}
}

func (w *WriteThroughPolicy[T, V]) SetConfig(ctx context.Context, cacheKey string, entityId T, value *V) error {

	w.LockManager.Lock(ctx, cacheKey)
	defer w.LockManager.Unlock(ctx)

	// Save in cache
	err := w.ConfigCacheOps.SaveConfig(ctx, cacheKey, value)
	if err != nil {
		return err
	}

	// Save in DB
	_, err = w.ConfigDbRepo.SaveConfig(ctx, &config.Config[T, V]{
		EntityId: entityId,
		Value:    config.ConfigValue[V]{Val: value},
	})

	if err != nil {
		return err
	}

	return nil

}
