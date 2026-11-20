// @ticket 00-project：统一响应封装，全项目唯一出参结构
package pkg

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 定义结构体 Response：用于承载一组相关的字段数据

type Response struct {
	// 结构体字段 Code，类型为 int，附带标签 `json:"code"`

	Code      int         `json:"code"`
	// 结构体字段 Message，类型为 string，附带标签 `json:"message"`

	Message   string      `json:"message"`
	// 执行 Data      interface{} `json:"data,omitempty"`，完成对应的业务处理

	Data      interface{} `json:"data,omitempty"`
	// 结构体字段 RequestID，类型为 string，附带标签 `json:"request_id,omitempty"`

	RequestID string      `json:"request_id,omitempty"`
// 当前代码块结束

}

// 声明常量 requestIDKey = "request_id"：固定不变的取值

const requestIDKey = "request_id"

// 定义函数 requestID：封装一段可复用的处理逻辑（参数：c *gin.Context）

func requestID(c *gin.Context) string {
	// 条件判断：v, ok := c.Get(requestIDKey); ok 时执行下面逻辑

	if v, ok := c.Get(requestIDKey); ok {
		// 条件判断：s, ok := v.(string); ok 时执行下面逻辑

		if s, ok := v.(string); ok {
			// 返回 s 给调用方，结束当前函数

			return s
		// 当前代码块结束

		}
	// 当前代码块结束

	}
	// 返回 "" 给调用方，结束当前函数

	return ""
// 当前代码块结束

}

// OK 成功响应
// 定义函数 OK：封装一段可复用的处理逻辑（参数：c *gin.Context, data interface{}）

func OK(c *gin.Context, data interface{}) {
	// 执行 c.JSON(http.StatusOK, Response{Code: 0, Message: "ok", Data: data, RequestID: requestID(c)})，完成对应的业务处理

	c.JSON(http.StatusOK, Response{Code: 0, Message: "ok", Data: data, RequestID: requestID(c)})
// 当前代码块结束

}

// Fail 业务错误响应（HTTP 200 + 业务码，前后端统一按 code 判断）
// 定义函数 Fail：封装一段可复用的处理逻辑（参数：c *gin.Context, code int）

func Fail(c *gin.Context, code int) {
	// 定义并初始化局部变量 e，取值为 GetErr(code)

	e := GetErr(code)
	// 执行 c.JSON(http.StatusOK, Response{Code: e.Code, Message: e.Message, RequestID: requestID(c)})，完成对应的业务处理

	c.JSON(http.StatusOK, Response{Code: e.Code, Message: e.Message, RequestID: requestID(c)})
// 当前代码块结束

}

// FailMsg 带自定义消息的响应
// 定义函数 FailMsg：封装一段可复用的处理逻辑（参数：c *gin.Context, code int, msg string）

func FailMsg(c *gin.Context, code int, msg string) {
	// 执行 c.JSON(http.StatusOK, Response{Code: code, Message: msg, RequestID: requestID(c)})，完成对应的业务处理

	c.JSON(http.StatusOK, Response{Code: code, Message: msg, RequestID: requestID(c)})
// 当前代码块结束

}

// AbortWithCode 中间件中断
// 定义函数 AbortWithCode：封装一段可复用的处理逻辑（参数：c *gin.Context, code int）

func AbortWithCode(c *gin.Context, code int) {
	// 执行 Fail(c, code)，完成对应的业务处理

	Fail(c, code)
	// 执行 c.Abort()，完成对应的业务处理

	c.Abort()
// 当前代码块结束

}
