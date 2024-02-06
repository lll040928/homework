package dao

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"homework/api/global"
	"homework/consts"
	"homework/model"
	"time"
)

func SelectWallet(username string) (*model.Wallet, error) {
	tx := global.MDB.Begin()
	if tx.Error != nil {
		global.Logger.Error("tx open failed")
		tx.Rollback()
		return nil, tx.Error
	}
	var w model.Wallet
	if err := tx.
		Where("username = ?", username).
		Find(&w).Error; err != nil {
		global.Logger.Error("mysql select failed" + err.Error())
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil { //提交事务并判断是否成功提交
		global.Logger.Error("tx close failed")
		tx.Rollback()
		return nil, err
	}
	return &w, nil
}
func UpdateWallet(wallet *model.Wallet) error {
	tx := global.MDB.Begin()
	if tx.Error != nil {
		global.Logger.Error("tx open failed")
		tx.Rollback()
		return tx.Error
	}
	if err := tx.
		Model(&model.Wallet{}).
		Where("username = ?", wallet.Username).
		Update(wallet).Error; err != nil {
		global.Logger.Error("mysql update failed" + err.Error())
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
func EmptyCart(carts []*model.Cart, username string, balance float64) error {
	tx := global.MDB.Begin()
	if tx.Error != nil {
		global.Logger.Error("tx open failed")
		tx.Rollback()
		return tx.Error
	}
	//创建订单并清空购物车
	for _, cart := range carts {
		order := model.Order{
			Username:  username,
			Gname:     cart.Gname,
			Price:     cart.Price,
			Count:     cart.Count,
			OrderTime: time.Now(),
		}
		if err := tx.
			Model(&model.Order{}).
			Create(&order).Error; err != nil {
			global.Logger.Error("order create failed" + err.Error())
			tx.Rollback()
			return err
		}
		if err := tx.
			Where("username = ? AND gid = ?", username, cart.Gid).
			Delete(&model.Cart{}).Error; err != nil {
			global.Logger.Error("cart delete failed" + err.Error())
			tx.Rollback()
			return err
		}
	}
	//更新钱包
	wallet := model.Wallet{
		Username: username,
		Balance:  balance,
	}
	if err := tx.
		Model(&model.Wallet{}).
		Where("username = ?", username).
		Update(wallet).Error; err != nil {
		global.Logger.Error("wallet update failed" + err.Error())
		tx.Rollback()
		return err
	}
	//提交事务并判断是否成功提交
	if err := tx.Commit().Error; err != nil {
		global.Logger.Error("tx close failed")
		tx.Rollback()
		return err
	}
	return nil
}
func GetOrder(username string) ([]*model.Order, error) {
	tx := global.MDB.Begin()
	if tx.Error != nil {
		global.Logger.Error("tx open failed")
		tx.Rollback()
		return nil, tx.Error
	}
	var orders []*model.Order
	if err := tx.
		Model(&model.Order{}).
		Where("username = ?", username).
		Find(&orders).Error; err != nil {
		global.Logger.Error("select order failed")
		tx.Rollback()
		return nil, tx.Error
	}
	if err := tx.Commit().Error; err != nil { //提交事务并判断是否成功提交
		global.Logger.Error("tx close failed")
		tx.Rollback()
		return nil, err
	}
	return orders, nil
}

//////////////////////////////////////////////////////////////////////

func SetWallet(c *gin.Context, wallet *model.Wallet) error {
	walletJSON, err := json.Marshal(wallet)
	if err != nil {
		global.Logger.Error("wallet marshal failed,err:" + err.Error())
		return err
	}
	_, err = global.RDB.Set(c, "wallet:"+wallet.Username, walletJSON, consts.RedisExpireDuration).Result()
	if err != nil {
		global.Logger.Error("wallet set failed,err:" + err.Error())
		return err
	}
	return nil
}
func GetWallet(c *gin.Context, username string) (*model.Wallet, error) {
	walletJSON, err := global.RDB.Get(c, "wallet:"+username).Result()
	if err != nil {
		global.Logger.Info("wallet get failed,err:" + err.Error())
		return nil, err
	}
	var wallet model.Wallet
	err = json.Unmarshal([]byte(walletJSON), &wallet)
	if err != nil {
		global.Logger.Error("wallet unmarshal failed,err:" + err.Error())
		return nil, err
	}
	return &wallet, nil
}
