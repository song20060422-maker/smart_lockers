// @ticket B端-告警处理 handler
package handler

import (
	"github.com/gin-gonic/gin"

	"smart_lockers_admin/internal/middleware"
	"smart_lockers_admin/internal/service"
	"smart_lockers_admin/internal/pkg"
)

// 定义结构体 AlarmHandler：用于承载一组相关的字段数据

type AlarmHandler struct{ svc *service.AlarmService }

// 定义函数 NewAlarm：封装一段可复用的处理逻辑（参数：svc *service.AlarmService）

func NewAlarm(svc *service.AlarmService) *AlarmHandler { return &AlarmHandler{svc: svc} }

// 定义方法 List：接收者 h *AlarmHandler，实现该类型对应的业务行为

func (h *AlarmHandler) List(c *gin.Context) {
	// 执行 respond(c, h.svc.List(middleware.User(c)), pkg.CodeOK)，完成对应的业务处理

	respond(c, h.svc.List(middleware.User(c)), pkg.CodeOK)
// 当前代码块结束

}

// 定义方法 Handle：接收者 h *AlarmHandler，实现该类型对应的业务行为

func (h *AlarmHandler) Handle(c *gin.Context) {
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
	// 声明变量 req struct{ Note string `json:"note"` }，用于保存运行过程中的数据

	var req struct{ Note string `json:"note"` }
	// 忽略该调用的返回值，仅执行其副作用

	_ = c.ShouldBindJSON(&req)
	// 定义并初始化局部变量 code，取值为 h.svc.Handle(middleware.User(c), id, req.Note, idemKey(c))

	code := h.svc.Handle(middleware.User(c), id, req.Note, idemKey(c))
	// 条件判断：code != pkg.CodeOK 时执行下面逻辑

	if code != pkg.CodeOK {
		// 执行 pkg.Fail(c, code)，完成对应的业务处理

		pkg.Fail(c, code)
		// 返回  给调用方，结束当前函数

		return
	// 当前代码块结束

	}
	// 执行 pkg.OK(c, nil)，完成对应的业务处理

	pkg.OK(c, nil)
// 当前代码块结束

}
