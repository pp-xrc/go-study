/*
第 26 课：JWT 认证原理 + Gin 集成

================================================================================
一、学习内容
================================================================================
  - JWT 结构：Header.Payload.Signature（无状态，服务端不存 token）
  - 使用 github.com/golang-jwt/jwt/v5 签发 / 校验
  - 过期时间建议 24 小时
  - 前端通常把 token 存 localStorage / cookie，请求头带 Authorization

================================================================================
二、实操任务
================================================================================
  1. GenerateToken(userID) / VerifyToken(token) → userID
  2. JWT 中间件：从 Authorization: Bearer <token> 取 token 并校验
  3. 给 /user/list、/user/:id 等需登录路由挂中间件

================================================================================
三、运行与测试
================================================================================
  go run ./26-JWT

  # 1) 登录拿 token（演示用，传 user_id 即可）
  curl -s -X POST http://localhost:8086/login \
    -H "Content-Type: application/json" \
    -d '{"user_id":1}'

  # 2) 无 token 访问受保护接口 → 401
  curl -i http://localhost:8086/user/list

  # 3) 带 token 访问
  TOKEN=这里粘贴上一步返回的token
  curl -s http://localhost:8086/user/list \
    -H "Authorization: Bearer $TOKEN"
  curl -s http://localhost:8086/user/1 \
    -H "Authorization: Bearer $TOKEN"

================================================================================
*/

package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"go-study/response"
)

// ======================== 1. JWT 工具 ========================

// 签名密钥（正式项目放配置 / 环境变量，不要写死提交仓库）
var jwtSecret = []byte("go-study-day26-secret-change-me")

// token 有效期：24 小时
const tokenTTL = 24 * time.Hour

// CustomClaims 自定义载荷：业务字段 + 标准字段（过期时间等）
type CustomClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken 根据用户 ID 生成 JWT 字符串
func GenerateToken(userID int64) (string, error) {
	now := time.Now()
	claims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenTTL)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(now),               // 签发时间
			NotBefore: jwt.NewNumericDate(now),               // 生效时间
			Issuer:    "go-study",                            // 签发者（可选）
		},
	}

	// Header 里算法 HS256；Payload 是 claims；用密钥生成 Signature
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// VerifyToken 校验 token，成功返回 userID
func VerifyToken(tokenString string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (any, error) {
		// 防止「算法混淆」攻击：只接受 HMAC
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意外的签名算法: %v", t.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return 0, errors.New("无效的 token")
	}
	if claims.UserID <= 0 {
		return 0, errors.New("token 中缺少有效 user_id")
	}
	return claims.UserID, nil
}

// ======================== 2. JWT 中间件 ========================

const ctxUserIDKey = "userID"

// JWTAuthMiddleware 从 Authorization: Bearer <token> 取 token 并校验
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			unauthorized(c, "缺少 Authorization 请求头")
			return
		}

		// 期望格式：Bearer xxxxx
		const prefix = "Bearer "
		if !strings.HasPrefix(auth, prefix) {
			unauthorized(c, "Authorization 格式错误，应为: Bearer <token>")
			return
		}
		tokenString := strings.TrimSpace(strings.TrimPrefix(auth, prefix))
		if tokenString == "" {
			unauthorized(c, "token 为空")
			return
		}

		userID, err := VerifyToken(tokenString)
		if err != nil {
			unauthorized(c, "token 无效或已过期: "+err.Error())
			return
		}

		// 校验通过：把 userID 放进上下文，后续 handler 可取用
		c.Set(ctxUserIDKey, userID)
		c.Next()
	}
}

func unauthorized(c *gin.Context, msg string) {
	response.Fail(c, http.StatusUnauthorized, 401, msg)
	c.Abort() // 中断后续 handler
}

// ======================== 3. 演示数据（内存） ========================

type UserVO struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

var demoUsers = []UserVO{
	{ID: 1, Username: "zhangsan", Nickname: "张三"},
	{ID: 2, Username: "lisi", Nickname: "李四"},
	{ID: 3, Username: "wangwu", Nickname: "王五"},
}

// ======================== 4. 路由与启动 ========================

func main() {
	r := gin.Default()

	// 公开接口：登录换 token（教学简化：只传 user_id）
	r.POST("/login", loginHandler)

	// 公开接口：健康检查（不需要登录）
	r.GET("/ping", func(c *gin.Context) {
		response.OK(c, "pong", gin.H{"auth": "JWT Day26"})
	})

	// 需登录的路由组：挂 JWT 中间件
	auth := r.Group("/")
	auth.Use(JWTAuthMiddleware())
	{
		auth.GET("/user/list", userListHandler) // 用户列表
		auth.GET("/user/:id", userDetailHandler) // 用户详情
	}

	addr := ":8086"
	fmt.Println("Day26 JWT demo 启动: http://localhost" + addr)
	fmt.Println("先 POST /login 拿 token，再带 Authorization: Bearer <token> 访问 /user/*")
	if err := r.Run(addr); err != nil {
		panic(err)
	}
}

type loginReq struct {
	UserID int64 `json:"user_id" binding:"required"`
}

func loginHandler(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil || req.UserID <= 0 {
		response.BadRequest(c, "请提供有效的 user_id")
		return
	}

	token, err := GenerateToken(req.UserID)
	if err != nil {
		response.ServerError(c, "生成 token 失败: "+err.Error())
		return
	}

	response.OK(c, "登录成功", gin.H{
		"token":      token,
		"token_type": "Bearer",
		"expires_in": int(tokenTTL.Seconds()), // 秒
		"user_id":    req.UserID,
	})
}

func userListHandler(c *gin.Context) {
	uid := c.GetInt64(ctxUserIDKey) // 中间件写入的当前登录用户
	response.OK(c, "获取列表成功", gin.H{
		"current_user_id": uid,
		"list":            demoUsers,
	})
}

func userDetailHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(c, "id 必须是正整数")
		return
	}

	for _, u := range demoUsers {
		if u.ID == id {
			response.OK(c, "获取详情成功", gin.H{
				"current_user_id": c.GetInt64(ctxUserIDKey),
				"user":            u,
			})
			return
		}
	}
	response.NotFound(c, "用户不存在")
}
