package cache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/rag594/konfigStore/config"
	configDb "github.com/rag594/konfigStore/db"
	"github.com/redis/go-redis/v9"
	"time"
)

type ConfigRedisRepo[T config.TenantId, V any] struct {
	RNonClusteredClient *redis.Client
	ConfigDbRepo        configDb.IConfigDbRepo[T, V]
	TTL                 time.Duration
}

func RegisterConfigForCacheOps[T config.TenantId, V any](client *redis.Client, configDbRepo configDb.IConfigDbRepo[T, V], ttl time.Duration) *ConfigRedisRepo[T, V] {
	return &ConfigRedisRepo[T, V]{
		RNonClusteredClient: client,
		TTL:                 ttl,
		ConfigDbRepo:        configDbRepo,
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

func (c *ConfigRedisRepo[T, V]) GetConfig(ctx context.Context, key string, entityId T) (*V, error) {
	v, err := c.GetConfigByKeyForEntity(ctx, key)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	// config present in cache
	if v != nil {
		return v, nil
	}

	return nil, nil
}

func (c *ConfigRedisRepo[T, V]) IsConfigCacheKeyPresent(ctx context.Context, key string) (bool, error) {
	val, err := c.RNonClusteredClient.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return val == 1, nil
}
