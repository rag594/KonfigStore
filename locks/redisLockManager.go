package locks

import (
	"context"
	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisLockManager struct {
	RNonClusteredClient *redis.Client
	RedisLockClient     *redislock.Client
	RedisLock           *redislock.Lock
}

func NewRedisLockManager(redisNClient *redis.Client) *RedisLockManager {
	return &RedisLockManager{
		RNonClusteredClient: redisNClient,
		RedisLockClient:     redislock.New(redisNClient),
	}
}

func (r *RedisLockManager) Lock(ctx context.Context, key string) error {
	lock, err := r.RedisLockClient.Obtain(ctx, key, time.Duration(1)*time.Second, nil)
	if err != nil {
		return err
	}
	r.RedisLock = lock
	return nil
}

func (r *RedisLockManager) Unlock(ctx context.Context) error {
	return r.RedisLock.Release(ctx)
}
