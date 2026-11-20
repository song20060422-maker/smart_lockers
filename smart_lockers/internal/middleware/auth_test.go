// @ticket 00-project + 06：JWT 生成/解析/篡改校验单元测试
package middleware

import (
	"strings"
	"testing"
)

// 定义函数 TestGenerateParseToken：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestGenerateParseToken(t *testing.T) {
	// 定义并初始化局部变量 tok, err，取值为 GenerateToken(123, 2)

	tok, err := GenerateToken(123, 2)
	// 条件判断：err != nil 时执行下面逻辑

	if err != nil {
		// 执行 t.Fatalf("generate token failed: %v", err)，完成对应的业务处理

		t.Fatalf("generate token failed: %v", err)
	// 当前代码块结束

	}
	// 定义并初始化局部变量 claims, err，取值为 ParseToken(tok)

	claims, err := ParseToken(tok)
	// 条件判断：err != nil 时执行下面逻辑

	if err != nil {
		// 执行 t.Fatalf("parse token failed: %v", err)，完成对应的业务处理

		t.Fatalf("parse token failed: %v", err)
	// 当前代码块结束

	}
	// 条件判断：claims.UserID != 123 || claims.TokenVer != 2 时执行下面逻辑

	if claims.UserID != 123 || claims.TokenVer != 2 {
		// 执行 t.Fatalf("unexpected claims: %+v", claims)，完成对应的业务处理

		t.Fatalf("unexpected claims: %+v", claims)
	// 当前代码块结束

	}
// 当前代码块结束

}

// 定义函数 TestParseTokenRejectTampered：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestParseTokenRejectTampered(t *testing.T) {
	// 定义并初始化局部变量 tok, _，取值为 GenerateToken(1, 0)

	tok, _ := GenerateToken(1, 0)
	// 定义并初始化局部变量 tampered，取值为 tok[:len(tok)-2] + "xx"

	tampered := tok[:len(tok)-2] + "xx"
	// 条件判断：_, err := ParseToken(tampered); err == nil 时执行下面逻辑

	if _, err := ParseToken(tampered); err == nil {
		// 执行 t.Fatal("tampered token should be rejected")，完成对应的业务处理

		t.Fatal("tampered token should be rejected")
	// 当前代码块结束

	}
	// 条件判断：_, err := ParseToken("not.a.token"); err == nil 时执行下面逻辑

	if _, err := ParseToken("not.a.token"); err == nil {
		// 执行 t.Fatal("malformed token should be rejected")，完成对应的业务处理

		t.Fatal("malformed token should be rejected")
	// 当前代码块结束

	}
	// 条件判断：_, err := ParseToken(strings.TrimSpace("  ")); err == nil 时执行下面逻辑

	if _, err := ParseToken(strings.TrimSpace("  ")); err == nil {
		// 执行 t.Fatal("empty token should be rejected")，完成对应的业务处理

		t.Fatal("empty token should be rejected")
	// 当前代码块结束

	}
// 当前代码块结束

}

// 定义函数 TestTokenCarriesDistinctVer：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestTokenCarriesDistinctVer(t *testing.T) {
	// 定义并初始化局部变量 a, _，取值为 GenerateToken(7, 1)

	a, _ := GenerateToken(7, 1)
	// 定义并初始化局部变量 b, _，取值为 GenerateToken(7, 2)

	b, _ := GenerateToken(7, 2)
	// 条件判断：a == b 时执行下面逻辑

	if a == b {
		// 执行 t.Fatal("tokens with different tokenVer must differ")，完成对应的业务处理

		t.Fatal("tokens with different tokenVer must differ")
	// 当前代码块结束

	}
// 当前代码块结束

}
