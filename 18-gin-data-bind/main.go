/*
第 18 课：Gin 数据绑定 —— 路径参数、查询参数、JSON 绑定

运行：
  go run ./18-gin-data-bind

测试：
  curl "http://localhost:8081/hello/张三?msg=你好"
  curl "http://localhost:8081/user/10001?role=admin&active=true"
  curl -X POST http://localhost:8081/user \
    -H "Content-Type: application/json" \
    -d '{"name":"张三","age":18,"phone":"13800138000"}'
  curl -X POST http://localhost:8081/user -H "Content-Type: application/json" -d '{"name":"","age":18}'
  curl -X POST http://localhost:8081/user -H "Content-Type: application/json" -d '{"name":"张三","age":0}'
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

// UserInfo 前端 JSON 对应的结构体
type UserInfo struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Phone string `json:"phone"`
}

func validateUser(user UserInfo) error {
	if strings.TrimSpace(user.Name) == "" {
		return fmt.Errorf("name 不能为空")
	}
	if user.Age <= 0 {
		return fmt.Errorf("age 必须大于 0，当前值: %d", user.Age)
	}
	if user.Phone != "" && len(user.Phone) < 11 {
		return fmt.Errorf("phone 格式不正确，长度至少 11 位")
	}
	return nil
}

func runServer(r *gin.Engine, addr string) {
	fmt.Println("⏳ 正在启动服务...")
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("❌ 启动失败:", err)
		return
	}
	fmt.Println("✅ 服务已就绪!")
	fmt.Printf("   http://localhost%s\n", addr)
	fmt.Println("   curl \"http://localhost" + addr + "/hello/张三?msg=你好\"")
	fmt.Println("   curl -X POST http://localhost" + addr + "/user \\")
	fmt.Println("     -H \"Content-Type: application/json\" \\")
	fmt.Println("     -d '{\"name\":\"张三\",\"age\":18,\"phone\":\"13800138000\"}'")
	fmt.Println("   按 Ctrl+C 停止")
	fmt.Println(strings.Repeat("-", 40))
	_ = http.Serve(listener, r)
}

func main() {
	r := gin.Default()

	// 路径参数 + 查询参数
	r.GET("/hello/:name", func(c *gin.Context) {
		name := c.Param("name")
		msg := c.DefaultQuery("msg", "Hello")
		lang := c.Query("lang")

		result := msg + ", " + name + "!"
		if lang != "" {
			result = result + " [" + lang + "]"
		}

		response.OK(c, "拼接成功", gin.H{
			"path_param":  name,
			"query_msg":   msg,
			"query_lang":  lang,
			"result":      result,
			"full_path":   c.FullPath(),
			"request_url": c.Request.URL.String(),
		})
	})

	r.GET("/user/:id", func(c *gin.Context) {
		id := c.Param("id")
		role := c.DefaultQuery("role", "guest")
		active := c.DefaultQuery("active", "false")

		response.OK(c, "获取用户成功", gin.H{
			"id":     id,
			"role":   role,
			"active": active,
			"desc":   "用户 " + id + "（角色:" + role + ", 激活:" + active + "）",
		})
	})

	// POST JSON + 校验
	r.POST("/user", func(c *gin.Context) {
		var user UserInfo

		if err := c.ShouldBindJSON(&user); err != nil {
			response.BadRequest(c, "JSON 绑定失败: "+err.Error())
			return
		}

		if err := validateUser(user); err != nil {
			response.BadRequest(c, "参数校验失败: "+err.Error())
			return
		}

		response.OK(c, "接收成功", user)
	})

	runServer(r, ":8081")
}
