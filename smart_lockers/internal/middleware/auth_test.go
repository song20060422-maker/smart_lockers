// @ticket 00-project + 06：JWT 生成/解析/篡改校验单元测试
package middleware

import (
	"strings"
	"testing"
)

func TestGenerateParseToken(t *testing.T) {
	tok, err := GenerateToken(123, 2)
	if err != nil {
		t.Fatalf("generate token failed: %v", err)
	}
	claims, err := ParseToken(tok)
	if err != nil {
		t.Fatalf("parse token failed: %v", err)
	}
	if claims.UserID != 123 || claims.TokenVer != 2 {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestParseTokenRejectTampered(t *testing.T) {
	tok, _ := GenerateToken(1, 0)
	tampered := tok[:len(tok)-2] + "xx"
	if _, err := ParseToken(tampered); err == nil {
		t.Fatal("tampered token should be rejected")
	}
	if _, err := ParseToken("not.a.token"); err == nil {
		t.Fatal("malformed token should be rejected")
	}
	if _, err := ParseToken(strings.TrimSpace("  ")); err == nil {
		t.Fatal("empty token should be rejected")
	}
}

func TestTokenCarriesDistinctVer(t *testing.T) {
	a, _ := GenerateToken(7, 1)
	b, _ := GenerateToken(7, 2)
	if a == b {
		t.Fatal("tokens with different tokenVer must differ")
	}
}
