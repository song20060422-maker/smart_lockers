// @ticket IoT-MQTT设备接入：纯函数单测（状态映射/重试策略/报文组装/报文解析）
package service

import (
	"strings"
	"testing"
)

// TestMapIoTStatusToLocal 柜格状态映射（方案 1/2/3 → 本地 0/1/2）
// 定义函数 TestMapIoTStatusToLocal：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestMapIoTStatusToLocal(t *testing.T) {
	// 定义并初始化局部变量 cases，取值为映射用例表

	cases := map[int]int{1: 0, 2: 1, 3: 2, 0: 2, 9: 2}
	// 循环处理：for in, want := range cases，逐轮迭代

	for in, want := range cases {
		// 条件判断：got := mapIoTStatusToLocal(in); got != want 时执行下面逻辑

		if got := mapIoTStatusToLocal(in); got != want {
			// 输出失败用例
			// 执行 t.Fatalf(...)，完成对应的业务处理

			t.Fatalf("mapIoTStatusToLocal(%d)=%d want=%d", in, got, want)
		// 当前代码块结束

		}
	// 当前代码块结束

	}
// 当前代码块结束

}

// TestMapOtaStatus OTA 状态字符串映射
// 定义函数 TestMapOtaStatus：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestMapOtaStatus(t *testing.T) {
	// 定义并初始化局部变量 cases，取值为映射用例表

	cases := map[string]int{"downloading": 1, "upgrading": 2, "success": 3, "failed": 4, "unknown": 0}
	// 循环处理：for in, want := range cases，逐轮迭代

	for in, want := range cases {
		// 条件判断：got := mapOtaStatus(in); got != want 时执行下面逻辑

		if got := mapOtaStatus(in); got != want {
			// 输出失败用例
			// 执行 t.Fatalf(...)，完成对应的业务处理

			t.Fatalf("mapOtaStatus(%q)=%d want=%d", in, got, want)
		// 当前代码块结束

		}
	// 当前代码块结束

	}
// 当前代码块结束

}

// TestCmdRetryPolicy 指令超时/重试策略（方案 5.2 超时策略表）
// 定义函数 TestCmdRetryPolicy：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestCmdRetryPolicy(t *testing.T) {
	// 定义并初始化局部变量 cases，取值为策略用例表：开门/关门 10s 重试 1、重启 30s 不重试、配置 60s 重试 2

	cases := map[string][2]int{"open_door": {10, 1}, "close_door": {10, 1}, "reboot": {30, 0}, "config": {60, 2}, "other": {10, 0}}
	// 循环处理：for in, want := range cases，逐轮迭代

	for in, want := range cases {
		// 定义并初始化局部变量 timeout, retry，取值为 cmdRetryPolicy(in)

		timeout, retry := cmdRetryPolicy(in)
		// 条件判断：timeout != want[0] || retry != want[1] 时执行下面逻辑

		if timeout != want[0] || retry != want[1] {
			// 输出失败用例
			// 执行 t.Fatalf(...)，完成对应的业务处理

			t.Fatalf("cmdRetryPolicy(%q)=(%d,%d) want=(%d,%d)", in, timeout, retry, want[0], want[1])
		// 当前代码块结束

		}
	// 当前代码块结束

	}
// 当前代码块结束

}

// TestBuildCommandPayload 指令报文组装（command_id/device_id/timeout 必含；扩展字段透传）
// 定义函数 TestBuildCommandPayload：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestBuildCommandPayload(t *testing.T) {
	// 组装开门指令报文
	// 定义并初始化局部变量 body，取值为 buildCommandPayload("CMD1", 1001, "A-01", "D123", 10, nil)

	body := buildCommandPayload("CMD1", 1001, "A-01", "D123", 10, nil)
	// 断言关键字段齐全
	// 条件判断：!strings.Contains(body, "\"command_id\":\"CMD1\"") 时执行下面逻辑

	if !strings.Contains(body, "\"command_id\":\"CMD1\"") {
		// 输出失败用例
		// 执行 t.Fatalf(...)，完成对应的业务处理

		t.Fatalf("payload missing command_id: %s", body)
	// 当前代码块结束

	}
	// 条件判断：!strings.Contains(body, "\"device_id\":1001") 时执行下面逻辑

	if !strings.Contains(body, "\"device_id\":1001") {
		// 输出失败用例
		// 执行 t.Fatalf(...)，完成对应的业务处理

		t.Fatalf("payload missing device_id: %s", body)
	// 当前代码块结束

	}
	// 条件判断：!strings.Contains(body, "\"locker_no\":\"A-01\"") || !strings.Contains(body, "\"order_id\":\"D123\"") 时执行下面逻辑

	if !strings.Contains(body, "\"locker_no\":\"A-01\"") || !strings.Contains(body, "\"order_id\":\"D123\"") {
		// 输出失败用例
		// 执行 t.Fatalf(...)，完成对应的业务处理

		t.Fatalf("payload missing locker_no/order_id: %s", body)
	// 当前代码块结束

	}
	// 条件判断：!strings.Contains(body, "\"timeout\":10") 时执行下面逻辑

	if !strings.Contains(body, "\"timeout\":10") {
		// 输出失败用例
		// 执行 t.Fatalf(...)，完成对应的业务处理

		t.Fatalf("payload missing timeout: %s", body)
	// 当前代码块结束

	}
// 当前代码块结束

}

// TestExtractJSONString 顶层字符串字段提取
// 定义函数 TestExtractJSONString：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestExtractJSONString(t *testing.T) {
	// 断言正常提取
	// 条件判断：got := extractJSONString(`{"alert_level":"high"}`, "alert_level"); got != "high" 时执行下面逻辑

	if got := extractJSONString(`{"alert_level":"high"}`, "alert_level"); got != "high" {
		// 输出失败用例
		// 执行 t.Fatalf(...)，完成对应的业务处理

		t.Fatalf("extract alert_level=%q want high", got)
	// 当前代码块结束

	}
	// 断言缺失字段返回空串
	// 条件判断：got := extractJSONString(`{}`, "ip_address"); got != "" 时执行下面逻辑

	if got := extractJSONString(`{}`, "ip_address"); got != "" {
		// 输出失败用例
		// 执行 t.Fatalf(...)，完成对应的业务处理

		t.Fatalf("extract missing key=%q want empty", got)
	// 当前代码块结束

	}
	// 断言非法 JSON 返回空串
	// 条件判断：got := extractJSONString("not-json", "k"); got != "" 时执行下面逻辑

	if got := extractJSONString("not-json", "k"); got != "" {
		// 输出失败用例
		// 执行 t.Fatalf(...)，完成对应的业务处理

		t.Fatalf("extract invalid json=%q want empty", got)
	// 当前代码块结束

	}
// 当前代码块结束

}
