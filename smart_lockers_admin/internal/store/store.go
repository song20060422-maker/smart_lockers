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
var ErrIdempotent = errors.New("idempotent hit")

// Store 内存仓储 @ticket B端-事务
type Store struct {
	mu      sync.RWMutex
	sites   map[uint]*model.Site
	devices map[uint]*model.Device
	recons  map[uint]*model.ReconcileRecord
	alarms  map[uint]*model.Alarm
	idem    map[string]bool
	nextID  uint
}

func New() *Store {
	s := &Store{
		sites: map[uint]*model.Site{}, devices: map[uint]*model.Device{},
		recons: map[uint]*model.ReconcileRecord{}, alarms: map[uint]*model.Alarm{},
		idem: map[string]bool{}, nextID: 1,
	}
	s.seed()
	return s
}

func (s *Store) allocID() uint { id := s.nextID; s.nextID++; return id } // 须持锁

// snapshot 深拷贝用于事务回滚 @ticket B端-事务（须持锁）
func (s *Store) snapshot() *Store {
	c := &Store{
		sites: map[uint]*model.Site{}, devices: map[uint]*model.Device{},
		recons: map[uint]*model.ReconcileRecord{}, alarms: map[uint]*model.Alarm{},
		idem: map[string]bool{}, nextID: s.nextID,
	}
	for k, v := range s.sites { x := *v; c.sites[k] = &x }
	for k, v := range s.devices { x := *v; c.devices[k] = &x }
	for k, v := range s.recons { x := *v; c.recons[k] = &x }
	for k, v := range s.alarms { x := *v; c.alarms[k] = &x }
	for k, v := range s.idem { c.idem[k] = v }
	return c
}

func (s *Store) restore(c *Store) { // 须持锁
	s.sites, s.devices, s.recons, s.alarms, s.idem, s.nextID = c.sites, c.devices, c.recons, c.alarms, c.idem, c.nextID
}

// Tx 事务：fn 返回错误则整体回滚。
// 注意：fn 内必须使用 *NoLock 系列方法，避免对同一 RWMutex 重复加锁造成死锁。@ticket B端-事务
func (s *Store) Tx(fn func() error) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	snap := s.snapshot()
	if err := fn(); err != nil {
		s.restore(snap)
		return err
	}
	return nil
}

// CheckIdemNoLock 幂等校验（须持锁，供 Tx 内调用）@ticket B端-幂等
func (s *Store) CheckIdemNoLock(key string) bool {
	if key == "" { return false }
	if s.idem[key] { return false }
	s.idem[key] = true
	return true
}

// CheckIdem 独立幂等校验（自行加锁）
func (s *Store) CheckIdem(key string) bool {
	s.mu.Lock(); defer s.mu.Unlock()
	return s.CheckIdemNoLock(key)
}

// ---- 网点 ----
func (s *Store) Sites() []model.Site {
	s.mu.RLock(); defer s.mu.RUnlock()
	return s.SitesNoLock()
}
func (s *Store) SitesNoLock() []model.Site {
	out := make([]model.Site, 0, len(s.sites))
	for _, v := range s.sites { out = append(out, *v) }
	return out
}
func (s *Store) Site(id uint) (model.Site, bool) {
	s.mu.RLock(); defer s.mu.RUnlock()
	return s.SiteNoLock(id)
}
func (s *Store) SiteNoLock(id uint) (model.Site, bool) {
	v, ok := s.sites[id]; if !ok { return model.Site{}, false }
	return *v, true
}
func (s *Store) AddSite(v *model.Site) {
	s.mu.Lock(); defer s.mu.Unlock()
	s.AddSiteNoLock(v)
}
func (s *Store) AddSiteNoLock(v *model.Site) { v.ID = s.allocID(); s.sites[v.ID] = v }
func (s *Store) UpdateSite(id uint, mutate func(*model.Site) bool) bool {
	s.mu.Lock(); defer s.mu.Unlock()
	return s.UpdateSiteNoLock(id, mutate)
}
func (s *Store) UpdateSiteNoLock(id uint, mutate func(*model.Site) bool) bool {
	v, ok := s.sites[id]; if !ok { return false }
	if mutate(v) { v.UpdatedAt = now(); return true }
	return false
}

