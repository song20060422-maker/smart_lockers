// @ticket B端-统一响应封装，与C端保持一致的全项目唯一出参结构
package pkg

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
}

func requestID(c *gin.Context) string {
	if v, ok := c.Get("request_id"); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// OK 成功响应
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: CodeOK, Message: "ok", Data: data, RequestID: requestID(c)})
}

// Fail 业务错误响应（HTTP 200 + 业务码）
func Fail(c *gin.Context, code int) {
	e := GetErr(code)
	c.JSON(http.StatusOK, Response{Code: e.Code, Message: e.Message, RequestID: requestID(c)})
}

// FailMsg 带自定义消息
func FailMsg(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{Code: code, Message: msg, RequestID: requestID(c)})
}

// AbortWithCode 中间件中断
func AbortWithCode(c *gin.Context, code int) {
	Fail(c, code)
	c.Abort()
}
