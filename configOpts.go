package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/rag594/konfigStore/readPolicy"
	"github.com/redis/go-redis/v9"
	"time"
)

type ConfigOpts struct {
	SqlxDbConn    *sqlx.DB
	RedisNCClient *redis.Client
	ReadPolicy    readPolicy.ReadPolicy
	TTL           time.Duration
}

type ConfigOptsOptions func(*ConfigOpts)

func WithSqlXDbConn(dbConn *sqlx.DB) ConfigOptsOptions {
	return func(c *ConfigOpts) {
		c.SqlxDbConn = dbConn
	}
}

func WithRedisNCClient(client *redis.Client) ConfigOptsOptions {
	return func(c *ConfigOpts) {
		c.RedisNCClient = client
	}
}

func WithTTL(ttl time.Duration) ConfigOptsOptions {
	return func(c *ConfigOpts) {
		c.TTL = ttl
	}
}

func WithReadPolicy(readPolicy readPolicy.ReadPolicy) ConfigOptsOptions {
	return func(c *ConfigOpts) {
		c.ReadPolicy = readPolicy
	}
}
