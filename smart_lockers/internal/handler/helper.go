// @ticket 00-project：handler 公共响应辅助
// @proto 原型对齐：无直接对应(公共响应辅助)；见 AGENT.MD 2.1
package handler

import (
	"github.com/gin-gonic/gin"

	"smart_lockers/internal/pkg"
)

// respond 统一处理 service 返回的 (data, code)
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

// idemKey 从请求头取幂等键 @ticket 00-project
// 定义函数 idemKey：封装一段可复用的处理逻辑（参数：c *gin.Context）

func idemKey(c *gin.Context) string {
	// 定义并初始化局部变量 k，取值为 c.GetHeader("Idempotency-Key")

	k := c.GetHeader("Idempotency-Key")
	// 条件判断：k == "" 时执行下面逻辑

	if k == "" {
		// 为 k 赋值：c.Query("idempotency_key")

		k = c.Query("idempotency_key")
	// 当前代码块结束

	}
	// 返回 k 给调用方，结束当前函数

	return k
// 当前代码块结束

}
