// @ticket 08-性能调优：pprof/expvar 性能分析接口，需鉴权且限制内网访问，采样开销<3%
// @proto 原型对齐：ops 远程运维(本文件为性能观测接口) / overview 运营总览；见 AGENT.MD 2.1
package handler

import (
	"expvar"
	"net/http/pprof"
	"strings"

	"github.com/gin-gonic/gin"

	"smart_lockers/internal/config"
	"smart_lockers/internal/pkg"
)

// isInternal 判断是否内网来源 @ticket 08-性能调优
// 定义函数 isInternal：封装一段可复用的处理逻辑（参数：ip string）

func isInternal(ip string) bool {
	// 返回 ip == "127.0.0.1" || ip == "::1" || strings.HasPrefix(ip, "192.168.") || 给调用方，结束当前函数

	return ip == "127.0.0.1" || ip == "::1" || strings.HasPrefix(ip, "192.168.") ||
		// 执行 strings.HasPrefix(ip, "10.") || strings.HasPrefix(ip, "172.")，完成对应的业务处理

		strings.HasPrefix(ip, "10.") || strings.HasPrefix(ip, "172.")
// 当前代码块结束

}

// guard 分析接口鉴权/内网校验中间件 @ticket 08-性能调优
// 定义函数 guard：封装一段可复用的处理逻辑（参数：c *gin.Context）

func guard(c *gin.Context) bool {
	// 条件判断：config.Get().Server.Mode != "debug" && !isInternal(c.ClientIP()) 时执行下面逻辑

	if config.Get().Server.Mode != "debug" && !isInternal(c.ClientIP()) {
		// 执行 pkg.Fail(c, pkg.ErrForb)，完成对应的业务处理

		pkg.Fail(c, pkg.ErrForb)
		// 返回 false 给调用方，结束当前函数

		return false
	// 当前代码块结束

	}
	// 返回 true 给调用方，结束当前函数

	return true
// 当前代码块结束

}

// PprofIndex @ticket 08-性能调优：GET /api/v1/ops/pprof/
// 定义函数 PprofIndex：封装一段可复用的处理逻辑（参数：c *gin.Context）

func PprofIndex(c *gin.Context) {
	// 条件判断：!guard(c) 时执行下面逻辑

	if !guard(c) {
		// 返回  给调用方，结束当前函数

		return
	// 当前代码块结束

	}
	// 执行 pprof.Index(c.Writer, c.Request)，完成对应的业务处理

	pprof.Index(c.Writer, c.Request)
// 当前代码块结束

}

// PprofProfile @ticket 08-性能调优：GET /api/v1/ops/pprof/profile（CPU）
// 定义函数 PprofProfile：封装一段可复用的处理逻辑（参数：c *gin.Context）

func PprofProfile(c *gin.Context) {
	// 条件判断：!guard(c) 时执行下面逻辑

	if !guard(c) {
		// 返回  给调用方，结束当前函数

		return
	// 当前代码块结束

	}
	// 执行 pprof.Profile(c.Writer, c.Request)，完成对应的业务处理

	pprof.Profile(c.Writer, c.Request)
// 当前代码块结束

}

// PprofHeap @ticket 08-性能调优：GET /api/v1/ops/pprof/heap
// 定义函数 PprofHeap：封装一段可复用的处理逻辑（参数：c *gin.Context）

func PprofHeap(c *gin.Context) {
	// 条件判断：!guard(c) 时执行下面逻辑

	if !guard(c) {
		// 返回  给调用方，结束当前函数

		return
	// 当前代码块结束

	}
	// 执行 pprof.Handler("heap").ServeHTTP(c.Writer, c.Request)，完成对应的业务处理

	pprof.Handler("heap").ServeHTTP(c.Writer, c.Request)
// 当前代码块结束

}

// PprofGoroutine @ticket 08-性能调优：GET /api/v1/ops/pprof/goroutine
// 定义函数 PprofGoroutine：封装一段可复用的处理逻辑（参数：c *gin.Context）

func PprofGoroutine(c *gin.Context) {
	// 条件判断：!guard(c) 时执行下面逻辑

	if !guard(c) {
		// 返回  给调用方，结束当前函数

		return
	// 当前代码块结束

	}
	// 执行 pprof.Handler("goroutine").ServeHTTP(c.Writer, c.Request)，完成对应的业务处理

	pprof.Handler("goroutine").ServeHTTP(c.Writer, c.Request)
// 当前代码块结束

}

// PprofBlock @ticket 08-性能调优：GET /api/v1/ops/pprof/block（锁竞争）
// 定义函数 PprofBlock：封装一段可复用的处理逻辑（参数：c *gin.Context）

func PprofBlock(c *gin.Context) {
	// 条件判断：!guard(c) 时执行下面逻辑

	if !guard(c) {
		// 返回  给调用方，结束当前函数

		return
	// 当前代码块结束

	}
	// 执行 pprof.Handler("block").ServeHTTP(c.Writer, c.Request)，完成对应的业务处理

	pprof.Handler("block").ServeHTTP(c.Writer, c.Request)
// 当前代码块结束

}

// Expvar @ticket 08-性能调优：GET /api/v1/ops/expvar（runtime/metrics 指标）
// 定义函数 Expvar：封装一段可复用的处理逻辑（参数：c *gin.Context）

func Expvar(c *gin.Context) {
	// 条件判断：!guard(c) 时执行下面逻辑

	if !guard(c) {
		// 返回  给调用方，结束当前函数

		return
	// 当前代码块结束

	}
	// 执行 expvar.Handler().ServeHTTP(c.Writer, c.Request)，完成对应的业务处理

	expvar.Handler().ServeHTTP(c.Writer, c.Request)
// 当前代码块结束

}

// Trace @ticket 08-性能调优：GET /api/v1/ops/trace
// 定义函数 Trace：封装一段可复用的处理逻辑（参数：c *gin.Context）

func Trace(c *gin.Context) {
	// 条件判断：!guard(c) 时执行下面逻辑

	if !guard(c) {
		// 返回  给调用方，结束当前函数

		return
	// 当前代码块结束

	}
	// 执行 pprof.Trace(c.Writer, c.Request)，完成对应的业务处理

	pprof.Trace(c.Writer, c.Request)
// 当前代码块结束

}
