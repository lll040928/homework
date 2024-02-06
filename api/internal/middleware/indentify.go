package middleware

import (
	"github.com/gin-gonic/gin"
)

func CheckRole1() func(c *gin.Context) {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "1" {
			c.JSON(200, gin.H{
				"code": 404,
				"msg":  "无权限访问",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
func CheckRole0() func(c *gin.Context) {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "0" {
			c.JSON(200, gin.H{
				"code": 404,
				"msg":  "无权访问",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
