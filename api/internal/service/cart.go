package service

import (
	"github.com/gin-gonic/gin"
	"homework/api/internal/dao"
	"homework/model"
	"strconv"
)

func GetCart(c *gin.Context) {
	username := c.GetString("username")
	carts, err := dao.SelectCartByUsername(username)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 0,
			"msg":  "未添加商品至购物车",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 1,
		"msg":  carts,
	})
}

func AddCart(c *gin.Context) {
	username := c.GetString("username")
	gidStr := c.PostForm("gid")
	countStr := c.PostForm("count")
	gid, _ := strconv.Atoi(gidStr)
	count, _ := strconv.Atoi(countStr)
	good, err := dao.GetGood(c, gid)
	if err != nil {
		good, err = dao.SearchGoodsByGid(gid)
		if err != nil {
			c.JSON(500, gin.H{
				"code": 0,
				"msg":  "获取商品信息失败",
			})
			return
		}
	}

	cart := model.Cart{
		Username: username,
		Gid:      gid,
		Gname:    good.Gname,
		Price:    good.Price,
		Count:    count,
	}
	err = dao.CreateCart(&cart)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 0,
			"msg":  "添加购物车失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 1,
		"msg":  "已添加至购物车",
	})
}

func DelCart(c *gin.Context) {
	username := c.GetString("username")
	gid, _ := strconv.Atoi(c.Query("gid"))
	err := dao.DeleteCart(username, gid)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 0,
			"msg":  "删除失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 1,
		"msg":  "删除成功",
	})
}
