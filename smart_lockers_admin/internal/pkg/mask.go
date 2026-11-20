// @ticket B端-脱敏：敏感字段输出前统一脱敏，禁止明文外泄
package pkg

import "strings"

// MaskPhone 手机号中间4位打码
func MaskPhone(s string) string {
	r := []rune(s)
	if len(r) <= 7 {
		return "***"
	}
	return string(r[:3]) + "****" + string(r[len(r)-4:])
}

// MaskName 姓名仅保留姓
func MaskName(s string) string {
	r := []rune(s)
	if len(r) == 0 {
		return ""
	}
	if len(r) == 1 {
		return string(r[0])
	}
	// 保留姓，其余每字以 * 替代（三字名"张三丰" => "张**"）
	return string(r[0]) + strings.Repeat("*", len(r)-1)
}

// MaskNo 编码/账号保留首尾各2位
func MaskNo(s string) string {
	r := []rune(s)
	if len(r) <= 4 {
		return "****"
	}
	return string(r[:2]) + "****" + string(r[len(r)-2:])
}
