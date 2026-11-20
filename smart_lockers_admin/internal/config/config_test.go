// @ticket B端-配置单元测试
package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaults(t *testing.T) {
	ResetForTest()
	c := defaults()
	if c.Server.Addr != ":18081" || c.Server.Mode != "debug" {
		t.Fatalf("bad server defaults: %+v", c.Server)
	}
	if len(c.DemoAdmins) != 3 {
		t.Fatalf("want 3 demo admins, got %d", len(c.DemoAdmins))
	}
	if c.JWT.ExpireHours <= 0 || c.Policy.AuditMaxRecords <= 0 {
		t.Fatalf("bad policy/jwt defaults: %+v %+v", c.JWT, c.Policy)
	}
}

func TestLoadWithConfigFile(t *testing.T) {
	ResetForTest()
	dir := t.TempDir()
	p := filepath.Join(dir, "c.yaml")
	body := "server:\n  addr: \":19999\"\n  mode: \"release\"\njwt:\n  secret: \"unit-test-secret\"\n  expire_hours: 2\n"
	if err := os.WriteFile(p, []byte(body), 0644); err != nil { t.Fatal(err) }
	t.Setenv("CONFIG_PATH", p)
	c := Load()
	if c.Server.Addr != ":19999" || c.Server.Mode != "release" || c.JWT.ExpireHours != 2 {
		t.Fatalf("config file not applied: %+v", c)
	}
	if Get() != c { t.Fatal("Get should return loaded singleton") }
	ResetForTest()
}

func TestLoadMissingFileUsesDefaults(t *testing.T) {
	ResetForTest()
	t.Setenv("CONFIG_PATH", filepath.Join(t.TempDir(), "none.yaml"))
	if c := Load(); c.Server.Addr != ":18081" { t.Fatalf("defaults not used: %+v", c.Server) }
	ResetForTest()
}
