// @ticket B端-公共响应与参数辅助
package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"smart_lockers_admin/internal/pkg"
)

// respond 统一处理 (data, code)
// 定义函数 respond：封装一段可复用的处理逻辑（参数：c *gin.Context, data interface{}, code int）

func respond(c *gin.Context, data interface{}, code int) {
	// 条件判断：code == pkg.CodeOK 时执行下面逻辑

	if code == pkg.CodeOK {
		// 执行 pkg.OK(c, data)，完成对应的业务处理

		pkg.OK(c, data)
	// 上述条件不成立时，执行 else 分支逻辑

	} else {
		// 执行 pkg.Fail(c, code)，完成对应的业务处理

		pkg.Fail(c, code)
	// 当前代码块结束

	}
// 当前代码块结束

}

// 定义函数 idemKey：封装一段可复用的处理逻辑（参数：c *gin.Context）

func idemKey(c *gin.Context) string {
	// 定义并初始化局部变量 k，取值为 c.GetHeader("Idempotency-Key")

	k := c.GetHeader("Idempotency-Key")
	// 条件判断：k == "" { k = c.Query("idempotency_key") } 时执行下面逻辑

	if k == "" { k = c.Query("idempotency_key") }
	// 返回 k 给调用方，结束当前函数

	return k
// 当前代码块结束

}

// 定义函数 uintParam：封装一段可复用的处理逻辑（参数：c *gin.Context, name string）

func uintParam(c *gin.Context, name string) (uint, bool) {
	// 定义并初始化局部变量 v, err，取值为 strconv.ParseUint(c.Param(name), 10, 64)

	v, err := strconv.ParseUint(c.Param(name), 10, 64)
	// 条件判断：err != nil || v == 0 { return 0, false } 时执行下面逻辑

	if err != nil || v == 0 { return 0, false }
	// 返回 uint(v), true 给调用方，结束当前函数

	return uint(v), true
// 当前代码块结束

}
