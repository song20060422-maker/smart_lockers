// @ticket B端-充值对账 handler
package handler

import (
	"github.com/gin-gonic/gin"

	"smart_lockers_admin/internal/middleware"
	"smart_lockers_admin/internal/service"
	"smart_lockers_admin/internal/pkg"
)

// 定义结构体 ReconcileHandler：用于承载一组相关的字段数据

type ReconcileHandler struct{ svc *service.ReconcileService }

// 定义函数 NewReconcile：封装一段可复用的处理逻辑（参数：svc *service.ReconcileService）

func NewReconcile(svc *service.ReconcileService) *ReconcileHandler { return &ReconcileHandler{svc: svc} }

// 定义方法 List：接收者 h *ReconcileHandler，实现该类型对应的业务行为

func (h *ReconcileHandler) List(c *gin.Context) {
	// 执行 respond(c, h.svc.List(middleware.User(c)), pkg.CodeOK)，完成对应的业务处理

	respond(c, h.svc.List(middleware.User(c)), pkg.CodeOK)
// 当前代码块结束

}

// 定义方法 Run：接收者 h *ReconcileHandler，实现该类型对应的业务行为

func (h *ReconcileHandler) Run(c *gin.Context) {
	// 声明变量 req struct {，用于保存运行过程中的数据

	var req struct {
		// 执行 BatchNo  string `json:"batch_no" binding:"required"`，完成对应的业务处理

		BatchNo  string `json:"batch_no" binding:"required"`
		// 执行 Expected int64  `json:"expected_amount"`，完成对应的业务处理

		Expected int64  `json:"expected_amount"`
		// 执行 Actual   int64  `json:"actual_amount"`，完成对应的业务处理

		Actual   int64  `json:"actual_amount"`
	// 当前代码块结束

	}
	// 条件判断：err := c.ShouldBindJSON(&req); err != nil 时执行下面逻辑

	if err := c.ShouldBindJSON(&req); err != nil {
		// 执行 pkg.Fail(c, pkg.ErrBad)，完成对应的业务处理

		pkg.Fail(c, pkg.ErrBad)
		// 返回  给调用方，结束当前函数

		return
	// 当前代码块结束

	}
	// 定义并初始化局部变量 code, rec，取值为 h.svc.Run(middleware.User(c), req.BatchNo, req.Expected, req.Actual, idemKey(c))

	code, rec := h.svc.Run(middleware.User(c), req.BatchNo, req.Expected, req.Actual, idemKey(c))
	// 条件判断：code != pkg.CodeOK 时执行下面逻辑

	if code != pkg.CodeOK {
		// 执行 pkg.Fail(c, code)，完成对应的业务处理

		pkg.Fail(c, code)
		// 返回  给调用方，结束当前函数

		return
	// 当前代码块结束

	}
	// 执行 pkg.OK(c, rec)，完成对应的业务处理

	pkg.OK(c, rec)
// 当前代码块结束

}
