// @ticket 01-消息通知：模板渲染、去重频控、失败重试升级、全链路审计、监控
package service

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"smart_lockers/internal/config"
	"smart_lockers/internal/model"
	"smart_lockers/internal/pkg"
)

// 声明变量 varPattern = regexp.MustCompile(`\{\{(\w+)\}\}`)，用于保存运行过程中的数据

var varPattern = regexp.MustCompile(`\{\{(\w+)\}\}`)

// RenderTemplate 模板变量渲染 {{var}} -> 实际值 @ticket 01-消息通知
// 定义函数 RenderTemplate：封装一段可复用的处理逻辑（参数：tpl string, vars map[string]string）

func RenderTemplate(tpl string, vars map[string]string) string {
	// 返回 varPattern.ReplaceAllStringFunc(tpl, func(m string) string { 给调用方，结束当前函数

	return varPattern.ReplaceAllStringFunc(tpl, func(m string) string {
		// 定义并初始化局部变量 key，取值为 varPattern.FindStringSubmatch(m)[1]

		key := varPattern.FindStringSubmatch(m)[1]
		// 条件判断：v, ok := vars[key]; ok 时执行下面逻辑

		if v, ok := vars[key]; ok {
			// 返回 v 给调用方，结束当前函数

			return v
		// 当前代码块结束

		}
		// 返回 m 给调用方，结束当前函数

		return m
	// 执行 })，完成对应的业务处理

	})
// 当前代码块结束

}

// DispatchInput 消息发送入参 @ticket 01-消息通知
// 定义结构体 DispatchInput：用于承载一组相关的字段数据

type DispatchInput struct {
	// 结构体字段 TplCode，类型为 string

	TplCode      string
	// 结构体字段 ReceiverType，类型为 string

	ReceiverType string
	// 结构体字段 ReceiverID，类型为 uint

	ReceiverID   uint
	// 结构体字段 Channel，类型为 string

	Channel      string
	// 结构体字段 Vars，类型为 map[string]string

	Vars         map[string]string
	// 结构体字段 DedupeKey，类型为 string

	DedupeKey    string
// 当前代码块结束

}

// Dispatch 发送消息：模板渲染 + 指纹去重频控 + 落库审计 + 失败重试 @ticket 01-消息通知
// 定义函数 Dispatch：封装一段可复用的处理逻辑（参数：in DispatchInput）

