package gou

import (
	"github.com/gin-gonic/gin"
)

// RspBase 代表返回的公共结构
type RspBase struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// RspBaseData 代表返回的公共结构
type RspBaseData struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ErrHandled 处理 err 错误，并且返回给 c
func ErrHandled(err error, c *gin.Context) bool {
	if err == nil {
		return false
	}

	c.JSON(200, RspBase{Status: 400, Message: err.Error()})
	return true
}
