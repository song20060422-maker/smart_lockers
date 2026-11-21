// @ticket B端-脱敏：敏感字段输出前统一脱敏，禁止明文外泄
package pkg

import "strings"

// MaskPhone 手机号中间4位打码
// 定义函数 MaskPhone：封装一段可复用的处理逻辑（参数：s string）

func MaskPhone(s string) string {
	// 定义并初始化局部变量 r，取值为 []rune(s)

	r := []rune(s)
	// 条件判断：len(r) <= 7 时执行下面逻辑

	if len(r) <= 7 {
		// 返回 "***" 给调用方，结束当前函数

		return "***"
	// 当前代码块结束

	}
	// 返回 string(r[:3]) + "****" + string(r[len(r)-4:]) 给调用方，结束当前函数

	return string(r[:3]) + "****" + string(r[len(r)-4:])
// 当前代码块结束

}

// MaskName 姓名仅保留姓
// 定义函数 MaskName：封装一段可复用的处理逻辑（参数：s string）

func MaskName(s string) string {
	// 定义并初始化局部变量 r，取值为 []rune(s)

	r := []rune(s)
	// 条件判断：len(r) == 0 时执行下面逻辑

	if len(r) == 0 {
		// 返回 "" 给调用方，结束当前函数

		return ""
	// 当前代码块结束

	}
	// 条件判断：len(r) == 1 时执行下面逻辑

	if len(r) == 1 {
		// 返回 string(r[0]) 给调用方，结束当前函数

		return string(r[0])
	// 当前代码块结束

	}
	// 保留姓，其余每字以 * 替代（三字名"张三丰" => "张**"）
	// 返回 string(r[0]) + strings.Repeat("*", len(r)-1) 给调用方，结束当前函数

	return string(r[0]) + strings.Repeat("*", len(r)-1)
// 当前代码块结束

}

// MaskNo 编码/账号保留首尾各2位
// 定义函数 MaskNo：封装一段可复用的处理逻辑（参数：s string）

func MaskNo(s string) string {
	// 定义并初始化局部变量 r，取值为 []rune(s)

	r := []rune(s)
	// 条件判断：len(r) <= 4 时执行下面逻辑

	if len(r) <= 4 {
		// 返回 "****" 给调用方，结束当前函数

		return "****"
	// 当前代码块结束

	}
	// 返回 string(r[:2]) + "****" + string(r[len(r)-2:]) 给调用方，结束当前函数

	return string(r[:2]) + "****" + string(r[len(r)-2:])
// 当前代码块结束

}
