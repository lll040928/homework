package middleware

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"homework/api/global"
	"homework/consts"
	"net/http"
	"strings"
	"time"
)

var secret = []byte(global.C.TokenSecret)

type MyClaims struct {
	jwt.StandardClaims
}

func GenToken(username string, role string) (string, error) {
	c := MyClaims{

		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(consts.TokenExpireDuration).Unix(),
			Issuer:    role,
			Subject:   username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(secret)
}
func ParseToken(tokenStr string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"msg":  "请求头中auth为空",
			})
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)

		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"msg":  "请求头中auth格式有误",
			})
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := ParseToken(parts[1])

		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"msg":  "无效的Token",
			})
			c.Abort()
			return
		}
		// 将当前请求的username信息保存到请求的上下文c上

		c.Set("username", mc.Subject)
		c.Set("role", mc.Issuer)
		c.Next()
	}
}
