// @ticket B端-告警处理 handler
package handler

import (
	"github.com/gin-gonic/gin"

	"smart_lockers_admin/internal/middleware"
	"smart_lockers_admin/internal/service"
	"smart_lockers_admin/internal/pkg"
)

type AlarmHandler struct{ svc *service.AlarmService }

func NewAlarm(svc *service.AlarmService) *AlarmHandler { return &AlarmHandler{svc: svc} }

func (h *AlarmHandler) List(c *gin.Context) {
	respond(c, h.svc.List(middleware.User(c)), pkg.CodeOK)
}

func (h *AlarmHandler) Handle(c *gin.Context) {
	id, ok := uintParam(c, "id")
	if !ok {
		pkg.Fail(c, pkg.ErrBad)
		return
	}
	var req struct{ Note string `json:"note"` }
	_ = c.ShouldBindJSON(&req)
	code := h.svc.Handle(middleware.User(c), id, req.Note, idemKey(c))
	if code != pkg.CodeOK {
		pkg.Fail(c, code)
		return
	}
	pkg.OK(c, nil)
}
