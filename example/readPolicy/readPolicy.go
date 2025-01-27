package main

import (
	"context"
	"fmt"
	"github.com/rag594/konfigStore/cache"
	"github.com/rag594/konfigStore/configRegister"
	"github.com/rag594/konfigStore/example"
	"github.com/rag594/konfigStore/konfigStore"
	"time"
)

type ComplexFeatureConfig struct {
	Enable string
}

func main() {
	// db/redis client used for testing
	dbConn := example.GetDbConnection()
	defer dbConn.Close()
	redisConn := example.NewRedisNonClusteredClient()

	/**
	Example on how to fetch the config
	*/

	// Read

	kStore := konfigStore.New(
		konfigStore.WithDatabase(&konfigStore.Database{Connection: dbConn}),
		konfigStore.WithRedisCache(&konfigStore.RedisCache{
			Connection:         redisConn,
			ClusterModeEnabled: false,
		}))

	// register your new configuration
	complexFeatConfigRegister := configRegister.RegisterConfig[int, ComplexFeatureConfig](
		kStore,
		configRegister.WithTTL(time.Minute),
	)

	// define your new cache key(it ius a function of entityId along with other options)
	cacheKeyForEntityA := cache.NewCacheKey[int, ComplexFeatureConfig](11)

	// get the config for any entity
	x, err := complexFeatConfigRegister.ReadPolicy.GetConfig(context.Background(), cacheKeyForEntityA.DefaultValue(), 11)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(x)

	cacheKeyForEntityB := cache.NewCacheKey[int, ComplexFeatureConfig](15)

	y, err := complexFeatConfigRegister.ReadPolicy.GetConfig(context.Background(), cacheKeyForEntityB.DefaultValue(), 15)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(y)
}
