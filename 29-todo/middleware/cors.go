package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS：浏览器走 Nginx :8090 同源反代时用不到；直接打 :8089（curl / 旧前端）时仍放行本机页面
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:8090")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
