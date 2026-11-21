// @ticket B端-存储层：内存仓储（事务快照+幂等键），生产可替换为 GORM 事务实现
// 说明：为在没有真实DB/凭据时仍可 go run 演示B端四项功能，采用内存实现；
//      字段语义与C端一致，种子数据均为演示占位（非真实业务指标）。
package store

import (
	"errors"
	"sync"

	"smart_lockers_admin/internal/model"
)

// ErrIdempotent 幂等命中 @ticket B端-幂等
// 声明变量 ErrIdempotent = errors.New("idempotent hit")，用于保存运行过程中的数据

var ErrIdempotent = errors.New("idempotent hit")

// Store 内存仓储 @ticket B端-事务
// 定义结构体 Store：用于承载一组相关的字段数据

type Store struct {
	// 结构体字段 mu，类型为 sync.RWMutex

	mu      sync.RWMutex
	// 结构体字段 sites，类型为 map[uint]*model.Site

	sites   map[uint]*model.Site
	// 结构体字段 devices，类型为 map[uint]*model.Device

	devices map[uint]*model.Device
	// 结构体字段 recons，类型为 map[uint]*model.ReconcileRecord

	recons  map[uint]*model.ReconcileRecord
	// 结构体字段 alarms，类型为 map[uint]*model.Alarm

	alarms  map[uint]*model.Alarm
	// 结构体字段 idem，类型为 map[string]bool

	idem    map[string]bool
	// 结构体字段 nextID，类型为 uint

	nextID  uint
// 当前代码块结束

}

// 定义函数 New：封装一段可复用的处理逻辑（参数：）

func New() *Store {
	// 定义并初始化局部变量 s，取值为 &Store{

	s := &Store{
		// 执行 sites: map[uint]*model.Site{}, devices: map[uint]*model.Device{},，完成对应的业务处理

		sites: map[uint]*model.Site{}, devices: map[uint]*model.Device{},
		// 执行 recons: map[uint]*model.ReconcileRecord{}, alarms: map[uint]*model.Alarm{},，完成对应的业务处理

		recons: map[uint]*model.ReconcileRecord{}, alarms: map[uint]*model.Alarm{},
		// 执行 idem: map[string]bool{}, nextID: 1,，完成对应的业务处理

		idem: map[string]bool{}, nextID: 1,
	// 当前代码块结束

	}
	// 执行 s.seed()，完成对应的业务处理

	s.seed()
	// 返回 s 给调用方，结束当前函数

	return s
// 当前代码块结束

}

// 定义方法 allocID：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) allocID() uint { id := s.nextID; s.nextID++; return id } // 须持锁

// snapshot 深拷贝用于事务回滚 @ticket B端-事务（须持锁）
// 定义方法 snapshot：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) snapshot() *Store {
	// 定义并初始化局部变量 c，取值为 &Store{

	c := &Store{
		// 执行 sites: map[uint]*model.Site{}, devices: map[uint]*model.Device{},，完成对应的业务处理

		sites: map[uint]*model.Site{}, devices: map[uint]*model.Device{},
		// 执行 recons: map[uint]*model.ReconcileRecord{}, alarms: map[uint]*model.Alarm{},，完成对应的业务处理

		recons: map[uint]*model.ReconcileRecord{}, alarms: map[uint]*model.Alarm{},
		// 执行 idem: map[string]bool{}, nextID: s.nextID,，完成对应的业务处理

		idem: map[string]bool{}, nextID: s.nextID,
	// 当前代码块结束

	}
	// 循环处理：for k, v := range s.sites { x := *v; c.sites[k] = &x }，逐轮迭代

	for k, v := range s.sites { x := *v; c.sites[k] = &x }
	// 循环处理：for k, v := range s.devices { x := *v; c.devices[k] = &x }，逐轮迭代

	for k, v := range s.devices { x := *v; c.devices[k] = &x }
	// 循环处理：for k, v := range s.recons { x := *v; c.recons[k] = &x }，逐轮迭代

	for k, v := range s.recons { x := *v; c.recons[k] = &x }
	// 循环处理：for k, v := range s.alarms { x := *v; c.alarms[k] = &x }，逐轮迭代

	for k, v := range s.alarms { x := *v; c.alarms[k] = &x }
	// 循环处理：for k, v := range s.idem { c.idem[k] = v }，逐轮迭代

	for k, v := range s.idem { c.idem[k] = v }
	// 返回 c 给调用方，结束当前函数

	return c
