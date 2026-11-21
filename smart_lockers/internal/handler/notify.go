// @ticket 01-消息通知：发送/记录/监控/场景 HTTP handler
// @proto 原型对齐：notify 通知中心 / alarm 告警中心(场景通知)；见 AGENT.MD 2.1
package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"smart_lockers/internal/middleware"
	"smart_lockers/internal/pkg"
	"smart_lockers/internal/service"
)

// Send @ticket 01-消息通知：POST /api/v1/notify/send
// 定义函数 Send：封装一段可复用的处理逻辑（参数：c *gin.Context）

func Send(c *gin.Context) {
	// 声明变量 req struct {，用于保存运行过程中的数据

	var req struct {
		// 执行 TplCode      string            `json:"tpl_code" binding:"required"`，完成对应的业务处理

		TplCode      string            `json:"tpl_code" binding:"required"`
		// 执行 ReceiverType string            `json:"receiver_type" binding:"required"`，完成对应的业务处理

		ReceiverType string            `json:"receiver_type" binding:"required"`
		// 执行 ReceiverID   uint              `json:"receiver_id" binding:"required"`，完成对应的业务处理

		ReceiverID   uint              `json:"receiver_id" binding:"required"`
		// 执行 Channel      string            `json:"channel"`，完成对应的业务处理

		Channel      string            `json:"channel"`
		// 执行 Vars         map[string]string `json:"vars"`，完成对应的业务处理

		Vars         map[string]string `json:"vars"`
		// 执行 DedupeKey    string            `json:"dedupe_key"`，完成对应的业务处理

		DedupeKey    string            `json:"dedupe_key"`
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
	// 定义并初始化局部变量 data, code，取值为 service.Dispatch(service.DispatchInput{

	data, code := service.Dispatch(service.DispatchInput{
		// 执行 TplCode: req.TplCode, ReceiverType: req.ReceiverType, ReceiverID: req.ReceiverID,，完成对应的业务处理

		TplCode: req.TplCode, ReceiverType: req.ReceiverType, ReceiverID: req.ReceiverID,
		// 执行 Channel: req.Channel, Vars: req.Vars, DedupeKey: req.DedupeKey,，完成对应的业务处理

		Channel: req.Channel, Vars: req.Vars, DedupeKey: req.DedupeKey,
	// 执行 })，完成对应的业务处理

	})
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}

// SendScene @ticket 01-消息通知：POST /api/v1/notify/scene
// 定义函数 SendScene：封装一段可复用的处理逻辑（参数：c *gin.Context）

func SendScene(c *gin.Context) {
	// 声明变量 req struct {，用于保存运行过程中的数据

	var req struct {
		// 执行 Scene        string            `json:"scene" binding:"required"`，完成对应的业务处理

		Scene        string            `json:"scene" binding:"required"`
		// 执行 ReceiverType string            `json:"receiver_type" binding:"required"`，完成对应的业务处理

		ReceiverType string            `json:"receiver_type" binding:"required"`
		// 执行 ReceiverID   uint              `json:"receiver_id" binding:"required"`，完成对应的业务处理

		ReceiverID   uint              `json:"receiver_id" binding:"required"`
		// 执行 Vars         map[string]string `json:"vars"`，完成对应的业务处理

		Vars         map[string]string `json:"vars"`
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
	// 定义并初始化局部变量 data, code，取值为 service.SceneNotify(req.Scene, req.ReceiverType, req.ReceiverID, req.Vars)

	data, code := service.SceneNotify(req.Scene, req.ReceiverType, req.ReceiverID, req.Vars)
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}

// Scenes @ticket 01-消息通知：GET /api/v1/notify/scenes（八类场景）
// 定义函数 Scenes：封装一段可复用的处理逻辑（参数：c *gin.Context）

func Scenes(c *gin.Context) {
	// 执行 pkg.OK(c, map[string]interface{}{"scenes": service.NotifyScene()})，完成对应的业务处理

	pkg.OK(c, map[string]interface{}{"scenes": service.NotifyScene()})
// 当前代码块结束

}

// Records @ticket 01-消息通知：GET /api/v1/notify/records
// 定义函数 Records：封装一段可复用的处理逻辑（参数：c *gin.Context）

func Records(c *gin.Context) {
	// 定义并初始化局部变量 receiverType，取值为 c.Query("receiver_type")

	receiverType := c.Query("receiver_type")
	// 定义并初始化局部变量 receiverID, _，取值为 strconv.ParseUint(c.Query("receiver_id"), 10, 64)

	receiverID, _ := strconv.ParseUint(c.Query("receiver_id"), 10, 64)
	// 定义并初始化局部变量 page, _，取值为 strconv.Atoi(c.DefaultQuery("page", "1"))

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	// 定义并初始化局部变量 pageSize, _，取值为 strconv.Atoi(c.DefaultQuery("page_size", "20"))

	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	// 定义并初始化局部变量 data, code，取值为 service.NotifyRecord(receiverType, uint(receiverID), page, pageSize)

	data, code := service.NotifyRecord(receiverType, uint(receiverID), page, pageSize)
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}

// Monitor @ticket 01-消息通知：GET /api/v1/notify/monitor
// 定义函数 Monitor：封装一段可复用的处理逻辑（参数：c *gin.Context）

func Monitor(c *gin.Context) {
	// 忽略该调用的返回值，仅执行其副作用

	_ = middleware.UID(c)
	// 定义并初始化局部变量 data, code，取值为 service.NotifyMonitor()

	data, code := service.NotifyMonitor()
	// 执行 respond(c, data, code)，完成对应的业务处理

	respond(c, data, code)
// 当前代码块结束

}
