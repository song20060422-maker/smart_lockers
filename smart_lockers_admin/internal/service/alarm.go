// @ticket B端-告警处理 service：列表 + 处理(幂等+状态校验+审计)
package service

import (
	"smart_lockers_admin/internal/middleware"
	"smart_lockers_admin/internal/model"
	"smart_lockers_admin/internal/pkg"
	"smart_lockers_admin/internal/store"
)

// 定义结构体 AlarmService：用于承载一组相关的字段数据

type AlarmService struct{ st *store.Store }

// 定义函数 NewAlarm：封装一段可复用的处理逻辑（参数：st *store.Store）

func NewAlarm(st *store.Store) *AlarmService { return &AlarmService{st: st} }

// List 告警列表 + 数据域
// 定义方法 List：接收者 s *AlarmService，实现该类型对应的业务行为

func (s *AlarmService) List(u middleware.CurrentUser) []model.Alarm {
	// 返回 s.st.Alarms(middleware.SiteScope(u)) 给调用方，结束当前函数

	return s.st.Alarms(middleware.SiteScope(u))
// 当前代码块结束

}

// Handle 处理告警：幂等键 + 状态校验 + 审计 @ticket B端-告警处理/幂等/审计
// 定义方法 Handle：接收者 s *AlarmService，实现该类型对应的业务行为

func (s *AlarmService) Handle(u middleware.CurrentUser, id uint, note, idemKey string) int {
	// 定义并初始化局部变量 alarm, ok，取值为 s.st.Alarm(id)

	alarm, ok := s.st.Alarm(id)
	// 条件判断：!ok 时执行下面逻辑

	if !ok {
		// 返回 pkg.ErrAlarmNF 给调用方，结束当前函数

		return pkg.ErrAlarmNF
	// 当前代码块结束

	}
	// 条件判断：alarm.Status != 0 时执行下面逻辑

	if alarm.Status != 0 {
		// 返回 pkg.ErrAlarmState 给调用方，结束当前函数

		return pkg.ErrAlarmState
	// 当前代码块结束

	}
	// 定义并初始化局部变量 err，取值为 s.st.Tx(func() error {

	err := s.st.Tx(func() error {
		// 条件判断：idemKey != "" && !s.st.CheckIdemNoLock("alarm:" + idemKey) 时执行下面逻辑

		if idemKey != "" && !s.st.CheckIdemNoLock("alarm:" + idemKey) {
			// 返回 store.ErrIdempotent 给调用方，结束当前函数

			return store.ErrIdempotent
		// 当前代码块结束

		}
		// 定义并初始化局部变量 ok，取值为 s.st.UpdateAlarmNoLock(id, func(a *model.Alarm) bool {

		ok := s.st.UpdateAlarmNoLock(id, func(a *model.Alarm) bool {
			// 条件判断：a.Status != 0 { return false } 时执行下面逻辑

			if a.Status != 0 { return false }
			// 为 a.Status 赋值：1

			a.Status = 1
			// 为 a.Handler 赋值：u.Name()

			a.Handler = u.Name()
			// 为 a.HandleNote 赋值：note

			a.HandleNote = note
			// 定义并初始化局部变量 now，取值为 store.Now()

			now := store.Now()
			// 为 a.HandledAt 赋值：&now

			a.HandledAt = &now
			// 返回 true 给调用方，结束当前函数

			return true
		// 执行 })，完成对应的业务处理

		})
		// 条件判断：!ok { return store.ErrIdempotent } 时执行下面逻辑

		if !ok { return store.ErrIdempotent }
		// 返回 nil 给调用方，结束当前函数

		return nil
	// 执行 })，完成对应的业务处理

	})
	// 条件判断：err == store.ErrIdempotent 时执行下面逻辑

	if err == store.ErrIdempotent {
		// 返回 pkg.ErrAlarmState 给调用方，结束当前函数

		return pkg.ErrAlarmState
	// 当前代码块结束

	}
	// 条件判断：err != nil 时执行下面逻辑

	if err != nil {
		// 返回 pkg.ErrInternal 给调用方，结束当前函数

		return pkg.ErrInternal
	// 当前代码块结束

	}
	// 执行 pkg.RecordAudit(u.Name(), string(u.Role), "alarm.handle", idstr(id), "success")，完成对应的业务处理

	pkg.RecordAudit(u.Name(), string(u.Role), "alarm.handle", idstr(id), "success")
	// 返回 pkg.CodeOK 给调用方，结束当前函数

	return pkg.CodeOK
// 当前代码块结束

}
