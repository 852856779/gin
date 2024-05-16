package global

import (
    "github.com/spf13/viper"
    "testproject/config"
    "gorm.io/gorm"
)

type Application struct {
    ConfigViper *viper.Viper
    Config config.Configuration
    DB *gorm.DB
}

var App = new(Application)