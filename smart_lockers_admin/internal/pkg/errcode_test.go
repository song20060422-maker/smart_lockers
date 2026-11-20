// @ticket B端-统一错误码单元测试
package pkg

import "testing"

// 定义函数 TestGetErr：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestGetErr(t *testing.T) {
	// 条件判断：got := GetErr(ErrForb).Message; got == "" || got == "系统内部错误" 时执行下面逻辑

	if got := GetErr(ErrForb).Message; got == "" || got == "系统内部错误" {
		// 执行 t.Fatalf("known code should have specific message, got %q", got)，完成对应的业务处理

		t.Fatalf("known code should have specific message, got %q", got)
	// 当前代码块结束

	}
	// 条件判断：got := GetErr(99999).Code; got != 99999 时执行下面逻辑

	if got := GetErr(99999).Code; got != 99999 {
		// 执行 t.Fatalf("unknown code pass-through failed: %d", got)，完成对应的业务处理

		t.Fatalf("unknown code pass-through failed: %d", got)
	// 当前代码块结束

	}
	// 条件判断：GetErr(CodeOK).Code != 0 时执行下面逻辑

	if GetErr(CodeOK).Code != 0 {
		// 执行 t.Fatal("ok code mismatch")，完成对应的业务处理

		t.Fatal("ok code mismatch")
	// 当前代码块结束

	}
// 当前代码块结束

}
