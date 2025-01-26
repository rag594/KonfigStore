package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/rag594/konfigStore/writePolicy"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

type ConfigOpts struct {
	SqlxDbConn    *sqlx.DB
	RedisNCClient *redis.Client
	WritePolicy   writePolicy.WritePolicy
	ConfigKey     string
	TTL           time.Duration
}

func (c *ConfigOpts) IsWriteAroundPolicy() bool {
	return len(c.WritePolicy) != 0 && strings.Compare(c.WritePolicy.Value(), writePolicy.WriteAround.Value()) == 0
}

func (c *ConfigOpts) IsWriteThroughPolicy() bool {
	return len(c.WritePolicy) != 0 && strings.Compare(c.WritePolicy.Value(), writePolicy.WriteThrough.Value()) == 0
}

func (c *ConfigOpts) IsWriteBackPolicy() bool {
	return len(c.WritePolicy) != 0 && strings.Compare(c.WritePolicy.Value(), writePolicy.WriteBack.Value()) == 0
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

func WithWritePolicy(writePolicy writePolicy.WritePolicy) ConfigOptsOptions {
	return func(c *ConfigOpts) {
		c.WritePolicy = writePolicy
	}
}

func WithConfigKey(configKey string) ConfigOptsOptions {
	return func(c *ConfigOpts) {
		c.ConfigKey = configKey
	}
}
