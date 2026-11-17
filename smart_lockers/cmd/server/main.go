// @ticket 00-project：服务入口；go run ./cmd/server 直接启动
package main

import (
	"log"
	"runtime"

	"smart_lockers/internal/config"
	"smart_lockers/internal/model"
	"smart_lockers/internal/pkg"
	"smart_lockers/internal/router"
)

func main() {
	cfg := config.Get()

	// 数据库/Redis（连接失败仅告警，保证 go run 可启动）@ticket 00-project
	pkg.InitDB(cfg.Database.DSN, cfg.Database.MaxOpenConns, cfg.Database.MaxIdleConns, cfg.Database.ConnMaxLifetimeS)
	pkg.InitRedis(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)

	// 性能分析采样：block profile（开销<3%）@ticket 08-性能调优
	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)

	// 自动建表迁移（DB 可用时）@ticket 00-project
	autoMigrate()

	r := router.Setup()
	log.Printf("[INFO] smart_lockers server starting on %s (mode=%s)", cfg.Server.Addr, cfg.Server.Mode)
	if err := r.Run(cfg.Server.Addr); err != nil {
		log.Fatalf("[FATAL] server start failed: %v", err)
	}
}

// autoMigrate 自动迁移表结构（生产建议改用 scripts/init.sql）@ticket 00-project
func autoMigrate() {
	if pkg.DB == nil {
		return
	}
	if err := pkg.DB.AutoMigrate(allModels()...); err != nil {
		pkg.Logf("[WARN] auto migrate: %v", err)
	}
}

// allModels 全量表模型 @ticket 00-project
func allModels() []interface{} {
	return []interface{}{
		&model.City{}, &model.LockerSite{}, &model.Guide{},
		&model.LockerCabinet{}, &model.LockerDoor{}, &model.PricingRule{}, &model.DepositOrder{},
		&model.TempOpenAuth{}, &model.MakeupRecord{},
		&model.Account{}, &model.BalanceLog{}, &model.RechargeOrder{},
		&model.UserProfile{}, &model.Coupon{}, &model.MyService{},
		&model.NotifyTemplate{}, &model.NotifyMessage{},
	}
}
