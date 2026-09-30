package router

import (
	"github.com/gin-gonic/gin"

	"go-study/29-todo/controller"
	"go-study/29-todo/middleware"
)

func Setup() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())

	user := r.Group("/user")
	{
		user.POST("/register", controller.Register)
		user.POST("/login", controller.Login)
		user.GET("/me", middleware.JWTAuth(), controller.Me)
	}

	todo := r.Group("/todo")
	todo.Use(middleware.JWTAuth())
	{
		todo.POST("", controller.Create)
		todo.GET("/list", controller.List)
		todo.GET("/:id", controller.Detail)
		todo.PUT("/:id", controller.Update)
		todo.DELETE("/:id", controller.Delete)
	}
	return r
}
