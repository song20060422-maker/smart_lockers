// @ticket B端-脱敏单元测试
package pkg

import "testing"

// 定义函数 TestMaskPhone：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestMaskPhone(t *testing.T) {
	// 定义并初始化局部变量 tests，取值为 []struct{ in, want string }{

	tests := []struct{ in, want string }{
		// 执行 {"13800138000", "138****8000"},，完成对应的业务处理

		{"13800138000", "138****8000"},
		// 执行 {"12345", "***"},，完成对应的业务处理

		{"12345", "***"},
	// 当前代码块结束

	}
	// 循环处理：for _, tt := range tests，逐轮迭代

	for _, tt := range tests {
		// 条件判断：got := MaskPhone(tt.in); got != tt.want 时执行下面逻辑

		if got := MaskPhone(tt.in); got != tt.want {
			// 执行 t.Fatalf("MaskPhone(%q)=%q want %q", tt.in, got, tt.want)，完成对应的业务处理

			t.Fatalf("MaskPhone(%q)=%q want %q", tt.in, got, tt.want)
		// 当前代码块结束

		}
	// 当前代码块结束

	}
// 当前代码块结束

}

// 定义函数 TestMaskNameNo：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestMaskNameNo(t *testing.T) {
	// 条件判断：got := MaskName("张三丰"); got != "张**" 时执行下面逻辑

	if got := MaskName("张三丰"); got != "张**" {
		// 执行 t.Fatalf("MaskName got %q", got)，完成对应的业务处理

		t.Fatalf("MaskName got %q", got)
	// 当前代码块结束

	}
	// 条件判断：got := MaskName("李"); got != "李" 时执行下面逻辑

	if got := MaskName("李"); got != "李" {
		// 执行 t.Fatalf("MaskName single got %q", got)，完成对应的业务处理

		t.Fatalf("MaskName single got %q", got)
	// 当前代码块结束

	}
	// 条件判断：got := MaskNo("ABCDEFGH"); got != "AB****GH" 时执行下面逻辑

	if got := MaskNo("ABCDEFGH"); got != "AB****GH" {
		// 执行 t.Fatalf("MaskNo got %q", got)，完成对应的业务处理

		t.Fatalf("MaskNo got %q", got)
	// 当前代码块结束

	}
// 当前代码块结束

}
