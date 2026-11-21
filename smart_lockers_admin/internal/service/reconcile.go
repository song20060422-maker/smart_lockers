// @ticket B端-充值对账 service：对账执行(事务+幂等+审计)、记录查询
package service

import (
	"fmt"

	"smart_lockers_admin/internal/config"
	"smart_lockers_admin/internal/middleware"
	"smart_lockers_admin/internal/model"
	"smart_lockers_admin/internal/pkg"
	"smart_lockers_admin/internal/store"
)

// 定义结构体 ReconcileService：用于承载一组相关的字段数据

type ReconcileService struct{ st *store.Store }

// 定义函数 NewReconcile：封装一段可复用的处理逻辑（参数：st *store.Store）

func NewReconcile(st *store.Store) *ReconcileService { return &ReconcileService{st: st} }

// List 对账记录 + 数据域
// 定义方法 List：接收者 s *ReconcileService，实现该类型对应的业务行为

func (s *ReconcileService) List(u middleware.CurrentUser) []model.ReconcileRecord {
	// 返回 s.st.Recons(middleware.SiteScope(u)) 给调用方，结束当前函数

	return s.st.Recons(middleware.SiteScope(u))
// 当前代码块结束

}

// Run 执行对账：幂等键 + 事务 + 容差配置 + 审计 @ticket B端-充值对账/幂等/事务/可配置化
// 定义方法 Run：接收者 s *ReconcileService，实现该类型对应的业务行为

func (s *ReconcileService) Run(u middleware.CurrentUser, batchNo string, expected, actual int64, idemKey string) (int, model.ReconcileRecord) {
	// 声明变量 rec model.ReconcileRecord，用于保存运行过程中的数据

	var rec model.ReconcileRecord
	// 条件判断：batchNo == "" 时执行下面逻辑

	if batchNo == "" {
		// 返回 pkg.ErrBad, rec 给调用方，结束当前函数

		return pkg.ErrBad, rec
	// 当前代码块结束

	}
	// 定义并初始化局部变量 tolerance，取值为 config.Get().Policy.ReconcileToleranceFen

	tolerance := config.Get().Policy.ReconcileToleranceFen
	// 定义并初始化局部变量 diff，取值为 actual - expected

	diff := actual - expected
	// 条件判断：diff < 0 { diff = -diff } 时执行下面逻辑

	if diff < 0 { diff = -diff }
	// 定义并初始化局部变量 status，取值为 1

	status := 1
	// 条件判断：diff > tolerance { status = 2 } 时执行下面逻辑

	if diff > tolerance { status = 2 }
	// 为 rec 赋值：model.ReconcileRecord{

	rec = model.ReconcileRecord{
		// 执行 BatchNo: batchNo, SiteID: middleware.SiteScope(u), ExpectedAmount: expected,，完成对应的业务处理

		BatchNo: batchNo, SiteID: middleware.SiteScope(u), ExpectedAmount: expected,
		// 执行 ActualAmount: actual, DiffAmount: actual - expected, Status: status,，完成对应的业务处理

		ActualAmount: actual, DiffAmount: actual - expected, Status: status,
		// 执行 IdempotencyKey: idemKey, OperatorID: u.ID,，完成对应的业务处理

		IdempotencyKey: idemKey, OperatorID: u.ID,
	// 当前代码块结束

	}
	// 定义并初始化局部变量 err，取值为 s.st.Tx(func() error {

	err := s.st.Tx(func() error {
		// 条件判断：idemKey != "" && !s.st.CheckIdemNoLock("recon:" + idemKey) 时执行下面逻辑

		if idemKey != "" && !s.st.CheckIdemNoLock("recon:" + idemKey) {
			// 返回 store.ErrIdempotent 给调用方，结束当前函数

			return store.ErrIdempotent
		// 当前代码块结束

		}
		// 循环处理：for _, x := range s.st.ReconsNoLock(0)，逐轮迭代

		for _, x := range s.st.ReconsNoLock(0) {
			// 条件判断：x.BatchNo == batchNo 时执行下面逻辑

			if x.BatchNo == batchNo {
				// 返回 store.ErrIdempotent 给调用方，结束当前函数

				return store.ErrIdempotent
			// 当前代码块结束

			}
		// 当前代码块结束

		}
		// 执行 s.st.AddReconNoLock(&rec)，完成对应的业务处理

		s.st.AddReconNoLock(&rec)
		// 返回 nil 给调用方，结束当前函数

		return nil
	// 执行 })，完成对应的业务处理

	})
	// 条件判断：err == store.ErrIdempotent 时执行下面逻辑

	if err == store.ErrIdempotent {
		// 返回 pkg.ErrReconcileDup, model.ReconcileRecord{} 给调用方，结束当前函数

		return pkg.ErrReconcileDup, model.ReconcileRecord{}
	// 当前代码块结束

	}
	// 条件判断：err != nil 时执行下面逻辑

	if err != nil {
		// 返回 pkg.ErrInternal, model.ReconcileRecord{} 给调用方，结束当前函数

		return pkg.ErrInternal, model.ReconcileRecord{}
	// 当前代码块结束

	}
	// 定义并初始化局部变量 result，取值为 "balanced"

	result := "balanced"
	// 条件判断：status == 2 { result = "diff" } 时执行下面逻辑

	if status == 2 { result = "diff" }
	// 执行 pkg.RecordAudit(u.Name(), string(u.Role), "reconcile.run", fmt.Sprintf("batch=%s result=%s", batchNo, result), "success")，完成对应的业务处理

	pkg.RecordAudit(u.Name(), string(u.Role), "reconcile.run", fmt.Sprintf("batch=%s result=%s", batchNo, result), "success")
	// 返回 pkg.CodeOK, rec 给调用方，结束当前函数

	return pkg.CodeOK, rec
// 当前代码块结束

}
