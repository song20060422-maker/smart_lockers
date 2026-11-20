// @ticket B端-服务入口：go run ./cmd/server 直接启动于 :18081
package main

import (
	"log"

	"smart_lockers_admin/internal/config"
	"smart_lockers_admin/internal/pkg"
	"smart_lockers_admin/internal/router"
	"smart_lockers_admin/internal/store"
)

func main() {
	cfg := config.Get()
	st := store.New()
	pkg.Logf("[INFO] store initialized (demo seed loaded) @ticket B端-网点管理")

	r := router.Setup(st)
	pkg.Logf("[INFO] smart_lockers_admin starting on %s (mode=%s)", cfg.Server.Addr, cfg.Server.Mode)
	if err := r.Run(cfg.Server.Addr); err != nil {
		log.Fatalf("[FATAL] server start failed: %v", err)
	}
}
