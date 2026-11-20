// @ticket B端-统一结构化日志输出
package pkg

import "log"

// 定义函数 Logf：封装一段可复用的处理逻辑（参数：format string, args ...interface{}）

func Logf(format string, args ...interface{}) {
	// 执行 log.Printf(format, args...)，完成对应的业务处理

	log.Printf(format, args...)
// 当前代码块结束

}
