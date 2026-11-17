// @ticket 00-project：通用中间件：RequestID、访问日志、Recovery、CORS
package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"smart_lockers/internal/pkg"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Request-ID")
		if id == "" {
			id = uuid.NewString()
		}
		c.Set("request_id", id)
		c.Next()
	}
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		pkg.Logf("[ACCESS] %s %s status=%d cost=%v rid=%v",
			c.Request.Method, c.Request.URL.Path, c.Writer.Status(), time.Since(start), c.GetString("request_id"))
	}
}

func Recovery() gin.HandlerFunc {
	return gin.CustomRecoveryWithWriter(nil, func(c *gin.Context, err interface{}) {
		pkg.Logf("[PANIC] %s %s panic=%v", c.Request.Method, c.Request.URL.Path, err)
		pkg.Fail(c, pkg.ErrInternal)
	})
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization,X-Request-Id,Idempotency-Key")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
