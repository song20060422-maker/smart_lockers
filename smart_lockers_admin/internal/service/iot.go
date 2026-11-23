// @ticket B端-IoT中台：指令下发/OTA 升级/中台总览服务层
// 对齐《智能快递柜IoT设备MQTT接入技术方案》：B端发布 platform/# 下行、订阅 device/response/# 与 device/alert/# 回执
package service

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"smart_lockers_admin/internal/middleware"
	"smart_lockers_admin/internal/model"
	"smart_lockers_admin/internal/pkg"
	"smart_lockers_admin/internal/store"
)

// IotService IoT 中台服务
// 定义结构体 IotService：用于承载一组相关的字段数据

type IotService struct{ st *store.Store }

// 定义函数 NewIot：封装一段可复用的处理逻辑（参数：st *store.Store）

func NewIot(st *store.Store) *IotService { return &IotService{st: st} }

// randomHex 生成 n 字节随机十六进制串（指令号/任务号生成）
// 定义函数 randomHex：封装一段可复用的处理逻辑（参数：n int）string

func randomHex(n int) string {
	// 定义并初始化局部变量 b，取值为 make([]byte, n)

	b := make([]byte, n)
	// 条件判断：_, err := rand.Read(b); err != nil 时执行下面逻辑

	if _, err := rand.Read(b); err != nil {
		// 兜底：时间戳（crypto/rand 理论上不失败）
		// 返回 fmt.Sprintf("%x", time.Now().UnixNano()) 给调用方，结束当前函数

		return fmt.Sprintf("%x", time.Now().UnixNano())
	// 当前代码块结束

	}
	// 返回 hex.EncodeToString(b) 给调用方，结束当前函数

	return hex.EncodeToString(b)
// 当前代码块结束

}

// DispatchCommand 下发控制指令：落库 + 发布 platform/command/{id}/{type}（原型 devices 视图"下发控制指令"）@ticket B端-IoT中台
// 定义方法 DispatchCommand：接收者 s *IotService，实现该类型对应的业务行为