// ---- 设备 ----
func (s *Store) Devices(siteID uint) []model.Device {
	s.mu.RLock(); defer s.mu.RUnlock()
	return s.DevicesNoLock(siteID)
}
func (s *Store) DevicesNoLock(siteID uint) []model.Device {
	out := []model.Device{}
	for _, v := range s.devices { if siteID == 0 || v.SiteID == siteID { out = append(out, *v) } }
	return out
}
func (s *Store) Device(id uint) (model.Device, bool) {
	s.mu.RLock(); defer s.mu.RUnlock()
	return s.DeviceNoLock(id)
}
func (s *Store) DeviceNoLock(id uint) (model.Device, bool) {
	v, ok := s.devices[id]; if !ok { return model.Device{}, false }
	return *v, true
}
func (s *Store) UpdateDevice(id uint, mutate func(*model.Device) bool) bool {
	s.mu.Lock(); defer s.mu.Unlock()
	return s.UpdateDeviceNoLock(id, mutate)
}
func (s *Store) UpdateDeviceNoLock(id uint, mutate func(*model.Device) bool) bool {
	v, ok := s.devices[id]; if !ok { return false }
	if mutate(v) { v.UpdatedAt = now(); return true }
	return false
}

// ---- 对账 ----
func (s *Store) Recons(siteID uint) []model.ReconcileRecord {
	s.mu.RLock(); defer s.mu.RUnlock()
	return s.ReconsNoLock(siteID)
}
func (s *Store) ReconsNoLock(siteID uint) []model.ReconcileRecord {
	out := []model.ReconcileRecord{}
	for _, v := range s.recons { if siteID == 0 || v.SiteID == siteID { out = append(out, *v) } }
	return out
}
func (s *Store) AddRecon(v *model.ReconcileRecord) {
	s.mu.Lock(); defer s.mu.Unlock()
	s.AddReconNoLock(v)
}
func (s *Store) AddReconNoLock(v *model.ReconcileRecord) { v.ID = s.allocID(); s.recons[v.ID] = v }

// ---- 告警 ----
func (s *Store) Alarms(siteID uint) []model.Alarm {
	s.mu.RLock(); defer s.mu.RUnlock()
	return s.AlarmsNoLock(siteID)
}
func (s *Store) AlarmsNoLock(siteID uint) []model.Alarm {
	out := []model.Alarm{}
	for _, v := range s.alarms { if siteID == 0 || v.SiteID == siteID { out = append(out, *v) } }
	return out
}
func (s *Store) Alarm(id uint) (model.Alarm, bool) {
	s.mu.RLock(); defer s.mu.RUnlock()
	return s.AlarmNoLock(id)
}
func (s *Store) AlarmNoLock(id uint) (model.Alarm, bool) {
	v, ok := s.alarms[id]; if !ok { return model.Alarm{}, false }
	return *v, true
}
func (s *Store) UpdateAlarm(id uint, mutate func(*model.Alarm) bool) bool {
	s.mu.Lock(); defer s.mu.Unlock()
	return s.UpdateAlarmNoLock(id, mutate)
}
func (s *Store) UpdateAlarmNoLock(id uint, mutate func(*model.Alarm) bool) bool {
	v, ok := s.alarms[id]; if !ok { return false }
	return mutate(v)
}

// seed 演示种子数据（占位示例，非真实业务指标）@ticket B端-网点管理
// 约定：网点ID 1001/1002 与 config 中 site_admin.site_id 对齐；设备/对账/告警 SiteID=1001
func (s *Store) seed() {
	t := now()
	s.sites[1001] = &model.Site{ID: 1001, Code: "SITE-1001", Name: "北京中关村网点",
		City: "北京", Address: "演示地址1", ManagerPhone: "13800138001", Status: 1, DoorTotal: 24, DoorFault: 0, CreatedAt: t, UpdatedAt: t}
	s.sites[1002] = &model.Site{ID: 1002, Code: "SITE-1002", Name: "上海陆家嘴网点",
		City: "上海", Address: "演示地址2", ManagerPhone: "13800138002", Status: 1, DoorTotal: 24, DoorFault: 1, CreatedAt: t, UpdatedAt: t}
	s.devices[2001] = &model.Device{ID: 2001, SiteID: 1001, Code: "CAB-001", Type: "cabinet", Status: 0, LastHeartbeat: t, UpdatedAt: t}
	s.devices[2002] = &model.Device{ID: 2002, SiteID: 1001, Code: "DOOR-001", Type: "door", Status: 2, LastHeartbeat: t, UpdatedAt: t}
	s.recons[3001] = &model.ReconcileRecord{ID: 3001, BatchNo: "RC-DEMO-001", SiteID: 1001, ExpectedAmount: 12800, ActualAmount: 12800, Status: 1, CreatedAt: t}
	s.alarms[4001] = &model.Alarm{ID: 4001, SiteID: 1001, DeviceID: 2002, Level: "warn", Content: "演示：柜门离线告警", Status: 0, CreatedAt: t}
	s.nextID = 4002
}