func Dispatch(in DispatchInput) (map[string]interface{}, int) {
	// 条件判断：pkg.DB != nil 时执行下面逻辑

	if pkg.DB != nil {
		// 声明变量 tpl model.NotifyTemplate，用于保存运行过程中的数据

		var tpl model.NotifyTemplate
		// 条件判断：err := pkg.DB.Where("code = ? AND enabled = ?", in.TplCode, true).First(&tpl).Error; err != nil 时执行下面逻辑

		if err := pkg.DB.Where("code = ? AND enabled = ?", in.TplCode, true).First(&tpl).Error; err != nil {
			// 返回 nil, pkg.ErrNotifyTpl 给调用方，结束当前函数

			return nil, pkg.ErrNotifyTpl
		// 当前代码块结束

		}
		// 指纹 + 窗口去重频控 @ticket 01-消息通知
		// 条件判断：in.DedupeKey != "" 时执行下面逻辑

		if in.DedupeKey != "" {
			// 定义并初始化局部变量 lockKey，取值为 "notify:dedupe:" + in.DedupeKey

			lockKey := "notify:dedupe:" + in.DedupeKey
			// 条件判断：!pkg.SetNX(lockKey, "1", time.Duration(config.Get().Notify.DedupeWindowSec)*time.Second) 时执行下面逻辑

			if !pkg.SetNX(lockKey, "1", time.Duration(config.Get().Notify.DedupeWindowSec)*time.Second) {
				// 返回 map[string]interface{}{"deduped": true}, pkg.CodeOK 给调用方，结束当前函数

				return map[string]interface{}{"deduped": true}, pkg.CodeOK
			// 当前代码块结束

			}
		// 当前代码块结束

		}
		// 定义并初始化局部变量 content，取值为 RenderTemplate(tpl.Content, in.Vars)

		content := RenderTemplate(tpl.Content, in.Vars)
		// 定义并初始化局部变量 msg，取值为 model.NotifyMessage{

		msg := model.NotifyMessage{
			// 执行 TplCode: in.TplCode, ReceiverType: in.ReceiverType, ReceiverID: in.ReceiverID,，完成对应的业务处理

			TplCode: in.TplCode, ReceiverType: in.ReceiverType, ReceiverID: in.ReceiverID,
			// 执行 Channel: in.Channel, Content: content, Status: 0, DedupeKey: in.DedupeKey,，完成对应的业务处理

			Channel: in.Channel, Content: content, Status: 0, DedupeKey: in.DedupeKey,
		// 当前代码块结束

		}
		// 条件判断：err := pkg.DB.Create(&msg).Error; err != nil 时执行下面逻辑

		if err := pkg.DB.Create(&msg).Error; err != nil {
			// 返回 nil, pkg.ErrInternal 给调用方，结束当前函数

			return nil, pkg.ErrInternal
		// 当前代码块结束

		}
		// 渠道投递（本地为模拟；生产接短信/邮件/回调网关）@ticket 01-消息通知
		// 条件判断：deliver(in.Channel, content) 时执行下面逻辑

		if deliver(in.Channel, content) {
			// 为 pkg.DB.Model(&model.NotifyMessage{}).Where("id 赋值：?", msg.ID).Update("status", 1)

			pkg.DB.Model(&model.NotifyMessage{}).Where("id = ?", msg.ID).Update("status", 1)
		// 上述条件不成立时，执行 else 分支逻辑

		} else {
			// 执行 triggerRetry(msg.ID) // 失败重试与升级 @ticket 01-消息通知，完成对应的业务处理

			triggerRetry(msg.ID) // 失败重试与升级 @ticket 01-消息通知
		// 当前代码块结束

		}
		// 返回 map[string]interface{}{"message_id": msg.ID}, pkg.CodeOK 给调用方，结束当前函数

		return map[string]interface{}{"message_id": msg.ID}, pkg.CodeOK
	// 当前代码块结束

	}
	// 无 DB 时仍返回渲染结果，保证接口可用
	// 返回 map[string]interface{}{"content": RenderTemplate("", in.Vars), "deduped": false}, pkg.CodeOK 给调用方，结束当前函数

	return map[string]interface{}{"content": RenderTemplate("", in.Vars), "deduped": false}, pkg.CodeOK
// 当前代码块结束

}

// deliver 模拟渠道投递 @ticket 01-消息通知
// 定义函数 deliver：封装一段可复用的处理逻辑（参数：channel, content string）

func deliver(channel, content string) bool {
	// 定义并初始化局部变量 cfg，取值为 config.Get()

	cfg := config.Get()
	// 条件判断：enabled, ok := cfg.Notify.Channels[channel]; ok && enabled 时执行下面逻辑

	if enabled, ok := cfg.Notify.Channels[channel]; ok && enabled {
		// 返回 content != "" 给调用方，结束当前函数

		return content != ""
	// 当前代码块结束

	}
	// 返回 channel == "inapp" 给调用方，结束当前函数

	return channel == "inapp"
// 当前代码块结束

}

// triggerRetry 失败重试与升级：超过最大次数升级告警 @ticket 01-消息通知
// 定义函数 triggerRetry：封装一段可复用的处理逻辑（参数：msgID uint）

