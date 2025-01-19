package main

import (
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/iancoleman/strcase"
	"github.com/rag594/konfigStore/cache"
	configDb "github.com/rag594/konfigStore/db"
	"github.com/rag594/konfigStore/locks"
	"github.com/rag594/konfigStore/model"
	"github.com/rag594/konfigStore/readPolicy"
	"strings"
)

type ConfigRegister[T model.TenantId, V any] struct {
	ConfigDbOps    configDb.IConfigDbRepo[T, V]
	ConfigCacheOps cache.IConfigCacheRepo[T, V]
	ReadPolicy     readPolicy.IReadPolicy[T, V]
	LockManager    locks.LockManager
}

func RegisterConfig[T model.TenantId, V any](configOptsOptions ...ConfigOptsOptions) *ConfigRegister[T, V] {
	// Registers a new config
	configRegister := &ConfigRegister[T, V]{}

	configOptions := &ConfigOpts{
		// sets the default configKey which is the struct name defined with snake case(A_B_C)
		ConfigKey: strcase.ToScreamingSnake(strings.Split(fmt.Sprintf("%T", *new(V)), ".")[1]),
	}

	for _, option := range configOptsOptions {
		option(configOptions)
	}

	// Registers Db ops for new config(this is registered by default)
	configDbOps := configDb.RegisterConfigForDbOps[T, V](configOptions.SqlxDbConn, configOptions.ConfigKey)
	configRegister.ConfigDbOps = configDbOps

	// Cache is optional for registration
	if configOptions.RedisNCClient != nil {
		// Registers Cache ops for new config
		configCacheOps := cache.RegisterConfigForCacheOps[T, V](configOptions.RedisNCClient, configDbOps, configOptions.TTL)
		configRegister.ConfigCacheOps = configCacheOps
		// using redSync for redis-locks
		pool := goredis.NewPool(configOptions.RedisNCClient)
		configRegister.LockManager = locks.NewRedisLockManager(redsync.New(pool))
	}

	// Register read policy - read-through
	if len(configOptions.ReadPolicy) != 0 && strings.Compare(configOptions.ReadPolicy.Value(), readPolicy.ReadThrough.Value()) == 0 {
		configRegister.ReadPolicy = readPolicy.NewReadThroughPolicy(configRegister.ConfigCacheOps)
	}

	// Register read policy - cache-aside
	if len(configOptions.ReadPolicy) != 0 && strings.Compare(configOptions.ReadPolicy.Value(), readPolicy.CacheAside.Value()) == 0 {
		configRegister.ReadPolicy = readPolicy.NewCacheAsidePolicy(configRegister.ConfigCacheOps, configDbOps, configRegister.LockManager)
	}

	return configRegister
}
