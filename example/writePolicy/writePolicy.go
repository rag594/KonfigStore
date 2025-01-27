package main

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rag594/konfigStore/cache"
	"github.com/rag594/konfigStore/configRegister"
	"github.com/rag594/konfigStore/example"
	"github.com/rag594/konfigStore/writePolicy"
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

	// register your new configuration
	smartFeatConfigRegister := configRegister.RegisterConfig[int, SmartFeatureConfig](
		configRegister.WithSqlXDbConn(dbConn),
		configRegister.WithRedisNCClient(redisConn),
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

	time.Sleep(time.Minute * 2)
}
