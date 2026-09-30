package utils

import "fmt"

const (
	CodeOK            = 0
	CodeUserNotFound  = 10001
	CodePasswordWrong = 10002
	CodePhoneExists   = 10003
	CodeForbidden     = 10004
	CodeUnauthorized  = 10005
	CodeTodoNotFound  = 10006 // 待办不存在或不属于当前用户
	CodeBadRequest    = 400
	CodeServer        = 500
)

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

func Unauthorized(msg string) *AppError {
	if msg == "" {
		msg = "请先登录"
	}
	return NewAppError(401, CodeUnauthorized, msg)
}

func TodoNotFound() *AppError {
	return NewAppError(404, CodeTodoNotFound, "待办不存在")
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
