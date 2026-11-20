// @ticket B端-网点管理 service：CRUD + 审计 + 脱敏
package service

import (
	"smart_lockers_admin/internal/middleware"
	"smart_lockers_admin/internal/model"
	"smart_lockers_admin/internal/pkg"
	"smart_lockers_admin/internal/store"
)

// 定义结构体 SiteService：用于承载一组相关的字段数据

type SiteService struct{ st *store.Store }

// 定义函数 NewSite：封装一段可复用的处理逻辑（参数：st *store.Store）

func NewSite(st *store.Store) *SiteService { return &SiteService{st: st} }

// List 网点列表，按数据域过滤 + 手机号脱敏 @ticket B端-网点管理/脱敏
// 定义方法 List：接收者 s *SiteService，实现该类型对应的业务行为

func (s *SiteService) List(u middleware.CurrentUser) []model.Site {
	// 定义并初始化局部变量 scope，取值为 middleware.SiteScope(u)

	scope := middleware.SiteScope(u)
	// 定义并初始化局部变量 all，取值为 s.st.Sites()

	all := s.st.Sites()
	// 定义并初始化局部变量 out，取值为 make([]model.Site, 0, len(all))

	out := make([]model.Site, 0, len(all))
	// 循环处理：for _, v := range all，逐轮迭代

	for _, v := range all {
		// 条件判断：scope != 0 && v.ID != scope 时执行下面逻辑

		if scope != 0 && v.ID != scope {
			// 执行 continue，完成对应的业务处理

			continue
		// 当前代码块结束

		}
		// 为 v.ManagerPhone 赋值：pkg.MaskPhone(v.ManagerPhone)

		v.ManagerPhone = pkg.MaskPhone(v.ManagerPhone)
		// 为 out 赋值：append(out, v)

		out = append(out, v)
	// 当前代码块结束

	}
	// 返回 out 给调用方，结束当前函数

	return out
// 当前代码块结束

}

// Create 新建网点：事务 + 审计 @ticket B端-网点管理/事务/审计
// 定义方法 Create：接收者 s *SiteService，实现该类型对应的业务行为

func (s *SiteService) Create(u middleware.CurrentUser, v *model.Site) int {
	// 条件判断：v.Code == "" || v.Name == "" 时执行下面逻辑

	if v.Code == "" || v.Name == "" {
		// 返回 pkg.ErrBad 给调用方，结束当前函数

		return pkg.ErrBad
	// 当前代码块结束

	}
	// 循环处理：for _, x := range s.st.Sites()，逐轮迭代

	for _, x := range s.st.Sites() {
		// 条件判断：x.Code == v.Code 时执行下面逻辑

		if x.Code == v.Code {
			// 返回 pkg.ErrSiteExists 给调用方，结束当前函数

			return pkg.ErrSiteExists
		// 当前代码块结束

		}
	// 当前代码块结束

	}
	// 定义并初始化局部变量 err，取值为 s.st.Tx(func() error {

	err := s.st.Tx(func() error {
		// 定义并初始化局部变量 now，取值为 store.Now()

		now := store.Now()
		// 为 v.Status, v.CreatedAt, v.UpdatedAt 赋值：1, now, now

		v.Status, v.CreatedAt, v.UpdatedAt = 1, now, now
		// 执行 s.st.AddSiteNoLock(v)，完成对应的业务处理

		s.st.AddSiteNoLock(v)
		// 返回 nil 给调用方，结束当前函数

		return nil
	// 执行 })，完成对应的业务处理

	})
	// 条件判断：err != nil 时执行下面逻辑

	if err != nil {
		// 返回 pkg.ErrInternal 给调用方，结束当前函数

		return pkg.ErrInternal
	// 当前代码块结束

	}
	// 执行 pkg.RecordAudit(u.Name(), string(u.Role), "site.create", v.Code, "success")，完成对应的业务处理

	pkg.RecordAudit(u.Name(), string(u.Role), "site.create", v.Code, "success")
	// 返回 pkg.CodeOK 给调用方，结束当前函数

	return pkg.CodeOK
// 当前代码块结束

}

// UpdateStatus 启停网点：状态变更鉴权 + 审计 @ticket B端-网点管理/审计
// 定义方法 UpdateStatus：接收者 s *SiteService，实现该类型对应的业务行为

func (s *SiteService) UpdateStatus(u middleware.CurrentUser, id uint, status int) int {
	// 条件判断：status != 0 && status != 1 时执行下面逻辑

	if status != 0 && status != 1 {
		// 返回 pkg.ErrBad 给调用方，结束当前函数

		return pkg.ErrBad
	// 当前代码块结束

	}
	// 条件判断：!s.st.UpdateSite(id, func(v *model.Site) bool 时执行下面逻辑

	if !s.st.UpdateSite(id, func(v *model.Site) bool {
		// 为 v.Status 赋值：status

		v.Status = status
		// 返回 true 给调用方，结束当前函数

		return true
	// 执行 }) {，完成对应的业务处理

	}) {
		// 返回 pkg.ErrNF 给调用方，结束当前函数

		return pkg.ErrNF
	// 当前代码块结束

	}
	// 执行 pkg.RecordAudit(u.Name(), string(u.Role), "site.status", idstr(id), "success")，完成对应的业务处理

	pkg.RecordAudit(u.Name(), string(u.Role), "site.status", idstr(id), "success")
	// 返回 pkg.CodeOK 给调用方，结束当前函数

	return pkg.CodeOK
// 当前代码块结束

}
