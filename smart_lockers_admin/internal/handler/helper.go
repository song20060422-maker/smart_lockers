// @ticket B端-公共响应与参数辅助
package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"smart_lockers_admin/internal/pkg"
)

// respond 统一处理 (data, code)
func respond(c *gin.Context, data interface{}, code int) {
	if code == pkg.CodeOK {
		pkg.OK(c, data)
	} else {
		pkg.Fail(c, code)
	}
}

func idemKey(c *gin.Context) string {
	k := c.GetHeader("Idempotency-Key")
	if k == "" { k = c.Query("idempotency_key") }
	return k
}

func uintParam(c *gin.Context, name string) (uint, bool) {
	v, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil || v == 0 { return 0, false }
	return uint(v), true
}
