package writePolicy

import (
	"context"
	"github.com/rag594/konfigStore/cache"
	"github.com/rag594/konfigStore/db"
	"github.com/rag594/konfigStore/locks"
	"github.com/rag594/konfigStore/model"
	"sync"
)

type WriteBackPolicy[T model.TenantId, V any] struct {
	ConfigCacheOps cache.IConfigCacheRepo[T, V]
	ConfigDbRepo   db.IConfigDbRepo[T, V]
	LockManager    locks.LockManager
	Wg             sync.WaitGroup
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

	// set the value in DB async
	w.setConfigAsyncInDb(ctx, cacheKey, entityId, value)

	return err

}

func (w *WriteBackPolicy[T, V]) setConfigAsyncInDb(ctx context.Context, cacheKey string, entityId T, value *V) {

	w.LockManager.Lock(ctx, cacheKey)
	defer func() {
		w.Wg.Done()
		w.LockManager.Unlock(ctx)

	}()

	w.Wg.Add(1)

	go w.ConfigDbRepo.SaveConfig(ctx, &model.Config[T, V]{
		EntityId: entityId,
		Value:    model.ConfigValue[V]{Val: value},
	})
}
