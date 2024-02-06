package service

import (
	"github.com/gin-gonic/gin"
	"homework/api/internal/dao"
	"homework/model"
	"strconv"
)

func Buy(c *gin.Context) {
	username := c.GetString("username")
	user, err := dao.GetUser(c, username)
	if err != nil {
		user, err = dao.GetUserByUsername(username)
		if err != nil {
			c.JSON(500, gin.H{
				"code": 0,
				"msg":  "查找用户出错",
			})
			return
		}
		_ = dao.SetUser(c, user)
	}
	if user.Address == "未填写" {
		c.JSON(500, gin.H{
			"code": 0,
			"msg":  "未填写收货地址",
		})
		return
	}
	carts, err := dao.SelectCartByUsername(username)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 0,
			"msg":  "查询购物车失败",
		})
		return
	}
	var sum float64 = 0
	for _, cart := range carts {
		sum += cart.Price * float64(cart.Count)
	}
	wallet, err := dao.GetWallet(c, username)
	if err != nil {
		wallet, err = dao.SelectWallet(username)
		if err != nil {
			c.JSON(500, gin.H{
				"code": 0,
				"msg":  "查询钱包失败",
			})
			return
		}
		_ = dao.SetWallet(c, wallet)
	}
	if sum > wallet.Balance {
		c.JSON(500, gin.H{
			"code": 0,
			"msg":  "余额不足",
		})
		return
	}
	now := wallet.Balance - sum
	err = dao.EmptyCart(carts, username, now)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 0,
			"msg":  "购买失败",
		})
	}

	c.JSON(200, gin.H{
		"code": 1,
		"msg":  "购买成功",
	})
}

func Recharge(c *gin.Context) {
	moneyStr := c.PostForm("money")
	username := c.GetString("username")
	money, err := strconv.ParseFloat(moneyStr, 64)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 0,
			"msg":  "金额错误",
		})
		return
	}
	wallet, err := dao.GetWallet(c, username)
	if err != nil {
		wallet, err = dao.SelectWallet(username)
		if err != nil {
			c.JSON(500, gin.H{
				"code": 0,
				"msg":  "查询钱包失败",
			})
			return
		}
		_ = dao.SetWallet(c, wallet)
	}
	now := wallet.Balance + money
	temp := model.Wallet{
		Username: username,
		Balance:  now,
	}
	err = dao.UpdateWallet(&temp)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 0,
			"msg":  "充值失败",
		})
		return
	}
	_ = dao.SetWallet(c, &temp)
	c.JSON(200, gin.H{
		"code": 1,
		"msg":  "充值成功",
	})
}

func Order(c *gin.Context) {
	username := c.GetString("username")
	orders, err := dao.GetOrder(username)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 0,
			"msg":  "获取订单失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 1,
		"msg":  orders,
	})
}
