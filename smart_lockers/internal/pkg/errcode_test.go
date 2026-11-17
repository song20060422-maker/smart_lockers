// @ticket 00-project：错误码表与响应码单元测试
package pkg

import "testing"

func TestGetErr(t *testing.T) {
	if m := GetErr(CodeOK).Message; m != "ok" {
		t.Fatalf("code 0 message = %q", m)
	}
	if m := GetErr(ErrAuth).Message; m != "登录已失效" {
		t.Fatalf("ErrAuth message = %q", m)
	}
	// 未知错误码兜底
	if m := GetErr(99999).Message; m != "系统内部错误" {
		t.Fatalf("unknown code message = %q", m)
	}
}

func TestErrCodesUnique(t *testing.T) {
	codes := []int{CodeOK, ErrBad, ErrAuth, ErrForb, ErrNF, ErrInternal, ErrLocFail,
		ErrBalanceLow, ErrLockerBusy, ErrOrderNF, ErrAvatarType, ErrNotifyTpl}
	seen := map[int]bool{}
	for _, c := range codes {
		if seen[c] {
			t.Fatalf("duplicate error code: %d", c)
		}
		seen[c] = true
	}
}