// 当前代码块结束

}

// 定义方法 restore：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) restore(c *Store) { // 须持锁
	// 为 s.sites, s.devices, s.recons, s.alarms, s.idem, s.nextID 赋值：c.sites, c.devices, c.recons, c.alarms, c.idem, c.nextID

	s.sites, s.devices, s.recons, s.alarms, s.idem, s.nextID = c.sites, c.devices, c.recons, c.alarms, c.idem, c.nextID
// 当前代码块结束

}

// Tx 事务：fn 返回错误则整体回滚。
// 注意：fn 内必须使用 *NoLock 系列方法，避免对同一 RWMutex 重复加锁造成死锁。@ticket B端-事务
// 定义方法 Tx：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) Tx(fn func() error) error {
	// 执行 s.mu.Lock()，完成对应的业务处理

	s.mu.Lock()
	// 延迟注册 s.mu.Unlock()，在当前函数返回前自动执行

	defer s.mu.Unlock()
	// 定义并初始化局部变量 snap，取值为 s.snapshot()

	snap := s.snapshot()
	// 条件判断：err := fn(); err != nil 时执行下面逻辑

	if err := fn(); err != nil {
		// 执行 s.restore(snap)，完成对应的业务处理

		s.restore(snap)
		// 返回 err 给调用方，结束当前函数

		return err
	// 当前代码块结束

	}
	// 返回 nil 给调用方，结束当前函数

	return nil
// 当前代码块结束

}

// CheckIdemNoLock 幂等校验（须持锁，供 Tx 内调用）@ticket B端-幂等
// 定义方法 CheckIdemNoLock：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) CheckIdemNoLock(key string) bool {
	// 条件判断：key == "" { return false } 时执行下面逻辑

	if key == "" { return false }
	// 条件判断：s.idem[key] { return false } 时执行下面逻辑

	if s.idem[key] { return false }
	// 为 s.idem[key] 赋值：true

	s.idem[key] = true
	// 返回 true 给调用方，结束当前函数

	return true
// 当前代码块结束

}

// CheckIdem 独立幂等校验（自行加锁）
// 定义方法 CheckIdem：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) CheckIdem(key string) bool {
	// 执行 s.mu.Lock(); defer s.mu.Unlock()，完成对应的业务处理

	s.mu.Lock(); defer s.mu.Unlock()
	// 返回 s.CheckIdemNoLock(key) 给调用方，结束当前函数

	return s.CheckIdemNoLock(key)
// 当前代码块结束

}

// ---- 网点 ----
// 定义方法 Sites：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) Sites() []model.Site {
	// 执行 s.mu.RLock(); defer s.mu.RUnlock()，完成对应的业务处理

	s.mu.RLock(); defer s.mu.RUnlock()
	// 返回 s.SitesNoLock() 给调用方，结束当前函数

	return s.SitesNoLock()
// 当前代码块结束

}
// 定义方法 SitesNoLock：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) SitesNoLock() []model.Site {
	// 定义并初始化局部变量 out，取值为 make([]model.Site, 0, len(s.sites))

	out := make([]model.Site, 0, len(s.sites))
	// 循环处理：for _, v := range s.sites { out = append(out, *v) }，逐轮迭代

	for _, v := range s.sites { out = append(out, *v) }
	// 返回 out 给调用方，结束当前函数

	return out
// 当前代码块结束

}
// 定义方法 Site：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) Site(id uint) (model.Site, bool) {
	// 执行 s.mu.RLock(); defer s.mu.RUnlock()，完成对应的业务处理

	s.mu.RLock(); defer s.mu.RUnlock()
	// 返回 s.SiteNoLock(id) 给调用方，结束当前函数

	return s.SiteNoLock(id)
