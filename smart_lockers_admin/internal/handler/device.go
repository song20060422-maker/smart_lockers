// @ticket B端-设备监控 handler
package handler

import (
	"github.com/gin-gonic/gin"

	"smart_lockers_admin/internal/middleware"
	"smart_lockers_admin/internal/service"
	"smart_lockers_admin/internal/pkg"
)

// 定义结构体 DeviceHandler：用于承载一组相关的字段数据

type DeviceHandler struct{ svc *service.DeviceService }

// 定义函数 NewDevice：封装一段可复用的处理逻辑（参数：svc *service.DeviceService）

func NewDevice(svc *service.DeviceService) *DeviceHandler { return &DeviceHandler{svc: svc} }

// 定义方法 List：接收者 h *DeviceHandler，实现该类型对应的业务行为

func (h *DeviceHandler) List(c *gin.Context) {
	// 执行 respond(c, h.svc.List(middleware.User(c)), pkg.CodeOK)，完成对应的业务处理

	respond(c, h.svc.List(middleware.User(c)), pkg.CodeOK)
// 当前代码块结束

}

// 定义方法 Monitor：接收者 h *DeviceHandler，实现该类型对应的业务行为

func (h *DeviceHandler) Monitor(c *gin.Context) {
	// 执行 respond(c, h.svc.Monitor(middleware.User(c)), pkg.CodeOK)，完成对应的业务处理

	respond(c, h.svc.Monitor(middleware.User(c)), pkg.CodeOK)
// 当前代码块结束

}

// 定义方法 Control：接收者 h *DeviceHandler，实现该类型对应的业务行为

func (h *DeviceHandler) Control(c *gin.Context) {
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
	// 执行 respond(c, nil, h.svc.Control(middleware.User(c), id, *req.Status))，完成对应的业务处理

	respond(c, nil, h.svc.Control(middleware.User(c), id, *req.Status))
// 当前代码块结束

}
