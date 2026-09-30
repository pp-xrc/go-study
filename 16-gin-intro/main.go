/*
第 16 课：Gin Web 框架入门

运行：
  go run ./16-gin-intro
  访问 http://localhost:8080
*/

package main

import (
	"github.com/gin-gonic/gin"

	"go-study/response"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		response.OK(c, "Hello, Gin!", gin.H{
			"module": "go-study",
		})
	})

	r.GET("/hello", func(c *gin.Context) {
		response.OK(c, "Hello, Gin!", nil)
	})

	r.GET("/ping", func(c *gin.Context) {
		response.OK(c, "pong", gin.H{"ping": "pong"})
	})

	r.Run(":8080")
}
