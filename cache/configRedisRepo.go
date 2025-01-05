package cache

import (
	"context"
	"encoding/json"
	"github.com/rag594/konfigStore/model"
	"github.com/redis/go-redis/v9"
	"time"
)

type ConfigRedisRepo[T model.TenantId, V any] struct {
	RNonClusteredClient *redis.Client
	TTL                 time.Duration
}

func RegisterConfigForCacheOps[T model.TenantId, V any](client *redis.Client, ttl time.Duration) *ConfigRedisRepo[T, V] {
	return &ConfigRedisRepo[T, V]{
		RNonClusteredClient: client,
		TTL:                 ttl,
	}
}

func (c *ConfigRedisRepo[T, V]) SaveConfig(ctx context.Context, key string, config *V) error {
	b, err := json.Marshal(config)
	if err != nil {
		return err
	}
	_, err = c.RNonClusteredClient.Set(ctx, key, b, c.TTL).Result()

	if err != nil {
		return nil
	}

	return nil
}

func (c *ConfigRedisRepo[T, V]) GetConfigByKeyForEntity(ctx context.Context, key string) (*V, error) {
	val, err := c.RNonClusteredClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	v := new(V)

	err = json.Unmarshal([]byte(val), &v)

	if err != nil {
		return nil, err
	}

	return v, nil
}
