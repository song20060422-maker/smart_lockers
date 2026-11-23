// @ticket 00-project：服务入口；go run ./cmd/server 直接启动
package main

import (
	"log"
	"runtime"

	"smart_lockers/internal/config"
	"smart_lockers/internal/model"
	"smart_lockers/internal/pkg"
	"smart_lockers/internal/router"
	"smart_lockers/internal/service"
)

// 定义函数 main：封装一段可复用的处理逻辑（参数：）

func main() {
	// 定义并初始化局部变量 cfg，取值为 config.Get()

	cfg := config.Get()

	// 数据库/Redis（连接失败仅告警，保证 go run 可启动）@ticket 00-project
	// 执行 pkg.InitDB(cfg.Database.DSN, cfg.Database.MaxOpenConns, cfg.Database.MaxIdleConns, cfg.Database.ConnMaxLifetimeS)，完成对应的业务处理

	pkg.InitDB(cfg.Database.DSN, cfg.Database.MaxOpenConns, cfg.Database.MaxIdleConns, cfg.Database.ConnMaxLifetimeS)
	// 执行 pkg.InitRedis(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)，完成对应的业务处理

	pkg.InitRedis(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	// 条件判断：err := pkg.InitAlipay(cfg.Alipay); err != nil 时执行下面逻辑

	if err := pkg.InitAlipay(cfg.Alipay); err != nil {
		// 执行 log.Printf("[WARN] alipay client init failed: %v", err)，完成对应的业务处理

		log.Printf("[WARN] alipay client init failed: %v", err)
	// 上述条件不成立时，执行 else 分支逻辑

	} else {
		// 执行 log.Printf("[INFO] alipay client ready (sandbox=%v)", cfg.Alipay.Sandbox)，完成对应的业务处理

		log.Printf("[INFO] alipay client ready (sandbox=%v)", cfg.Alipay.Sandbox)
	// 当前代码块结束

	}

	// 性能分析采样：block profile（开销<3%）@ticket 08-性能调优
	// 执行 runtime.SetBlockProfileRate(1)，完成对应的业务处理

	runtime.SetBlockProfileRate(1)
	// 执行 runtime.SetMutexProfileFraction(1)，完成对应的业务处理

	runtime.SetMutexProfileFraction(1)

	// 自动建表迁移（DB 可用时）@ticket 00-project
	// 执行 autoMigrate()，完成对应的业务处理

	autoMigrate()
	// 数字人独立账号种子：表为空时写入 admin/123456（bcrypt 哈希），确保账号真实存在于数据库 @ticket 07-数字人Coze
	// 执行 service.SeedDigitalHumanAccount()，完成对应的业务处理

	service.SeedDigitalHumanAccount()

	// IoT 设备 MQTT 接入：连接 EMQX 并订阅 device/# 上行（失败降级演示模式）@ticket IoT-MQTT设备接入
	// 执行 pkg.InitMQTT(pkg.NewMqttConfig(cfg.MQTT), service.RouteMessage)，完成对应的业务处理

	pkg.InitMQTT(pkg.NewMqttConfig(cfg.MQTT), service.RouteMessage)
	// 条件判断：pkg.MqttReady() 时执行下面逻辑

	if pkg.MqttReady() {
		// 订阅全部上行 Topic（方案 3.4 ACL：后端允许订阅 device/#）@ticket IoT-MQTT设备接入
		// 条件判断：err := pkg.MqttSubscribe("device/#", nil); err != nil 时执行下面逻辑

		if err := pkg.MqttSubscribe("device/#", nil); err != nil {
			// 执行 log.Printf("[WARN] mqtt subscribe device/# failed: %v", err)，完成对应的业务处理

			log.Printf("[WARN] mqtt subscribe device/# failed: %v", err)
		// 当前代码块结束

		}
	// 当前代码块结束

	}

	// 定义并初始化局部变量 r，取值为 router.Setup()

	r := router.Setup()
	// 执行 log.Printf("[INFO] smart_lockers server starting on %s (mode=%s)", cfg.Server.Addr, cfg.Server.Mode)，完成对应的业务处理

	log.Printf("[INFO] smart_lockers server starting on %s (mode=%s)", cfg.Server.Addr, cfg.Server.Mode)
	// 条件判断：err := r.Run(cfg.Server.Addr); err != nil 时执行下面逻辑

	if err := r.Run(cfg.Server.Addr); err != nil {
		// 执行 log.Fatalf("[FATAL] server start failed: %v", err)，完成对应的业务处理

		log.Fatalf("[FATAL] server start failed: %v", err)
	// 当前代码块结束

	}
// 当前代码块结束

}

// autoMigrate 自动迁移表结构（生产建议改用 scripts/init.sql）@ticket 00-project
// 定义函数 autoMigrate：封装一段可复用的处理逻辑（参数：）

func autoMigrate() {
	// 条件判断：pkg.DB == nil 时执行下面逻辑

	if pkg.DB == nil {
		// 返回  给调用方，结束当前函数

		return
	// 当前代码块结束

	}
	// 条件判断：err := pkg.DB.AutoMigrate(allModels()...); err != nil 时执行下面逻辑

	if err := pkg.DB.AutoMigrate(allModels()...); err != nil {
		// 执行 pkg.Logf("[WARN] auto migrate: %v", err)，完成对应的业务处理

		pkg.Logf("[WARN] auto migrate: %v", err)
	// 当前代码块结束

	}
// 当前代码块结束

}

// allModels 全量表模型 @ticket 00-project
// 定义函数 allModels：封装一段可复用的处理逻辑（参数：）

func allModels() []interface{} {
	// 返回 []interface{}{ 给调用方，结束当前函数

	return []interface{}{
		// 执行 &model.City{}, &model.LockerSite{}, &model.Guide{},，完成对应的业务处理

		&model.City{}, &model.LockerSite{}, &model.Guide{},
		// 执行 &model.LockerCabinet{}, &model.LockerDoor{}, &model.PricingRule{}, &model.DepositOrder{},，完成对应的业务处理

		&model.LockerCabinet{}, &model.LockerDoor{}, &model.PricingRule{}, &model.DepositOrder{},
		// 执行 &model.TempOpenAuth{}, &model.MakeupRecord{},，完成对应的业务处理

		&model.TempOpenAuth{}, &model.MakeupRecord{},
		// 执行 &model.Account{}, &model.BalanceLog{}, &model.RechargeOrder{},，完成对应的业务处理

		&model.Account{}, &model.BalanceLog{}, &model.RechargeOrder{},
		// 执行 &model.UserProfile{}, &model.Coupon{}, &model.MyService{},，完成对应的业务处理

		&model.UserProfile{}, &model.Coupon{}, &model.MyService{},
		// 执行 &model.NotifyTemplate{}, &model.NotifyMessage{},，完成对应的业务处理

		&model.NotifyTemplate{}, &model.NotifyMessage{},
		// IoT 设备接入表（工单 IoT-MQTT）@ticket IoT-MQTT设备接入
		// 执行 &model.IoTDevice{}, &model.DeviceToken{}, &model.DeviceCommand{},，完成对应的业务处理

		&model.IoTDevice{}, &model.DeviceToken{}, &model.DeviceCommand{},
		// 执行 &model.DeviceEvent{}, &model.OtaTask{},，完成对应的业务处理

		&model.DeviceEvent{}, &model.OtaTask{},
		// 数字人独立登录账号表（工单07）@ticket 07-数字人Coze
		// 执行 &model.DigitalHumanAccount{}，完成对应的业务处理

		&model.DigitalHumanAccount{},
	// 当前代码块结束

	}
// 当前代码块结束

}
