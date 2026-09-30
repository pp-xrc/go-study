/*
第 27 课：Gin + GORM + JWT 整合 —— 用户模块实战

================================================================================
一、串联流程
================================================================================
  请求进来：Gin 路由
       ↓
  公开接口：注册 / 登录
  受保护接口：JWT 中间件先解析 user_id，再 GORM 查库

  注册：校验参数 → bcrypt 加密密码 → GORM 插入
  登录：校验参数 → GORM 查用户 → bcrypt 比对密码 → 签发 JWT
  鉴权接口：从 JWT 取出 user_id → GORM 查库 → 按角色做简单权限控制

================================================================================
二、接口
================================================================================
  POST /user/register     注册（公开）
  POST /user/login        登录（公开，返回 token）
  GET  /user/me           当前登录用户（需 JWT）
  GET  /user/list         用户列表（需 JWT，仅 admin）
  GET  /user/:id          用户详情（需 JWT；本人或 admin）

================================================================================
三、运行与测试
================================================================================
  go run ./27-Gin-GORM-JWT

  # 注册
  curl -s -X POST http://localhost:8087/user/register \
    -H "Content-Type: application/json" \
    -d '{"username":"zhangsan","password":"123456","nickname":"张三"}'

  # 登录（记下返回的 token）
  curl -s -X POST http://localhost:8087/user/login \
    -H "Content-Type: application/json" \
    -d '{"username":"admin","password":"admin123"}'

  # 带 token 访问
  TOKEN=这里粘贴token
  curl -s http://localhost:8087/user/me -H "Authorization: Bearer $TOKEN"
  curl -s http://localhost:8087/user/list -H "Authorization: Bearer $TOKEN"

  内置管理员：username=admin  password=admin123

================================================================================
*/

package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"go-study/response"
)

// db 全局数据库连接，handler 里直接用（教学 demo；正式项目建议依赖注入）
var db *gorm.DB

const (
	ctxUserIDKey = "userID"
	roleAdmin    = "admin"
	roleUser     = "user"
	tokenTTL     = 24 * time.Hour
)

// 签名密钥：正式项目放到环境变量，不要写死在代码里
var jwtSecret = []byte("go-study-day27-secret-change-me")

// ======================== 1. 模型 ========================

// User 映射 go_db.user 表
// 本课在原表基础上补充 username / nickname / role，AutoMigrate 会自动加列
type User struct {
	ID         uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"column:name;size:64;not null" json:"-"` // 旧列，非空；注册时与 nickname 同步
	Username   string    `gorm:"column:username;size:64;not null;default:''" json:"username"`
	Nickname   string    `gorm:"column:nickname;size:64;not null;default:''" json:"nickname"`
	Password   string    `gorm:"column:password;size:255;not null" json:"-"` // json:"-" 绝不返回给前端
	Role       string    `gorm:"column:role;size:16;not null;default:user" json:"role"`
	Age        uint      `gorm:"column:age;not null;default:0" json:"age"`
	Phone      string    `gorm:"column:phone;size:20;not null;default:''" json:"phone"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"create_time"`
}

func (User) TableName() string { return "user" }

// UserVO 对外展示（不含密码）
type UserVO struct {
	ID         uint      `json:"id"`
	Username   string    `json:"username"`
	Nickname   string    `json:"nickname"`
	Role       string    `json:"role"`
	Age        uint      `json:"age"`
	Phone      string    `json:"phone"`
	CreateTime time.Time `json:"create_time"`
}

func toVO(u User) UserVO {
	return UserVO{
		ID:         u.ID,
		Username:   u.Username,
		Nickname:   u.Nickname,
		Role:       u.Role,
		Age:        u.Age,
		Phone:      u.Phone,
		CreateTime: u.CreateTime,
	}
}

type RegisterReq struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	Nickname string `json:"nickname"`
}

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ======================== 2. bcrypt：加密 / 比对（不能解密） ========================
//
// bcrypt 是单向哈希：注册时加密入库，登录时拿明文和哈希比对。
// 课件说的「解密」在工程上其实是 Compare，哈希不可逆，无法还原明文。

// HashPassword 注册时调用：把明文密码变成哈希
func HashPassword(plain string) (string, error) {
	// DefaultCost=10，数值越大越慢越安全
	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// CheckPassword 登录时调用：明文 vs 库里的哈希，匹配返回 nil
func CheckPassword(hashed, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}

// ======================== 3. JWT：签发 / 校验 ========================

type CustomClaims struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"` // 放进 token，中间件可做粗粒度权限判断
	jwt.RegisteredClaims
}

