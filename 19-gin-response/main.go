/*
第 19 课：Gin 响应处理 + 统一返回格式

response 包：
  response.Response(c, httpStatus, code, message, data)  // 通用
  response.OK / Created / BadRequest / NotFound / ServerError  // 快捷方法

运行：
  go run ./19-gin-response

测试：
  curl http://localhost:8082/json
  curl http://localhost:8082/string
  curl -i -X POST http://localhost:8082/created
  curl -i http://localhost:8082/not-found
  curl -i http://localhost:8082/bad-request
  curl -i http://localhost:8082/error
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
	fmt.Println("⏳ 正在启动服务...")
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("❌ 启动失败:", err)
		return
	}
	fmt.Println("✅ 服务已就绪!")
	fmt.Printf("   http://localhost%s\n", addr)
	fmt.Println("   curl http://localhost" + addr + "/json")
	fmt.Println("   curl -i -X POST http://localhost" + addr + "/created")
	fmt.Println("   curl -i http://localhost" + addr + "/not-found")
	fmt.Println("   按 Ctrl+C 停止")
	fmt.Println(strings.Repeat("-", 40))
	_ = http.Serve(listener, r)
}

func main() {
	r := gin.Default()

	// ---------- 1. 返回 JSON（统一格式）----------
	r.GET("/json", func(c *gin.Context) {
		response.OK(c, "ok", gin.H{
			"name": "Gin",
			"type": "json",
		})
	})

	// 完全自定义：HTTP 状态码 + 业务码 + 提示 + 数据
	r.GET("/custom", func(c *gin.Context) {
		response.Response(c, http.StatusOK, 0, "自定义成功", gin.H{
			"id":   10001,
			"name": "张三",
		})
	})

	// ---------- 2. 返回字符串 ----------
	r.GET("/string", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, Gin! 这是纯文本响应")
	})

	r.GET("/string/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "你好, %s!", name)
	})

	// ---------- 3. 指定状态码（统一格式）----------
	// 注意：/created 是 POST，要用 curl -X POST
	r.POST("/created", func(c *gin.Context) {
		response.Created(c, "创建成功", gin.H{"id": 1})
	})

	r.GET("/bad-request", func(c *gin.Context) {
		response.BadRequest(c, "参数错误")
	})

	r.GET("/not-found", func(c *gin.Context) {
		response.NotFound(c, "资源不存在")
	})

	r.GET("/error", func(c *gin.Context) {
		response.ServerError(c, "服务器内部错误")
	})

	r.GET("/string/forbidden", func(c *gin.Context) {
		c.String(http.StatusForbidden, "禁止访问（403）")
	})

	runServer(r, ":8082")
}