// 当前代码块结束

}
// 定义方法 SiteNoLock：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) SiteNoLock(id uint) (model.Site, bool) {
	// 定义并初始化局部变量 v, ok，取值为 s.sites[id]; if !ok { return model.Site{}, false }

	v, ok := s.sites[id]; if !ok { return model.Site{}, false }
	// 返回 *v, true 给调用方，结束当前函数

	return *v, true
// 当前代码块结束

}
// 定义方法 AddSite：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) AddSite(v *model.Site) {
	// 执行 s.mu.Lock(); defer s.mu.Unlock()，完成对应的业务处理

	s.mu.Lock(); defer s.mu.Unlock()
	// 执行 s.AddSiteNoLock(v)，完成对应的业务处理

	s.AddSiteNoLock(v)
// 当前代码块结束

}
// 定义方法 AddSiteNoLock：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) AddSiteNoLock(v *model.Site) { v.ID = s.allocID(); s.sites[v.ID] = v }
// 定义方法 UpdateSite：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) UpdateSite(id uint, mutate func(*model.Site) bool) bool {
	// 执行 s.mu.Lock(); defer s.mu.Unlock()，完成对应的业务处理

	s.mu.Lock(); defer s.mu.Unlock()
	// 返回 s.UpdateSiteNoLock(id, mutate) 给调用方，结束当前函数

	return s.UpdateSiteNoLock(id, mutate)
// 当前代码块结束

}
// 定义方法 UpdateSiteNoLock：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) UpdateSiteNoLock(id uint, mutate func(*model.Site) bool) bool {
	// 定义并初始化局部变量 v, ok，取值为 s.sites[id]; if !ok { return false }

	v, ok := s.sites[id]; if !ok { return false }
	// 条件判断：mutate(v) { v.UpdatedAt = now(); return true } 时执行下面逻辑

	if mutate(v) { v.UpdatedAt = now(); return true }
	// 返回 false 给调用方，结束当前函数

	return false
// 当前代码块结束

}

// ---- 设备 ----
// 定义方法 Devices：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) Devices(siteID uint) []model.Device {
	// 执行 s.mu.RLock(); defer s.mu.RUnlock()，完成对应的业务处理

	s.mu.RLock(); defer s.mu.RUnlock()
	// 返回 s.DevicesNoLock(siteID) 给调用方，结束当前函数

	return s.DevicesNoLock(siteID)
// 当前代码块结束

}
// 定义方法 DevicesNoLock：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) DevicesNoLock(siteID uint) []model.Device {
	// 定义并初始化局部变量 out，取值为 []model.Device{}

	out := []model.Device{}
	// 循环处理：for _, v := range s.devices { if siteID == 0 || v.SiteID == siteID { out = append(out, *v) } }，逐轮迭代

	for _, v := range s.devices { if siteID == 0 || v.SiteID == siteID { out = append(out, *v) } }
	// 返回 out 给调用方，结束当前函数

	return out
// 当前代码块结束

}
// 定义方法 Device：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) Device(id uint) (model.Device, bool) {
	// 执行 s.mu.RLock(); defer s.mu.RUnlock()，完成对应的业务处理

	s.mu.RLock(); defer s.mu.RUnlock()
	// 返回 s.DeviceNoLock(id) 给调用方，结束当前函数

	return s.DeviceNoLock(id)
// 当前代码块结束

}
// 定义方法 DeviceNoLock：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) DeviceNoLock(id uint) (model.Device, bool) {
	// 定义并初始化局部变量 v, ok，取值为 s.devices[id]; if !ok { return model.Device{}, false }

	v, ok := s.devices[id]; if !ok { return model.Device{}, false }
	// 返回 *v, true 给调用方，结束当前函数

	return *v, true
// 当前代码块结束

}
// 定义方法 UpdateDevice：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) UpdateDevice(id uint, mutate func(*model.Device) bool) bool {
	// 执行 s.mu.Lock(); defer s.mu.Unlock()，完成对应的业务处理

	s.mu.Lock(); defer s.mu.Unlock()
	// 返回 s.UpdateDeviceNoLock(id, mutate) 给调用方，结束当前函数

	return s.UpdateDeviceNoLock(id, mutate)