func (s *IotService) DispatchCommand(u middleware.CurrentUser, deviceID uint, cmdType, lockerNo string) (map[string]interface{}, int) {
	// 指令类型白名单校验（方案 3.3）
	// 定义并初始化局部变量 valid，取值为 map[string]bool{"open_door": true, "close_door": true, "reboot": true, "config": true}

	valid := map[string]bool{"open_door": true, "close_door": true, "reboot": true, "config": true}
	// 条件判断：!valid[cmdType] 时执行下面逻辑

	if !valid[cmdType] {
		// 返回 nil, pkg.ErrCmdInvalid 给调用方，结束当前函数

		return nil, pkg.ErrCmdInvalid
	// 当前代码块结束

	}
	// 条件判断：d, ok := s.st.Device(deviceID); !ok 时执行下面逻辑

	d, ok := s.st.Device(deviceID)
	// 条件判断：!ok 时执行下面逻辑

	if !ok {
		// 返回 nil, pkg.ErrDeviceNF 给调用方，结束当前函数

		return nil, pkg.ErrDeviceNF
	// 当前代码块结束

	}
	// 数据域校验：网点管理员仅可操作本网点设备 @ticket B端-RBAC
	// 条件判断：scope := middleware.SiteScope(u); scope != 0 && d.SiteID != scope 时执行下面逻辑

	if scope := middleware.SiteScope(u); scope != 0 && d.SiteID != scope {
		// 返回 nil, pkg.ErrForb 给调用方，结束当前函数

		return nil, pkg.ErrForb
	// 当前代码块结束

	}
	// 生成全局唯一指令号（command_id 幂等依据，方案 5.3）
	// 定义并初始化局部变量 commandID，取值为 fmt.Sprintf("CMD%s%s", time.Now().Format("20060102150405"), strings.ToUpper(randomHex(4)))

	commandID := fmt.Sprintf("CMD%s%s", time.Now().Format("20060102150405"), strings.ToUpper(randomHex(4)))
	// 组装下行报文（方案 5.1 open_door 报文格式）
	// 定义并初始化局部变量 payloadBytes，取值为 json.Marshal

	payloadBytes, _ := json.Marshal(map[string]interface{}{
		// 执行 "command_id": commandID,，完成对应的业务处理

		"command_id": commandID,
		// 执行 "device_id": d.ID,，完成对应的业务处理

		"device_id":  d.ID,
		// 执行 "timestamp": time.Now().UnixMilli(),，完成对应的业务处理

		"timestamp":  time.Now().UnixMilli(),
		// 执行 "locker_no": lockerNo,，完成对应的业务处理

		"locker_no":  lockerNo,
		// 执行 "timeout": 30,，完成对应的业务处理

		"timeout":    30,
	// 当前代码块结束

	})
	// 发布指令（Broker 不可用时仅落库，published=false 供前端提示演示状态）
	// 定义并初始化局部变量 pubErr，取值为 pkg.MqttPublish(fmt.Sprintf("platform/command/%d/%s", d.ID, cmdType), 1, payloadBytes)

	pubErr := pkg.MqttPublish(fmt.Sprintf("platform/command/%d/%s", d.ID, cmdType), 1, payloadBytes)
	// 指令落库（含幂等指令号/审计触发人）@ticket B端-IoT中台
	// 定义并初始化局部变量 rec，取值为 &model.DeviceCommand{...}

	rec := &model.DeviceCommand{
		// 执行 CommandID: commandID,，完成对应的业务处理

		CommandID: commandID,
		// 执行 DeviceID: d.ID,，完成对应的业务处理

		DeviceID:  d.ID,
		// 执行 DeviceSN: d.DeviceSN,，完成对应的业务处理

		DeviceSN:  d.DeviceSN,
		// 执行 CmdType: cmdType,，完成对应的业务处理

		CmdType:   cmdType,
		// 执行 LockerNo: lockerNo,，完成对应的业务处理

		LockerNo:  lockerNo,
		// 执行 Payload: string(payloadBytes),，完成对应的业务处理

		Payload:   string(payloadBytes),
		// 执行 Status: 0,，完成对应的业务处理

		Status:    0,
		// 执行 Operator: u.Name(),，完成对应的业务处理

		Operator:  u.Name(),
		// 执行 CreatedAt: time.Now(),，完成对应的业务处理

		CreatedAt: time.Now(),
	// 当前代码块结束

	}
	// 条件判断：pubErr == nil 时执行下面逻辑（发布成功即进入"已下发"）

	if pubErr == nil {
		// 为 rec.Status 赋值：0（已下发待响应）

		rec.Status = 0
	// 上述条件不成立时，执行 else 分支逻辑

	} else {
		// Broker 不可达：标记超时态，前端可见"未发布"
		// 为 rec.Status 赋值：3

		rec.Status = 3
		// 执行 pkg.Logf("[WARN] @ticket B端-IoT command publish deferred: cmd=%s err=%v", commandID, pubErr)，完成对应的业务处理

		pkg.Logf("[WARN] @ticket B端-IoT command publish deferred: cmd=%s err=%v", commandID, pubErr)
	// 当前代码块结束

	}
	// 执行 s.st.AddCommand(rec)，完成对应的业务处理

	s.st.AddCommand(rec)
	// 审计留痕（设备指令属风险动作）@ticket B端-审计
	// 执行 pkg.RecordAudit(u.Name(), string(u.Role), "iot.command.dispatch", idstr(deviceID)+" "+cmdType, "success")，完成对应的业务处理

	pkg.RecordAudit(u.Name(), string(u.Role), "iot.command.dispatch", idstr(deviceID)+" "+cmdType, "success")
	// 返回指令回执给调用方，结束当前函数

	return map[string]interface{}{
		// 执行 "command_id": commandID,，完成对应的业务处理

		"command_id": commandID,
		// 执行 "device_id": d.ID,，完成对应的业务处理

		"device_id":  d.ID,
		// 执行 "cmd_type": cmdType,，完成对应的业务处理

		"cmd_type":   cmdType,
		// 执行 "published": pubErr == nil,，完成对应的业务处理

		"published":  pubErr == nil,
		// 执行 "status": rec.Status,，完成对应的业务处理

		"status":     rec.Status,
	// 当前代码块结束

	}, pkg.CodeOK
// 当前代码块结束

}

