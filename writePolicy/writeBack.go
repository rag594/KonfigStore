package writePolicy

import (
	"context"
	"github.com/rag594/konfigStore/cache"
	"github.com/rag594/konfigStore/db"
	"github.com/rag594/konfigStore/locks"
	"github.com/rag594/konfigStore/model"
)

type WriteBackPolicy[T model.TenantId, V any] struct {
	ConfigCacheOps cache.IConfigCacheRepo[T, V]
	ConfigDbRepo   db.IConfigDbRepo[T, V]
	LockManager    locks.LockManager
}

func NewWriteBackPolicy[T model.TenantId, V any](configCacheOps cache.IConfigCacheRepo[T, V], configDbOps db.IConfigDbRepo[T, V], locksManager locks.LockManager) *WriteBackPolicy[T, V] {
	return &WriteBackPolicy[T, V]{
		ConfigCacheOps: configCacheOps,
		ConfigDbRepo:   configDbOps,
		LockManager:    locksManager,
	}
}

func (w *WriteBackPolicy[T, V]) SetConfig(ctx context.Context, cacheKey string, entityId T, value *V) error {

	w.LockManager.Lock(ctx, cacheKey)
	defer w.LockManager.Unlock(ctx)

	// Save in cache
	err := w.ConfigCacheOps.SaveConfig(ctx, cacheKey, value)
	if err != nil {
		return err
	}

	// TODO: async update the primary DB

	return err

}
