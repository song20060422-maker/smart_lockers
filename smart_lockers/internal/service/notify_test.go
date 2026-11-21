// @ticket 01-消息通知：模板渲染单元测试
package service

import "testing"

// 定义函数 TestRenderTemplate：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestRenderTemplate(t *testing.T) {
	// 定义并初始化局部变量 vars，取值为 map[string]string{"order_no": "SO123", "name": "张三"}

	vars := map[string]string{"order_no": "SO123", "name": "张三"}
	// 定义并初始化局部变量 cases，取值为 []struct {

	cases := []struct {
		// 执行 tpl  string，完成对应的业务处理

		tpl  string
		// 执行 want string，完成对应的业务处理

		want string
	// 执行 }{，完成对应的业务处理

	}{
		// 执行 {"订单{{order_no}}已创建", "订单SO123已创建"},，完成对应的业务处理

		{"订单{{order_no}}已创建", "订单SO123已创建"},
		// 执行 {"{{name}}，您的订单{{order_no}}超时", "张三，您的订单SO123超时"},，完成对应的业务处理

		{"{{name}}，您的订单{{order_no}}超时", "张三，您的订单SO123超时"},
		// 执行 {"无变量文本", "无变量文本"},，完成对应的业务处理

		{"无变量文本", "无变量文本"},
		// 执行 {"未知{{missing}}变量", "未知{{missing}}变量"}, // 缺失变量保留原样，完成对应的业务处理

		{"未知{{missing}}变量", "未知{{missing}}变量"}, // 缺失变量保留原样
	// 当前代码块结束

	}
	// 循环处理：for _, c := range cases，逐轮迭代

	for _, c := range cases {
		// 条件判断：got := RenderTemplate(c.tpl, vars); got != c.want 时执行下面逻辑

		if got := RenderTemplate(c.tpl, vars); got != c.want {
			// 为 t.Fatalf("RenderTemplate(%q) 赋值：%q, want %q", c.tpl, got, c.want)

			t.Fatalf("RenderTemplate(%q) = %q, want %q", c.tpl, got, c.want)
		// 当前代码块结束

		}
	// 当前代码块结束

	}
// 当前代码块结束

}

// 定义函数 TestNotifySceneCount：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestNotifySceneCount(t *testing.T) {
	// 条件判断：got := len(NotifyScene()); got != 8 时执行下面逻辑

	if got := len(NotifyScene()); got != 8 {
		// 执行 t.Fatalf("expected 8 notify scenes, got %d", got)，完成对应的业务处理

		t.Fatalf("expected 8 notify scenes, got %d", got)
	// 当前代码块结束

	}
// 当前代码块结束

}
