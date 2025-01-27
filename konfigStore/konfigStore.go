package konfigStore

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Database struct {
	Connection *sqlx.DB
}

type RedisCache struct {
	Connection         *redis.Client
	ClusterModeEnabled bool
}

type KonfigStore struct {
	Database   *Database
	RedisCache *RedisCache
}

type Options func(*KonfigStore)

func New(opts ...Options) *KonfigStore {
	k := &KonfigStore{}

	for _, o := range opts {
		o(k)
	}

	return k
}

func WithDatabase(d *Database) Options {
	return func(store *KonfigStore) {
		store.Database = d
	}
}

func WithRedisCache(c *RedisCache) Options {
	return func(store *KonfigStore) {
		store.RedisCache = c
	}
}
