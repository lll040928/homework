package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"time"
)

func LimitRate(limiter *rate.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 20*time.Second) //等待20秒
		defer cancel()
		if err := limiter.Wait(ctx); err != nil {
			c.JSON(429, gin.H{
				"code": 0,
				"msg":  "请求超时",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
