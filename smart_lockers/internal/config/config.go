// @ticket 00-project：集中式配置加载，保证 go run 直接可用
package config

import (
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	JWT      JWTConfig      `yaml:"jwt"`
	Notify   NotifyConfig   `yaml:"notify"`
	Billing  BillingConfig  `yaml:"billing"`
	Alipay   AlipayConfig   `yaml:"alipay"`
}

type ServerConfig struct {
	Addr string `yaml:"addr"`
	Mode string `yaml:"mode"`
}

type DatabaseConfig struct {
	DSN              string `yaml:"dsn"`
	MaxOpenConns     int    `yaml:"max_open_conns"`
	MaxIdleConns     int    `yaml:"max_idle_conns"`
	ConnMaxLifetimeS int    `yaml:"conn_max_lifetime_sec"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type JWTConfig struct {
	Secret      string `yaml:"secret"`
	ExpireHours int    `yaml:"expire_hours"`
}

type NotifyConfig struct {
	DedupeWindowSec int             `yaml:"dedupe_window_sec"`
	MaxRetry        int             `yaml:"max_retry"`
	Channels        map[string]bool `yaml:"channels"`
}

type BillingConfig struct {
	Currency string `yaml:"currency"`
}

// AlipayConfig 支付宝支付配置：私钥/公钥仅从被 gitignore 排除的 config.yaml 读取 @ticket 02-支付宝支付
type AlipayConfig struct {
	AppID      string `yaml:"app_id"`
	PrivateKey string `yaml:"private_key"`
	PublicKey  string `yaml:"public_key"`
	NotifyURL  string `yaml:"notify_url"`
	ReturnURL  string `yaml:"return_url"`
	Sandbox    bool   `yaml:"sandbox"`
}

var (
	once sync.Once
	cfg  *Config
)

// Load 读取 config/config.yaml；支持 CONFIG_PATH 覆盖；缺省值保证直接 go run 可启动
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
		Server:   ServerConfig{Addr: ":8080", Mode: "debug"},
		Database: DatabaseConfig{MaxOpenConns: 50, MaxIdleConns: 10, ConnMaxLifetimeS: 3600},
		JWT:      JWTConfig{Secret: "smart_lockers_local_dev_secret", ExpireHours: 24},
		Notify:   NotifyConfig{DedupeWindowSec: 60, MaxRetry: 3, Channels: map[string]bool{"inapp": true}},
		Billing:  BillingConfig{Currency: "CNY"},
		Alipay:   AlipayConfig{Sandbox: true},
	}
}

func Get() *Config { return Load() }
