package router

import (
	"github.com/gin-gonic/gin"
	"homework/api/global"
	"homework/api/internal/middleware"
	"homework/api/internal/service"
)

func InitRouter() {
	r := gin.Default()

	r.POST("/register", service.Register)
	r.POST("/login", service.Login)
	r.POST("/logOut", service.LogOut)
	r.POST("/forget", service.Forget)
	s := r.Group("/store")
	{
		s.Use(middleware.JWT(), middleware.CheckRole1())
		s.POST("/push", service.PushGood)
		s.GET("/all", service.GetOwnGood)

	}
	u := r.Group("/user")
	{
		u.Use(middleware.JWT(), middleware.CheckRole0())
		u.GET("/info", service.GetInfo)
		u.POST("/info", service.ChangeInfo)
		u.POST("/recharge", service.Recharge)
		g := u.Group("/good")
		{
			g.GET("/all", service.GetGoods)
			g.GET("/search", service.SearchGoods)
			g.POST("/cart", service.AddCart)
			g.GET("/cart", service.GetCart)
			g.DELETE("/cart", service.DelCart)
			g.POST("/trade", service.Buy)
			g.GET("/order", service.Order)
		}
	}

	err := r.Run(":8080")
	if err != nil {
		global.Logger.Warn("router init failed")
		panic(err)
	}
}
