package main

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/rag594/konfigStore/readPolicy"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type ComplexFeatureConfig struct {
	Enable string
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

	complexFeatConfig := &ComplexFeatureConfig{Enable: "fewgvfwrEGRW"}

	dbConn := GetDbConnection()

	defer dbConn.Close()

	redisConn := NewRedisNonClusteredClient()

	complexFeatConfigRegister := RegisterConfig[int, ComplexFeatureConfig](
		complexFeatConfig,
		WithSqlXDbConn(dbConn),
		WithRedisNCClient(redisConn),
		WithTTL(time.Minute),
		WithReadPolicy(readPolicy.CacheAside),
	)

	complexFeatConfigRegister.Config.EntityId = 11

	x, err := complexFeatConfigRegister.ReadPolicy.GetConfig(context.Background(), complexFeatConfigRegister.Config.GetKey(), complexFeatConfigRegister.Config.EntityId)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(x)

}
