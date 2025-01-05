package main

import (
	"github.com/rag594/konfigStore/cache"
	configDb "github.com/rag594/konfigStore/db"
	"github.com/rag594/konfigStore/model"
	"github.com/rag594/konfigStore/readPolicy"
	"strings"
)

type ConfigRegister[T model.TenantId, V any] struct {
	Config         *model.Config[T, V]
	ConfigDbOps    configDb.IConfigDbRepo[T, V]
	ConfigCacheOps cache.IConfigCacheRepo[T, V]
	ReadPolicy     readPolicy.IReadPolicy[T, V]
}

func RegisterConfig[T model.TenantId, V any](value *V, configOptsOptions ...ConfigOptsOptions) *ConfigRegister[T, V] {
	// Registers a new config
	config := model.NewConfig[T, V](value)

	configRegister := &ConfigRegister[T, V]{
		Config: config,
	}

	configOptions := &ConfigOpts{}

	for _, option := range configOptsOptions {
		option(configOptions)
	}

	// Registers Db ops for new config(this is registered by default)
	configDbOps := configDb.RegisterConfigForDbOps[T, V](configOptions.SqlxDbConn)
	configRegister.ConfigDbOps = configDbOps

	// Cache is optional for registration
	if configOptions.RedisNCClient != nil {
		// Registers Cache ops for new config
		configCacheOps := cache.RegisterConfigForCacheOps[T, V](configOptions.RedisNCClient, configOptions.TTL)
		configRegister.ConfigCacheOps = configCacheOps
	}

	// Register read policy
	if len(configOptions.ReadPolicy) != 0 && strings.Compare(configOptions.ReadPolicy.Value(), ReadThrough.Value()) == 0 {
		configRegister.ReadPolicy = readPolicy.NewReadThroughPolicy(configDbOps, configRegister.ConfigCacheOps)
	}

	return configRegister
}
