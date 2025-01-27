# KonfigStore - A cache backed multi-tenant config store[WIP]


A multi tenant configuration management. Configuration are like key/value pairs for a tenant.These configurations are specific to feature/workflow that a tenant uses in your system.

**Note**: Configurations in this case are not application specific configurations, these are configurations specific to a tenant(tenant can be any entity in the ecosystem)

### Features:
- [x] Values are of JSON format
- [x] Developers/users can define their custom values as structs.
- [x] default TTL and custom TTL for each config
- [x] Register/De-Register a configuration using a hook based mechanism. Registration of a configuration will simply register a configuration of specific type. 
- [x] Caching of configuration in redis
- [x] Persistent storage backed by rdbms/nosql databases(mysql)
- [x] Support of different cache write policies(write-through, write-around, write-back)
- [x] Cache stampede protection
- [ ] Settings/Options at each configuration
  - [x] ttl of the config in cache
  - [ ] distributed cache mode
  - [ ] db persistence mode
  - [ ] db timeout
  - [ ] cache timeout
  - [ ] eager refresh
  - [x] write policy
  - [x] custom configKey
  - [ ] custom cacheKey
- [ ] Monitoring
- [ ] Logging
- [ ] Current state of a group/category of configuration in cache/db.
- [ ] Grouping/Categorisation of config to ease the fetch of config of similar category/groups(do we need multiple groups per config?)
- [ ] Lineage of a configuration with time(will include changes in configuration, timestamp, updatedBy etc)
- [ ] Web based UI for managing the configurations per tenant(RBAC for listing, editing, viewing configurations per tenant). It is optional.


### How to use it

> **_NOTE:_**  Currently for database MYSQL v8.0 or greater and Redis as a distributed cache

> **_NOTE:_**  sqlx(https://github.com/jmoiron/sqlx) is currently supported for Db and go-redis/v9(https://github.com/redis/go-redis) for redis

#### Initialise the KonfigStore

```go
konfigStore := konfigStore.New(
		konfigStore.WithDatabase(&konfigStore.Database{Connection: dbConn}),
		konfigStore.WithRedisCache(&konfigStore.RedisCache{
			Connection:         redisConn,
			ClusterModeEnabled: false,
		}))
```

#### Write

```go
type IWritePolicy[T config.TenantId, V any] interface {
	SetConfig(ctx context.Context, cacheKey string, entityId T, value *V) error
}
```
##### Suppose our configuration of below type per tenant

```go
type SmartFeatureConfig struct {
    SmartSetting string
}
```

##### Register your new configuration

Three write policies are supported, we can select as per the needs per configuration level:
1. write-through
2. write-around
3. write-back

```go
smartFeatConfigRegister := configRegister.RegisterConfig[int, SmartFeatureConfig](
      kStore,
      configRegister.WithWritePolicy(writePolicy.WriteBack),
      configRegister.WithTTL(time.Minute),
)
```

##### Define your cache key

> **_NOTE:_**  Currently cacheKey is a function of entityId, in coming release we can provide custom func as well

```go
// define your new cache key(it is a function of entityId along with other options)
    cacheKeyForEntityC := cache.NewCacheKey[int, SmartFeatureConfig](20)

```

##### Write/Set the config

```go
    err := smartFeatConfigRegister.WritePolicy.SetConfig(context.Background(), cacheKeyForEntityC.DefaultValue(), 20, val)
    
    if err != nil {
        // handle error
    }
```



#### Read

```go
type IReadPolicy[T config.TenantId, V any] interface {
	GetConfig(ctx context.Context, cacheKey string, entityId T) (*V, error)
}
```

##### Suppose our configuration of below type per tenant

```go
type ComplexFeatureConfig struct {
	Enable string
}
```

##### Register your new configuration

```go
complexFeatConfigRegister := configRegister.RegisterConfig[int, ComplexFeatureConfig](
	kStore, 
	configRegister.WithTTL(time.Minute), 
	)
```

##### Define your cache key

> **_NOTE:_**  Currently cacheKey is a function of entityId, in coming release we can provide custom func as well

```go
// define your new cache key(it ius a function of entityId along with other options)
	cacheKeyForEntityA := cache.NewCacheKey[int, ComplexFeatureConfig](11)

```

##### Read the config

```go
// get the config for any entity
	config, err := complexFeatConfigRegister.ReadPolicy.GetConfig(context.Background(), cacheKeyForEntityA.DefaultValue(), 11)
	if err != nil {
		// handle error
	}
```

