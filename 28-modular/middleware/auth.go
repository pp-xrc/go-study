package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"go-study/28-modular/utils"
)

const CtxUserID = "userID"
const CtxRole = "role"

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
		raw := strings.TrimSpace(strings.TrimPrefix(auth, prefix))
		claims, err := utils.VerifyToken(raw)
		if err != nil {
			utils.FailApp(c, utils.Unauthorized("token 无效或已过期"))
			c.Abort()
			return
		}
		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxRole, claims.Role)
		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString(CtxRole) != "admin" {
			utils.FailApp(c, utils.Forbidden("仅管理员可查看用户列表"))
			c.Abort()
			return
		}
		c.Next()
	}
}