// GenerateToken 登录成功后签发，有效期 24 小时
func GenerateToken(userID int64, role string) (string, error) {
	now := time.Now()
	claims := CustomClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "go-study-day27",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// VerifyToken 校验签名和过期时间，成功返回 claims
func VerifyToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意外的签名算法: %v", t.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("无效的 token")
	}
	if claims.UserID <= 0 {
		return nil, errors.New("token 中缺少有效 user_id")
	}
	return claims, nil
}

// ======================== 4. 中间件 ========================

// JWTAuthMiddleware 从 Authorization: Bearer <token> 解析出 user_id / role
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			unauthorized(c, "缺少 Authorization 请求头，请先登录")
			return
		}
		const prefix = "Bearer "
		if !strings.HasPrefix(auth, prefix) {
			unauthorized(c, "Authorization 格式错误，应为: Bearer <token>")
			return
		}
		raw := strings.TrimSpace(strings.TrimPrefix(auth, prefix))
		if raw == "" {
			unauthorized(c, "token 为空")
			return
		}

		claims, err := VerifyToken(raw)
		if err != nil {
			unauthorized(c, "token 无效或已过期: "+err.Error())
			return
		}

		c.Set(ctxUserIDKey, claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// AdminOnly 简单权限：仅 admin 可通过（挂在 /user/list 上）
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != roleAdmin {
			response.Fail(c, http.StatusForbidden, 403, "权限不足：仅管理员可查看用户列表")
			c.Abort()
			return
		}
		c.Next()
	}
}

func unauthorized(c *gin.Context, msg string) {
	response.Fail(c, http.StatusUnauthorized, 401, msg)
	c.Abort()
}

// ======================== 5. Handler ========================

// registerHandler 注册：校验 → 查重 → bcrypt → GORM Create
func registerHandler(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误：username 至少 3 位，password 至少 6 位")
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	req.Nickname = strings.TrimSpace(req.Nickname)
	if req.Nickname == "" {
		req.Nickname = req.Username
	}

	var exist User
	err := db.Where("username = ?", req.Username).First(&exist).Error
	if err == nil {
		response.Fail(c, http.StatusConflict, 409, "用户名已存在")
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		response.ServerError(c, "查询用户失败")
		return
	}

	hashed, err := HashPassword(req.Password)
	if err != nil {
		response.ServerError(c, "密码加密失败")
		return
	}

	u := User{
		Name:     req.Nickname, // 旧表 name 非空，与昵称保持一致
		Username: req.Username,
		Nickname: req.Nickname,
		Password: hashed, // 存哈希，不存明文
		Role:     roleUser,
	}
	if err := db.Create(&u).Error; err != nil {
		response.ServerError(c, "创建用户失败: "+err.Error())
		return
	}

	response.Created(c, "注册成功", gin.H{"user": toVO(u)})
}

// loginHandler 登录：校验 → GORM 查用户 → bcrypt 比对 → 签发 JWT
func loginHandler(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请提供 username 和 password")
		return
	}

	var u User
	err := db.Where("username = ?", strings.TrimSpace(req.Username)).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusUnauthorized, 401, "用户名或密码错误")
		return
	}
	if err != nil {
		response.ServerError(c, "查询用户失败")
		return
	}

	// 比对失败不要区分「用户不存在 / 密码错」，避免被枚举用户名
	if err := CheckPassword(u.Password, req.Password); err != nil {
		response.Fail(c, http.StatusUnauthorized, 401, "用户名或密码错误")
		return
	}

	token, err := GenerateToken(int64(u.ID), u.Role)
	if err != nil {
		response.ServerError(c, "生成 token 失败")
		return
	}

	response.OK(c, "登录成功", gin.H{
		"token":      token,
		"token_type": "Bearer",
		"expires_in": int(tokenTTL.Seconds()),
		"user":       toVO(u),
	})
}

