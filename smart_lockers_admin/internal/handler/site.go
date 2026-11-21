// @ticket B端-网点管理 handler
package handler

import (
	"github.com/gin-gonic/gin"

	"smart_lockers_admin/internal/middleware"
	"smart_lockers_admin/internal/model"
	"smart_lockers_admin/internal/service"
	"smart_lockers_admin/internal/pkg"
)

// 定义结构体 SiteHandler：用于承载一组相关的字段数据

type SiteHandler struct{ svc *service.SiteService }

// 定义函数 NewSite：封装一段可复用的处理逻辑（参数：svc *service.SiteService）

func NewSite(svc *service.SiteService) *SiteHandler { return &SiteHandler{svc: svc} }

// 定义方法 List：接收者 h *SiteHandler，实现该类型对应的业务行为

func (h *SiteHandler) List(c *gin.Context) {
	// 执行 respond(c, h.svc.List(middleware.User(c)), pkg.CodeOK)，完成对应的业务处理

	respond(c, h.svc.List(middleware.User(c)), pkg.CodeOK)
// 当前代码块结束

}

// 定义方法 Create：接收者 h *SiteHandler，实现该类型对应的业务行为

func (h *SiteHandler) Create(c *gin.Context) {
	// 声明变量 v model.Site，用于保存运行过程中的数据

	var v model.Site
	// 条件判断：err := c.ShouldBindJSON(&v); err != nil 时执行下面逻辑

	if err := c.ShouldBindJSON(&v); err != nil {
		// 执行 pkg.Fail(c, pkg.ErrBad)，完成对应的业务处理

		pkg.Fail(c, pkg.ErrBad)
		// 返回  给调用方，结束当前函数

		return
	// 当前代码块结束

	}
	// 执行 respond(c, nil, h.svc.Create(middleware.User(c), &v))，完成对应的业务处理

	respond(c, nil, h.svc.Create(middleware.User(c), &v))
// 当前代码块结束

}

// 定义方法 UpdateStatus：接收者 h *SiteHandler，实现该类型对应的业务行为

func (h *SiteHandler) UpdateStatus(c *gin.Context) {
	// 定义并初始化局部变量 id, ok，取值为 uintParam(c, "id")

	id, ok := uintParam(c, "id")
	// 条件判断：!ok 时执行下面逻辑

	if !ok {
		// 执行 pkg.Fail(c, pkg.ErrBad)，完成对应的业务处理

		pkg.Fail(c, pkg.ErrBad)
		// 返回  给调用方，结束当前函数

		return
	// 当前代码块结束

	}
	// 声明变量 req struct{ Status *int `json:"status" binding:"required"` }，用于保存运行过程中的数据

	var req struct{ Status *int `json:"status" binding:"required"` }
	// 条件判断：err := c.ShouldBindJSON(&req); err != nil 时执行下面逻辑

	if err := c.ShouldBindJSON(&req); err != nil {
		// 执行 pkg.Fail(c, pkg.ErrBad)，完成对应的业务处理

		pkg.Fail(c, pkg.ErrBad)
		// 返回  给调用方，结束当前函数

		return
	// 当前代码块结束

	}
	// 执行 respond(c, nil, h.svc.UpdateStatus(middleware.User(c), id, *req.Status))，完成对应的业务处理

	respond(c, nil, h.svc.UpdateStatus(middleware.User(c), id, *req.Status))
// 当前代码块结束

}
