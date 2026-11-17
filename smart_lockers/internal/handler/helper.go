// @ticket 00-project：handler 公共响应辅助
package handler

import (
	"github.com/gin-gonic/gin"

	"smart_lockers/internal/pkg"
)

// respond 统一处理 service 返回的 (data, code)
func respond(c *gin.Context, data interface{}, code int) {
	if code == pkg.CodeOK {
		pkg.OK(c, data)
	} else {
		pkg.Fail(c, code)
	}
}

// idemKey 从请求头取幂等键 @ticket 00-project
func idemKey(c *gin.Context) string {
	k := c.GetHeader("Idempotency-Key")
	if k == "" {
		k = c.Query("idempotency_key")
	}
	return k
}
