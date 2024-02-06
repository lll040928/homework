package global

import (
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

var C Config

var Logger *zap.Logger

var MDB *gorm.DB

var RDB *redis.Client
