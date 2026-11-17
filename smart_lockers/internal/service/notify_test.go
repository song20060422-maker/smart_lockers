// @ticket 01-消息通知：模板渲染单元测试
package service

import "testing"

func TestRenderTemplate(t *testing.T) {
	vars := map[string]string{"order_no": "SO123", "name": "张三"}
	cases := []struct {
		tpl  string
		want string
	}{
		{"订单{{order_no}}已创建", "订单SO123已创建"},
		{"{{name}}，您的订单{{order_no}}超时", "张三，您的订单SO123超时"},
		{"无变量文本", "无变量文本"},
		{"未知{{missing}}变量", "未知{{missing}}变量"}, // 缺失变量保留原样
	}
	for _, c := range cases {
		if got := RenderTemplate(c.tpl, vars); got != c.want {
			t.Fatalf("RenderTemplate(%q) = %q, want %q", c.tpl, got, c.want)
		}
	}
}

func TestNotifySceneCount(t *testing.T) {
	if got := len(NotifyScene()); got != 8 {
		t.Fatalf("expected 8 notify scenes, got %d", got)
	}
}