// meHandler 从 JWT 取出 user_id，再 GORM 查当前用户
func meHandler(c *gin.Context) {
	uid := uint(c.GetInt64(ctxUserIDKey))
	u, err := findUserByID(uid)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}
	response.OK(c, "获取当前用户成功", gin.H{"user": toVO(u)})
}

// listHandler 用户列表（仅 admin，由 AdminOnly 中间件拦截）
func listHandler(c *gin.Context) {
	var list []User
	if err := db.Where("username <> ?", "").Order("id asc").Find(&list).Error; err != nil {
		response.ServerError(c, "查询列表失败")
		return
	}
	vos := make([]UserVO, 0, len(list))
	for _, u := range list {
		vos = append(vos, toVO(u))
	}
	response.OK(c, "获取列表成功", gin.H{
		"current_user_id": c.GetInt64(ctxUserIDKey),
		"count":           len(vos),
		"list":            vos,
	})
}

// detailHandler 用户详情：admin 可看任何人；普通用户只能看自己
func detailHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.BadRequest(c, "id 必须是正整数")
		return
	}

	currentID := uint(c.GetInt64(ctxUserIDKey))
	role := c.GetString("role")
	if role != roleAdmin && currentID != uint(id) {
		response.Fail(c, http.StatusForbidden, 403, "权限不足：只能查看自己的信息")
		return
	}

	u, err := findUserByID(uint(id))
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}
	response.OK(c, "获取详情成功", gin.H{"user": toVO(u)})
}

func findUserByID(id uint) (User, error) {
	var u User
	err := db.First(&u, id).Error
	return u, err
}

// ======================== 6. 启动 ========================

func main() {
	db = mustConnect()
	if err := db.AutoMigrate(&User{}); err != nil {
		log.Fatalf("迁移表结构失败: %v", err)
	}
	seedAdmin()

	r := gin.Default()

	user := r.Group("/user")
	{
		// 公开：注册、登录
		user.POST("/register", registerHandler)
		user.POST("/login", loginHandler)

		// 需登录
		auth := user.Group("/")
		auth.Use(JWTAuthMiddleware())
		{
			auth.GET("/me", meHandler)
			// 静态路径 /list 必须写在 /:id 前面，否则 "list" 会被当成 id
			auth.GET("/list", AdminOnly(), listHandler)
			auth.GET("/:id", detailHandler)
		}
	}

	addr := ":8087"
	fmt.Println("Day27 Gin+GORM+JWT 启动: http://localhost" + addr)
	fmt.Println("管理员账号: admin / admin123")
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}

func mustConnect() *gorm.DB {
	dsn := "root:@tcp(127.0.0.1:3306)/go_db?charset=utf8mb4&parseTime=True&loc=Local"
	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("GORM 连接失败: %v", err)
	}
	sqlDB, err := conn.DB()
	if err != nil {
		log.Fatal(err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Ping 失败: %v（请确认 MySQL 已启动、go_db 存在）", err)
	}
	fmt.Println("✅ GORM 已连接 go_db")
	return conn
}

// seedAdmin 写入内置管理员，方便测试 /user/list
func seedAdmin() {
	var count int64
	db.Model(&User{}).Where("username = ?", "admin").Count(&count)
	if count > 0 {
		return
	}
	hashed, err := HashPassword("admin123")
	if err != nil {
		log.Fatalf("管理员密码加密失败: %v", err)
	}
	admin := User{
		Name:     "管理员",
		Username: "admin",
		Nickname: "管理员",
		Password: hashed,
		Role:     roleAdmin,
	}
	if err := db.Create(&admin).Error; err != nil {
		log.Fatalf("创建管理员失败: %v", err)
	}
	fmt.Println("📦 已创建管理员 admin / admin123")
}