// 当前代码块结束

}
// 定义方法 UpdateDeviceNoLock：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) UpdateDeviceNoLock(id uint, mutate func(*model.Device) bool) bool {
	// 定义并初始化局部变量 v, ok，取值为 s.devices[id]; if !ok { return false }

	v, ok := s.devices[id]; if !ok { return false }
	// 条件判断：mutate(v) { v.UpdatedAt = now(); return true } 时执行下面逻辑

	if mutate(v) { v.UpdatedAt = now(); return true }
	// 返回 false 给调用方，结束当前函数

	return false
// 当前代码块结束

}

// ---- 对账 ----
// 定义方法 Recons：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) Recons(siteID uint) []model.ReconcileRecord {
	// 执行 s.mu.RLock(); defer s.mu.RUnlock()，完成对应的业务处理

	s.mu.RLock(); defer s.mu.RUnlock()
	// 返回 s.ReconsNoLock(siteID) 给调用方，结束当前函数

	return s.ReconsNoLock(siteID)
// 当前代码块结束

}
// 定义方法 ReconsNoLock：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) ReconsNoLock(siteID uint) []model.ReconcileRecord {
	// 定义并初始化局部变量 out，取值为 []model.ReconcileRecord{}

	out := []model.ReconcileRecord{}
	// 循环处理：for _, v := range s.recons { if siteID == 0 || v.SiteID == siteID { out = append(out, *v) } }，逐轮迭代

	for _, v := range s.recons { if siteID == 0 || v.SiteID == siteID { out = append(out, *v) } }
	// 返回 out 给调用方，结束当前函数

	return out
// 当前代码块结束

}
// 定义方法 AddRecon：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) AddRecon(v *model.ReconcileRecord) {
	// 执行 s.mu.Lock(); defer s.mu.Unlock()，完成对应的业务处理

	s.mu.Lock(); defer s.mu.Unlock()
	// 执行 s.AddReconNoLock(v)，完成对应的业务处理

	s.AddReconNoLock(v)
// 当前代码块结束

}
// 定义方法 AddReconNoLock：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) AddReconNoLock(v *model.ReconcileRecord) { v.ID = s.allocID(); s.recons[v.ID] = v }

// ---- 告警 ----
// 定义方法 Alarms：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) Alarms(siteID uint) []model.Alarm {
	// 执行 s.mu.RLock(); defer s.mu.RUnlock()，完成对应的业务处理

	s.mu.RLock(); defer s.mu.RUnlock()
	// 返回 s.AlarmsNoLock(siteID) 给调用方，结束当前函数

	return s.AlarmsNoLock(siteID)
// 当前代码块结束

}
// 定义方法 AlarmsNoLock：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) AlarmsNoLock(siteID uint) []model.Alarm {
	// 定义并初始化局部变量 out，取值为 []model.Alarm{}

	out := []model.Alarm{}
	// 循环处理：for _, v := range s.alarms { if siteID == 0 || v.SiteID == siteID { out = append(out, *v) } }，逐轮迭代

	for _, v := range s.alarms { if siteID == 0 || v.SiteID == siteID { out = append(out, *v) } }
	// 返回 out 给调用方，结束当前函数

	return out
// 当前代码块结束

}
// 定义方法 Alarm：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) Alarm(id uint) (model.Alarm, bool) {
	// 执行 s.mu.RLock(); defer s.mu.RUnlock()，完成对应的业务处理

	s.mu.RLock(); defer s.mu.RUnlock()
	// 返回 s.AlarmNoLock(id) 给调用方，结束当前函数

	return s.AlarmNoLock(id)
// 当前代码块结束

}
// 定义方法 AlarmNoLock：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) AlarmNoLock(id uint) (model.Alarm, bool) {
	// 定义并初始化局部变量 v, ok，取值为 s.alarms[id]; if !ok { return model.Alarm{}, false }

	v, ok := s.alarms[id]; if !ok { return model.Alarm{}, false }
	// 返回 *v, true 给调用方，结束当前函数

	return *v, true
// 当前代码块结束

}
// 定义方法 UpdateAlarm：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) UpdateAlarm(id uint, mutate func(*model.Alarm) bool) bool {
	// 执行 s.mu.Lock(); defer s.mu.Unlock()，完成对应的业务处理

	s.mu.Lock(); defer s.mu.Unlock()
	// 返回 s.UpdateAlarmNoLock(id, mutate) 给调用方，结束当前函数

	return s.UpdateAlarmNoLock(id, mutate)
