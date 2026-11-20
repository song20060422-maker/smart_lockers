// @ticket B端-设备监控 service：状态聚合 + 远程控制(鉴权+审计)
package service

import (
	"smart_lockers_admin/internal/middleware"
	"smart_lockers_admin/internal/model"
	"smart_lockers_admin/internal/pkg"
	"smart_lockers_admin/internal/store"
)

type DeviceService struct{ st *store.Store }

func NewDevice(st *store.Store) *DeviceService { return &DeviceService{st: st} }

// List 设备列表 + 数据域
func (s *DeviceService) List(u middleware.CurrentUser) []model.Device {
	return s.st.Devices(middleware.SiteScope(u))
}

// Monitor 监控概览：在线/离线/故障统计 @ticket B端-设备监控
func (s *DeviceService) Monitor(u middleware.CurrentUser) map[string]interface{} {
	devs := s.st.Devices(middleware.SiteScope(u))
	online, offline, fault := 0, 0, 0
	for _, d := range devs {
		switch d.Status {
		case 0: online++
		case 1: offline++
		case 2: fault++
		}
	}
	return map[string]interface{}{"total": len(devs), "online": online, "offline": offline, "fault": fault}
}

// Control 设备远程控制：状态变更 + 审计 @ticket B端-设备监控/审计
func (s *DeviceService) Control(u middleware.CurrentUser, id uint, status int) int {
	if status < 0 || status > 2 {
		return pkg.ErrBad
	}
	if _, ok := s.st.Device(id); !ok {
		return pkg.ErrDeviceNF
	}
	if !s.st.UpdateDevice(id, func(d *model.Device) bool { d.Status = status; return true }) {
		return pkg.ErrDeviceNF
	}
	pkg.RecordAudit(u.Name(), string(u.Role), "device.control", idstr(id), "success")
	return pkg.CodeOK
}
