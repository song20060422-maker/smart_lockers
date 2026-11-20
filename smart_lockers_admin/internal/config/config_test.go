// @ticket B端-配置单元测试
package config

import (
	"os"
	"path/filepath"
	"testing"
)

// 定义函数 TestDefaults：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestDefaults(t *testing.T) {
	// 执行 ResetForTest()，完成对应的业务处理

	ResetForTest()
	// 定义并初始化局部变量 c，取值为 defaults()

	c := defaults()
	// 条件判断：c.Server.Addr != ":18081" || c.Server.Mode != "debug" 时执行下面逻辑

	if c.Server.Addr != ":18081" || c.Server.Mode != "debug" {
		// 执行 t.Fatalf("bad server defaults: %+v", c.Server)，完成对应的业务处理

		t.Fatalf("bad server defaults: %+v", c.Server)
	// 当前代码块结束

	}
	// 条件判断：len(c.DemoAdmins) != 3 时执行下面逻辑

	if len(c.DemoAdmins) != 3 {
		// 执行 t.Fatalf("want 3 demo admins, got %d", len(c.DemoAdmins))，完成对应的业务处理

		t.Fatalf("want 3 demo admins, got %d", len(c.DemoAdmins))
	// 当前代码块结束

	}
	// 条件判断：c.JWT.ExpireHours <= 0 || c.Policy.AuditMaxRecords <= 0 时执行下面逻辑

	if c.JWT.ExpireHours <= 0 || c.Policy.AuditMaxRecords <= 0 {
		// 执行 t.Fatalf("bad policy/jwt defaults: %+v %+v", c.JWT, c.Policy)，完成对应的业务处理

		t.Fatalf("bad policy/jwt defaults: %+v %+v", c.JWT, c.Policy)
	// 当前代码块结束

	}
// 当前代码块结束

}

// 定义函数 TestLoadWithConfigFile：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestLoadWithConfigFile(t *testing.T) {
	// 执行 ResetForTest()，完成对应的业务处理

	ResetForTest()
	// 定义并初始化局部变量 dir，取值为 t.TempDir()

	dir := t.TempDir()
	// 定义并初始化局部变量 p，取值为 filepath.Join(dir, "c.yaml")

	p := filepath.Join(dir, "c.yaml")
	// 定义并初始化局部变量 body，取值为 "server:\n  addr: \":19999\"\n  mode: \"release\"\njwt:\n  secret: \"unit-test-secret\"\n  expire_hours: 2\n"

	body := "server:\n  addr: \":19999\"\n  mode: \"release\"\njwt:\n  secret: \"unit-test-secret\"\n  expire_hours: 2\n"
	// 条件判断：err := os.WriteFile(p, []byte(body), 0644); err != nil { t.Fatal(err) } 时执行下面逻辑

	if err := os.WriteFile(p, []byte(body), 0644); err != nil { t.Fatal(err) }
	// 执行 t.Setenv("CONFIG_PATH", p)，完成对应的业务处理

	t.Setenv("CONFIG_PATH", p)
	// 定义并初始化局部变量 c，取值为 Load()

	c := Load()
	// 条件判断：c.Server.Addr != ":19999" || c.Server.Mode != "release" || c.JWT.ExpireHours != 2 时执行下面逻辑

	if c.Server.Addr != ":19999" || c.Server.Mode != "release" || c.JWT.ExpireHours != 2 {
		// 执行 t.Fatalf("config file not applied: %+v", c)，完成对应的业务处理

		t.Fatalf("config file not applied: %+v", c)
	// 当前代码块结束

	}
	// 条件判断：Get() != c { t.Fatal("Get should return loaded singleton") } 时执行下面逻辑

	if Get() != c { t.Fatal("Get should return loaded singleton") }
	// 执行 ResetForTest()，完成对应的业务处理

	ResetForTest()
// 当前代码块结束

}

// 定义函数 TestLoadMissingFileUsesDefaults：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestLoadMissingFileUsesDefaults(t *testing.T) {
	// 执行 ResetForTest()，完成对应的业务处理

	ResetForTest()
	// 执行 t.Setenv("CONFIG_PATH", filepath.Join(t.TempDir(), "none.yaml"))，完成对应的业务处理

	t.Setenv("CONFIG_PATH", filepath.Join(t.TempDir(), "none.yaml"))
	// 条件判断：c := Load(); c.Server.Addr != ":18081" { t.Fatalf("defaults not used: %+v", c.Server) } 时执行下面逻辑

	if c := Load(); c.Server.Addr != ":18081" { t.Fatalf("defaults not used: %+v", c.Server) }
	// 执行 ResetForTest()，完成对应的业务处理

	ResetForTest()
// 当前代码块结束

}
