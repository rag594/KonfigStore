package configRegister

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/rag594/konfig-store/cache"
	"github.com/rag594/konfig-store/config"
	configDb "github.com/rag594/konfig-store/db"
	"github.com/rag594/konfig-store/konfigStore"
	"github.com/rag594/konfig-store/locks"
	"github.com/rag594/konfig-store/readPolicy"
	"github.com/rag594/konfig-store/writePolicy"
	"strings"
)

type ConfigRegister[T config.TenantId, V any] struct {
	configDbOps    configDb.IConfigDbRepo[T, V]
	configCacheOps cache.IConfigCacheRepo[T, V]
	ReadPolicy     readPolicy.IReadPolicy[T, V]
	WritePolicy    writePolicy.IWritePolicy[T, V]
	lockManager    locks.LockManager
	KonfigStore    *konfigStore.KonfigStore
}

func RegisterConfig[T config.TenantId, V any](konfigStore *konfigStore.KonfigStore, configOptsOptions ...ConfigOptsOptions) *ConfigRegister[T, V] {
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
	configDbOps := configDb.RegisterConfigForDbOps[T, V](konfigStore.Database.Connection, configOptions.ConfigKey)
	configRegister.configDbOps = configDbOps

	// Cache is optional for registration
	if konfigStore.RedisCache.Connection != nil {
		// Registers Cache ops for new config
		configCacheOps := cache.RegisterConfigForCacheOps[T, V](konfigStore.RedisCache.Connection, configDbOps, configOptions.TTL)
		configRegister.configCacheOps = configCacheOps
		configRegister.lockManager = locks.NewRedisLockManager(konfigStore.RedisCache.Connection)
	}

	// Default read policy
	configRegister.ReadPolicy = readPolicy.NewDefaultReadPolicy(configRegister.configCacheOps, configDbOps)

	// Register write policy - write-around
	if configOptions.IsWriteAroundPolicy() {
		configRegister.WritePolicy = writePolicy.NewWriteAroundPolicy(configRegister.configCacheOps, configDbOps, configRegister.lockManager)
	}

	// Register write policy - write-through
	if configOptions.IsWriteThroughPolicy() {
		configRegister.WritePolicy = writePolicy.NewWriteThroughPolicy(configRegister.configCacheOps, configDbOps, configRegister.lockManager)
	}

	// Register write policy - write-back
	if configOptions.IsWriteBackPolicy() {
		configRegister.WritePolicy = writePolicy.NewWriteBackPolicy(configRegister.configCacheOps, configDbOps, configRegister.lockManager)
	}

	return configRegister
}
