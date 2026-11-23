// @ticket B端-服务入口：go run ./cmd/server 直接启动于 :18081
package main

import (
	"log"

	"smart_lockers_admin/internal/config"
	"smart_lockers_admin/internal/pkg"
	"smart_lockers_admin/internal/router"
	"smart_lockers_admin/internal/service"
	"smart_lockers_admin/internal/store"
)

// 定义函数 main：封装一段可复用的处理逻辑（参数：）

func main() {
	// 定义并初始化局部变量 cfg，取值为 config.Get()

	cfg := config.Get()
	// 定义仓储变量 st，按 db.enable 配置选择真实 GORM 后端或内存后端

	var st *store.Store
	// 条件判断：配置启用了真实数据库

	if cfg.DB.Enable {
		// 依据 DSN 与连接池参数打开数据库
		db, err := store.OpenDBWithPool(cfg.DB.DSN, cfg.DB.MaxOpenConns, cfg.DB.MaxIdleConns, cfg.DB.ConnMaxLifetimeSec)
		// 条件判断：连接或建表或种子写入失败
		if err != nil {
			// 打印告警并回退内存仓储，保证服务仍可启动
			pkg.Logf("[WARN] GORM backend init failed: %v; fallback to in-memory store", err)
			// 回退到内存仓储
			st = store.New()
			// 条件判断：初始化成功
		} else {
			// 用真实数据库句柄构造带 GORM 后端的仓储
			st, err = store.NewWithDB(db)
			// 条件判断：建表或种子写入失败
			if err != nil {
				// 打印告警并回退内存仓储
				pkg.Logf("[WARN] GORM migrate/seed failed: %v; fallback to in-memory store", err)
				// 回退到内存仓储
				st = store.New()
				// 当前代码块结束
			} else {
				// 打印真实后端初始化成功日志
				pkg.Logf("[INFO] store initialized with GORM backend (real DB) @ticket B端-真实对接")
				// 当前代码块结束
			}
			// 当前代码块结束
		}
		// 条件判断：未启用真实数据库
	} else {
		// 使用内存仓储（默认演示模式）
		st = store.New()
		// 打印内存后端日志
		pkg.Logf("[INFO] store initialized (in-memory demo seed loaded) @ticket B端-网点管理")
		// 当前代码块结束
	}

	// IoT 中台 MQTT 接入：连接 EMQX 并订阅指令回执/设备告警上行（失败降级演示模式）@ticket B端-IoT中台
	// 执行 pkg.InitMQTT(pkg.NewMqttConfig(cfg.MQTT), func(topic string, payload []byte) { service.RouteMessage(st, topic, payload) })，完成对应的业务处理

	pkg.InitMQTT(pkg.NewMqttConfig(cfg.MQTT), func(topic string, payload []byte) {
		// 转发到 IoT 中台上行消息路由
		// 执行 service.RouteMessage(st, topic, payload)，完成对应的业务处理

		service.RouteMessage(st, topic, payload)
	// 当前代码块结束

	})
	// 条件判断：pkg.MqttReady() 时执行下面逻辑

	if pkg.MqttReady() {
		// 订阅指令响应回执（方案 3.4 ACL：后端允许订阅 device/#）@ticket B端-IoT中台
		// 忽略该调用的返回值，仅执行其副作用

		_ = pkg.MqttSubscribe("device/response/#")
		// 订阅设备告警上报（方案 4.4 告警流程）
		// 忽略该调用的返回值，仅执行其副作用

		_ = pkg.MqttSubscribe("device/alert/#")
	// 当前代码块结束

	}

	// 定义并初始化局部变量 r，取值为 router.Setup(st)

	r := router.Setup(st)
	// 执行 pkg.Logf("[INFO] smart_lockers_admin starting on %s (mode=%s)", cfg.Server.Addr, cfg.Server.Mode)，完成对应的业务处理

	pkg.Logf("[INFO] smart_lockers_admin starting on %s (mode=%s)", cfg.Server.Addr, cfg.Server.Mode)
	// 条件判断：err := r.Run(cfg.Server.Addr); err != nil 时执行下面逻辑

	if err := r.Run(cfg.Server.Addr); err != nil {
		// 执行 log.Fatalf("[FATAL] server start failed: %v", err)，完成对应的业务处理

		log.Fatalf("[FATAL] server start failed: %v", err)
		// 当前代码块结束

	}
	// 当前代码块结束

}
