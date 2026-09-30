/*
第 21 课：Gin 文件上传 + 路由重定向

学习内容：
  - 单文件上传：c.FormFile() + c.SaveUploadedFile()
  - 多文件上传：c.MultipartForm()，遍历保存
  - 路由重定向：c.Redirect()（临时 302 / 永久 301）

实操任务：
  1. POST /upload/single  上传单张图片到项目 upload 目录
  2. POST /upload/multi   批量上传多张图片
  3. GET  /old → 重定向到 /new

复盘要点：
  - 注意文件大小限制、文件类型校验（本课末尾有简单补充）

运行：
  go run ./21-gin-upload

测试单文件：
  curl -F "file=@./README.md" http://localhost:8084/upload/single

测试多文件：
  curl -F "files=@./README.md" -F "files=@./go.mod" http://localhost:8084/upload/multi

测试重定向（-L 跟随跳转）：
  curl -i http://localhost:8084/old
  curl -iL http://localhost:8084/old
*/

package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"go-study/response"
)

const (
	uploadDir   = "upload"       // 保存目录（相对当前工作目录，建议在项目根目录运行）
	maxMemory   = 8 << 20        // Multipart 内存缓冲：8MB
	maxFileSize = 10 << 20       // 单文件最大：10MB（复盘补充）
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
	fmt.Println("   curl -F \"file=@./README.md\" http://localhost" + addr + "/upload/single")
	fmt.Println("   curl -i http://localhost" + addr + "/old")
	fmt.Println("   按 Ctrl+C 停止")
	fmt.Println(strings.Repeat("-", 40))
	_ = http.Serve(listener, r)
}

// ensureUploadDir 确保 upload 目录存在
func ensureUploadDir() error {
	return os.MkdirAll(uploadDir, 0755)
}

// savePath 生成保存路径：upload/原始文件名
func savePath(filename string) string {
	// 只用文件名，防止路径穿越（如 ../../etc/passwd）
	base := filepath.Base(filename)
	return filepath.Join(uploadDir, base)
}

func main() {
	if err := ensureUploadDir(); err != nil {
		fmt.Println("创建 upload 目录失败:", err)
		return
	}

	r := gin.Default()

	// 限制 multipart 表单内存（超出部分写临时文件）
	r.MaxMultipartMemory = maxMemory

	// =====================================================================
	// 任务 1：单文件上传
	// 表单字段名：file
	// 使用：c.FormFile("file") + c.SaveUploadedFile(file, dst)
	// =====================================================================
	r.POST("/upload/single", func(c *gin.Context) {
		// FormFile：从 multipart/form-data 中取出名为 file 的单个文件
		file, err := c.FormFile("file")
		if err != nil {
			response.BadRequest(c, "请使用表单字段 file 上传文件: "+err.Error())
			return
		}

		// 复盘补充：简单校验大小
		if file.Size > maxFileSize {
			response.BadRequest(c, fmt.Sprintf("文件过大，限制 %dMB", maxFileSize>>20))
			return
		}

		dst := savePath(file.Filename)
		// SaveUploadedFile：把上传内容保存到本地路径
		if err := c.SaveUploadedFile(file, dst); err != nil {
			response.ServerError(c, "保存失败: "+err.Error())
			return
		}

		response.OK(c, "单文件上传成功", gin.H{
			"filename": file.Filename,
			"size":     file.Size,
			"saved_as": dst,
		})
	})

	// =====================================================================
	// 任务 2：多文件上传
	// 表单字段名：files（可多个同名字段）
	// 使用：c.MultipartForm()，再遍历 form.File["files"]
	// =====================================================================
	r.POST("/upload/multi", func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			response.BadRequest(c, "解析 multipart 失败: "+err.Error())
			return
		}

		files := form.File["files"]
		if len(files) == 0 {
			response.BadRequest(c, "请使用表单字段 files 上传一个或多个文件")
			return
		}

		var saved []gin.H
		for _, file := range files {
			if file.Size > maxFileSize {
				response.BadRequest(c, fmt.Sprintf("%s 过大，限制 %dMB", file.Filename, maxFileSize>>20))
				return
			}

			dst := savePath(file.Filename)
			if err := c.SaveUploadedFile(file, dst); err != nil {
				response.ServerError(c, "保存 "+file.Filename+" 失败: "+err.Error())
				return
			}

			saved = append(saved, gin.H{
				"filename": file.Filename,
				"size":     file.Size,
				"saved_as": dst,
			})
		}

		response.OK(c, "多文件上传成功", gin.H{
			"count": len(saved),
			"files": saved,
		})
	})

	// =====================================================================
	// 任务 3：路由重定向
	// 临时重定向 302：浏览器下次仍可能访问 /old
	// 永久重定向 301：告诉客户端「资源已永久搬家」到 /new
	// =====================================================================
	r.GET("/old", func(c *gin.Context) {
		// http.StatusFound = 302 临时重定向
		c.Redirect(http.StatusFound, "/new")
	})

	r.GET("/old-permanent", func(c *gin.Context) {
		// http.StatusMovedPermanently = 301 永久重定向
		c.Redirect(http.StatusMovedPermanently, "/new")
	})

	r.GET("/new", func(c *gin.Context) {
		response.OK(c, "这是新地址 /new", gin.H{
			"tip": "你可能是从 /old 或 /old-permanent 跳转过来的",
		})
	})

	runServer(r, ":8084")
}
