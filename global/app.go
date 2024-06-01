package global

import (
    "github.com/spf13/viper"
    "testproject/config"
    "gorm.io/gorm"
    "github.com/redis/go-redis/v9"
)

type Application struct {
    ConfigViper *viper.Viper
    Config config.Configuration
    DB *gorm.DB
    Redis *redis.Client
    RedisClusterClient *redis.ClusterClient
}

var App = new(Application)