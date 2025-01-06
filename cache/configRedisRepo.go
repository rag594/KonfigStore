package cache

import (
	"context"
	"encoding/json"
	"errors"
	configDb "github.com/rag594/konfigStore/db"
	"github.com/rag594/konfigStore/model"
	"github.com/redis/go-redis/v9"
	"time"
)

type ConfigRedisRepo[T model.TenantId, V any] struct {
	RNonClusteredClient *redis.Client
	ConfigDbRepo        configDb.IConfigDbRepo[T, V]
	TTL                 time.Duration
}

func RegisterConfigForCacheOps[T model.TenantId, V any](client *redis.Client, configDbRepo configDb.IConfigDbRepo[T, V], ttl time.Duration) *ConfigRedisRepo[T, V] {
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
	// TODO: synchronization mechanisms between cache and primary data store
	v, err := c.GetConfigByKeyForEntity(ctx, key)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	// config present in cache
	if v != nil {
		return v, nil
	}

	// fetch the config from primary data source
	vDb, err := c.ConfigDbRepo.GetConfigByKeyForEntity(ctx, entityId)

	if err != nil {
		return nil, err
	}

	// set the config in cache
	err = c.SaveConfig(ctx, key, vDb)

	if err != nil {
		return nil, err
	}

	return vDb, nil
}
