// @ticket 00-project：错误码表与响应码单元测试
package pkg

import "testing"

// 定义函数 TestGetErr：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestGetErr(t *testing.T) {
	// 条件判断：m := GetErr(CodeOK).Message; m != "ok" 时执行下面逻辑

	if m := GetErr(CodeOK).Message; m != "ok" {
		// 为 t.Fatalf("code 0 message 赋值：%q", m)

		t.Fatalf("code 0 message = %q", m)
	// 当前代码块结束

	}
	// 条件判断：m := GetErr(ErrAuth).Message; m != "登录已失效" 时执行下面逻辑

	if m := GetErr(ErrAuth).Message; m != "登录已失效" {
		// 为 t.Fatalf("ErrAuth message 赋值：%q", m)

		t.Fatalf("ErrAuth message = %q", m)
	// 当前代码块结束

	}
	// 未知错误码兜底
	// 条件判断：m := GetErr(99999).Message; m != "系统内部错误" 时执行下面逻辑

	if m := GetErr(99999).Message; m != "系统内部错误" {
		// 为 t.Fatalf("unknown code message 赋值：%q", m)

		t.Fatalf("unknown code message = %q", m)
	// 当前代码块结束

	}
// 当前代码块结束

}

// 定义函数 TestErrCodesUnique：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestErrCodesUnique(t *testing.T) {
	// 定义并初始化局部变量 codes，取值为 []int{CodeOK, ErrBad, ErrAuth, ErrForb, ErrNF, ErrInternal, ErrLocFail,

	codes := []int{CodeOK, ErrBad, ErrAuth, ErrForb, ErrNF, ErrInternal, ErrLocFail,
		// 执行 ErrBalanceLow, ErrLockerBusy, ErrOrderNF, ErrAvatarType, ErrNotifyTpl}，完成对应的业务处理

		ErrBalanceLow, ErrLockerBusy, ErrOrderNF, ErrAvatarType, ErrNotifyTpl}
	// 定义并初始化局部变量 seen，取值为 map[int]bool{}

	seen := map[int]bool{}
	// 循环处理：for _, c := range codes，逐轮迭代

	for _, c := range codes {
		// 条件判断：seen[c] 时执行下面逻辑

		if seen[c] {
			// 执行 t.Fatalf("duplicate error code: %d", c)，完成对应的业务处理

			t.Fatalf("duplicate error code: %d", c)
		// 当前代码块结束

		}
		// 为 seen[c] 赋值：true

		seen[c] = true
	// 当前代码块结束

	}
// 当前代码块结束

}
