package bootstrap

import (
    "context"
    "github.com/redis/go-redis/v9"
    "go.uber.org/zap"
    "testproject/global"
	"fmt"
	"sync"
	"time"
)
var redisClient *redis.Client
func InitializeRedis() *redis.Client {
	// redisCluster
	var sync sync.Once
	sync.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr:     global.App.Config.Redis.Host + ":" + global.App.Config.Redis.Port,
			Password: global.App.Config.Redis.Password, // no password set
			DB:       global.App.Config.Redis.DB,       // use default DB
		})
		redisClient = client
	})
	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println(zap.Any("err", err));
		// global.App.Log.Error("Redis connect ping failed, err:", zap.Any("err", err))
		return nil
	}
	fmt.Println("redis connect sucess");
    return redisClient

}
var redisCluster *redis.ClusterClient
var redisClusterUrl []string = []string{
	"192.168.31.110:6380",
	"192.168.31.110:6381",
	"192.168.31.110:6382",
	"192.168.31.110:6383",
	"192.168.31.110:6384",
	"192.168.31.110:6385",
}
//连接redis集群
func InitializeRedisCluster() *redis.ClusterClient {
	// redisCluster
	var sync sync.Once
	sync.Do(func() {
		cluster := redis.NewClusterClient(&redis.ClusterOptions{
			// Addr:     global.App.Config.Redis.Host + ":" + global.App.Config.Redis.Port,
			Addrs:     redisClusterUrl,
			Password: global.App.Config.Redis.Password, // no password set
			RouteRandomly: true, 
			DialTimeout:  30 * time.Second,   //连接超时时间
			ReadTimeout:  30 * time.Second,   //读取超时时间
			WriteTimeout: 30 * time.Second,   //写入超时时间
		})
		redisCluster = cluster
	})
	_, err := redisCluster.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println(err);
		// global.App.Log.Error("Redis connect ping failed, err:", zap.Any("err", err))
		return nil
	}
	fmt.Println("redis connect cluster sucess");
    return redisCluster

}