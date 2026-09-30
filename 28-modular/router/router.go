package router

import (
	"github.com/gin-gonic/gin"

	"go-study/28-modular/controller"
	"go-study/28-modular/middleware"
)

// Setup 只注册路由：公开接口 vs 需登录接口
func Setup() *gin.Engine {
	r := gin.Default()

	user := r.Group("/user")
	{
		user.POST("/register", controller.Register)
		user.POST("/login", controller.Login)

		auth := user.Group("/")
		auth.Use(middleware.JWTAuth())
		{
			auth.GET("/me", controller.Me)
			auth.GET("/list", middleware.AdminOnly(), controller.List)
			auth.GET("/:id", controller.Detail)
		}
	}
	return r
}