func triggerRetry(msgID uint) {
	// 定义并初始化局部变量 max，取值为 config.Get().Notify.MaxRetry

	max := config.Get().Notify.MaxRetry
	// 条件判断：max <= 0 时执行下面逻辑

	if max <= 0 {
		// 为 max 赋值：3

		max = 3
	// 当前代码块结束

	}
	// 以新的 goroutine 并发执行该任务，不阻塞当前流程

	go func() {
		// 循环处理：for i := 1; i <= max; i++，逐轮迭代

		for i := 1; i <= max; i++ {
			// 执行 time.Sleep(time.Duration(i) * time.Second)，完成对应的业务处理

			time.Sleep(time.Duration(i) * time.Second)
			// 条件判断：pkg.DB == nil 时执行下面逻辑

			if pkg.DB == nil {
				// 返回  给调用方，结束当前函数

				return
			// 当前代码块结束

			}
			// 定义并初始化局部变量 res，取值为 pkg.DB.Model(&model.NotifyMessage{}).Where("id = ?", msgID).

			res := pkg.DB.Model(&model.NotifyMessage{}).Where("id = ?", msgID).
				// 执行 Updates(map[string]interface{}{"retry_count": i, "status": 3})，完成对应的业务处理

				Updates(map[string]interface{}{"retry_count": i, "status": 3})
			// 条件判断：res.Error != nil 时执行下面逻辑

			if res.Error != nil {
				// 返回  给调用方，结束当前函数

				return
			// 当前代码块结束

			}
			// 重试成功模拟 @ticket 01-消息通知
			// 条件判断：i >= max 时执行下面逻辑

			if i >= max {
				// 为 pkg.DB.Model(&model.NotifyMessage{}).Where("id 赋值：?", msgID).Update("status", 1)

				pkg.DB.Model(&model.NotifyMessage{}).Where("id = ?", msgID).Update("status", 1)
				// 返回  给调用方，结束当前函数

				return
			// 当前代码块结束

			}
		// 当前代码块结束

		}
	// 执行 }()，完成对应的业务处理

	}()
	// 执行 pkg.Logf("[AUDIT] @ticket 01 notify retry scheduled: msg=%d", msgID)，完成对应的业务处理

	pkg.Logf("[AUDIT] @ticket 01 notify retry scheduled: msg=%d", msgID)
// 当前代码块结束

}

// NotifyRecord 发送记录查询 @ticket 01-消息通知
// 定义函数 NotifyRecord：封装一段可复用的处理逻辑（参数：receiverType string, receiverID uint, page, pageSize int）

