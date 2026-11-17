// @ticket 01-消息通知：发送/记录/监控/场景 HTTP handler
package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"smart_lockers/internal/middleware"
	"smart_lockers/internal/pkg"
	"smart_lockers/internal/service"
)

// Send @ticket 01-消息通知：POST /api/v1/notify/send
func Send(c *gin.Context) {
	var req struct {
		TplCode      string            `json:"tpl_code" binding:"required"`
		ReceiverType string            `json:"receiver_type" binding:"required"`
		ReceiverID   uint              `json:"receiver_id" binding:"required"`
		Channel      string            `json:"channel"`
		Vars         map[string]string `json:"vars"`
		DedupeKey    string            `json:"dedupe_key"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, pkg.ErrBad)
		return
	}
	data, code := service.Dispatch(service.DispatchInput{
		TplCode: req.TplCode, ReceiverType: req.ReceiverType, ReceiverID: req.ReceiverID,
		Channel: req.Channel, Vars: req.Vars, DedupeKey: req.DedupeKey,
	})
	respond(c, data, code)
}

// SendScene @ticket 01-消息通知：POST /api/v1/notify/scene
func SendScene(c *gin.Context) {
	var req struct {
		Scene        string            `json:"scene" binding:"required"`
		ReceiverType string            `json:"receiver_type" binding:"required"`
		ReceiverID   uint              `json:"receiver_id" binding:"required"`
		Vars         map[string]string `json:"vars"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Fail(c, pkg.ErrBad)
		return
	}
	data, code := service.SceneNotify(req.Scene, req.ReceiverType, req.ReceiverID, req.Vars)
	respond(c, data, code)
}

// Scenes @ticket 01-消息通知：GET /api/v1/notify/scenes（八类场景）
func Scenes(c *gin.Context) {
	pkg.OK(c, map[string]interface{}{"scenes": service.NotifyScene()})
}

// Records @ticket 01-消息通知：GET /api/v1/notify/records
func Records(c *gin.Context) {
	receiverType := c.Query("receiver_type")
	receiverID, _ := strconv.ParseUint(c.Query("receiver_id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	data, code := service.NotifyRecord(receiverType, uint(receiverID), page, pageSize)
	respond(c, data, code)
}

// Monitor @ticket 01-消息通知：GET /api/v1/notify/monitor
func Monitor(c *gin.Context) {
	_ = middleware.UID(c)
	data, code := service.NotifyMonitor()
	respond(c, data, code)
}
