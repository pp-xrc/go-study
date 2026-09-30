// Package response 统一 API 返回格式
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Body 统一响应体结构
type Body struct {
	Code    int    `json:"code"`    // 业务状态码：0 成功，非 0 失败
	Message string `json:"message"` // 提示信息
	Data    any    `json:"data"`    // 业务数据，可为 nil
}

// Response 统一返回函数：自定义 HTTP 状态码、提示信息、数据
//
//	response.Response(c, http.StatusOK, 0, "ok", data)
//	response.Response(c, http.StatusBadRequest, 400, "参数错误", nil)
func Response(c *gin.Context, httpStatus int, code int, message string, data any) {
	c.JSON(httpStatus, Body{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// ---------- 常用快捷方法 ----------

// OK 成功（HTTP 200，业务码 0）
func OK(c *gin.Context, message string, data any) {
	Response(c, http.StatusOK, 0, message, data)
}

// Created 创建成功（HTTP 201）
func Created(c *gin.Context, message string, data any) {
	Response(c, http.StatusCreated, 0, message, data)
}

// Fail 失败：自定义 HTTP 状态码和业务码
func Fail(c *gin.Context, httpStatus int, code int, message string) {
	Response(c, httpStatus, code, message, nil)
}

// BadRequest 参数错误（HTTP 400）
func BadRequest(c *gin.Context, message string) {
	Fail(c, http.StatusBadRequest, 400, message)
}

// NotFound 资源不存在（HTTP 404）
func NotFound(c *gin.Context, message string) {
	Fail(c, http.StatusNotFound, 404, message)
}

// ServerError 服务器错误（HTTP 500）
func ServerError(c *gin.Context, message string) {
	Fail(c, http.StatusInternalServerError, 500, message)
}
