// @ticket B端-告警处理 service：列表 + 处理(幂等+状态校验+审计)
package service

import (
	"smart_lockers_admin/internal/middleware"
	"smart_lockers_admin/internal/model"
	"smart_lockers_admin/internal/pkg"
	"smart_lockers_admin/internal/store"
)

type AlarmService struct{ st *store.Store }

func NewAlarm(st *store.Store) *AlarmService { return &AlarmService{st: st} }

// List 告警列表 + 数据域
func (s *AlarmService) List(u middleware.CurrentUser) []model.Alarm {
	return s.st.Alarms(middleware.SiteScope(u))
}

// Handle 处理告警：幂等键 + 状态校验 + 审计 @ticket B端-告警处理/幂等/审计
func (s *AlarmService) Handle(u middleware.CurrentUser, id uint, note, idemKey string) int {
	alarm, ok := s.st.Alarm(id)
	if !ok {
		return pkg.ErrAlarmNF
	}
	if alarm.Status != 0 {
		return pkg.ErrAlarmState
	}
	err := s.st.Tx(func() error {
		if idemKey != "" && !s.st.CheckIdemNoLock("alarm:" + idemKey) {
			return store.ErrIdempotent
		}
		ok := s.st.UpdateAlarmNoLock(id, func(a *model.Alarm) bool {
			if a.Status != 0 { return false }
			a.Status = 1
			a.Handler = u.Name()
			a.HandleNote = note
			now := store.Now()
			a.HandledAt = &now
			return true
		})
		if !ok { return store.ErrIdempotent }
		return nil
	})
	if err == store.ErrIdempotent {
		return pkg.ErrAlarmState
	}
	if err != nil {
		return pkg.ErrInternal
	}
	pkg.RecordAudit(u.Name(), string(u.Role), "alarm.handle", idstr(id), "success")
	return pkg.CodeOK
}
