package main

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rag594/konfigStore/cache"
	configDb "github.com/rag594/konfigStore/db"
	"time"
)

type UpiOptimalRouting struct {
	Enable string
}

func main() {

	upiConfig := &UpiOptimalRouting{Enable: "fewgvfwrEGRW"}

	dbConn := configDb.GetDbConnection()

	defer dbConn.Close()

	redisConn := cache.NewRedisNonClusteredClient()

	upiOptimalRoutingConfigRegister := RegisterConfig[int, UpiOptimalRouting](
		upiConfig,
		WithSqlXDbConn(dbConn),
		WithRedisNCClient(redisConn),
		WithTTL(time.Minute),
		WithReadPolicy(ReadThrough),
	)

	upiOptimalRoutingConfigRegister.Config.EntityId = 11

	//x, err := upiOptimalRoutingConfigRegister.ConfigDbOps.SaveConfig(context.Background(), upiOptimalRoutingConfigRegister.Config)
	//
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//fmt.Println(x)
	//
	//d, err := upiOptimalRoutingConfigRegister.ConfigDbOps.GetConfigByKeyForEntity(context.Background(), upiOptimalRoutingConfigRegister.Config.GetKey(), upiOptimalRoutingConfigRegister.Config.EntityId)
	//
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//fmt.Println(d)

	//err := upiOptimalRoutingConfigRegister.ConfigCacheOps.SaveConfig(context.Background(), upiOptimalRoutingConfigRegister.Config.GetKey(), &upiOptimalRoutingConfigRegister.Config.Value)
	//
	//if err != nil {
	//	fmt.Println("error in storing", err)
	//}
	//
	//x, err := upiOptimalRoutingConfigRegister.ConfigCacheOps.GetConfigByKeyForEntity(context.Background(), upiOptimalRoutingConfigRegister.Config.GetKey())
	//
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//fmt.Println(x)

	x, err := upiOptimalRoutingConfigRegister.ReadPolicy.GetConfig(context.Background(), upiOptimalRoutingConfigRegister.Config.GetKey(), upiOptimalRoutingConfigRegister.Config.EntityId)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(x)

}
