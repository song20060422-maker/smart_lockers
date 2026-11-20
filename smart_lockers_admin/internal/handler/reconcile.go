// @ticket B端-充值对账 handler
package handler

import (
	"github.com/gin-gonic/gin"

	"smart_lockers_admin/internal/middleware"
	"smart_lockers_admin/internal/service"
	"smart_lockers_admin/internal/pkg"
)

type ReconcileHandler struct{ svc *service.ReconcileService }

func NewReconcile(svc *service.ReconcileService) *ReconcileHandler { return &ReconcileHandler{svc: svc} }

func (h *ReconcileHandler) List(c *gin.Context) {
	respond(c, h.svc.List(middleware.User(c)), pkg.CodeOK)
}

func (h *ReconcileHandler) Run(c *gin.Context) {
	var req struct {
		BatchNo  string `json:"batch_no" binding:"required"`
		Expected int64  `json:"expected_amount"`
		Actual   int64  `json:"actual_amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, pkg.ErrBad)
		return
	}
	code, rec := h.svc.Run(middleware.User(c), req.BatchNo, req.Expected, req.Actual, idemKey(c))
	if code != pkg.CodeOK {
		pkg.Fail(c, code)
		return
	}
	pkg.OK(c, rec)
}
