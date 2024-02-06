package service

import (
	"github.com/gin-gonic/gin"
	"homework/api/internal/dao"
	"homework/model"
)

func GetGoods(c *gin.Context) {
	var goods []*model.Good
	var err error
	goods, err = dao.GetAllGoods(c)
	if err != nil || goods == nil {
		goods, err = dao.PutAllGoods()
		if err != nil {
			c.JSON(500, gin.H{
				"code": "0",
				"msg":  "获取商品失败",
			})
			return
		}
		for _, v := range goods {
			_ = dao.SetGood(c, v)
		}
	}
	var good []model.Good
	for _, v := range goods {
		good = append(good, *v)
	}

	c.JSON(200, gin.H{
		"code": "1",
		"msg":  good,
	})
}

func SearchGoods(c *gin.Context) {
	keyword := c.Query("keyword")
	goods, err := dao.SearchGoods(keyword)
	if err != nil {
		c.JSON(500, gin.H{
			"code": "0",
			"msg":  "无法获取商品信息",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": "1",
		"msg":  goods,
	})
}
