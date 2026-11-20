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

type ReconcileService struct{ st *store.Store }

func NewReconcile(st *store.Store) *ReconcileService { return &ReconcileService{st: st} }

// List 对账记录 + 数据域
func (s *ReconcileService) List(u middleware.CurrentUser) []model.ReconcileRecord {
	return s.st.Recons(middleware.SiteScope(u))
}

// Run 执行对账：幂等键 + 事务 + 容差配置 + 审计 @ticket B端-充值对账/幂等/事务/可配置化
func (s *ReconcileService) Run(u middleware.CurrentUser, batchNo string, expected, actual int64, idemKey string) (int, model.ReconcileRecord) {
	var rec model.ReconcileRecord
	if batchNo == "" {
		return pkg.ErrBad, rec
	}
	tolerance := config.Get().Policy.ReconcileToleranceFen
	diff := actual - expected
	if diff < 0 { diff = -diff }
	status := 1
	if diff > tolerance { status = 2 }
	rec = model.ReconcileRecord{
		BatchNo: batchNo, SiteID: middleware.SiteScope(u), ExpectedAmount: expected,
		ActualAmount: actual, DiffAmount: actual - expected, Status: status,
		IdempotencyKey: idemKey, OperatorID: u.ID,
	}
	err := s.st.Tx(func() error {
		if idemKey != "" && !s.st.CheckIdemNoLock("recon:" + idemKey) {
			return store.ErrIdempotent
		}
		for _, x := range s.st.ReconsNoLock(0) {
			if x.BatchNo == batchNo {
				return store.ErrIdempotent
			}
		}
		s.st.AddReconNoLock(&rec)
		return nil
	})
	if err == store.ErrIdempotent {
		return pkg.ErrReconcileDup, model.ReconcileRecord{}
	}
	if err != nil {
		return pkg.ErrInternal, model.ReconcileRecord{}
	}
	result := "balanced"
	if status == 2 { result = "diff" }
	pkg.RecordAudit(u.Name(), string(u.Role), "reconcile.run", fmt.Sprintf("batch=%s result=%s", batchNo, result), "success")
	return pkg.CodeOK, rec
}
