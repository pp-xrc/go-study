/*
第 22 课：Gin 进阶 —— 用户模块综合练习（无数据库）

================================================================================
一、学习内容（对应截图）
================================================================================
  - 复习：路由、参数、响应、中间件
  - 综合练习：用户模块（注册 / 登录 / 列表 / 详情 / 更新）
  - 数据存内存 map，重启后清空（无数据库）
  - 统一返回格式：go-study/response
  - 全部接口挂 CORS + 日志中间件
  - 可用 Postman / curl 测正常与异常请求

================================================================================
二、接口一览（RESTful）
================================================================================
  POST /user/register     注册
  POST /user/login        登录
  GET  /user/list         用户列表
  GET  /user/:id          用户详情
  PUT  /user/:id          更新用户

================================================================================
三、运行与测试
================================================================================
  go run ./22-gin-advance

  # 注册
  curl -X POST http://localhost:8085/user/register \
    -H "Content-Type: application/json" \
    -d '{"username":"zhangsan","password":"123456","nickname":"张三"}'

  # 登录
  curl -X POST http://localhost:8085/user/login \
    -H "Content-Type: application/json" \
    -d '{"username":"zhangsan","password":"123456"}'

  # 列表 / 详情 / 更新
  curl http://localhost:8085/user/list
  curl http://localhost:8085/user/1
  curl -X PUT http://localhost:8085/user/1 \
    -H "Content-Type: application/json" \
    -d '{"nickname":"张三丰","age":20}'

  # 异常：缺参数、用户不存在、用户名重复、密码错误
  curl -X POST http://localhost:8085/user/register -H "Content-Type: application/json" -d '{}'
  curl http://localhost:8085/user/999
  curl -X POST http://localhost:8085/user/login \
    -H "Content-Type: application/json" \
    -d '{"username":"zhangsan","password":"wrong"}'
*/

package main

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"go-study/response"
)

// ======================== 1. 数据模型 ========================

// User 用户实体（内存存储）
// Password 仅内部使用，返回给前端时用 UserVO 脱敏
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"` // json:"-" 序列化时忽略，防止密码泄露
	Nickname string `json:"nickname"`
	Age      int    `json:"age"`
}

// UserVO 对外展示的用户视图（不含密码）
type UserVO struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Age      int    `json:"age"`
}

func toVO(u User) UserVO {
	return UserVO{
		ID:       u.ID,
		Username: u.Username,
		Nickname: u.Nickname,
		Age:      u.Age,
	}
}

// 请求体结构（绑定 JSON）
type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateReq struct {
	Nickname string `json:"nickname"`
	Age      int    `json:"age"`
}

// ======================== 2. 内存仓库（模拟数据库）========================

type UserStore struct {
	mu      sync.RWMutex
	users   map[int]User
	byName  map[string]int // username → id，便于登录查找
	nextID  int
}

func NewUserStore() *UserStore {
	return &UserStore{
		users:  make(map[int]User),
		byName: make(map[string]int),
		nextID: 1,
	}
}

func (s *UserStore) Register(username, password, nickname string) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.byName[username]; ok {
		return User{}, fmt.Errorf("用户名已存在")
	}

	u := User{
		ID:       s.nextID,
		Username: username,
		Password: password, // 教学 demo 明文；正式项目必须哈希（bcrypt）
		Nickname: nickname,
		Age:      0,
	}
	s.users[u.ID] = u
	s.byName[username] = u.ID
	s.nextID++
	return u, nil
}

func (s *UserStore) Login(username, password string) (User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	id, ok := s.byName[username]
	if !ok {
		return User{}, fmt.Errorf("用户不存在")
	}
	u := s.users[id]
	if u.Password != password {
		return User{}, fmt.Errorf("密码错误")
	}
	return u, nil
}

func (s *UserStore) List() []UserVO {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := make([]UserVO, 0, len(s.users))
	for _, u := range s.users {
		list = append(list, toVO(u))
	}
	return list
}

func (s *UserStore) Get(id int) (UserVO, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	u, ok := s.users[id]
	if !ok {
		return UserVO{}, fmt.Errorf("用户不存在")
	}
	return toVO(u), nil
}

func (s *UserStore) Update(id int, nickname string, age int) (UserVO, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	u, ok := s.users[id]
	if !ok {
		return UserVO{}, fmt.Errorf("用户不存在")
	}
	if nickname != "" {
		u.Nickname = nickname
	}
	if age > 0 {
		u.Age = age
	}
	s.users[id] = u
	return toVO(u), nil
}

// ======================== 3. 中间件 ========================

