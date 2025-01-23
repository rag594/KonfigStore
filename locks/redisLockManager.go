package locks

import (
	"context"
	"github.com/go-redsync/redsync/v4"
)

type RedisLockManager struct {
	RedSync *redsync.Redsync
	Mutex   *redsync.Mutex
}

func NewRedisLockManager(redSync *redsync.Redsync) *RedisLockManager {

	return &RedisLockManager{
		RedSync: redSync,
	}
}

func (r *RedisLockManager) Lock(ctx context.Context, key string) error {
	mutex := r.RedSync.NewMutex(key)
	r.Mutex = mutex
	return r.Mutex.LockContext(ctx)

}

func (r *RedisLockManager) Unlock(ctx context.Context) (bool, error) {
	return r.Mutex.UnlockContext(ctx)
}
