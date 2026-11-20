// @ticket B端-网点管理 handler
package handler

import (
	"github.com/gin-gonic/gin"

	"smart_lockers_admin/internal/middleware"
	"smart_lockers_admin/internal/model"
	"smart_lockers_admin/internal/service"
	"smart_lockers_admin/internal/pkg"
)

type SiteHandler struct{ svc *service.SiteService }

func NewSite(svc *service.SiteService) *SiteHandler { return &SiteHandler{svc: svc} }

func (h *SiteHandler) List(c *gin.Context) {
	respond(c, h.svc.List(middleware.User(c)), pkg.CodeOK)
}

func (h *SiteHandler) Create(c *gin.Context) {
	var v model.Site
	if err := c.ShouldBindJSON(&v); err != nil {
		pkg.Fail(c, pkg.ErrBad)
		return
	}
	respond(c, nil, h.svc.Create(middleware.User(c), &v))
}

func (h *SiteHandler) UpdateStatus(c *gin.Context) {
	id, ok := uintParam(c, "id")
	if !ok {
		pkg.Fail(c, pkg.ErrBad)
		return
	}
	var req struct{ Status *int `json:"status" binding:"required"` }
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, pkg.ErrBad)
		return
	}
	respond(c, nil, h.svc.UpdateStatus(middleware.User(c), id, *req.Status))
}
