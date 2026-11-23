// @ticket 07-数字人Coze：本地演示日志工具，统一输出带时间戳的演示日志
// 定义包 service：为 HTTP / WebSocket handler 提供可复用的业务能力

package service

import (
	// 引入 fmt：格式化日志模板与参数

	"fmt"
	// 引入 log：输出带时间戳的日志到标准输出

	"log"
)

// LogDemo 输出本地演示日志，格式与标准库 log.Printf 一致
// 定义函数 LogDemo：接收格式串与可变参数并打印日志（参数：format string, args ...interface{}）

func LogDemo(format string, args ...interface{}) {
	// 调用 log.Printf 按给定格式输出日志

	log.Printf(format, args...)
	// 当前代码块结束

}

// 保留 fmt 的引用以避免未使用导入，formatDemo 供内部日志调试扩展使用
// 定义函数 formatDemo：返回格式化后的日志字符串（参数：format string, args ...interface{}）

func formatDemo(format string, args ...interface{}) string {
	// 返回 fmt.Sprintf 的格式化结果给调用方，结束当前函数

	return fmt.Sprintf(format, args...)
	// 当前代码块结束

}
