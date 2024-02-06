package service

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"homework/api/global"
	"homework/api/internal/dao"
	"homework/api/internal/middleware"
	"homework/model"
	"homework/utils"
	"strconv"
)

func Register(c *gin.Context) {
	var user model.UserInfo
	err := c.ShouldBind(&user)
	if err != nil {
		global.Logger.Warn("should bind failed")
		return
	}

	if user.Username == "" || user.Password == "" {
		c.JSON(500, gin.H{
			"code": 0,
			"msg":  "用户名或密码为空",
		})
		return
	}
	if len(user.Phone) != 11 {
		c.JSON(500, gin.H{
			"code": 0,
			"msg":  "手机号不合法",
		})
		return
	}

	_, err = dao.GetUserByUsername(user.Username)
	if err != gorm.ErrRecordNotFound {
		c.JSON(500, gin.H{
			"code": 0,
			"msg":  "用户已存在",
		})
		return
	}
	err = dao.CreateUser(&user)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 0,
			"msg":  "注册失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 1,
		"msg":  "注册成功",
	})
}

func LogIn(c *gin.Context) {
	var user model.UserInfo
	err := c.ShouldBind(&user)
	if err != nil {
		global.Logger.Warn("shouldBind failed:" + err.Error())
		return
	}
	if user.Username == "" || user.Password == "" {
		c.JSON(500, gin.H{
			"code":  0,
			"msg":   "用户名或密码为空",
			"token": "null",
		})
		return
	}
	var temp *model.UserInfo
	temp, err = dao.GetUser(c, user.Username)
	if err != nil {
		temp, err = dao.GetUserByUsername(user.Username)
		if err != nil {
			c.JSON(500, gin.H{
				"code": 0,
				"msg":  "用户不存在",
			})
			return
		}
		_ = dao.SetUser(c, temp)
	}

	password, err := utils.Decrypt(temp.Password)
	if err != nil {
		global.Logger.Warn("解密错误：" + err.Error())
		return
	}
	if password != user.Password {
		c.JSON(500, gin.H{
			"code":  0,
			"msg":   "密码错误",
			"token": "null",
		})
		return
	}

	role := strconv.Itoa(temp.Role)
	token, err := middleware.GenToken(user.Username, role)
	if err != nil {
		c.JSON(500, gin.H{
			"code":  0,
			"msg":   "无法生成token",
			"token": "null",
		})
		return
	}
	c.JSON(200, gin.H{
		"code":  1,
		"msg":   "登录成功",
		"token": token,
	})
}

func Forget(c *gin.Context) {
	newPass := c.PostForm("newpassword")
	var user model.UserInfo
	err := c.ShouldBind(&user)
	if err != nil {
		global.Logger.Warn("shouldBind failed:" + err.Error())
	}
	if user.Username == "" {
		c.JSON(500, gin.H{
			"code": 0,
			"msg":  "用户名或密码为空",
		})
		return
	}
	if len(user.Phone) != 11 {
		c.JSON(500, gin.H{
			"code": 0,
			"msg":  "手机号不合法",
		})
		return
	}
	usr, err := dao.GetUser(c, user.Username)
	if err != nil {
		usr, err = dao.GetUserByUsername(user.Username)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(500, gin.H{
					"code": 0,
					"msg":  "用户不存在",
				})
				return
			}
			c.JSON(500, gin.H{
				"code": 0,
				"msg":  "查找用户出错",
			})
			return
		}
		_ = dao.SetUser(c, usr)
	}

	phone, err := utils.Decrypt(usr.Phone)
	if err != nil {
		global.Logger.Warn("解密错误：" + err.Error())
		return
	}
	if phone != usr.Phone {
		c.JSON(500, gin.H{
			"code": 0,
			"msg":  "验证失败",
		})
		return
	}
	newPassStr, err := utils.Encrypt(newPass)
	if err != nil {
		global.Logger.Warn("加密错误：" + err.Error())
		return
	}
	temp := &model.UserInfo{
		Uid:      usr.Uid,
		Username: usr.Username,
		Password: newPassStr,
		Phone:    usr.Phone,
	}
	err = dao.UpdateUser(temp)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 0,
			"msg":  "修改失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 1,
		"msg":  "修改成功",
	})
}

func LogOut(c *gin.Context) {
	c.JSON(200, gin.H{
		"code": 1,
		"msg":  "clear jwt-token",
		"user": c.GetString("username"),
	})
}

func GetInfo(c *gin.Context) {

	username := c.GetString("username")
	user, err := dao.GetUser(c, username)
	if err != nil {
		user, err = dao.GetUserByUsername(username)
		if err != nil {
			c.JSON(500, gin.H{
				"code": 0,
				"msg":  "获取用户信息失败",
			})
			return
		}
		_ = dao.SetUser(c, user)
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
	role := "顾客"
	if user.Role != 0 {
		role = "店铺"
	}
	c.JSON(200, gin.H{
		"code": 1,
		"msg": gin.H{
			"Uid":      user.Uid,
			"Role":     role,
			"Username": user.Username,
			"Address":  user.Address,
			"Balance":  wallet.Balance,
		},
	})

}

func ChangeInfo(c *gin.Context) {
	address := c.PostForm("address")
	username := c.GetString("username")
	usr, err := dao.GetUser(c, username)
	if err != nil {
		usr, err = dao.GetUserByUsername(username)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(500, gin.H{
					"code": 0,
					"msg":  "用户不存在",
				})
				return
			}
			c.JSON(500, gin.H{
				"code": 0,
				"msg":  "查找用户出错",
			})
			return
		}
		_ = dao.SetUser(c, usr)
	}
	temp := &model.UserInfo{
		Uid:      usr.Uid,
		Username: usr.Username,
		Password: usr.Password,
		Phone:    usr.Phone,
		Address:  address,
	}
	err = dao.UpdateUser(temp)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "更改用户个人信息失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 1,
		"msg":  "更改用户个人信息成功",
	})
}
