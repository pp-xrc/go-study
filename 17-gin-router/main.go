/*
第 17 课：Gin 路由 —— GET/POST/PUT/DELETE、路径参数、路由分组

运行：
  go run ./17-gin-router

测试：
  curl http://localhost:8080/api
  curl -X POST http://localhost:8080/api
  curl http://localhost:8080/user/list
  curl -X POST http://localhost:8080/user/create
  curl http://localhost:8080/user/10001
  curl http://localhost:8080/admin/list
  curl -X POST http://localhost:8080/admin/login
*/

package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"go-study/response"
)

func runServer(r *gin.Engine, addr string) {
	fmt.Println("⏳ 正在编译/启动服务（首次 go run 可能需等待数秒）...")
	fmt.Printf("⏳ 正在绑定端口 %s ...\n", addr)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("❌ 启动失败（端口可能被占用）:", err)
		fmt.Println("   提示: lsof -i", addr, " 查看占用进程")
		return
	}

	fmt.Println("✅ 服务已就绪!")
	fmt.Printf("   本地访问: http://localhost%s\n", addr)
	fmt.Println("   测试命令:")
	fmt.Println("     curl http://localhost" + addr + "/api")
	fmt.Println("     curl http://localhost" + addr + "/user/list")
	fmt.Println("     curl http://localhost" + addr + "/admin/list")
	fmt.Println("   按 Ctrl+C 停止服务")
	fmt.Println(strings.Repeat("-", 40))

	if err := http.Serve(listener, r); err != nil {
		fmt.Println("❌ 服务异常退出:", err)
	}
}

func main() {
	r := gin.Default()

	r.GET("/api", func(c *gin.Context) {
		response.OK(c, "查询数据", gin.H{"method": "GET"})
	})

	r.POST("/api", func(c *gin.Context) {
		response.OK(c, "创建数据", gin.H{"method": "POST"})
	})

	r.PUT("/api", func(c *gin.Context) {
		response.OK(c, "更新数据", gin.H{"method": "PUT"})
	})

	r.DELETE("/api", func(c *gin.Context) {
		response.OK(c, "删除数据", gin.H{"method": "DELETE"})
	})

	r.GET("/get", func(c *gin.Context) {
		response.OK(c, "ok", gin.H{"method": "GET"})
	})
	r.POST("/post", func(c *gin.Context) {
		response.OK(c, "ok", gin.H{"method": "POST"})
	})
	r.PUT("/put", func(c *gin.Context) {
		response.OK(c, "ok", gin.H{"method": "PUT"})
	})
	r.DELETE("/delete", func(c *gin.Context) {
		response.OK(c, "ok", gin.H{"method": "DELETE"})
	})

	userGroup := r.Group("/user")
	{
		userGroup.GET("/list", func(c *gin.Context) {
			response.OK(c, "用户列表", gin.H{
				"group":  "user",
				"method": "GET",
				"action": "list",
				"list":   []string{"张三", "李四"},
			})
		})

		userGroup.POST("/create", func(c *gin.Context) {
			response.Created(c, "用户创建成功", gin.H{
				"group":  "user",
				"method": "POST",
				"action": "create",
			})
		})

		userGroup.GET("/:id", func(c *gin.Context) {
			id := c.Param("id")
			response.OK(c, "获取用户成功", gin.H{
				"group":  "user",
				"method": "GET",
				"id":     id,
			})
		})
	}

	adminGroup := r.Group("/admin")
	{
		adminGroup.GET("/list", func(c *gin.Context) {
			response.OK(c, "管理员列表", gin.H{
				"group":  "admin",
				"method": "GET",
				"action": "list",
				"list":   []string{"admin1", "admin2"},
			})
		})

		adminGroup.POST("/login", func(c *gin.Context) {
			response.OK(c, "管理员登录成功", gin.H{
				"group":  "admin",
				"method": "POST",
				"action": "login",
			})
		})
	}

	runServer(r, ":8080")
}
