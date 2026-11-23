// @ticket B端-依赖装配
package router

import (
	"smart_lockers_admin/internal/service"
	"smart_lockers_admin/internal/store"
)

// 定义函数 serviceSite：封装一段可复用的处理逻辑（参数：st *store.Store）

func serviceSite(st *store.Store) *service.SiteService     { return service.NewSite(st) }
// 定义函数 serviceDevice：封装一段可复用的处理逻辑（参数：st *store.Store）

func serviceDevice(st *store.Store) *service.DeviceService { return service.NewDevice(st) }
// 定义函数 serviceRecon：封装一段可复用的处理逻辑（参数：st *store.Store）

func serviceRecon(st *store.Store) *service.ReconcileService { return service.NewReconcile(st) }
// 定义函数 serviceAlarm：封装一段可复用的处理逻辑（参数：st *store.Store）

func serviceAlarm(st *store.Store) *service.AlarmService   { return service.NewAlarm(st) }
// 定义函数 serviceIot：封装一段可复用的处理逻辑（参数：st *store.Store）@ticket B端-IoT中台

func serviceIot(st *store.Store) *service.IotService       { return service.NewIot(st) }