// 当前代码块结束

}
// 定义方法 UpdateAlarmNoLock：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) UpdateAlarmNoLock(id uint, mutate func(*model.Alarm) bool) bool {
	// 定义并初始化局部变量 v, ok，取值为 s.alarms[id]; if !ok { return false }

	v, ok := s.alarms[id]; if !ok { return false }
	// 返回 mutate(v) 给调用方，结束当前函数

	return mutate(v)
// 当前代码块结束

}

// seed 演示种子数据（占位示例，非真实业务指标）@ticket B端-网点管理
// 约定：网点ID 1001/1002 与 config 中 site_admin.site_id 对齐；设备/对账/告警 SiteID=1001
// 定义方法 seed：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) seed() {
	// 定义并初始化局部变量 t，取值为 now()

	t := now()
	// 为 s.sites[1001] 赋值：&model.Site{ID: 1001, Code: "SITE-1001", Name: "北京中关村网点",

	s.sites[1001] = &model.Site{ID: 1001, Code: "SITE-1001", Name: "北京中关村网点",
		// 执行 City: "北京", Address: "演示地址1", ManagerPhone: "13800138001", Status: 1, DoorTotal: 24, DoorFault: 0, CreatedAt: t, UpdatedAt: t}，完成对应的业务处理

		City: "北京", Address: "演示地址1", ManagerPhone: "13800138001", Status: 1, DoorTotal: 24, DoorFault: 0, CreatedAt: t, UpdatedAt: t}
	// 为 s.sites[1002] 赋值：&model.Site{ID: 1002, Code: "SITE-1002", Name: "上海陆家嘴网点",

	s.sites[1002] = &model.Site{ID: 1002, Code: "SITE-1002", Name: "上海陆家嘴网点",
		// 执行 City: "上海", Address: "演示地址2", ManagerPhone: "13800138002", Status: 1, DoorTotal: 24, DoorFault: 1, CreatedAt: t, UpdatedAt: t}，完成对应的业务处理

		City: "上海", Address: "演示地址2", ManagerPhone: "13800138002", Status: 1, DoorTotal: 24, DoorFault: 1, CreatedAt: t, UpdatedAt: t}
	// 为 s.devices[2001] 赋值：&model.Device{ID: 2001, SiteID: 1001, Code: "CAB-001", Type: "cabinet", Status: 0, LastHeartbeat: t, UpdatedAt: t}

	s.devices[2001] = &model.Device{ID: 2001, SiteID: 1001, Code: "CAB-001", Type: "cabinet", Status: 0, LastHeartbeat: t, UpdatedAt: t}
	// 为 s.devices[2002] 赋值：&model.Device{ID: 2002, SiteID: 1001, Code: "DOOR-001", Type: "door", Status: 2, LastHeartbeat: t, UpdatedAt: t}

	s.devices[2002] = &model.Device{ID: 2002, SiteID: 1001, Code: "DOOR-001", Type: "door", Status: 2, LastHeartbeat: t, UpdatedAt: t}
	// 为 s.recons[3001] 赋值：&model.ReconcileRecord{ID: 3001, BatchNo: "RC-DEMO-001", SiteID: 1001, ExpectedAmount: 12800, ActualAmount: 12800, Status: 1, CreatedAt: t}

	s.recons[3001] = &model.ReconcileRecord{ID: 3001, BatchNo: "RC-DEMO-001", SiteID: 1001, ExpectedAmount: 12800, ActualAmount: 12800, Status: 1, CreatedAt: t}
	// 为 s.alarms[4001] 赋值：&model.Alarm{ID: 4001, SiteID: 1001, DeviceID: 2002, Level: "warn", Content: "演示：柜门离线告警", Status: 0, CreatedAt: t}

	s.alarms[4001] = &model.Alarm{ID: 4001, SiteID: 1001, DeviceID: 2002, Level: "warn", Content: "演示：柜门离线告警", Status: 0, CreatedAt: t}
	// 为 s.nextID 赋值：4002

	s.nextID = 4002
// 当前代码块结束

}
