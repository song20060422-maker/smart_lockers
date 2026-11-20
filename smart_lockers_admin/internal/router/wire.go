// @ticket B端-依赖装配
package router

import (
	"smart_lockers_admin/internal/service"
	"smart_lockers_admin/internal/store"
)

func serviceSite(st *store.Store) *service.SiteService     { return service.NewSite(st) }
func serviceDevice(st *store.Store) *service.DeviceService { return service.NewDevice(st) }
func serviceRecon(st *store.Store) *service.ReconcileService { return service.NewReconcile(st) }
func serviceAlarm(st *store.Store) *service.AlarmService   { return service.NewAlarm(st) }