// RequestLogger 全局日志：方法、路径、时间、状态码、耗时
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		fmt.Printf("[日志] %s | %s %s | status=%d | cost=%v\n",
			start.Format("2006-01-02 15:04:05"),
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			time.Since(start),
		)
	}
}

// CORS 全局跨域：允许所有源、方法、请求头（方便前端联调）
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "*")
		c.Header("Access-Control-Allow-Headers", "*")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// ======================== 4. 启动辅助 ========================

func runServer(r *gin.Engine, addr string) {
	fmt.Println("⏳ 正在启动服务...")
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("❌ 启动失败:", err)
		return
	}
	fmt.Println("✅ 用户模块服务已就绪!")
	fmt.Printf("   http://localhost%s\n", addr)
	fmt.Println("   POST /user/register  注册")
	fmt.Println("   POST /user/login     登录")
	fmt.Println("   GET  /user/list      列表")
	fmt.Println("   GET  /user/:id       详情")
	fmt.Println("   PUT  /user/:id       更新")
	fmt.Println("   按 Ctrl+C 停止")
	fmt.Println(strings.Repeat("-", 40))
	_ = http.Serve(listener, r)
}

// ======================== 5. 主程序 & 路由 ========================

func main() {
	store := NewUserStore()

	r := gin.New()
	r.Use(gin.Recovery())
	// 题目要求：全部接口挂 CORS + 日志
	r.Use(CORS(), RequestLogger())

	user := r.Group("/user")
	{
		// ---------- POST /user/register 注册 ----------
		user.POST("/register", func(c *gin.Context) {
			var req RegisterReq
			if err := c.ShouldBindJSON(&req); err != nil {
				response.BadRequest(c, "JSON 格式错误: "+err.Error())
				return
			}
			req.Username = strings.TrimSpace(req.Username)
			req.Password = strings.TrimSpace(req.Password)
			if req.Username == "" || req.Password == "" {
				response.BadRequest(c, "username 和 password 不能为空")
				return
			}
			if req.Nickname == "" {
				req.Nickname = req.Username
			}

			u, err := store.Register(req.Username, req.Password, req.Nickname)
			if err != nil {
				response.BadRequest(c, err.Error())
				return
			}
			response.Created(c, "注册成功", toVO(u))
		})

		// ---------- POST /user/login 登录 ----------
		user.POST("/login", func(c *gin.Context) {
			var req LoginReq
			if err := c.ShouldBindJSON(&req); err != nil {
				response.BadRequest(c, "JSON 格式错误: "+err.Error())
				return
			}
			if req.Username == "" || req.Password == "" {
				response.BadRequest(c, "username 和 password 不能为空")
				return
			}

			u, err := store.Login(req.Username, req.Password)
			if err != nil {
				// 登录失败常用 401
				response.Fail(c, http.StatusUnauthorized, 401, err.Error())
				return
			}

			// 教学用简易 token（正式项目用 JWT）
			token := fmt.Sprintf("token-%d-%d", u.ID, time.Now().Unix())
			response.OK(c, "登录成功", gin.H{
				"token": token,
				"user":  toVO(u),
			})
		})

		// ---------- GET /user/list 列表 ----------
		// 注意：静态路径 /list 必须写在 /:id 前面，否则 "list" 会被当成 id
		user.GET("/list", func(c *gin.Context) {
			list := store.List()
			response.OK(c, "获取列表成功", gin.H{
				"count": len(list),
				"list":  list,
			})
		})

		// ---------- GET /user/:id 详情 ----------
		user.GET("/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil || id <= 0 {
				response.BadRequest(c, "id 必须是正整数")
				return
			}
			u, err := store.Get(id)
			if err != nil {
				response.NotFound(c, err.Error())
				return
			}
			response.OK(c, "获取详情成功", u)
		})

		// ---------- PUT /user/:id 更新 ----------
		user.PUT("/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil || id <= 0 {
				response.BadRequest(c, "id 必须是正整数")
				return
			}

			var req UpdateReq
			if err := c.ShouldBindJSON(&req); err != nil {
				response.BadRequest(c, "JSON 格式错误: "+err.Error())
				return
			}
			if strings.TrimSpace(req.Nickname) == "" && req.Age <= 0 {
				response.BadRequest(c, "请至少提供 nickname 或 age")
				return
			}

			u, err := store.Update(id, strings.TrimSpace(req.Nickname), req.Age)
			if err != nil {
				response.NotFound(c, err.Error())
				return
			}
			response.OK(c, "更新成功", u)
		})
	}

	runServer(r, ":8085")
}
