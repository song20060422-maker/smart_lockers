// @ticket B端-设备监控 service：状态聚合 + 远程控制(鉴权+审计)
package service

import (
	"smart_lockers_admin/internal/middleware"
	"smart_lockers_admin/internal/model"
	"smart_lockers_admin/internal/pkg"
	"smart_lockers_admin/internal/store"
)

// 定义结构体 DeviceService：用于承载一组相关的字段数据

type DeviceService struct{ st *store.Store }

// 定义函数 NewDevice：封装一段可复用的处理逻辑（参数：st *store.Store）

func NewDevice(st *store.Store) *DeviceService { return &DeviceService{st: st} }

// List 设备列表 + 数据域
// 定义方法 List：接收者 s *DeviceService，实现该类型对应的业务行为

func (s *DeviceService) List(u middleware.CurrentUser) []model.Device {
	// 返回 s.st.Devices(middleware.SiteScope(u)) 给调用方，结束当前函数

	return s.st.Devices(middleware.SiteScope(u))
// 当前代码块结束

}

// Monitor 监控概览：在线/离线/故障统计 @ticket B端-设备监控
// 定义方法 Monitor：接收者 s *DeviceService，实现该类型对应的业务行为

func (s *DeviceService) Monitor(u middleware.CurrentUser) map[string]interface{} {
	// 定义并初始化局部变量 devs，取值为 s.st.Devices(middleware.SiteScope(u))

	devs := s.st.Devices(middleware.SiteScope(u))
	// 定义并初始化局部变量 online, offline, fault，取值为 0, 0, 0

	online, offline, fault := 0, 0, 0
	// 循环处理：for _, d := range devs，逐轮迭代

	for _, d := range devs {
		// 多分支判断：根据 d.Status 的值选择分支

		switch d.Status {
		// 当取值匹配 0 时进入该分支

		case 0: online++
		// 当取值匹配 1 时进入该分支

		case 1: offline++
		// 当取值匹配 2 时进入该分支

		case 2: fault++
		// 当前代码块结束

		}
	// 当前代码块结束

	}
	// 返回 map[string]interface{}{"total": len(devs), "online": online, "offline": offline, "fault": fault} 给调用方，结束当前函数

	return map[string]interface{}{"total": len(devs), "online": online, "offline": offline, "fault": fault}
// 当前代码块结束

}

// Control 设备远程控制：状态变更 + 审计 @ticket B端-设备监控/审计
// 定义方法 Control：接收者 s *DeviceService，实现该类型对应的业务行为

func (s *DeviceService) Control(u middleware.CurrentUser, id uint, status int) int {
	// 条件判断：status < 0 || status > 2 时执行下面逻辑

	if status < 0 || status > 2 {
		// 返回 pkg.ErrBad 给调用方，结束当前函数

		return pkg.ErrBad
	// 当前代码块结束

	}
	// 条件判断：_, ok := s.st.Device(id); !ok 时执行下面逻辑

	if _, ok := s.st.Device(id); !ok {
		// 返回 pkg.ErrDeviceNF 给调用方，结束当前函数

		return pkg.ErrDeviceNF
	// 当前代码块结束

	}
	// 条件判断：!s.st.UpdateDevice(id, func(d *model.Device) bool { d.Status = status; return true }) 时执行下面逻辑

	if !s.st.UpdateDevice(id, func(d *model.Device) bool { d.Status = status; return true }) {
		// 返回 pkg.ErrDeviceNF 给调用方，结束当前函数

		return pkg.ErrDeviceNF
	// 当前代码块结束

	}
	// 执行 pkg.RecordAudit(u.Name(), string(u.Role), "device.control", idstr(id), "success")，完成对应的业务处理

	pkg.RecordAudit(u.Name(), string(u.Role), "device.control", idstr(id), "success")
	// 返回 pkg.CodeOK 给调用方，结束当前函数

	return pkg.CodeOK
// 当前代码块结束

}