// CommandRecords 指令响应记录列表（原型 devices 视图"指令响应记录"）@ticket B端-IoT中台
// 定义方法 CommandRecords：接收者 s *IotService，实现该类型对应的业务行为

func (s *IotService) CommandRecords(u middleware.CurrentUser, deviceID uint) []model.DeviceCommand {
	// 返回 s.st.Commands(deviceID) 给调用方，结束当前函数

	return s.st.Commands(deviceID)
// 当前代码块结束

}

// CreateOta 创建并下发 OTA 升级任务（原型 ota 视图"新建升级任务"）@ticket B端-IoT中台
// 定义方法 CreateOta：接收者 s *IotService，实现该类型对应的业务行为

func (s *IotService) CreateOta(u middleware.CurrentUser, deviceID uint, version, url, md5 string, size int64, force bool) (map[string]interface{}, int) {
	// 条件判断：version == "" || url == "" 时执行下面逻辑

	if version == "" || url == "" {
		// 返回 nil, pkg.ErrBad 给调用方，结束当前函数

		return nil, pkg.ErrBad
	// 当前代码块结束

	}
	// 条件判断：d, ok := s.st.Device(deviceID); !ok 时执行下面逻辑

	d, ok := s.st.Device(deviceID)
	// 条件判断：!ok 时执行下面逻辑

	if !ok {
		// 返回 nil, pkg.ErrDeviceNF 给调用方，结束当前函数

		return nil, pkg.ErrDeviceNF
	// 当前代码块结束

	}
	// 数据域校验
	// 条件判断：scope := middleware.SiteScope(u); scope != 0 && d.SiteID != scope 时执行下面逻辑

	if scope := middleware.SiteScope(u); scope != 0 && d.SiteID != scope {
		// 返回 nil, pkg.ErrForb 给调用方，结束当前函数

		return nil, pkg.ErrForb
	// 当前代码块结束

	}
	// 生成 OTA 任务号（方案 3.3）
	// 定义并初始化局部变量 otaID，取值为 fmt.Sprintf("OTA%s%s", time.Now().Format("20060102150405"), strings.ToUpper(randomHex(3)))

	otaID := fmt.Sprintf("OTA%s%s", time.Now().Format("20060102150405"), strings.ToUpper(randomHex(3)))
	// 任务落库
	// 定义并初始化局部变量 task，取值为 &model.OtaTask{...}

	task := &model.OtaTask{
		// 执行 OtaID: otaID,，完成对应的业务处理

		OtaID:           otaID,
		// 执行 DeviceID: d.ID,，完成对应的业务处理

		DeviceID:        d.ID,
		// 执行 DeviceSN: d.DeviceSN,，完成对应的业务处理

		DeviceSN:        d.DeviceSN,
		// 执行 FirmwareVersion: version,，完成对应的业务处理

		FirmwareVersion: version,
		// 执行 DownloadURL: url,，完成对应的业务处理

		DownloadURL:     url,
		// 执行 MD5: md5,，完成对应的业务处理

		MD5:             md5,
		// 执行 FileSize: size,，完成对应的业务处理

		FileSize:        size,
		// 执行 ForceUpgrade: force,，完成对应的业务处理

		ForceUpgrade:    force,
		// 执行 Status: 0,，完成对应的业务处理

		Status:          0,
		// 执行 Operator: u.Name(),，完成对应的业务处理

		Operator:        u.Name(),
		// 执行 CreatedAt: time.Now(), UpdatedAt: time.Now(),，完成对应的业务处理

		CreatedAt:       time.Now(), UpdatedAt: time.Now(),
	// 当前代码块结束

	}
	// 组装 OTA 下行报文（方案 3.3 platform/ota/{id}/upgrade）
	// 定义并初始化局部变量 payloadBytes，取值为 json.Marshal

	payloadBytes, _ := json.Marshal(map[string]interface{}{
		// 执行 "ota_id": otaID,，完成对应的业务处理

		"ota_id":           otaID,
		// 执行 "device_id": d.ID,，完成对应的业务处理

		"device_id":        d.ID,
		// 执行 "firmware_version": version,，完成对应的业务处理

		"firmware_version": version,
		// 执行 "download_url": url,，完成对应的业务处理

		"download_url":     url,
		// 执行 "md5": md5,，完成对应的业务处理

		"md5":              md5,
		// 执行 "file_size": size,，完成对应的业务处理

		"file_size":        size,
		// 执行 "force_upgrade": force,，完成对应的业务处理

		"force_upgrade":    force,
	// 当前代码块结束

	})
	// 发布 OTA 通知（QoS1；300s 超时不重试，等待设备主动拉取）
	// 定义并初始化局部变量 pubErr，取值为 pkg.MqttPublish(fmt.Sprintf("platform/ota/%d/upgrade", d.ID), 1, payloadBytes)

	pubErr := pkg.MqttPublish(fmt.Sprintf("platform/ota/%d/%s", d.ID, "upgrade"), 1, payloadBytes)
	// 发布失败时任务标记失败态，前端可见
	// 条件判断：pubErr != nil 时执行下面逻辑

	if pubErr != nil {
		// 为 task.Status 赋值：4

		task.Status = 4
		// 为 task.Message 赋值："broker unavailable"

		task.Message = "broker unavailable"
		// 执行 pkg.Logf("[WARN] @ticket B端-IoT ota publish deferred: ota=%s err=%v", otaID, pubErr)，完成对应的业务处理

		pkg.Logf("[WARN] @ticket B端-IoT ota publish deferred: ota=%s err=%v", otaID, pubErr)
	// 当前代码块结束

	}
	// 执行 s.st.AddOta(task)，完成对应的业务处理

	s.st.AddOta(task)
	// 审计留痕（OTA 属高风险动作）@ticket B端-审计
	// 执行 pkg.RecordAudit(u.Name(), string(u.Role), "iot.ota.create", idstr(deviceID)+" "+version, "success")，完成对应的业务处理

	pkg.RecordAudit(u.Name(), string(u.Role), "iot.ota.create", idstr(deviceID)+" "+version, "success")
	// 返回任务回执给调用方，结束当前函数

	return map[string]interface{}{
		// 执行 "ota_id": otaID,，完成对应的业务处理

		"ota_id":        otaID,
		// 执行 "device_id": d.ID,，完成对应的业务处理

		"device_id":     d.ID,
		// 执行 "firmware_version": version,，完成对应的业务处理

		"firmware_version": version,
		// 执行 "published": pubErr == nil,，完成对应的业务处理

		"published":     pubErr == nil,
		// 执行 "status": task.Status,，完成对应的业务处理

		"status":        task.Status,
	// 当前代码块结束

	}, pkg.CodeOK
// 当前代码块结束

}

