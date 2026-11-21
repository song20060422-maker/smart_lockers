// @ticket B端-可配置化：集中式配置加载，CONFIG_PATH 覆盖，缺省值保证 go run 直接可用
package config

import (
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

// 定义结构体 Config：用于承载一组相关的字段数据

type Config struct {
	// 结构体字段 Server，类型为 ServerConfig，附带标签 `yaml:"server"`

	Server     ServerConfig   `yaml:"server"`
	// 结构体字段 JWT，类型为 JWTConfig，附带标签 `yaml:"jwt"`

	JWT        JWTConfig      `yaml:"jwt"`
	// 结构体字段 DemoAdmins，类型为 []DemoAdmin，附带标签 `yaml:"demo_admins"`

	DemoAdmins []DemoAdmin    `yaml:"demo_admins"`
	// 结构体字段 Policy，类型为 PolicyConfig，附带标签 `yaml:"policy"`

	Policy     PolicyConfig   `yaml:"policy"`
// 当前代码块结束

}

// 定义结构体 ServerConfig：用于承载一组相关的字段数据

type ServerConfig struct {
	// 结构体字段 Addr，类型为 string，附带标签 `yaml:"addr"`

	Addr string `yaml:"addr"`
	// 结构体字段 Mode，类型为 string，附带标签 `yaml:"mode"`

	Mode string `yaml:"mode"`
// 当前代码块结束

}

// 定义结构体 JWTConfig：用于承载一组相关的字段数据

type JWTConfig struct {
	// 结构体字段 Secret，类型为 string，附带标签 `yaml:"secret"`

	Secret      string `yaml:"secret"`
	// 结构体字段 ExpireHours，类型为 int，附带标签 `yaml:"expire_hours"`

	ExpireHours int    `yaml:"expire_hours"`
// 当前代码块结束

}

// DemoAdmin 演示账号占位配置 @ticket B端-RBAC
// 定义结构体 DemoAdmin：用于承载一组相关的字段数据

type DemoAdmin struct {
	// 结构体字段 Username，类型为 string，附带标签 `yaml:"username"`

	Username string `yaml:"username"`
	// 结构体字段 Password，类型为 string，附带标签 `yaml:"password"`

	Password string `yaml:"password"`
	// 结构体字段 Role，类型为 string，附带标签 `yaml:"role"`

	Role     string `yaml:"role"`
	// 结构体字段 SiteID，类型为 uint，附带标签 `yaml:"site_id"`

	SiteID   uint   `yaml:"site_id"`
// 当前代码块结束

}

// PolicyConfig 可配置化业务策略 @ticket B端-可配置化
// 定义结构体 PolicyConfig：用于承载一组相关的字段数据

type PolicyConfig struct {
	// 结构体字段 AlarmLevels，类型为 string，附带标签 `yaml:"alarm_levels"`

	AlarmLevels           string `yaml:"alarm_levels"`
	// 结构体字段 ReconcileToleranceFen，类型为 int64，附带标签 `yaml:"reconcile_tolerance_fen"`

	ReconcileToleranceFen int64  `yaml:"reconcile_tolerance_fen"`
	// 结构体字段 AuditMaxRecords，类型为 int，附带标签 `yaml:"audit_max_records"`

	AuditMaxRecords       int    `yaml:"audit_max_records"`
// 当前代码块结束

}

// 开始定义一组包级/局部变量的 var 声明块

var (
	// 执行 once sync.Once，完成对应的业务处理

	once sync.Once
	// 执行 cfg  *Config，完成对应的业务处理

	cfg  *Config
// 结束参数/表达式列表

)

// Load 读取 config/config.yaml；缺省值保证直接 go run 可启动
// 定义函数 Load：封装一段可复用的处理逻辑（参数：）

func Load() *Config {
	// 执行 once.Do(func() {，完成对应的业务处理

	once.Do(func() {
		// 定义并初始化局部变量 c，取值为 defaults()

		c := defaults()
		// 定义并初始化局部变量 path，取值为 os.Getenv("CONFIG_PATH")

		path := os.Getenv("CONFIG_PATH")
		// 条件判断：path == "" 时执行下面逻辑

		if path == "" {
			// 为 path 赋值："config/config.yaml"

			path = "config/config.yaml"
		// 当前代码块结束

		}
		// 条件判断：b, err := os.ReadFile(path); err == nil 时执行下面逻辑

		if b, err := os.ReadFile(path); err == nil {
			// 忽略该调用的返回值，仅执行其副作用

			_ = yaml.Unmarshal(b, c)
		// 当前代码块结束

		}
		// 为 cfg 赋值：c

		cfg = c
	// 执行 })，完成对应的业务处理

	})
	// 返回 cfg 给调用方，结束当前函数

	return cfg
// 当前代码块结束

}

// 定义函数 defaults：封装一段可复用的处理逻辑（参数：）

func defaults() *Config {
	// 返回 &Config{ 给调用方，结束当前函数

	return &Config{
		// 执行 Server: ServerConfig{Addr: ":18081", Mode: "debug"},，完成对应的业务处理

		Server: ServerConfig{Addr: ":18081", Mode: "debug"},
		// 执行 JWT:    JWTConfig{Secret: "smart_lockers_admin_local_dev_secret", ExpireHours: 24},，完成对应的业务处理

		JWT:    JWTConfig{Secret: "smart_lockers_admin_local_dev_secret", ExpireHours: 24},
		// 执行 Policy: PolicyConfig{，完成对应的业务处理

		Policy: PolicyConfig{
			// 执行 AlarmLevels:           "info,warn,critical",，完成对应的业务处理

			AlarmLevels:           "info,warn,critical",
			// 执行 ReconcileToleranceFen: 0,，完成对应的业务处理

			ReconcileToleranceFen: 0,
			// 执行 AuditMaxRecords:       500,，完成对应的业务处理

			AuditMaxRecords:       500,
		// 定义结束

		},
		// 执行 DemoAdmins: []DemoAdmin{，完成对应的业务处理

		DemoAdmins: []DemoAdmin{
			// 执行 {Username: "operator", Password: "example", Role: "operator", SiteID: 0},，完成对应的业务处理

			{Username: "operator", Password: "example", Role: "operator", SiteID: 0},
			// 执行 {Username: "site_admin", Password: "example", Role: "site_admin", SiteID: 1001},，完成对应的业务处理

			{Username: "site_admin", Password: "example", Role: "site_admin", SiteID: 1001},
			// 执行 {Username: "user", Password: "example", Role: "user", SiteID: 0},，完成对应的业务处理

			{Username: "user", Password: "example", Role: "user", SiteID: 0},
		// 定义结束

		},
	// 当前代码块结束

	}
// 当前代码块结束

}

// 定义函数 Get：封装一段可复用的处理逻辑（参数：）

func Get() *Config { return Load() }

// ResetForTest 仅供单测重置单例
// 定义函数 ResetForTest：封装一段可复用的处理逻辑（参数：）

func ResetForTest() { once = sync.Once{}; cfg = nil }
