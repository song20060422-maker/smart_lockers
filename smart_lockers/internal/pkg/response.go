// @ticket 00-project：统一响应封装，全项目唯一出参结构
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

const requestIDKey = "request_id"

func requestID(c *gin.Context) string {
	if v, ok := c.Get(requestIDKey); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// OK 成功响应
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: 0, Message: "ok", Data: data, RequestID: requestID(c)})
}

// Fail 业务错误响应（HTTP 200 + 业务码，前后端统一按 code 判断）
func Fail(c *gin.Context, code int) {
	e := GetErr(code)
	c.JSON(http.StatusOK, Response{Code: e.Code, Message: e.Message, RequestID: requestID(c)})
}

// FailMsg 带自定义消息的响应
func FailMsg(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{Code: code, Message: msg, RequestID: requestID(c)})
}

// AbortWithCode 中间件中断
func AbortWithCode(c *gin.Context, code int) {
	Fail(c, code)
	c.Abort()
}
