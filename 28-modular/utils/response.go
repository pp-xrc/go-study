package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Body 统一响应：HTTP 状态码管传输，code 管业务
type Body struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func JSON(c *gin.Context, httpStatus, code int, message string, data any) {
	c.JSON(httpStatus, Body{Code: code, Message: message, Data: data})
}

func OK(c *gin.Context, message string, data any) {
	JSON(c, http.StatusOK, CodeOK, message, data)
}

func Created(c *gin.Context, message string, data any) {
	JSON(c, http.StatusCreated, CodeOK, message, data)
}

// FailApp 把自定义错误写成统一 JSON，并 Abort（给中间件用）
func FailApp(c *gin.Context, err *AppError) {
	JSON(c, err.HTTP, err.Code, err.Message, nil)
}

// Handle 控制器统一出口：成功走 OK；失败识别 AppError / 其它错误
func Handle(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	if ae, ok := err.(*AppError); ok {
		FailApp(c, ae)
		return true
	}
	FailApp(c, Server(err.Error()))
	return true
}
