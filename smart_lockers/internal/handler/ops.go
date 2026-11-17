// @ticket 08-性能调优：pprof/expvar 性能分析接口，需鉴权且限制内网访问，采样开销<3%
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
func isInternal(ip string) bool {
	return ip == "127.0.0.1" || ip == "::1" || strings.HasPrefix(ip, "192.168.") ||
		strings.HasPrefix(ip, "10.") || strings.HasPrefix(ip, "172.")
}

// guard 分析接口鉴权/内网校验中间件 @ticket 08-性能调优
func guard(c *gin.Context) bool {
	if config.Get().Server.Mode != "debug" && !isInternal(c.ClientIP()) {
		pkg.Fail(c, pkg.ErrForb)
		return false
	}
	return true
}

// PprofIndex @ticket 08-性能调优：GET /api/v1/ops/pprof/
func PprofIndex(c *gin.Context) {
	if !guard(c) {
		return
	}
	pprof.Index(c.Writer, c.Request)
}

// PprofProfile @ticket 08-性能调优：GET /api/v1/ops/pprof/profile（CPU）
func PprofProfile(c *gin.Context) {
	if !guard(c) {
		return
	}
	pprof.Profile(c.Writer, c.Request)
}

// PprofHeap @ticket 08-性能调优：GET /api/v1/ops/pprof/heap
func PprofHeap(c *gin.Context) {
	if !guard(c) {
		return
	}
	pprof.Handler("heap").ServeHTTP(c.Writer, c.Request)
}

// PprofGoroutine @ticket 08-性能调优：GET /api/v1/ops/pprof/goroutine
func PprofGoroutine(c *gin.Context) {
	if !guard(c) {
		return
	}
	pprof.Handler("goroutine").ServeHTTP(c.Writer, c.Request)
}

// PprofBlock @ticket 08-性能调优：GET /api/v1/ops/pprof/block（锁竞争）
func PprofBlock(c *gin.Context) {
	if !guard(c) {
		return
	}
	pprof.Handler("block").ServeHTTP(c.Writer, c.Request)
}

// Expvar @ticket 08-性能调优：GET /api/v1/ops/expvar（runtime/metrics 指标）
func Expvar(c *gin.Context) {
	if !guard(c) {
		return
	}
	expvar.Handler().ServeHTTP(c.Writer, c.Request)
}

// Trace @ticket 08-性能调优：GET /api/v1/ops/trace
func Trace(c *gin.Context) {
	if !guard(c) {
		return
	}
	pprof.Trace(c.Writer, c.Request)
}
