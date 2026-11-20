// @ticket B端-网点管理 service：CRUD + 审计 + 脱敏
package service

import (
	"smart_lockers_admin/internal/middleware"
	"smart_lockers_admin/internal/model"
	"smart_lockers_admin/internal/pkg"
	"smart_lockers_admin/internal/store"
)

type SiteService struct{ st *store.Store }

func NewSite(st *store.Store) *SiteService { return &SiteService{st: st} }

// List 网点列表，按数据域过滤 + 手机号脱敏 @ticket B端-网点管理/脱敏
func (s *SiteService) List(u middleware.CurrentUser) []model.Site {
	scope := middleware.SiteScope(u)
	all := s.st.Sites()
	out := make([]model.Site, 0, len(all))
	for _, v := range all {
		if scope != 0 && v.ID != scope {
			continue
		}
		v.ManagerPhone = pkg.MaskPhone(v.ManagerPhone)
		out = append(out, v)
	}
	return out
}

// Create 新建网点：事务 + 审计 @ticket B端-网点管理/事务/审计
func (s *SiteService) Create(u middleware.CurrentUser, v *model.Site) int {
	if v.Code == "" || v.Name == "" {
		return pkg.ErrBad
	}
	for _, x := range s.st.Sites() {
		if x.Code == v.Code {
			return pkg.ErrSiteExists
		}
	}
	err := s.st.Tx(func() error {
		now := store.Now()
		v.Status, v.CreatedAt, v.UpdatedAt = 1, now, now
		s.st.AddSiteNoLock(v)
		return nil
	})
	if err != nil {
		return pkg.ErrInternal
	}
	pkg.RecordAudit(u.Name(), string(u.Role), "site.create", v.Code, "success")
	return pkg.CodeOK
}

// UpdateStatus 启停网点：状态变更鉴权 + 审计 @ticket B端-网点管理/审计
func (s *SiteService) UpdateStatus(u middleware.CurrentUser, id uint, status int) int {
	if status != 0 && status != 1 {
		return pkg.ErrBad
	}
	if !s.st.UpdateSite(id, func(v *model.Site) bool {
		v.Status = status
		return true
	}) {
		return pkg.ErrNF
	}
	pkg.RecordAudit(u.Name(), string(u.Role), "site.status", idstr(id), "success")
	return pkg.CodeOK
}
