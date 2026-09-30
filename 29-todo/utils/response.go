package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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

func FailApp(c *gin.Context, err *AppError) {
	JSON(c, err.HTTP, err.Code, err.Message, nil)
}
