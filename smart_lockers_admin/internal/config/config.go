// @ticket B端-可配置化：集中式配置加载，CONFIG_PATH 覆盖，缺省值保证 go run 直接可用
package config

import (
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server     ServerConfig   `yaml:"server"`
	JWT        JWTConfig      `yaml:"jwt"`
	DemoAdmins []DemoAdmin    `yaml:"demo_admins"`
	Policy     PolicyConfig   `yaml:"policy"`
}

type ServerConfig struct {
	Addr string `yaml:"addr"`
	Mode string `yaml:"mode"`
}

type JWTConfig struct {
	Secret      string `yaml:"secret"`
	ExpireHours int    `yaml:"expire_hours"`
}

// DemoAdmin 演示账号占位配置 @ticket B端-RBAC
type DemoAdmin struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Role     string `yaml:"role"`
	SiteID   uint   `yaml:"site_id"`
}

// PolicyConfig 可配置化业务策略 @ticket B端-可配置化
type PolicyConfig struct {
	AlarmLevels           string `yaml:"alarm_levels"`
	ReconcileToleranceFen int64  `yaml:"reconcile_tolerance_fen"`
	AuditMaxRecords       int    `yaml:"audit_max_records"`
}

var (
	once sync.Once
	cfg  *Config
)

// Load 读取 config/config.yaml；缺省值保证直接 go run 可启动
func Load() *Config {
	once.Do(func() {
		c := defaults()
		path := os.Getenv("CONFIG_PATH")
		if path == "" {
			path = "config/config.yaml"
		}
		if b, err := os.ReadFile(path); err == nil {
			_ = yaml.Unmarshal(b, c)
		}
		cfg = c
	})
	return cfg
}

func defaults() *Config {
	return &Config{
		Server: ServerConfig{Addr: ":18081", Mode: "debug"},
		JWT:    JWTConfig{Secret: "smart_lockers_admin_local_dev_secret", ExpireHours: 24},
		Policy: PolicyConfig{
			AlarmLevels:           "info,warn,critical",
			ReconcileToleranceFen: 0,
			AuditMaxRecords:       500,
		},
		DemoAdmins: []DemoAdmin{
			{Username: "operator", Password: "example", Role: "operator", SiteID: 0},
			{Username: "site_admin", Password: "example", Role: "site_admin", SiteID: 1001},
			{Username: "user", Password: "example", Role: "user", SiteID: 0},
		},
	}
}

func Get() *Config { return Load() }

// ResetForTest 仅供单测重置单例
func ResetForTest() { once = sync.Once{}; cfg = nil }
