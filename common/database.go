package common

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"homework/api/global"
	"homework/model"
)

func InitDatabase() {
	InitMySQL()
	InitRedis()
}
func InitMySQL() {
	config := global.C.MysqlInfo
	s := "%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf(s, config.User, config.Password, config.Host, config.DBName)
	var err error
	global.MDB, err = gorm.Open("mysql", dsn)
	if err != nil {
		global.Logger.Error("MySQL initial failed:" + err.Error())
		return
	}
	global.MDB.AutoMigrate(&model.Good{})
	global.MDB.AutoMigrate(&model.Wallet{})
	global.MDB.AutoMigrate(&model.Cart{})
	global.MDB.AutoMigrate(&model.UserInfo{})
	global.MDB.AutoMigrate(&model.Order{})
	global.Logger.Info("MySQL initial success")
}

func InitRedis() {
	config := global.C.RedisInfo
	global.RDB = redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       0,
	})
	global.Logger.Info("Redis initial success")
}
