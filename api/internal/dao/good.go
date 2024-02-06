package dao

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"homework/api/global"
	"homework/model"
	"strconv"
)

func PutAllGoods() ([]*model.Good, error) {
	tx := global.MDB.Begin()
	if tx.Error != nil {
		global.Logger.Error("tx open failed")
		tx.Rollback()
		return nil, tx.Error
	}
	var goods []*model.Good
	ret := tx.Find(&goods)
	if ret.Error != nil {
		global.Logger.Error("select goods failed")
		tx.Rollback()
		return nil, ret.Error
	}
	if err := tx.Commit().Error; err != nil { //提交事务并判断是否成功提交
		global.Logger.Error("tx close failed")
		tx.Rollback()
		return nil, err
	}
	return goods, nil
}
func SearchGoods(keyword string) ([]*model.Good, error) {
	tx := global.MDB.Begin()
	if tx.Error != nil {
		global.Logger.Error("tx open failed")
		tx.Rollback()
		return nil, tx.Error
	}
	var goods []*model.Good
	ret := tx.Where("gname LIKE ?", "%"+keyword+"%").Find(&goods)
	if ret.Error != nil {
		global.Logger.Info("select goods failed")
		tx.Rollback()
		return nil, ret.Error
	}
	if err := tx.Commit().Error; err != nil { //提交事务并判断是否成功提交
		global.Logger.Error("tx close failed")
		tx.Rollback()
		return nil, err
	}
	return goods, nil
}
func SearchGoodsByOid(oid int) ([]*model.Good, error) {
	tx := global.MDB.Begin()
	if tx.Error != nil {
		global.Logger.Error("tx open failed")
		tx.Rollback()
		return nil, tx.Error
	}
	var goods []*model.Good
	ret := tx.Where("owner_id = ?", oid).Find(&goods)
	if ret.Error != nil {
		global.Logger.Info("select goods failed")
		tx.Rollback()
		return nil, ret.Error
	}
	if err := tx.Commit().Error; err != nil { //提交事务并判断是否成功提交
		global.Logger.Error("tx close failed")
		tx.Rollback()
		return nil, err
	}
	return goods, nil
}
func AddGood(good *model.Good) error {
	tx := global.MDB.Begin()
	if tx.Error != nil {
		global.Logger.Error("tx open failed")
		tx.Rollback()
		return tx.Error
	}
	err := tx.Create(&good).Error
	if err != nil {
		global.Logger.Error("mysql insert failed" + err.Error())
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
func SearchGoodsByGid(gid int) (*model.Good, error) {
	tx := global.MDB.Begin()
	if tx.Error != nil {
		global.Logger.Error("tx open failed")
		tx.Rollback()
		return nil, tx.Error
	}
	var good model.Good
	ret := tx.Where("gid = ?", gid).Find(&good)
	if ret.Error != nil {
		global.Logger.Info("select goods failed")
		tx.Rollback()
		return nil, ret.Error
	}
	if err := tx.Commit().Error; err != nil { //提交事务并判断是否成功提交
		global.Logger.Error("tx close failed")
		tx.Rollback()
		return nil, err
	}
	return &good, nil
}

/////////////////////////////////////////////////////////////

func SetGood(c *gin.Context, good *model.Good) error {
	goodJSON, err := json.Marshal(good)
	if err != nil {
		global.Logger.Error("good marshal failed,err:" + err.Error())
		return err
	}
	_, err = global.RDB.HSet(c, "goods", strconv.Itoa(good.Gid), goodJSON).Result()
	if err != nil {
		if err != nil {
			global.Logger.Error("good set failed,err:" + err.Error())
			return err
		}
	}
	return nil
}
func GetGood(c *gin.Context, gid int) (*model.Good, error) {
	goodJSON, err := global.RDB.HGet(c, "goods", strconv.Itoa(gid)).Result()
	if err != nil {
		global.Logger.Error("good get failed,err:" + err.Error())
		return nil, err
	}
	var good model.Good
	err = json.Unmarshal([]byte(goodJSON), &good)
	if err != nil {
		global.Logger.Error("good unmarshal failed,err:" + err.Error())
		return nil, err
	}
	return &good, nil
}
func GetAllGoods(c *gin.Context) ([]*model.Good, error) {
	goodsJSON, err := global.RDB.HGetAll(c, "goods").Result()
	if err != nil {
		global.Logger.Error("good get failed,err:" + err.Error())
		return nil, err
	}
	var goods []*model.Good
	for _, v := range goodsJSON {
		var good model.Good
		err = json.Unmarshal([]byte(v), &good)
		if err != nil {
			global.Logger.Error("good unmarshal failed,err:" + err.Error())
			return nil, err
		}
		goods = append(goods, &good)
	}
	return goods, nil
}