// OtaTasks OTA 任务列表
// 定义方法 OtaTasks：接收者 s *IotService，实现该类型对应的业务行为

func (s *IotService) OtaTasks(deviceID uint) []model.OtaTask {
	// 返回 s.st.OtaTasks(deviceID) 给调用方，结束当前函数

	return s.st.OtaTasks(deviceID)
// 当前代码块结束

}

// OnResponse 处理 device/response/{id}/{type} 回执：更新指令终态（原型 devices 视图"指令响应记录"）@ticket B端-IoT中台
// 定义函数 OnResponse：封装一段可复用的处理逻辑（参数：st *store.Store, deviceID uint, payload string）

func OnResponse(st *store.Store, deviceID uint, payload string) {
	// 定义并初始化局部变量 rp，取值为 map[string]interface{}{} 并解析报文

	rp := map[string]interface{}{}
	// 条件判断：err := json.Unmarshal([]byte(payload), &rp); err != nil 时执行下面逻辑

	if err := json.Unmarshal([]byte(payload), &rp); err != nil {
		// 执行 pkg.Logf("[WARN] @ticket B端-IoT response invalid: device=%d err=%v", deviceID, err)，完成对应的业务处理

		pkg.Logf("[WARN] @ticket B端-IoT response invalid: device=%d err=%v", deviceID, err)
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 提取指令号与结果
	// 定义并初始化局部变量 commandID, _，取值为 rp["command_id"].(string)

	commandID, _ := rp["command_id"].(string)
	// 定义并初始化局部变量 result, _，取值为 rp["result"].(string)

	result, _ := rp["result"].(string)
	// 条件判断：commandID == "" 时执行下面逻辑

	if commandID == "" {
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 按 commandID 定位指令并更新终态
	// 条件判断：rec, ok := st.CommandByID(commandID); ok 时执行下面逻辑

	if rec, ok := st.CommandByID(commandID); ok {
		// 定义并初始化局部变量 nowT，取值为 time.Now()

		nowT := time.Now()
		// 更新终态：result=success → 1 成功，否则 2 失败
		// 执行 st.UpdateCommand(rec.ID, func(m *model.DeviceCommand) bool {，完成对应的业务处理

		st.UpdateCommand(rec.ID, func(m *model.DeviceCommand) bool {
			// 条件判断：result == "success" 时执行下面逻辑

			if result == "success" {
				// 为 m.Status 赋值：1

				m.Status = 1
			// 上述条件不成立时，执行 else 分支逻辑

			} else {
				// 为 m.Status 赋值：2

				m.Status = 2
			// 当前代码块结束

			}
			// 记录响应报文与时间
			// 为 m.Response, m.RespondedAt 赋值：payload, &nowT

			m.Response, m.RespondedAt = payload, &nowT
			// 返回 true 给调用方，结束当前函数

			return true
		// 当前代码块结束

		})
	// 当前代码块结束

	}
// 当前代码块结束

}

// OnAlert 处理 device/alert/{id}/{action} 告警：联动 B端告警中心（方案 4.4 告警流程）@ticket B端-IoT中台
// 定义函数 OnAlert：封装一段可复用的处理逻辑（参数：st *store.Store, deviceID uint, topic, payload string）

func OnAlert(st *store.Store, deviceID uint, topic, payload string) {
	// 告警级别从报文提取（方案 3.2 alert_level），缺省 warn
	// 定义并初始化局部变量 rp，取值为 map[string]interface{}{}

	rp := map[string]interface{}{}
	// 忽略该调用的返回值，仅执行其副作用

	_ = json.Unmarshal([]byte(payload), &rp)
	// 定义并初始化局部变量 level, _，取值为 rp["alert_level"].(string)

	level, _ := rp["alert_level"].(string)
	// 条件判断：level == "" 时执行下面逻辑

	if level == "" {
		// 为 level 赋值："warn"

		level = "warn"
	// 当前代码块结束

	}
	// 级别映射：方案 low/medium/high/critical → B端 info/warn/critical
	// 定义并初始化局部变量 mapped，取值为 mapAlarmLevel(level)

	mapped := mapAlarmLevel(level)
	// 查设备取站点（告警挂到设备所属网点）
	// 条件判断：d, ok := st.Device(deviceID); ok 时执行下面逻辑

	if d, ok := st.Device(deviceID); ok {
		// 创建告警记录（方案 4.4 步骤 5：创建告警工单）
		// 执行 st.AddAlarm(&model.Alarm{...})，完成对应的业务处理

		st.AddAlarm(&model.Alarm{
			// 执行 SiteID: d.SiteID,，完成对应的业务处理

			SiteID:    d.SiteID,
			// 执行 DeviceID: d.ID,，完成对应的业务处理

			DeviceID:  d.ID,
			// 执行 Level: mapped,，完成对应的业务处理

			Level:     mapped,
			// 执行 Content: topic + " " + payload,，完成对应的业务处理

			Content:   topic + " " + payload,
			// 执行 Status: 0,，完成对应的业务处理

			Status:    0,
			// 执行 CreatedAt: time.Now(),，完成对应的业务处理

			CreatedAt: time.Now(),
		// 当前代码块结束

		})
		// 执行 pkg.Logf("[WARN] @ticket B端-IoT alert recorded: device=%d level=%s topic=%s", deviceID, mapped, topic)，完成对应的业务处理

		pkg.Logf("[WARN] @ticket B端-IoT alert recorded: device=%d level=%s topic=%s", deviceID, mapped, topic)
	// 当前代码块结束

	}
// 当前代码块结束

}

// mapAlarmLevel 方案告警级别 → B端告警级别映射
// 定义函数 mapAlarmLevel：封装一段可复用的处理逻辑（参数：level string）string

func mapAlarmLevel(level string) string {
	// 多分支判断：根据 level 的值选择分支

	switch level {
	// 当取值匹配 low 时进入该分支
	case "low":
		// 返回 "info" 给调用方，结束当前函数

		return "info"
	// 当取值匹配 medium 时进入该分支
	case "medium":
		// 返回 "warn" 给调用方，结束当前函数

		return "warn"
	// 当取值匹配 high/critical 时进入该分支
	case "high", "critical":
		// 返回 "critical" 给调用方，结束当前函数

		return "critical"
	// 当前代码块结束

	}
	// 兜底 warn
	// 返回 "warn" 给调用方，结束当前函数

	return "warn"
// 当前代码块结束

}

// RouteMessage B端上行消息统一入口：device/response/# 更新指令回执，device/alert/# 联动告警 @ticket B端-IoT中台
// 定义函数 RouteMessage：封装一段可复用的处理逻辑（参数：st *store.Store, topic string, payload []byte）

func RouteMessage(st *store.Store, topic string, payload []byte) {
	// 按 / 切分 Topic：device/{kind}/{device_id}/{action}
	// 定义并初始化局部变量 parts，取值为 strings.Split(topic, "/")

	parts := strings.Split(topic, "/")
	// 条件判断：len(parts) < 3 || parts[0] != "device" 时执行下面逻辑

	if len(parts) < 3 || parts[0] != "device" {
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 解析设备 ID
	// 定义并初始化局部变量 deviceID, err，取值为 strconv.ParseUint(parts[2], 10, 64)

	deviceID, err := strconv.ParseUint(parts[2], 10, 64)
	// 条件判断：err != nil 时执行下面逻辑

	if err != nil {
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 多分支判断：根据 parts[1] 的值选择分支

	switch parts[1] {
	// 当取值匹配 response 时进入该分支
	case "response":
		// 执行 OnResponse(st, uint(deviceID), string(payload))，完成对应的业务处理

		OnResponse(st, uint(deviceID), string(payload))
	// 当取值匹配 alert 时进入该分支
	case "alert":
		// 执行 OnAlert(st, uint(deviceID), topic, string(payload))，完成对应的业务处理

		OnAlert(st, uint(deviceID), topic, string(payload))
	// 当前代码块结束

	}
// 当前代码块结束

}

// TopicSpec MQTT Topic 设计规范数据（原型 mqtt 视图"Topic 设计规范"静态展示数据源）@ticket B端-IoT中台
// 定义函数 TopicSpec：封装一段可复用的处理逻辑（参数：）[]map[string]string

func TopicSpec() []map[string]string {
	// 返回规范条目列表给调用方，结束当前函数

	return []map[string]string{
		// 设备上线
		{"topic": "device/status/{id}/online", "qos": "1", "direction": "上行", "desc": "设备上线消息（连接成功时）"},
		// 设备离线（遗嘱）
		{"topic": "device/status/{id}/offline", "qos": "1", "direction": "上行", "desc": "设备离线消息（遗嘱消息）"},
		// 心跳
		{"topic": "device/status/{id}/heartbeat", "qos": "0", "direction": "上行", "desc": "设备心跳（每5分钟）"},
		// 柜格状态
		{"topic": "device/status/{id}/locker", "qos": "1", "direction": "上行", "desc": "柜格状态变更"},
		// 开门指令
		{"topic": "platform/command/{id}/open_door", "qos": "1", "direction": "下行", "desc": "开门指令（响应 device/response/{id}/open_door）"},
		// 关门指令
		{"topic": "platform/command/{id}/close_door", "qos": "1", "direction": "下行", "desc": "关门指令"},
		// 重启指令
		{"topic": "platform/command/{id}/reboot", "qos": "1", "direction": "下行", "desc": "重启指令（30s超时不重试）"},
		// 配置下发
		{"topic": "platform/config/{id}/price", "qos": "1", "direction": "下行", "desc": "价格配置下发"},
		// OTA 升级
		{"topic": "platform/ota/{id}/upgrade", "qos": "1", "direction": "下行", "desc": "OTA 升级通知"},
		// OTA 进度
		{"topic": "device/ota/{id}/progress", "qos": "1", "direction": "上行", "desc": "OTA 升级进度上报"},
	// 当前代码块结束

	}
// 当前代码块结束

}

// AclSpec MQTT ACL 权限规则数据（原型 mqtt 视图"ACL 权限规则"静态展示数据源）@ticket B端-IoT中台
// 定义函数 AclSpec：封装一段可复用的处理逻辑（参数：）[]map[string]string

func AclSpec() []map[string]string {
	// 返回 ACL 条目列表给调用方，结束当前函数

	return []map[string]string{
		// 设备发布权限
		{"role": "设备(SL20240101001)", "perm": "publish", "rule": "device/+/1001/#", "desc": "设备只能发布自己的上行 Topic"},
		// 设备订阅权限
		{"role": "设备(SL20240101001)", "perm": "subscribe", "rule": "platform/+/1001/#", "desc": "设备只能订阅自己的下行 Topic"},
		// 设备禁止项
		{"role": "设备(SL20240101001)", "perm": "deny", "rule": "device/#", "desc": "禁止设备订阅其他设备的 Topic"},
		// 后端服务权限
		{"role": "backend-service", "perm": "publish", "rule": "platform/#", "desc": "后端允许发布全部下行 Topic"},
		// 后端服务订阅权限
		{"role": "backend-service", "perm": "subscribe", "rule": "device/#", "desc": "后端允许订阅全部上行 Topic"},
	// 当前代码块结束

	}
// 当前代码块结束

}

// Overview IoT 中台总览：设备状态分布 + EMQX 健康 + Topic/ACL 规范（原型 overview 视图数据源）@ticket B端-IoT中台
// 定义方法 Overview：接收者 s *IotService，实现该类型对应的业务行为

func (s *IotService) Overview(u middleware.CurrentUser) map[string]interface{} {
	// 设备状态统计（复用设备监控口径：0在线 1离线 2故障）
	// 定义并初始化局部变量 devs，取值为 s.st.Devices(middleware.SiteScope(u))

	devs := s.st.Devices(middleware.SiteScope(u))
	// 定义并初始化局部变量 online, offline, fault，取值为 0, 0, 0

	online, offline, fault := 0, 0, 0
	// 循环处理：for _, d := range devs，逐轮迭代

	for _, d := range devs {
		// 多分支判断：根据 d.Status 的值选择分支

		switch d.Status {
		// 当取值匹配 0 时进入该分支
		case 0:
			// 执行 online++，完成对应的业务处理

			online++
		// 当取值匹配 1 时进入该分支
		case 1:
			// 执行 offline++，完成对应的业务处理

			offline++
		// 当取值匹配 2 时进入该分支
		case 2:
			// 执行 fault++，完成对应的业务处理

			fault++
		// 当前代码块结束

		}
	// 当前代码块结束

	}
	// 指令/OTA 统计
	// 定义并初始化局部变量 cmds，取值为 s.st.Commands(0)

	cmds := s.st.Commands(0)
	// 定义并初始化局部变量 otas，取值为 s.st.OtaTasks(0)

	otas := s.st.OtaTasks(0)
	// EMQX 健康状态（connected=Broker 可达；本地演示模式为 standby）
	// 定义并初始化局部变量 emqxStatus，取值为 "standby"

	emqxStatus := "standby"
	// 条件判断：pkg.MqttReady() 时执行下面逻辑

	if pkg.MqttReady() {
		// 为 emqxStatus 赋值："connected"

		emqxStatus = "connected"
	// 当前代码块结束

	}
	// 返回总览数据给调用方，结束当前函数

	return map[string]interface{}{
		// 执行 "devices": map[string]interface{}{"total": len(devs), "online": online, "offline": offline, "fault": fault},，完成对应的业务处理

		"devices": map[string]interface{}{"total": len(devs), "online": online, "offline": offline, "fault": fault},
		// 执行 "emqx": map[string]interface{}{"status": emqxStatus, "broker": pkg.BrokerURL()},，完成对应的业务处理

		"emqx":    map[string]interface{}{"status": emqxStatus, "broker": pkg.BrokerURL()},
		// 执行 "commands_total": len(cmds), "ota_total": len(otas),，完成对应的业务处理

		"commands_total": len(cmds), "ota_total": len(otas),
		// 执行 "topic_spec": TopicSpec(), "acl_spec": AclSpec(),，完成对应的业务处理

		"topic_spec": TopicSpec(), "acl_spec": AclSpec(),
	// 当前代码块结束

	}
// 当前代码块结束

}
