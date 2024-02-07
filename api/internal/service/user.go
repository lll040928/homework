package service

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"homework/api/global"
	"homework/api/internal/dao"
	"homework/api/internal/middleware"
	"homework/model"
	"homework/utils"
	"net/http"
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"statues": 500,
			"message": "username or password are empty",
		})
		return
	}
	if len(user.Phone) != 11 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statues": 500,
			"message": "numbers are failed",
		})
		return
	}

	_, err = dao.GetUserByUsername(user.Username)
	if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statues": 500,
			"message": "user already exists",
		})
		return
	}
	err = dao.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statues": 500,
			"message": "register failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"statues": 200,
		"message": "add user successful",
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
		c.JSON(http.StatusOK, gin.H{
			"statues": 0,
			"message": "username or password are empty",
			"token":   "null",
		})
		return
	}
	var temp *model.UserInfo
	temp, err = dao.GetUser(c, user.Username)
	if err != nil {
		temp, err = dao.GetUserByUsername(user.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 0,
				"msg":  "user doesn't exit",
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"statues": 500,
			"message": "password wrong",
			"token":   "null",
		})
		return
	}

	role := strconv.Itoa(temp.Role)
	token, err := middleware.GenToken(user.Username, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statues": 500,
			"message": "not token",
			"token":   "null",
		})
		return
	}
	c.JSON(200, gin.H{
		"statues": 200,
		"message": "login successful",
		"token":   token,
	})
	//正确登录成功，设置cookie
	c.SetCookie("gin_demo_cookie", "test", 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "login successful",
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "username or password are empty",
		})
		return
	}
	if len(user.Phone) != 11 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "numbers are failed",
		})
		return
	}
	usr, err := dao.GetUser(c, user.Username)
	if err != nil {
		usr, err = dao.GetUserByUsername(user.Username)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusInternalServerError, gin.H{
					"statues": 500,
					"message": "user doesn't exist",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"statues": 500,
				"message": "search user is failed",
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"statues": 500,
			"message": "verification failed",
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"statues": 500,
			"message": "fix is failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"statues": 200,
		"message": "modification successful",
	})
}

func LogOut(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"statues": 200,
		"message": "clear jwt-token",
		"user":    c.GetString("username"),
	})
}

func GetInfo(c *gin.Context) {

	username := c.GetString("username")
	user, err := dao.GetUser(c, username)
	if err != nil {
		user, err = dao.GetUserByUsername(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"statues": 500,
				"message": "failed to retrieve user information",
			})
			return
		}
		_ = dao.SetUser(c, user)
	}
	wallet, err := dao.GetWallet(c, username)
	if err != nil {
		wallet, err = dao.SelectWallet(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"statues": 500,
				"message": "failed to query the wallet",
			})
			return
		}
		_ = dao.SetWallet(c, wallet)
	}
	role := "customer"
	if user.Role != 0 {
		role = "store"
	}
	c.JSON(http.StatusOK, gin.H{
		"statues": 200,
		"message": gin.H{
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
				c.JSON(http.StatusInternalServerError, gin.H{
					"statues": 500,
					"message": "user not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"statues": 500,
				"message": "failed to find user",
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"statues": 500,
			"message": "failed to update user personal information",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"statues": 200,
		"message": "user's personal information updated successfully",
	})
}
