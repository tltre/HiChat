package global

import (
	"HiChat/config"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var DB *gorm.DB
var RedisDB *redis.Client
var ServiceConfig *config.ServiceConfig
