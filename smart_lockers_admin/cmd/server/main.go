// @ticket B端-服务入口：go run ./cmd/server 直接启动于 :18081
package main

import (
	"log"

	"smart_lockers_admin/internal/config"
	"smart_lockers_admin/internal/pkg"
	"smart_lockers_admin/internal/router"
	"smart_lockers_admin/internal/store"
)

// 定义函数 main：封装一段可复用的处理逻辑（参数：）

func main() {
	// 定义并初始化局部变量 cfg，取值为 config.Get()

	cfg := config.Get()
	// 定义并初始化局部变量 st，取值为 store.New()

	st := store.New()
	// 执行 pkg.Logf("[INFO] store initialized (demo seed loaded) @ticket B端-网点管理")，完成对应的业务处理

	pkg.Logf("[INFO] store initialized (demo seed loaded) @ticket B端-网点管理")

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
