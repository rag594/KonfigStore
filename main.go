package main

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/rag594/konfigStore/cache"
	"github.com/rag594/konfigStore/writePolicy"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type ComplexFeatureConfig struct {
	Enable string
}

type SmartFeatureConfig struct {
	SmartSetting string
}

func NewRedisNonClusteredClient() *redis.Client {
	// Initialize Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Adjust for your Redis setup
	})

	return rdb
}

func GetDbConnection() *sqlx.DB {
	// Replace with your database credentials
	dsn := "rabby:rabby123@tcp(localhost:3306)/configDb"

	// Open the database connection
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}

	// Test the database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v\n", err)
	}

	fmt.Println("Successfully connected to MySQL!")

	return db
}

func main() {

	// db/redis client used for testing
	dbConn := GetDbConnection()
	defer dbConn.Close()
	redisConn := NewRedisNonClusteredClient()

	/**
	Example on how to fetch the config
	*/

	// Read

	// register your new configuration
	complexFeatConfigRegister := RegisterConfig[int, ComplexFeatureConfig](
		WithSqlXDbConn(dbConn),
		WithRedisNCClient(redisConn),
		WithTTL(time.Minute),
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

	// Write

	// register your new configuration
	smartFeatConfigRegister := RegisterConfig[int, SmartFeatureConfig](
		WithSqlXDbConn(dbConn),
		WithRedisNCClient(redisConn),
		WithWritePolicy(writePolicy.WriteBack),
		WithTTL(time.Minute),
	)

	// define your new cache key(it ius a function of entityId along with other options)
	cacheKeyForEntityC := cache.NewCacheKey[int, SmartFeatureConfig](20)

	val := &SmartFeatureConfig{SmartSetting: "gsgsgsfhg"}

	err = smartFeatConfigRegister.WritePolicy.SetConfig(context.Background(), cacheKeyForEntityC.DefaultValue(), 20, val)

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
