package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"go-study/29-todo/utils"
)

const CtxUserID = "userID"

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			utils.FailApp(c, utils.Unauthorized("缺少 Authorization 请求头"))
			c.Abort()
			return
		}
		const prefix = "Bearer "
		if !strings.HasPrefix(auth, prefix) {
			utils.FailApp(c, utils.Unauthorized("Authorization 格式应为 Bearer <token>"))
			c.Abort()
			return
		}
		claims, err := utils.VerifyToken(strings.TrimSpace(strings.TrimPrefix(auth, prefix)))
		if err != nil {
			utils.FailApp(c, utils.Unauthorized("token 无效或已过期"))
			c.Abort()
			return
		}
		c.Set(CtxUserID, claims.UserID)
		c.Next()
	}
}

func UID(c *gin.Context) uint {
	return uint(c.GetInt64(CtxUserID))
}
