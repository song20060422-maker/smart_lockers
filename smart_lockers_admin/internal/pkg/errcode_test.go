// @ticket B端-统一错误码单元测试
package pkg

import "testing"

func TestGetErr(t *testing.T) {
	if got := GetErr(ErrForb).Message; got == "" || got == "系统内部错误" {
		t.Fatalf("known code should have specific message, got %q", got)
	}
	if got := GetErr(99999).Code; got != 99999 {
		t.Fatalf("unknown code pass-through failed: %d", got)
	}
	if GetErr(CodeOK).Code != 0 {
		t.Fatal("ok code mismatch")
	}
}
