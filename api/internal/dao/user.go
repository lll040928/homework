package dao

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"homework/api/global"
	"homework/consts"
	"homework/model"
	"homework/utils"
)

func CreateUser(user *model.UserInfo) error {
	tx := global.MDB.Begin()
	if tx.Error != nil {
		global.Logger.Error("tx open failed")
		tx.Rollback()
		return tx.Error
	}
	var temp model.UserInfo
	var err error
	password, err := utils.Encrypt(user.Password)
	if err != nil {
		global.Logger.Info("crypt failed:" + err.Error())
		return err
	}
	phone, err := utils.Encrypt(user.Phone)
	if err != nil {
		global.Logger.Info("crypt failed:" + err.Error())
		return err
	}
	err = tx.Where("username = ?", user.Username).First(&temp).Error
	if err != gorm.ErrRecordNotFound && err != nil {
		tx.Rollback()
		return errors.New(consts.MySQLExist)
	}
	temp = model.UserInfo{
		Username: user.Username,
		Password: password,
		Phone:    phone,
		Role:     user.Role,
	}
	err = tx.Create(&temp).Error
	if err != nil {
		global.Logger.Error("mysql insert failed" + err.Error())
		tx.Rollback()
		return err
	}
	w := model.Wallet{
		Username: user.Username,
	}
	if err = tx.Create(&w).Error; err != nil {
		global.Logger.Error("mysql insert failed" + err.Error())
		tx.Rollback()
		return err
	}

	if err = tx.Commit().Error; err != nil { //提交事务并判断是否成功提交
		global.Logger.Error("tx close failed")
		tx.Rollback()
		return err
	}
	return nil
}
func GetUserByUsername(username string) (*model.UserInfo, error) {
	tx := global.MDB.Begin() //开启事务

	if tx.Error != nil { //检查事务是否正常开启
		global.Logger.Error("tx open failed")
		tx.Rollback()
		return nil, tx.Error
	}
	var user model.UserInfo
	if err := tx.Where("username = ?", username).First(&user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		global.Logger.Error("tx close failed")
		tx.Rollback()
		return nil, err
	}
	return &user, nil
}
func UpdateUser(user *model.UserInfo) error {
	tx := global.MDB.Begin()
	if tx.Error != nil {
		global.Logger.Error("tx open failed")
		tx.Rollback()
		return tx.Error
	}
	if err := tx.
		Model(&model.UserInfo{}).
		Updates(user).
		Where("username = ?", user.Username).Error; err != nil {
		global.Logger.Warn("update user failed")
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil { //提交事务并判断是否成功提交
		global.Logger.Error("tx close failed")
		tx.Rollback()
		return err
	}
	return nil
}
func SetUser(c *gin.Context, user *model.UserInfo) error {
	userJSON, err := json.Marshal(user)
	if err != nil {
		global.Logger.Error("user_info marshal failed,err:" + err.Error())
		return err
	}
	_, err = global.RDB.Set(c, "user:"+user.Username, userJSON, consts.RedisExpireDuration).Result()
	if err != nil {
		global.Logger.Error("user_info set failed,err:" + err.Error())
		return err
	}
	return nil
}
func GetUser(c *gin.Context, username string) (*model.UserInfo, error) {
	userJSON, err := global.RDB.Get(c, "user:"+username).Result()
	if err != nil {
		global.Logger.Info("user_info get failed,err:" + err.Error())
		return nil, err
	}
	var user model.UserInfo
	err = json.Unmarshal([]byte(userJSON), &user)
	if err != nil {
		global.Logger.Error("user_info unmarshal failed,err:" + err.Error())
		return nil, err
	}
	return &user, nil
}
