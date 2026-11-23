// @ticket 00-project：集中式配置加载，保证 go run 直接可用
package config

import (
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

// 定义结构体 Config：用于承载一组相关的字段数据

type Config struct {
	// 结构体字段 Server，类型为 ServerConfig，附带标签 `yaml:"server"`

	Server   ServerConfig   `yaml:"server"`
	// 结构体字段 Database，类型为 DatabaseConfig，附带标签 `yaml:"database"`

	Database DatabaseConfig `yaml:"database"`
	// 结构体字段 Redis，类型为 RedisConfig，附带标签 `yaml:"redis"`

	Redis    RedisConfig    `yaml:"redis"`
	// 结构体字段 JWT，类型为 JWTConfig，附带标签 `yaml:"jwt"`

	JWT      JWTConfig      `yaml:"jwt"`
	// 结构体字段 Notify，类型为 NotifyConfig，附带标签 `yaml:"notify"`

	Notify   NotifyConfig   `yaml:"notify"`
	// 结构体字段 Billing，类型为 BillingConfig，附带标签 `yaml:"billing"`

	Billing  BillingConfig  `yaml:"billing"`
	// 结构体字段 Alipay，类型为 AlipayConfig，附带标签 `yaml:"alipay"`

	Alipay   AlipayConfig   `yaml:"alipay"`
	// 结构体字段 MQTT，类型为 MqttConfig，附带标签 `yaml:"mqtt"`；IoT 设备接入（工单 IoT-MQTT）

	MQTT     MqttConfig     `yaml:"mqtt"`
// 当前代码块结束

}

// MqttConfig MQTT 设备接入配置 @ticket IoT-MQTT设备接入
// Enabled=false 或 Broker 不可达时自动降级为本地演示模式，保证 go run 可启动
// 定义结构体 MqttConfig：用于承载一组相关的字段数据

type MqttConfig struct {
	// 结构体字段 Enabled，类型为 bool，附带标签 `yaml:"enabled"`；true 时连接 EMQX

	Enabled bool `yaml:"enabled"`
	// 结构体字段 Broker，类型为 string，附带标签 `yaml:"broker"`；EMQX 地址 tcp://host:1883

	Broker string `yaml:"broker"`
	// 结构体字段 ClientID，类型为 string，附带标签 `yaml:"client_id"`；后端服务客户端 ID

	ClientID string `yaml:"client_id"`
	// 结构体字段 Username，类型为 string，附带标签 `yaml:"username"`；ACL 后端服务用户名

	Username string `yaml:"username"`
	// 结构体字段 Password，类型为 string，附带标签 `yaml:"password"`；凭据仅写入 config.yaml（不入库）

	Password string `yaml:"password"`
	// 结构体字段 TokenExpireHours，类型为 int，附带标签 `yaml:"token_expire_hours"`；设备令牌有效期（方案默认 24h）

	TokenExpireHours int `yaml:"token_expire_hours"`
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

// 定义结构体 DatabaseConfig：用于承载一组相关的字段数据

type DatabaseConfig struct {
	// 结构体字段 DSN，类型为 string，附带标签 `yaml:"dsn"`

	DSN              string `yaml:"dsn"`
	// 结构体字段 MaxOpenConns，类型为 int，附带标签 `yaml:"max_open_conns"`

	MaxOpenConns     int    `yaml:"max_open_conns"`
	// 结构体字段 MaxIdleConns，类型为 int，附带标签 `yaml:"max_idle_conns"`

	MaxIdleConns     int    `yaml:"max_idle_conns"`
	// 结构体字段 ConnMaxLifetimeS，类型为 int，附带标签 `yaml:"conn_max_lifetime_sec"`

	ConnMaxLifetimeS int    `yaml:"conn_max_lifetime_sec"`
// 当前代码块结束

}

// 定义结构体 RedisConfig：用于承载一组相关的字段数据

type RedisConfig struct {
	// 结构体字段 Addr，类型为 string，附带标签 `yaml:"addr"`

	Addr     string `yaml:"addr"`
	// 结构体字段 Password，类型为 string，附带标签 `yaml:"password"`

	Password string `yaml:"password"`
	// 结构体字段 DB，类型为 int，附带标签 `yaml:"db"`

	DB       int    `yaml:"db"`
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

// 定义结构体 NotifyConfig：用于承载一组相关的字段数据

type NotifyConfig struct {
	// 结构体字段 DedupeWindowSec，类型为 int，附带标签 `yaml:"dedupe_window_sec"`

	DedupeWindowSec int             `yaml:"dedupe_window_sec"`
	// 结构体字段 MaxRetry，类型为 int，附带标签 `yaml:"max_retry"`

	MaxRetry        int             `yaml:"max_retry"`
	// 结构体字段 Channels，类型为 map[string]bool，附带标签 `yaml:"channels"`

	Channels        map[string]bool `yaml:"channels"`
// 当前代码块结束

}

// 定义结构体 BillingConfig：用于承载一组相关的字段数据

type BillingConfig struct {
	// 结构体字段 Currency，类型为 string，附带标签 `yaml:"currency"`

	Currency string `yaml:"currency"`
// 当前代码块结束

}

// AlipayConfig 支付宝支付配置：私钥/公钥仅从被 gitignore 排除的 config.yaml 读取 @ticket 02-支付宝支付
// 定义结构体 AlipayConfig：用于承载一组相关的字段数据

type AlipayConfig struct {
	// 结构体字段 AppID，类型为 string，附带标签 `yaml:"app_id"`

	AppID      string `yaml:"app_id"`
	// 结构体字段 PrivateKey，类型为 string，附带标签 `yaml:"private_key"`

	PrivateKey string `yaml:"private_key"`
	// 结构体字段 PublicKey，类型为 string，附带标签 `yaml:"public_key"`

	PublicKey  string `yaml:"public_key"`
	// 结构体字段 NotifyURL，类型为 string，附带标签 `yaml:"notify_url"`

	NotifyURL  string `yaml:"notify_url"`
	// 结构体字段 ReturnURL，类型为 string，附带标签 `yaml:"return_url"`

	ReturnURL  string `yaml:"return_url"`
	// 结构体字段 Sandbox，类型为 bool，附带标签 `yaml:"sandbox"`

	Sandbox    bool   `yaml:"sandbox"`
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

// Load 读取 config/config.yaml；支持 CONFIG_PATH 覆盖；缺省值保证直接 go run 可启动
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
		// 执行 Server:   ServerConfig{Addr: ":8080", Mode: "debug"},，完成对应的业务处理

		Server:   ServerConfig{Addr: ":8080", Mode: "debug"},
		// 执行 Database: DatabaseConfig{MaxOpenConns: 50, MaxIdleConns: 10, ConnMaxLifetimeS: 3600},，完成对应的业务处理

		Database: DatabaseConfig{MaxOpenConns: 50, MaxIdleConns: 10, ConnMaxLifetimeS: 3600},
		// 执行 JWT:      JWTConfig{Secret: "smart_lockers_local_dev_secret", ExpireHours: 24},，完成对应的业务处理

		JWT:      JWTConfig{Secret: "smart_lockers_local_dev_secret", ExpireHours: 24},
		// 执行 Notify:   NotifyConfig{DedupeWindowSec: 60, MaxRetry: 3, Channels: map[string]bool{"inapp": true}},，完成对应的业务处理

		Notify:   NotifyConfig{DedupeWindowSec: 60, MaxRetry: 3, Channels: map[string]bool{"inapp": true}},
		// 执行 Billing:  BillingConfig{Currency: "CNY"},，完成对应的业务处理

		Billing:  BillingConfig{Currency: "CNY"},
		// 执行 Alipay:   AlipayConfig{Sandbox: true},，完成对应的业务处理

		Alipay:   AlipayConfig{Sandbox: true},
		// 执行 MQTT:     MqttConfig{Enabled: false, Broker: "tcp://127.0.0.1:1883", ClientID: "smart-lockers-backend", Username: "backend-service", TokenExpireHours: 24},，完成对应的业务处理

		MQTT:     MqttConfig{Enabled: false, Broker: "tcp://127.0.0.1:1883", ClientID: "smart-lockers-backend", Username: "backend-service", TokenExpireHours: 24},
	// 当前代码块结束

	}
// 当前代码块结束

}

// 定义函数 Get：封装一段可复用的处理逻辑（参数：）

func Get() *Config { return Load() }
