// @ticket B端-通用中间件：RequestID、访问日志、Recovery、CORS
package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"smart_lockers_admin/internal/pkg"
)

// 定义函数 RequestID：封装一段可复用的处理逻辑（参数：）

func RequestID() gin.HandlerFunc {
	// 返回 func(c *gin.Context) { 给调用方，结束当前函数

	return func(c *gin.Context) {
		// 定义并初始化局部变量 id，取值为 c.GetHeader("X-Request-ID")

		id := c.GetHeader("X-Request-ID")
		// 条件判断：id == "" { id = uuid.NewString() } 时执行下面逻辑

		if id == "" { id = uuid.NewString() }
		// 执行 c.Set("request_id", id)，完成对应的业务处理

		c.Set("request_id", id)
		// 执行 c.Next()，完成对应的业务处理

		c.Next()
	// 当前代码块结束

	}
// 当前代码块结束

}

// 定义函数 Logger：封装一段可复用的处理逻辑（参数：）

func Logger() gin.HandlerFunc {
	// 返回 func(c *gin.Context) { 给调用方，结束当前函数

	return func(c *gin.Context) {
		// 定义并初始化局部变量 start，取值为 time.Now()

		start := time.Now()
		// 执行 c.Next()，完成对应的业务处理

		c.Next()
		// 执行 pkg.Logf("[ACCESS] %s %s status=%d cost=%v rid=%v role=%v",，完成对应的业务处理

		pkg.Logf("[ACCESS] %s %s status=%d cost=%v rid=%v role=%v",
			// 执行 c.Request.Method, c.Request.URL.Path, c.Writer.Status(), time.Since(start),，完成对应的业务处理

			c.Request.Method, c.Request.URL.Path, c.Writer.Status(), time.Since(start),
			// 执行 c.GetString("request_id"), c.GetString("role"))，完成对应的业务处理

			c.GetString("request_id"), c.GetString("role"))
	// 当前代码块结束

	}
// 当前代码块结束

}

// 定义函数 Recovery：封装一段可复用的处理逻辑（参数：）

func Recovery() gin.HandlerFunc {
	// 返回 gin.CustomRecoveryWithWriter(nil, func(c *gin.Context, err interface{}) { 给调用方，结束当前函数

	return gin.CustomRecoveryWithWriter(nil, func(c *gin.Context, err interface{}) {
		// 执行 pkg.Logf("[PANIC] %s %s panic=%v", c.Request.Method, c.Request.URL.Path, err)，完成对应的业务处理

		pkg.Logf("[PANIC] %s %s panic=%v", c.Request.Method, c.Request.URL.Path, err)
		// 执行 pkg.Fail(c, pkg.ErrInternal)，完成对应的业务处理

		pkg.Fail(c, pkg.ErrInternal)
	// 执行 })，完成对应的业务处理

	})
// 当前代码块结束

}

// 定义函数 CORS：封装一段可复用的处理逻辑（参数：）

func CORS() gin.HandlerFunc {
	// 返回 func(c *gin.Context) { 给调用方，结束当前函数

	return func(c *gin.Context) {
		// 执行 c.Header("Access-Control-Allow-Origin", "*")，完成对应的业务处理

		c.Header("Access-Control-Allow-Origin", "*")
		// 执行 c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")，完成对应的业务处理

		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		// 执行 c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization,X-Request-Id,Idempotency-Key")，完成对应的业务处理

		c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization,X-Request-Id,Idempotency-Key")
		// 条件判断：c.Request.Method == http.MethodOptions 时执行下面逻辑

		if c.Request.Method == http.MethodOptions {
			// 执行 c.AbortWithStatus(http.StatusNoContent)，完成对应的业务处理

			c.AbortWithStatus(http.StatusNoContent)
			// 返回  给调用方，结束当前函数

			return
		// 当前代码块结束

		}
		// 执行 c.Next()，完成对应的业务处理

		c.Next()
	// 当前代码块结束

	}
// 当前代码块结束

}
