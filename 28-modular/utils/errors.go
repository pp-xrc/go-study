package utils

import "fmt"

// 业务错误码（成功为 0，写在响应的 code 字段里）
const (
	CodeOK            = 0
	CodeUserNotFound  = 10001 // 用户不存在
	CodePasswordWrong = 10002 // 密码错误
	CodePhoneExists   = 10003 // 手机号已注册
	CodeForbidden     = 10004 // 权限不足
	CodeUnauthorized  = 10005 // 未登录 / token 无效
	CodeBadRequest    = 400
	CodeServer        = 500
)

// AppError 可带 HTTP 状态码 + 业务码，controller 统一转成 JSON
type AppError struct {
	HTTP    int
	Code    int
	Message string
}

func (e *AppError) Error() string { return e.Message }

func NewAppError(httpStatus, code int, msg string) *AppError {
	return &AppError{HTTP: httpStatus, Code: code, Message: msg}
}

func UserNotFound() *AppError {
	return NewAppError(401, CodeUserNotFound, "用户不存在")
}

func PasswordWrong() *AppError {
	return NewAppError(401, CodePasswordWrong, "密码错误")
}

func PhoneExists() *AppError {
	return NewAppError(409, CodePhoneExists, "手机号已注册")
}

func Forbidden(msg string) *AppError {
	if msg == "" {
		msg = "权限不足"
	}
	return NewAppError(403, CodeForbidden, msg)
}

func Unauthorized(msg string) *AppError {
	if msg == "" {
		msg = "请先登录"
	}
	return NewAppError(401, CodeUnauthorized, msg)
}

func BadRequest(msg string) *AppError {
	return NewAppError(400, CodeBadRequest, msg)
}

func Server(msg string) *AppError {
	if msg == "" {
		msg = "服务器内部错误"
	}
	return NewAppError(500, CodeServer, msg)
}

func WrapDB(err error) *AppError {
	return Server(fmt.Sprintf("数据库错误: %v", err))
}
