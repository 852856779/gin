package services

import (
	"context"
	// "log"
	// "sync"
	// "time"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/bson"
	"testproject/global"
	"github.com/IBM/sarama"
	"fmt"
)


type redisHander struct {

}

var RedisHander = new(redisHander)




func (redisHander *redisHander) DelCache(message *sarama.ConsumerMessage){
	fmt.Println("code start")
	// fmt.Println(string(message.Value))
	userIdKey := string(message.Value)
	fmt.Println(string(message.Value))
	
	ctx := context.Background()
    redis := global.App.RedisClusterClient
	
	fmt.Println(redis)
	// userIdKey = "userId:2003"
	redis.Del(ctx, userIdKey)
	fmt.Println("code end")
}