// @ticket B端-IoT中台 handler：总览/指令下发与记录/OTA 任务
package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"smart_lockers_admin/internal/middleware"
	"smart_lockers_admin/internal/pkg"
	"smart_lockers_admin/internal/service"
)

// 定义结构体 IotHandler：用于承载一组相关的字段数据

type IotHandler struct{ svc *service.IotService }

// 定义函数 NewIot：封装一段可复用的处理逻辑（参数：svc *service.IotService）

func NewIot(svc *service.IotService) *IotHandler { return &IotHandler{svc: svc} }

// Overview IoT 中台总览 @ticket B端-IoT中台
// 定义方法 Overview：接收者 h *IotHandler，实现该类型对应的业务行为

func (h *IotHandler) Overview(c *gin.Context) {
	// 执行 respond(c, h.svc.Overview(middleware.User(c)), pkg.CodeOK)，完成对应的业务处理

	respond(c, h.svc.Overview(middleware.User(c)), pkg.CodeOK)
// 当前代码块结束

}

// Dispatch 下发控制指令 @ticket B端-IoT中台
// 定义方法 Dispatch：接收者 h *IotHandler，实现该类型对应的业务行为

func (h *IotHandler) Dispatch(c *gin.Context) {
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
	// 声明变量 req，用于保存运行过程中的数据

	var req struct {
		// 结构体字段 CmdType，类型为 string，附带标签 `json:"cmd_type" binding:"required"`

		CmdType  string `json:"cmd_type" binding:"required"`
		// 结构体字段 LockerNo，类型为 string，附带标签 `json:"locker_no"`

		LockerNo string `json:"locker_no"`
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
	// 执行 respond(c, h.svc.DispatchCommand(middleware.User(c), id, req.CmdType, req.LockerNo))，完成对应的业务处理

	data, code := h.svc.DispatchCommand(middleware.User(c), id, req.CmdType, req.LockerNo)
	// 统一响应
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}

// CommandRecords 指令响应记录 @ticket B端-IoT中台
// 定义方法 CommandRecords：接收者 h *IotHandler，实现该类型对应的业务行为

func (h *IotHandler) CommandRecords(c *gin.Context) {
	// 解析可选 device_id 查询参数
	// 定义并初始化局部变量 deviceID，取值为 uintQuery(c, "device_id")

	deviceID := uintQuery(c, "device_id")
	// 执行 respond(c, h.svc.CommandRecords(middleware.User(c), deviceID), pkg.CodeOK)，完成对应的业务处理

	respond(c, h.svc.CommandRecords(middleware.User(c), deviceID), pkg.CodeOK)
// 当前代码块结束

}

// OtaCreate 新建 OTA 升级任务 @ticket B端-IoT中台
// 定义方法 OtaCreate：接收者 h *IotHandler，实现该类型对应的业务行为

func (h *IotHandler) OtaCreate(c *gin.Context) {
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
	// 声明变量 req，用于保存运行过程中的数据

	var req struct {
		// 结构体字段 FirmwareVersion，类型为 string，附带标签 `json:"firmware_version" binding:"required"`

		FirmwareVersion string `json:"firmware_version" binding:"required"`
		// 结构体字段 DownloadURL，类型为 string，附带标签 `json:"download_url" binding:"required"`

		DownloadURL     string `json:"download_url" binding:"required"`
		// 结构体字段 MD5，类型为 string，附带标签 `json:"md5"`

		MD5             string `json:"md5"`
		// 结构体字段 FileSize，类型为 int64，附带标签 `json:"file_size"`

		FileSize        int64  `json:"file_size"`
		// 结构体字段 ForceUpgrade，类型为 bool，附带标签 `json:"force_upgrade"`

		ForceUpgrade    bool   `json:"force_upgrade"`
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
	// 调用服务层创建任务
	// 执行 data, code := h.svc.CreateOta(...)，完成对应的业务处理

	data, code := h.svc.CreateOta(middleware.User(c), id, req.FirmwareVersion, req.DownloadURL, req.MD5, req.FileSize, req.ForceUpgrade)
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}

// OtaTasks OTA 任务列表 @ticket B端-IoT中台
// 定义方法 OtaTasks：接收者 h *IotHandler，实现该类型对应的业务行为

func (h *IotHandler) OtaTasks(c *gin.Context) {
	// 解析可选 device_id 查询参数
	// 定义并初始化局部变量 deviceID，取值为 uintQuery(c, "device_id")

	deviceID := uintQuery(c, "device_id")
	// 执行 respond(c, h.svc.OtaTasks(deviceID), pkg.CodeOK)，完成对应的业务处理

	respond(c, h.svc.OtaTasks(deviceID), pkg.CodeOK)
// 当前代码块结束

}

// uintQuery 解析 query 数字参数，缺省 0
// 定义函数 uintQuery：封装一段可复用的处理逻辑（参数：c *gin.Context, key string）uint

func uintQuery(c *gin.Context, key string) uint {
	// 定义并初始化局部变量 n, err，取值为 strconv.ParseUint(c.Query(key), 10, 64)

	n, err := strconv.ParseUint(c.Query(key), 10, 64)
	// 条件判断：err != nil 时执行下面逻辑

	if err != nil {
		// 返回 0 给调用方，结束当前函数

		return 0
	// 当前代码块结束

	}
	// 返回 uint(n) 给调用方，结束当前函数

	return uint(n)
// 当前代码块结束

}
