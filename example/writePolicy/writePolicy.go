package main

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rag594/konfig-store/cache"
	"github.com/rag594/konfig-store/configRegister"
	"github.com/rag594/konfig-store/example"
	"github.com/rag594/konfig-store/konfigStore"
	"github.com/rag594/konfig-store/writePolicy"
	"time"
)

type SmartFeatureConfig struct {
	SmartSetting string
}

func main() {

	// db/redis client used for testing
	dbConn := example.GetDbConnection()
	defer dbConn.Close()
	redisConn := example.NewRedisNonClusteredClient()

	// Write

	kStore := konfigStore.New(
		konfigStore.WithDatabase(&konfigStore.Database{Connection: dbConn}),
		konfigStore.WithRedisCache(&konfigStore.RedisCache{
			Connection:         redisConn,
			ClusterModeEnabled: false,
		}))

	// register your new configuration
	smartFeatConfigRegister := configRegister.RegisterConfig[int, SmartFeatureConfig](
		kStore,
		configRegister.WithWritePolicy(writePolicy.WriteBack),
		configRegister.WithTTL(time.Minute),
	)

	// define your new cache key(it ius a function of entityId along with other options)
	cacheKeyForEntityC := cache.NewCacheKey[int, SmartFeatureConfig](20)

	val := &SmartFeatureConfig{SmartSetting: "vrwqrgrqe"}

	err := smartFeatConfigRegister.WritePolicy.SetConfig(context.Background(), cacheKeyForEntityC.DefaultValue(), 20, val)

	if err != nil {
		fmt.Println(err)
	}

	smartVal, err := smartFeatConfigRegister.ReadPolicy.GetConfig(context.Background(), cacheKeyForEntityC.DefaultValue(), 20)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(smartVal)
}
