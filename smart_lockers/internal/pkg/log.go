// @ticket 00-project：统一结构化日志输出
package pkg

import "log"

func Logf(format string, args ...interface{}) {
	log.Printf(format, args...)
}
