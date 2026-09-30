/*
第 20 课：Gin 中间件（Middleware）

================================================================================
一、什么是中间件？
================================================================================
中间件 = 请求到达「业务 Handler」之前 / 之后，统一执行的拦截逻辑。

典型用途：
  - 日志：记录每个请求的方法、路径、时间、耗时
  - 鉴权：校验 Token / Cookie，未登录则拦截
  - 跨域（CORS）：给浏览器加 Access-Control-* 响应头
  - 限流、埋点、统一错误处理……

================================================================================
二、请求执行流程（洋葱模型）
================================================================================

  客户端请求
      │
      ▼
  ┌─────────────┐
  │ 中间件 A 前半 │  ← c.Next() 之前
  └──────┬──────┘
         ▼
  ┌─────────────┐
  │ 中间件 B 前半 │
  └──────┬──────┘
         ▼
  ┌─────────────┐
  │  业务 Handler │  ← 真正的接口逻辑
  └──────┬──────┘
         ▼
  ┌─────────────┐
  │ 中间件 B 后半 │  ← c.Next() 之后
  └──────┬──────┘
         ▼
  ┌─────────────┐
  │ 中间件 A 后半 │
  └──────┬──────┘
         ▼
     响应客户端

关键 API：
  c.Next()   继续执行下一个中间件 / Handler
  c.Abort()  中断后续中间件和 Handler（常用于鉴权失败、OPTIONS 预检）

================================================================================
三、注册方式
================================================================================
  r.Use(中间件)                 全局：所有路由都经过
  group.Use(中间件)             分组：只作用于该组路由
  r.GET("/x", 中间件, 处理函数)  单路由：只作用于这一条

本题要求：
  1. 自定义请求日志中间件，打印方法、路径、时间
  2. 实现 CORS 中间件（允许所有源、方法、请求头）
  3. 日志 → 全局；CORS → /user 路由分组

扩展补充（教学加分）：
  4. 鉴权中间件 AuthRequired（演示 c.Abort / c.Set / c.Get）
  5. /private 分组挂鉴权；单路由挂鉴权

================================================================================
四、运行与测试
================================================================================
  go run ./20-gin-middleware

  # /user 分组：响应头应有 Access-Control-*
  curl -i http://localhost:8083/user/list
  curl -i -X OPTIONS http://localhost:8083/user/list
  curl -i -X POST http://localhost:8083/user/create

  # 非 /user：没有 CORS 头（对比）
  curl -i http://localhost:8083/public/hello

  # 扩展：鉴权
  curl -i http://localhost:8083/private/profile
  curl -i -H "Authorization: Bearer demo-token" http://localhost:8083/private/profile
  curl -i -H "Authorization: Bearer demo-token" http://localhost:8083/admin
*/

package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"go-study/response"
)

func runServer(r *gin.Engine, addr string) {
	fmt.Println("⏳ 正在启动服务...")
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("❌ 启动失败:", err)
		return
	}
	fmt.Println("✅ 服务已就绪!")
	fmt.Printf("   http://localhost%s\n", addr)
	fmt.Println("   curl -i http://localhost" + addr + "/user/list")
	fmt.Println("   curl -i http://localhost" + addr + "/public/hello")
	fmt.Println("   curl -i -H \"Authorization: Bearer demo-token\" http://localhost" + addr + "/private/profile")
	fmt.Println("   按 Ctrl+C 停止")
	fmt.Println(strings.Repeat("-", 40))
	_ = http.Serve(listener, r)
}

// ---------------------------------------------------------------------------
// ① 自定义请求日志中间件
// 返回类型必须是 gin.HandlerFunc，即 func(*gin.Context)
// 打印：请求时间、方法、路径（耗时作为加分观察项）
// ---------------------------------------------------------------------------
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now() // 请求进入时的时间

		// 把控制权交给后续中间件 / Handler
		c.Next()

		// ---- 下面是「后半段」：Handler 执行完才会跑到这里 ----
		fmt.Printf("[日志] 时间=%s | 方法=%s | 路径=%s | 耗时=%v\n",
			start.Format("2006-01-02 15:04:05"),
			c.Request.Method,
			c.Request.URL.Path,
			time.Since(start),
		)
	}
}

// ---------------------------------------------------------------------------
// ② CORS 跨域中间件
// 浏览器跨域请求时，需要服务端返回 Access-Control-* 响应头
// 本题要求：允许所有源、所有方法、所有请求头
// ---------------------------------------------------------------------------
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")  // 所有源
		c.Header("Access-Control-Allow-Methods", "*") // 所有方法
		c.Header("Access-Control-Allow-Headers", "*") // 所有请求头

		// 预检请求（OPTIONS）直接返回，不必进业务 Handler
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// ---------------------------------------------------------------------------
// ④ 扩展：简易鉴权中间件
// 要求请求头：Authorization: Bearer demo-token
// 校验失败则 Abort，后面的 Handler 不会执行
// ---------------------------------------------------------------------------
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth != "Bearer demo-token" {
			response.Fail(c, http.StatusUnauthorized, 401, "未授权：请携带 Authorization: Bearer demo-token")
			c.Abort() // 中断：不再执行后续中间件和业务 Handler
			return
		}

		// 鉴权通过：把用户信息塞进上下文，供后面 Handler 使用
		c.Set("user_id", "10001")
		c.Set("user_name", "demo-user")

		c.Next() // 继续往后走
	}
}

func main() {
	// gin.Default() 自带 Logger + Recovery
	// 这里用 gin.New()，方便完全自己控制中间件
	r := gin.New()
	r.Use(gin.Recovery()) // panic 时自动恢复，避免整个进程崩掉

	// ---------- ③-1 日志中间件：全局注册 ----------
	// 所有路由（/public、/user）都会经过 RequestLogger
	r.Use(RequestLogger())

	// ---------- 非 /user 路由：没有 CORS（用于对比）----------
	r.GET("/public/hello", func(c *gin.Context) {
		response.OK(c, "公开接口（无 CORS 中间件）", gin.H{
			"tip": "对比 /user/list 的响应头，这里没有 Access-Control-*",
		})
	})

	// ---------- ③-2 CORS 中间件：只挂在 /user 分组 ----------
	// group.Use(...) 只对本分组内的路由生效
	user := r.Group("/user")
	user.Use(CORS())
	{
		user.GET("/list", func(c *gin.Context) {
			response.OK(c, "用户列表（带 CORS）", gin.H{
				"list": []string{"张三", "李四"},
			})
		})

		user.POST("/create", func(c *gin.Context) {
			response.Created(c, "创建用户成功（带 CORS）", gin.H{
				"id": 10001,
			})
		})
	}

	// ---------- ④ 扩展：分组级鉴权中间件 ----------
	private := r.Group("/private")
	private.Use(AuthRequired()) // 只对 /private 下的路由生效
	{
		private.GET("/profile", func(c *gin.Context) {
			// 从鉴权中间件写入的上下文中取值
			userID, _ := c.Get("user_id")
			userName, _ := c.Get("user_name")

			response.OK(c, "鉴权通过", gin.H{
				"user_id":   userID,
				"user_name": userName,
			})
		})

		private.GET("/secret", func(c *gin.Context) {
			response.OK(c, "机密数据", gin.H{"secret": "42"})
		})
	}

	// ---------- ⑤ 扩展：单路由中间件 ----------
	// 写法：r.GET(路径, 中间件, 处理函数)
	r.GET("/admin", AuthRequired(), func(c *gin.Context) {
		response.OK(c, "管理员页", gin.H{"role": "admin"})
	})

	runServer(r, ":8083")
}
