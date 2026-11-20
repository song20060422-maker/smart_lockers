// @ticket B端-设备监控 handler
package handler

import (
	"github.com/gin-gonic/gin"

	"smart_lockers_admin/internal/middleware"
	"smart_lockers_admin/internal/service"
	"smart_lockers_admin/internal/pkg"
)

type DeviceHandler struct{ svc *service.DeviceService }

func NewDevice(svc *service.DeviceService) *DeviceHandler { return &DeviceHandler{svc: svc} }

func (h *DeviceHandler) List(c *gin.Context) {
	respond(c, h.svc.List(middleware.User(c)), pkg.CodeOK)
}

func (h *DeviceHandler) Monitor(c *gin.Context) {
	respond(c, h.svc.Monitor(middleware.User(c)), pkg.CodeOK)
}

func (h *DeviceHandler) Control(c *gin.Context) {
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
	respond(c, nil, h.svc.Control(middleware.User(c), id, *req.Status))
}