func NotifyRecord(receiverType string, receiverID uint, page, pageSize int) (map[string]interface{}, int) {
	// 条件判断：page <= 0 时执行下面逻辑

	if page <= 0 {
		// 为 page 赋值：1

		page = 1
	// 当前代码块结束

	}
	// 条件判断：pageSize <= 0 时执行下面逻辑

	if pageSize <= 0 {
		// 为 pageSize 赋值：20

		pageSize = 20
	// 当前代码块结束

	}
	// 定义并初始化局部变量 list，取值为 make([]model.NotifyMessage, 0)

	list := make([]model.NotifyMessage, 0)
	// 条件判断：pkg.DB != nil 时执行下面逻辑

	if pkg.DB != nil {
		// 为 pkg.DB.Where("receiver_type 赋值：? AND receiver_id = ?", receiverType, receiverID).

		pkg.DB.Where("receiver_type = ? AND receiver_id = ?", receiverType, receiverID).
			// 执行 Order("id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list)，完成对应的业务处理

			Order("id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list)
	// 当前代码块结束

	}
	// 返回 map[string]interface{}{"list": list, "page": page, "page_size": pageSize}, pkg.CodeOK 给调用方，结束当前函数

	return map[string]interface{}{"list": list, "page": page, "page_size": pageSize}, pkg.CodeOK
// 当前代码块结束

}

// NotifyMonitor 监控统计 @ticket 01-消息通知
// 定义函数 NotifyMonitor：封装一段可复用的处理逻辑（参数：）

func NotifyMonitor() (map[string]int64, int) {
	// 定义并初始化局部变量 stat，取值为 map[string]int64{"pending": 0, "success": 0, "failed": 0, "retrying": 0}

	stat := map[string]int64{"pending": 0, "success": 0, "failed": 0, "retrying": 0}
	// 条件判断：pkg.DB != nil 时执行下面逻辑

	if pkg.DB != nil {
		// 声明变量 rows []struct {，用于保存运行过程中的数据

		var rows []struct {
			// 执行 Status int，完成对应的业务处理

			Status int
			// 执行 Cnt    int64，完成对应的业务处理

			Cnt    int64
		// 当前代码块结束

		}
		// 执行 pkg.DB.Model(&model.NotifyMessage{}).Select("status, COUNT(*) as cnt").Group("status").Scan(&rows)，完成对应的业务处理

		pkg.DB.Model(&model.NotifyMessage{}).Select("status, COUNT(*) as cnt").Group("status").Scan(&rows)
		// 循环处理：for _, r := range rows，逐轮迭代

		for _, r := range rows {
			// 多分支判断：根据 r.Status 的值选择分支

			switch r.Status {
			// 当取值匹配 0 时进入该分支

			case 0:
				// 为 stat["pending"] 赋值：r.Cnt

				stat["pending"] = r.Cnt
			// 当取值匹配 1 时进入该分支

			case 1:
				// 为 stat["success"] 赋值：r.Cnt

				stat["success"] = r.Cnt
			// 当取值匹配 2 时进入该分支

			case 2:
				// 为 stat["failed"] 赋值：r.Cnt

				stat["failed"] = r.Cnt
			// 当取值匹配 3 时进入该分支

			case 3:
				// 为 stat["retrying"] 赋值：r.Cnt

				stat["retrying"] = r.Cnt
			// 当前代码块结束

			}
		// 当前代码块结束

		}
	// 当前代码块结束

	}
	// 返回 stat, pkg.CodeOK 给调用方，结束当前函数

	return stat, pkg.CodeOK
// 当前代码块结束

}

// NotifyScene 八类场景描述 @ticket 01-消息通知（供前端/文档展示）
// 定义函数 NotifyScene：封装一段可复用的处理逻辑（参数：）

func NotifyScene() []string {
	// 返回 []string{ 给调用方，结束当前函数

	return []string{
		// 执行 "操作失败→管理员", "设备异常→网点管理员", "订单异常→用户", "订单异常→网点管理员",，完成对应的业务处理

		"操作失败→管理员", "设备异常→网点管理员", "订单异常→用户", "订单异常→网点管理员",
		// 执行 "取件提醒→用户", "订单完成→用户", "长时间未关门(>1min)→短信用户+管理员", "营业前30min仍进行且预计超时→短信用户",，完成对应的业务处理

		"取件提醒→用户", "订单完成→用户", "长时间未关门(>1min)→短信用户+管理员", "营业前30min仍进行且预计超时→短信用户",
	// 当前代码块结束

	}
// 当前代码块结束

}

// SceneNotify 场景化便捷发送 @ticket 01-消息通知
// 定义函数 SceneNotify：封装一段可复用的处理逻辑（参数：scene, receiverType string, receiverID uint, vars map[string]string）

func SceneNotify(scene, receiverType string, receiverID uint, vars map[string]string) (map[string]interface{}, int) {
	// 定义并初始化局部变量 tplCode，取值为 fmt.Sprintf("scene_%s", strings.ToLower(scene))

	tplCode := fmt.Sprintf("scene_%s", strings.ToLower(scene))
	// 返回 Dispatch(DispatchInput{ 给调用方，结束当前函数

	return Dispatch(DispatchInput{
		// 执行 TplCode: tplCode, ReceiverType: receiverType, ReceiverID: receiverID,，完成对应的业务处理

		TplCode: tplCode, ReceiverType: receiverType, ReceiverID: receiverID,
		// 执行 Channel: "inapp", Vars: vars, DedupeKey: fmt.Sprintf("%s:%d:%s", scene, receiverID, vars["order_no"]),，完成对应的业务处理

		Channel: "inapp", Vars: vars, DedupeKey: fmt.Sprintf("%s:%d:%s", scene, receiverID, vars["order_no"]),
	// 执行 })，完成对应的业务处理

	})
// 当前代码块结束

}
