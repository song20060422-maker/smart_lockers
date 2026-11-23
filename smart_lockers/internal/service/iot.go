// @ticket IoT-MQTT设备接入：设备认证/上下行消息处理/指令下发/OTA 业务层
// 严格对齐《智能快递柜IoT设备MQTT接入技术方案》v1.0：
//
//	2.1 设备生命周期（出厂→已注册→待激活→在线/离线→维护/停用/报废）
//	2.2 认证流程（SN+Secret 换 24h Token；EMQX HTTP Auth Hook 校验）
//	3.2/3.3 Topic 规范（device/# 上行、platform/# 下行）
//	5.1~5.3 指令下发（command_id 幂等、超时重试：开门/关门 10s 重试 1、重启 30s 不重试、配置 60s 重试 2）
//	3.3 OTA 升级（platform/ota/{id}/upgrade 下发、device/ota/{id}/progress 进度上报）
package service

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"

	"smart_lockers/internal/config"
	"smart_lockers/internal/model"
	"smart_lockers/internal/pkg"
)

// mqttTopicDevicePrefix 上行 Topic 前缀（设备→云端）
// 声明变量 mqttTopicDevicePrefix = "device/"，用于保存运行过程中的数据

const mqttTopicDevicePrefix = "device/"

// mqttTopicPlatformPrefix 下行 Topic 前缀（云端→设备）
// 声明变量 mqttTopicPlatformPrefix = "platform/"，用于保存运行过程中的数据

const mqttTopicPlatformPrefix = "platform/"

// iotTokenTTL 设备 MQTT 令牌有效期：默认 24 小时（方案 2.2），可由 config.mqtt.token_expire_hours 覆盖
// 声明变量 iotTokenTTL，用于保存运行过程中的数据

var iotTokenTTL = 24 * time.Hour

// cmdTimer 定时器条目：用于指令超时监控（方案 5.2）
// 定义结构体 cmdTimer：用于承载一组相关的字段数据

type cmdTimer struct {
	// 结构体字段 timer，类型为 *time.Timer

	timer *time.Timer
// 当前代码块结束

}

// 声明变量 cmdTimers/cmdTimersMu，用于保存运行过程中的数据

var (
	// 执行 cmdTimers = make(map[string]*cmdTimer)，完成对应的业务处理

	cmdTimers = make(map[string]*cmdTimer)
	// 执行 cmdTimersMu sync.Mutex，完成对应的业务处理

	cmdTimersMu sync.Mutex
// 结束参数/表达式列表

)

// randomHex 生成 n 字节随机数的十六进制串（Token/指令号生成用）
// 定义函数 randomHex：封装一段可复用的处理逻辑（参数：n int）string

func randomHex(n int) string {
	// 定义并初始化局部变量 b，取值为 make([]byte, n)

	b := make([]byte, n)
	// 条件判断：_, err := rand.Read(b); err != nil 时执行下面逻辑

	if _, err := rand.Read(b); err != nil {
		// 兜底：时间戳毫秒值（理论上 crypto/rand 不失败）
		// 返回 fmt.Sprintf("%x", time.Now().UnixNano()) 给调用方，结束当前函数

		return fmt.Sprintf("%x", time.Now().UnixNano())
	// 当前代码块结束

	}
	// 返回 hex.EncodeToString(b) 给调用方，结束当前函数

	return hex.EncodeToString(b)
// 当前代码块结束

}

// deviceBySN 按 SN 查询设备
// 定义函数 deviceBySN：封装一段可复用的处理逻辑（参数：sn string）(*model.IoTDevice, error)

func deviceBySN(sn string) (*model.IoTDevice, error) {
	// 条件判断：pkg.DB == nil 时执行下面逻辑

	if pkg.DB == nil {
		// 返回 nil, fmt.Errorf("database unavailable") 给调用方，结束当前函数

		return nil, fmt.Errorf("database unavailable")
	// 当前代码块结束

	}
	// 声明变量 d model.IoTDevice，用于保存运行过程中的数据

	var d model.IoTDevice
	// 条件判断：err := pkg.DB.Where("device_sn = ?", sn).First(&d).Error; err != nil 时执行下面逻辑

	if err := pkg.DB.Where("device_sn = ?", sn).First(&d).Error; err != nil {
		// 返回 nil, err 给调用方，结束当前函数

		return nil, err
	// 当前代码块结束

	}
	// 返回 &d, nil 给调用方，结束当前函数

	return &d, nil
// 当前代码块结束

}

// DeviceAuth 设备请求 MQTT 凭证：SN+出厂密钥校验通过后签发 24h Token（方案 2.2 认证接口）
// 定义函数 DeviceAuth：封装一段可复用的处理逻辑（参数：sn, secret, firmware, ip string）(map[string]interface{}, int)

func DeviceAuth(sn, secret, firmware, ip string) (map[string]interface{}, int) {
	// 条件判断：sn == "" || secret == "" 时执行下面逻辑

	if sn == "" || secret == "" {
		// 返回 nil, pkg.ErrBad 给调用方，结束当前函数

		return nil, pkg.ErrBad
	// 当前代码块结束

	}
	// 定义并初始化局部变量 d, err，取值为 deviceBySN(sn)

	d, err := deviceBySN(sn)
	// 条件判断：err != nil 时执行下面逻辑

	if err != nil {
		// 返回 nil, pkg.ErrDeviceNF 给调用方，结束当前函数

		return nil, pkg.ErrDeviceNF
	// 当前代码块结束

	}
	// 出厂密钥校验 @ticket IoT-MQTT设备接入（方案 2.2 步骤 2）
	// 条件判断：d.DeviceSecret == "" || d.DeviceSecret != secret 时执行下面逻辑

	if d.DeviceSecret == "" || d.DeviceSecret != secret {
		// 执行 pkg.Logf("[AUDIT] @ticket IoT device auth secret mismatch: sn=%s ip=%s", sn, ip)，完成对应的业务处理

		pkg.Logf("[AUDIT] @ticket IoT device auth secret mismatch: sn=%s ip=%s", sn, ip)
		// 返回 nil, pkg.ErrDeviceSecret 给调用方，结束当前函数

		return nil, pkg.ErrDeviceSecret
	// 当前代码块结束

	}
	// 条件判断：d.Status == model.IoTDevDisabled || d.Status == model.IoTDevScrapped 时执行下面逻辑

	if d.Status == model.IoTDevDisabled || d.Status == model.IoTDevScrapped {
		// 已停用/已报废设备拒绝签发凭证（方案 2.1 生命周期）
		// 返回 nil, pkg.ErrDeviceState 给调用方，结束当前函数

		return nil, pkg.ErrDeviceState
	// 当前代码块结束

	}
	// 定义并初始化局部变量 token，取值为 randomHex(32)

	token := randomHex(32)
	// 令牌有效期：配置 token_expire_hours 优先，缺省 24h（方案 2.2）@ticket IoT-MQTT设备接入
	// 定义并初始化局部变量 ttl，取值为 iotTokenTTL

	ttl := iotTokenTTL
	// 条件判断：h := config.Get().MQTT.TokenExpireHours; h > 0 时执行下面逻辑

	if h := config.Get().MQTT.TokenExpireHours; h > 0 {
		// 为 ttl 赋值：time.Duration(h) * time.Hour

		ttl = time.Duration(h) * time.Hour
	// 当前代码块结束

	}
	// 定义并初始化局部变量 expireAt，取值为 time.Now().Add(ttl)

	expireAt := time.Now().Add(ttl)
	// 保存 Token 记录（方案 2.2 步骤 5）@ticket IoT-MQTT设备接入
	// 条件判断：err := pkg.DB.Create(&model.DeviceToken{DeviceID: d.ID, Token: token, ExpireAt: expireAt}).Error; err != nil 时执行下面逻辑

	if err := pkg.DB.Create(&model.DeviceToken{DeviceID: d.ID, Token: token, ExpireAt: expireAt}).Error; err != nil {
		// 执行 pkg.Logf("[WARN] @ticket IoT device token save failed: %v", err)，完成对应的业务处理

		pkg.Logf("[WARN] @ticket IoT device token save failed: %v", err)
		// 返回 nil, pkg.ErrInternal 给调用方，结束当前函数

		return nil, pkg.ErrInternal
	// 当前代码块结束

	}
	// 记录固件版本与最近 IP（可选字段，方案 2.2 请求参数）
	// 条件判断：firmware != "" 时执行下面逻辑

	if firmware != "" {
		// 执行 pkg.DB.Model(&model.IoTDevice{}).Where("id = ?", d.ID).Updates(map[string]interface{}{"firmware_version": firmware, "ip_address": ip, "updated_at": time.Now()})，完成对应的业务处理

		pkg.DB.Model(&model.IoTDevice{}).Where("id = ?", d.ID).Updates(map[string]interface{}{"firmware_version": firmware, "ip_address": ip, "updated_at": time.Now()})
	// 当前代码块结束

	}
	// 执行 pkg.Logf("[AUDIT] @ticket IoT device token issued: sn=%s device_id=%d expire=%s", sn, d.ID, expireAt.Format(time.RFC3339))，完成对应的业务处理

	pkg.Logf("[AUDIT] @ticket IoT device token issued: sn=%s device_id=%d expire=%s", sn, d.ID, expireAt.Format(time.RFC3339))
	// 返回凭证报文（broker_url/client_id/username/password/token_expire_at，方案 2.2 响应参数）给调用方，结束当前函数

	return map[string]interface{}{
		// 执行 "broker_url": cfgMQTTBrokerURL(),，完成对应的业务处理

		"broker_url": cfgMQTTBrokerURL(),
		// 执行 "client_id": fmt.Sprintf("device_%d", d.ID)，完成对应的业务处理

		"client_id":  fmt.Sprintf("device_%d", d.ID),
		// 执行 "username": sn，完成对应的业务处理

		"username":   sn,
		// 执行 "password": token，完成对应的业务处理

		"password":   token,
		// 执行 "token_expire_at": expireAt.Unix()，完成对应的业务处理

		"token_expire_at": expireAt.Unix(),
		// 执行 "device_id": d.ID，完成对应的业务处理

		"device_id":  d.ID,
	// 当前代码块结束

	}, pkg.CodeOK
// 当前代码块结束

}

// cfgMQTTBrokerURL 从配置拼接 broker_url（方案 2.2 响应字段）
// 定义函数 cfgMQTTBrokerURL：封装一段可复用的处理逻辑（参数：）string

func cfgMQTTBrokerURL() string {
	// 返回 config.Get().MQTT.Broker 给调用方，结束当前函数

	return config.Get().MQTT.Broker
// 当前代码块结束

}

// VerifyEMQXToken 校验设备 Token（EMQX HTTP Auth Hook 调用入口；username=SN，password=Token）
// 定义函数 VerifyEMQXToken：封装一段可复用的处理逻辑（参数：username, password string）bool

func VerifyEMQXToken(username, password string) bool {
	// 条件判断：pkg.DB == nil || username == "" || password == "" 时执行下面逻辑

	if pkg.DB == nil || username == "" || password == "" {
		// 返回 false 给调用方，结束当前函数

		return false
	// 当前代码块结束

	}
	// 声明变量 t model.DeviceToken，用于保存运行过程中的数据

	var t model.DeviceToken
	// 条件判断：err := pkg.DB.Where("token = ?", password).First(&t).Error; err != nil 时执行下面逻辑

	if err := pkg.DB.Where("token = ?", password).First(&t).Error; err != nil {
		// 返回 false 给调用方，结束当前函数

		return false
	// 当前代码块结束

	}
	// Token 归属校验：username 必须是该 Token 所属设备的 SN @ticket IoT-MQTT设备接入
	// 定义并初始化局部变量 d, err，取值为 deviceBySN(username)

	d, err := deviceBySN(username)
	// 条件判断：err != nil || d.ID != t.DeviceID 时执行下面逻辑

	if err != nil || d.ID != t.DeviceID {
		// 返回 false 给调用方，结束当前函数

		return false
	// 当前代码块结束

	}
	// 有效期与吊销校验（方案 2.2 步骤 9~10）
	// 返回 (!t.Revoked && time.Now().Before(t.ExpireAt)) 给调用方，结束当前函数

	return !t.Revoked && time.Now().Before(t.ExpireAt)
// 当前代码块结束

}

// HandleConnected 设备上线处理（EMQX webhook client.connected 或上行 online 消息）@ticket IoT-MQTT设备接入
// 定义函数 HandleConnected：封装一段可复用的处理逻辑（参数：deviceID uint, ip string）

func HandleConnected(deviceID uint, ip string) {
	// 条件判断：pkg.DB == nil 时执行下面逻辑

	if pkg.DB == nil {
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 定义并初始化局部变量 nowT，取值为 time.Now()

	nowT := time.Now()
	// 更新设备状态为在线（方案 2.5 状态同步机制）
	// 执行 pkg.DB.Model(&model.IoTDevice{}).Where("id = ?", deviceID).Updates(...)，完成对应的业务处理

	pkg.DB.Model(&model.IoTDevice{}).Where("id = ?", deviceID).Updates(map[string]interface{}{
		// 执行 "status": model.IoTDevOnline,，完成对应的业务处理

		"status": model.IoTDevOnline,
		// 执行 "last_online_at": nowT,，完成对应的业务处理

		"last_online_at": nowT,
		// 执行 "updated_at": nowT,，完成对应的业务处理

		"updated_at": nowT,
	// 当前代码块结束

	})
	// 条件判断：ip != "" 时执行下面逻辑

	if ip != "" {
		// 执行 pkg.DB.Model(&model.IoTDevice{}).Where("id = ?", deviceID).Update("ip_address", ip)，完成对应的业务处理

		pkg.DB.Model(&model.IoTDevice{}).Where("id = ?", deviceID).Update("ip_address", ip)
	// 当前代码块结束

	}
	// 执行 pkg.Logf("[INFO] @ticket IoT device online: device_id=%d ip=%s", deviceID, ip)，完成对应的业务处理

	pkg.Logf("[INFO] @ticket IoT device online: device_id=%d ip=%s", deviceID, ip)
// 当前代码块结束

}

// HandleDisconnected 设备离线处理（EMQX webhook client.disconnected 或上行 offline 消息）@ticket IoT-MQTT设备接入
// 定义函数 HandleDisconnected：封装一段可复用的处理逻辑（参数：deviceID uint）

func HandleDisconnected(deviceID uint) {
	// 条件判断：pkg.DB == nil 时执行下面逻辑

	if pkg.DB == nil {
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 定义并初始化局部变量 nowT，取值为 time.Now()

	nowT := time.Now()
	// 更新设备状态为离线（方案 2.5：离线超过 5 分钟由监控侧触发告警）

	pkg.DB.Model(&model.IoTDevice{}).Where("id = ?", deviceID).Updates(map[string]interface{}{
		// 执行 "status": model.IoTDevOffline,，完成对应的业务处理

		"status": model.IoTDevOffline,
		// 执行 "updated_at": nowT,，完成对应的业务处理

		"updated_at": nowT,
	// 当前代码块结束

	})
	// 执行 pkg.Logf("[INFO] @ticket IoT device offline: device_id=%d", deviceID)，完成对应的业务处理

	pkg.Logf("[INFO] @ticket IoT device offline: device_id=%d", deviceID)
// 当前代码块结束

}

// HandleHeartbeat 心跳上报：刷新 last_heartbeat（方案 4.2，每 5 分钟 QoS0）
// 定义函数 HandleHeartbeat：封装一段可复用的处理逻辑（参数：deviceID uint, payload string）

func HandleHeartbeat(deviceID uint, payload string) {
	// 条件判断：pkg.DB == nil 时执行下面逻辑

	if pkg.DB == nil {
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 定义并初始化局部变量 nowT，取值为 time.Now()

	nowT := time.Now()
	// 更新最近心跳时间（方案 4.2 步骤 3）
	// 执行 pkg.DB.Model(&model.IoTDevice{}).Where("id = ?", deviceID).Updates(...)，完成对应的业务处理

	pkg.DB.Model(&model.IoTDevice{}).Where("id = ?", deviceID).Updates(map[string]interface{}{
		// 执行 "last_heartbeat": nowT,，完成对应的业务处理

		"last_heartbeat": nowT,
		// 执行 "updated_at": nowT,，完成对应的业务处理

		"updated_at": nowT,
	// 当前代码块结束

	})
	// 心跳报文留存为事件（供 B端 中台查询）
	// 执行 saveDeviceEvent(deviceID, "", "device/status/"+strconv.FormatUint(uint64(deviceID), 10)+"/heartbeat", "status", "info", payload)，完成对应的业务处理

	saveDeviceEvent(deviceID, "", "device/status/"+strconv.FormatUint(uint64(deviceID), 10)+"/heartbeat", "status", "info", payload)
// 当前代码块结束

}

// lockerPayload 柜格状态上报报文（方案 3.2 device/status/{id}/locker）
// 定义结构体 lockerPayload：用于承载一组相关的字段数据

type lockerPayload struct {
	// 结构体字段 Timestamp，类型为 int64，附带标签 `json:"timestamp"`

	Timestamp int64 `json:"timestamp"`
	// 结构体字段 Lockers，类型为 []lockerItem，附带标签 `json:"lockers"`

	Lockers   []lockerItem `json:"lockers"`
// 当前代码块结束

}

// lockerItem 单柜格状态项
// 定义结构体 lockerItem：用于承载一组相关的字段数据

type lockerItem struct {
	// 结构体字段 LockerNo，类型为 string，附带标签 `json:"locker_no"`

	LockerNo string `json:"locker_no"`
	// 结构体字段 Status，类型为 int，附带标签 `json:"status"`；1空闲 2占用 3故障（方案枚举）

	Status   int    `json:"status"`
	// 结构体字段 OrderID，类型为 json.Number，附带标签 `json:"order_id"`；可为 null，兼容整型/字符串

	OrderID  json.Number `json:"order_id"`
// 当前代码块结束

}

// HandleLockerReport 柜格状态上报落库：同步 LockerDoor 与柜体可用数（方案 4.1 状态上报流程）
// 定义函数 HandleLockerReport：封装一段可复用的处理逻辑（参数：deviceID uint, payload string）

func HandleLockerReport(deviceID uint, payload string) {
	// 条件判断：pkg.DB == nil 时执行下面逻辑

	if pkg.DB == nil {
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 定义并初始化局部变量 lp lockerPayload，取值为零值并解析报文

	var lp lockerPayload
	// 条件判断：err := json.Unmarshal([]byte(payload), &lp); err != nil 时执行下面逻辑

	if err := json.Unmarshal([]byte(payload), &lp); err != nil {
		// 执行 pkg.Logf("[WARN] @ticket IoT locker payload invalid: device_id=%d err=%v", deviceID, err)，完成对应的业务处理

		pkg.Logf("[WARN] @ticket IoT locker payload invalid: device_id=%d err=%v", deviceID, err)
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 声明变量 d model.IoTDevice，用于保存运行过程中的数据

	var d model.IoTDevice
	// 条件判断：err := pkg.DB.First(&d, deviceID).Error; err != nil 时执行下面逻辑

	if err := pkg.DB.First(&d, deviceID).Error; err != nil {
		// 执行 pkg.Logf("[WARN] @ticket IoT locker report device not found: device_id=%d", deviceID)，完成对应的业务处理

		pkg.Logf("[WARN] @ticket IoT locker report device not found: device_id=%d", deviceID)
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 逐柜格同步状态：方案枚举 1空闲/2占用/3故障 → 本库 0空闲/1占用/2故障 @ticket IoT-MQTT设备接入
	// 循环处理：for _, it := range lp.Lockers，逐轮迭代

	for _, it := range it_status(&lp) {
		// 映射方案状态到本地柜门状态
		// 定义并初始化局部变量 doorStatus，取值为 mapIoTStatusToLocal(it.Status)

		doorStatus := mapIoTStatusToLocal(it.Status)
		// 按柜体绑定 + 柜格编号定位柜门记录
		// 定义并初始化局部变量 res，取值为 pkg.DB.Model(&model.LockerDoor{}).Where(...)，完成对应的业务处理

		res := pkg.DB.Model(&model.LockerDoor{}).
			Where("cabinet_id = ? AND code = ?", d.CabinetID, it.LockerNo).
			Updates(map[string]interface{}{"status": doorStatus, "updated_at": time.Now()})
		// 条件判断：res.Error != nil 时执行下面逻辑

		if res.Error != nil {
			// 执行 pkg.Logf("[WARN] @ticket IoT locker update failed: door=%s err=%v", it.LockerNo, res.Error)，完成对应的业务处理

			pkg.Logf("[WARN] @ticket IoT locker update failed: door=%s err=%v", it.LockerNo, res.Error)
		// 当前代码块结束

		}
	// 当前代码块结束

	}
	// 刷新柜体可用数：可用 = 空闲柜门数 @ticket IoT-MQTT设备接入
	// 执行 refreshCabinetAvailable(d.CabinetID)，完成对应的业务处理

	refreshCabinetAvailable(d.CabinetID)
	// 报文留存为事件
	// 执行 saveDeviceEvent(deviceID, d.DeviceSN, "device/status/"+strconv.FormatUint(uint64(deviceID), 10)+"/locker", "status", "info", payload)，完成对应的业务处理

	saveDeviceEvent(deviceID, d.DeviceSN, "device/status/"+strconv.FormatUint(uint64(deviceID), 10)+"/locker", "status", "info", payload)
// 当前代码块结束

}

// it_status 兼容辅助：遍历柜格列表（独立函数便于单测桩替换）
// 定义函数 it_status：封装一段可复用的处理逻辑（参数：lp *lockerPayload）[]lockerItem

func it_status(lp *lockerPayload) []lockerItem {
	// 返回 lp.Lockers 给调用方，结束当前函数

	return lp.Lockers
// 当前代码块结束

}

// mapIoTStatusToLocal 方案柜格状态枚举 → 本地柜门状态映射
// 定义函数 mapIoTStatusToLocal：封装一段可复用的处理逻辑（参数：s int）int

func mapIoTStatusToLocal(s int) int {
	// 多分支判断：根据 s 的值选择分支

	switch s {
	// 当取值匹配 1 时进入该分支（方案 1 空闲 → 本地 0）

	case 1:
		// 返回 0 给调用方，结束当前函数

		return 0
	// 当取值匹配 2 时进入该分支（方案 2 占用 → 本地 1）

	case 2:
		// 返回 1 给调用方，结束当前函数

		return 1
	// 当取值匹配 3 时进入该分支（方案 3 故障 → 本地 2）

	case 3:
		// 返回 2 给调用方，结束当前函数

		return 2
// 当前代码块结束

	}
	// 未知状态按故障处理，避免误开放
	// 返回 2 给调用方，结束当前函数

	return 2
// 当前代码块结束

}

// refreshCabinetAvailable 重算柜体可用数
// 定义函数 refreshCabinetAvailable：封装一段可复用的处理逻辑（参数：cabinetID uint）

func refreshCabinetAvailable(cabinetID uint) {
	// 条件判断：cabinetID == 0 时执行下面逻辑

	if cabinetID == 0 {
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 声明变量 avail int64，用于保存运行过程中的数据

	var avail int64
	// 统计空闲柜门数
	// 执行 pkg.DB.Model(&model.LockerDoor{}).Where("cabinet_id = ? AND status = 0", cabinetID).Count(&avail)，完成对应的业务处理

	pkg.DB.Model(&model.LockerDoor{}).Where("cabinet_id = ? AND status = 0", cabinetID).Count(&avail)
	// 声明变量 total int64，用于保存运行过程中的数据

	var total int64
	// 统计柜门总数
	// 执行 pkg.DB.Model(&model.LockerDoor{}).Where("cabinet_id = ?", cabinetID).Count(&total)，完成对应的业务处理

	pkg.DB.Model(&model.LockerDoor{}).Where("cabinet_id = ?", cabinetID).Count(&total)
	// 执行 pkg.DB.Model(&model.LockerCabinet{}).Where("id = ?", cabinetID).Updates(map[string]interface{}{"available": avail, "total": total, "updated_at": time.Now()})，完成对应的业务处理

	pkg.DB.Model(&model.LockerCabinet{}).Where("id = ?", cabinetID).Updates(map[string]interface{}{"available": avail, "total": total, "updated_at": time.Now()})
// 当前代码块结束

}

// HandleDeviceLog 设备日志上报留存（方案 4.3：INFO/WARN/ERROR 实时上报；ERROR 触发告警）
// 定义函数 HandleDeviceLog：封装一段可复用的处理逻辑（参数：deviceID uint, topic, payload string, isError bool）

func HandleDeviceLog(deviceID uint, topic, payload string, isError bool) {
	// 定义并初始化局部变量 lvl，取值为 "info"

	lvl := "info"
	// 条件判断：isError 时执行下面逻辑

	if isError {
		// 为 lvl 赋值："warn"

		lvl = "warn"
	// 当前代码块结束

	}
	// 留存日志事件
	// 执行 saveDeviceEvent(deviceID, "", topic, "log", lvl, payload)，完成对应的业务处理

	saveDeviceEvent(deviceID, "", topic, "log", lvl, payload)
	// 错误日志触发设备异常场景通知（方案 4.3 流程 + 工单01 设备异常→网点管理员）@ticket IoT-MQTT设备接入
	// 条件判断：isError 时执行下面逻辑

	if isError {
		// 执行 sceneDeviceAlert(deviceID, topic, payload)，完成对应的业务处理

		sceneDeviceAlert(deviceID, topic, payload)
	// 当前代码块结束

	}
// 当前代码块结束

}

// HandleDeviceAlert 设备告警上报：留事件 + 高/严重级别联动场景通知（方案 4.4 告警流程）
// 定义函数 HandleDeviceAlert：封装一段可复用的处理逻辑（参数：deviceID uint, topic, payload string, level string）

func HandleDeviceAlert(deviceID uint, topic, payload string, level string) {
	// 留存告警事件
	// 执行 saveDeviceEvent(deviceID, "", topic, "alert", level, payload)，完成对应的业务处理
	// 告警级别联动：高(high)/严重(critical) 发送紧急场景通知，中(warn) 短信（此处统一走场景通知由渠道开关控制）@ticket IoT-MQTT设备接入

	saveDeviceEvent(deviceID, "", topic, "alert", level, payload)
	// 执行 sceneDeviceAlert(deviceID, topic, payload)，完成对应的业务处理

	sceneDeviceAlert(deviceID, topic, payload)
// 当前代码块结束

}

// sceneDeviceAlert 联动工单01 设备异常→网点管理员 场景通知 @ticket IoT-MQTT设备接入
// 定义函数 sceneDeviceAlert：封装一段可复用的处理逻辑（参数：deviceID uint, topic, payload string）

func sceneDeviceAlert(deviceID uint, topic, payload string) {
	// 条件判断：pkg.DB == nil 时执行下面逻辑

	if pkg.DB == nil {
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 声明变量 d model.IoTDevice，用于保存运行过程中的数据

	var d model.IoTDevice
	// 条件判断：err := pkg.DB.First(&d, deviceID).Error; err != nil 时执行下面逻辑

	if err := pkg.DB.First(&d, deviceID).Error; err != nil {
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 定义并初始化局部变量 _, code，取值为 SceneNotify("设备异常→网点管理员", "admin", d.SiteID, map[string]string{"device_sn": d.DeviceSN, "topic": topic, "message": payload})

	_, code := SceneNotify("设备异常→网点管理员", "admin", d.SiteID, map[string]string{"device_sn": d.DeviceSN, "topic": topic, "message": payload})
	// 条件判断：code != pkg.CodeOK 时执行下面逻辑

	if code != pkg.CodeOK {
		// 执行 pkg.Logf("[WARN] @ticket IoT scene notify failed: device_id=%d code=%d", deviceID, code)，完成对应的业务处理

		pkg.Logf("[WARN] @ticket IoT scene notify failed: device_id=%d code=%d", deviceID, code)
	// 当前代码块结束

	}
// 当前代码块结束

}

// saveDeviceEvent 留存设备上行事件
// 定义函数 saveDeviceEvent：封装一段可复用的处理逻辑（参数：deviceID uint, sn, topic, kind, level, payload string）

func saveDeviceEvent(deviceID uint, sn, topic, kind, level, payload string) {
	// 条件判断：pkg.DB == nil 时执行下面逻辑

	if pkg.DB == nil {
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 条件判断：err := pkg.DB.Create(&model.DeviceEvent{...}).Error; err != nil 时执行下面逻辑

	if err := pkg.DB.Create(&model.DeviceEvent{
		// 执行 DeviceID: deviceID,，完成对应的业务处理

		DeviceID: deviceID,
		// 执行 DeviceSN: sn,，完成对应的业务处理

		DeviceSN: sn,
		// 执行 Topic: topic,，完成对应的业务处理

		Topic:    topic,
		// 执行 Kind: kind,，完成对应的业务处理

		Kind:     kind,
		// 执行 Level: level,，完成对应的业务处理

		Level:    level,
		// 执行 Content: payload,，完成对应的业务处理

		Content:  payload,
		// 执行 CreatedAt: time.Now(),，完成对应的业务处理

		CreatedAt: time.Now(),
	// 当前代码块结束

	}).Error; err != nil {
		// 执行 pkg.Logf("[WARN] @ticket IoT event save failed: %v", err)，完成对应的业务处理

		pkg.Logf("[WARN] @ticket IoT event save failed: %v", err)
	// 当前代码块结束

	}
// 当前代码块结束

}

// cmdRetryPolicy 指令超时/重试策略（方案 5.2 超时策略表）@ticket IoT-MQTT设备接入
// 定义函数 cmdRetryPolicy：封装一段可复用的处理逻辑（参数：cmdType string）(timeoutSec, maxRetry int)

func cmdRetryPolicy(cmdType string) (int, int) {
	// 多分支判断：根据 cmdType 的值选择分支

	switch cmdType {
	// 开门/关门：10s 超时，重试 1 次

	case "open_door", "close_door":
		// 返回 10, 1 给调用方，结束当前函数

		return 10, 1
	// 重启：30s 超时，不重试

	case "reboot":
		// 返回 30, 0 给调用方，结束当前函数

		return 30, 0
	// 配置下发：60s 超时，重试 2 次

	case "config":
		// 返回 60, 2 给调用方，结束当前函数

		return 60, 2
// 当前代码块结束

	}
	// 兜底：10s 超时，不重试
	// 返回 10, 0 给调用方，结束当前函数

	return 10, 0
// 当前代码块结束

}

// buildCommandPayload 组装下行指令报文（方案 5.1 open_door 报文格式）
// 定义函数 buildCommandPayload：封装一段可复用的处理逻辑（参数：commandID string, deviceID uint, lockerNo, orderNo string, timeoutSec int, extra map[string]interface{}）string

func buildCommandPayload(commandID string, deviceID uint, lockerNo, orderNo string, timeoutSec int, extra map[string]interface{}) string {
	// 定义并初始化局部变量 m，取值为 map[string]interface{}{...}

	m := map[string]interface{}{
		// 执行 "command_id": commandID,，完成对应的业务处理

		"command_id": commandID,
		// 执行 "device_id": deviceID,，完成对应的业务处理

		"device_id":  deviceID,
		// 执行 "timestamp": time.Now().UnixMilli(),，完成对应的业务处理

		"timestamp":  time.Now().UnixMilli(),
		// 执行 "timeout": timeoutSec,，完成对应的业务处理

		"timeout":    timeoutSec,
	// 当前代码块结束

	}
	// 条件判断：lockerNo != "" 时执行下面逻辑

	if lockerNo != "" {
		// 执行 m["locker_no"] = lockerNo，完成对应的业务处理

		m["locker_no"] = lockerNo
	// 当前代码块结束

	}
	// 条件判断：orderNo != "" 时执行下面逻辑

	if orderNo != "" {
		// 执行 m["order_id"] = orderNo，完成对应的业务处理

		m["order_id"] = orderNo
	// 当前代码块结束

	}
	// 循环处理：for k, v := range extra，逐轮迭代（合并扩展字段，如配置下发 price_config）

	for k, v := range extra {
		// 执行 m[k] = v，完成对应的业务处理

		m[k] = v
	// 当前代码块结束

	}
	// 定义并初始化局部变量 b, err，取值为 json.Marshal(m)

	b, err := json.Marshal(m)
	// 条件判断：err != nil 时执行下面逻辑

	if err != nil {
		// 返回 "{}" 给调用方，结束当前函数

		return "{}"
	// 当前代码块结束

	}
	// 返回 string(b) 给调用方，结束当前函数

	return string(b)
// 当前代码块结束

}

// DispatchCommand 下发平台指令：落库 + 发布 platform/command/{id}/{type} + 启动超时监控（方案 5.1~5.3）
// 定义函数 DispatchCommand：封装一段可复用的处理逻辑（参数：sn, cmdType, lockerNo, orderNo, operator string, extra map[string]interface{}）(map[string]interface{}, int)

func DispatchCommand(sn, cmdType, lockerNo, orderNo, operator string, extra map[string]interface{}) (map[string]interface{}, int) {
	// 定义并初始化局部变量 d, err，取值为 deviceBySN(sn)

	d, err := deviceBySN(sn)
	// 条件判断：err != nil 时执行下面逻辑

	if err != nil {
		// 返回 nil, pkg.ErrDeviceNF 给调用方，结束当前函数

		return nil, pkg.ErrDeviceNF
	// 当前代码块结束

	}
	// 设备在线校验：离线设备拒绝下发（方案 2.5 状态管理）@ticket IoT-MQTT设备接入
	// 条件判断：d.Status != model.IoTDevOnline 时执行下面逻辑

	if d.Status != model.IoTDevOnline {
		// 返回 nil, pkg.ErrDeviceState 给调用方，结束当前函数

		return nil, pkg.ErrDeviceState
	// 当前代码块结束

	}
	// 读取超时/重试策略
	// 定义并初始化局部变量 timeoutSec, maxRetry，取值为 cmdRetryPolicy(cmdType)

	timeoutSec, maxRetry := cmdRetryPolicy(cmdType)
	// 生成全局唯一指令号（command_id 幂等依据，方案 5.3）
	// 定义并初始化局部变量 commandID，取值为 fmt.Sprintf("CMD%s%s", time.Now().Format("20060102150405"), strings.ToUpper(randomHex(4)))

	commandID := fmt.Sprintf("CMD%s%s", time.Now().Format("20060102150405"), strings.ToUpper(randomHex(4)))
	// 组装报文
	// 定义并初始化局部变量 payload，取值为 buildCommandPayload(commandID, d.ID, lockerNo, orderNo, timeoutSec, extra)

	payload := buildCommandPayload(commandID, d.ID, lockerNo, orderNo, timeoutSec, extra)
	// 指令落库 @ticket IoT-MQTT设备接入（资金/柜门动作接口必须有事务保护与审计）
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
		// 执行 OrderNo: orderNo,，完成对应的业务处理

		OrderNo:   orderNo,
		// 执行 Payload: payload,，完成对应的业务处理

		Payload:   payload,
		// 执行 TimeoutSec: timeoutSec,，完成对应的业务处理

		TimeoutSec: timeoutSec,
		// 执行 MaxRetry: maxRetry,，完成对应的业务处理

		MaxRetry:  maxRetry,
		// 执行 Status: model.CmdPending,，完成对应的业务处理

		Status:    model.CmdPending,
		// 执行 Operator: operator,，完成对应的业务处理

		Operator:  operator,
		// 执行 CreatedAt: time.Now(),，完成对应的业务处理

		CreatedAt: time.Now(),
	// 当前代码块结束

	}
	// 条件判断：err := pkg.DB.Create(rec).Error; err != nil 时执行下面逻辑

	if err := pkg.DB.Create(rec).Error; err != nil {
		// 执行 pkg.Logf("[WARN] @ticket IoT command save failed: %v", err)，完成对应的业务处理

		pkg.Logf("[WARN] @ticket IoT command save failed: %v", err)
		// 返回 nil, pkg.ErrInternal 给调用方，结束当前函数

		return nil, pkg.ErrInternal
	// 当前代码块结束

	}
	// 发布指令（Broker 不可用时指令保持待下发，可由定时补偿重发；演示环境可用 HTTP 兼容通道触发）
	// 定义并初始化局部变量 topic，取值为 fmt.Sprintf("%scommand/%d/%s", mqttTopicPlatformPrefix, d.ID, cmdType)

	topic := fmt.Sprintf("%scommand/%d/%s", mqttTopicPlatformPrefix, d.ID, cmdType)
	// 定义并初始化局部变量 pubErr，取值为 pkg.MqttPublish(topic, 1, []byte(payload))（QoS1，方案 5.1）

	pubErr := pkg.MqttPublish(topic, 1, []byte(payload))
	// 条件判断：pubErr == nil 时执行下面逻辑

	if pubErr == nil {
		// 执行 markCommandSent(rec)，完成对应的业务处理

		markCommandSent(rec)
		// 执行 scheduleCommandTimeout(rec)，完成对应的业务处理

		scheduleCommandTimeout(rec)
	// 上述条件不成立时，执行 else 分支逻辑

	} else {
		// 执行 pkg.Logf("[WARN] @ticket IoT command publish deferred: cmd=%s err=%v", commandID, pubErr)，完成对应的业务处理

		pkg.Logf("[WARN] @ticket IoT command publish deferred: cmd=%s err=%v", commandID, pubErr)
	// 当前代码块结束

	}
	// 执行 pkg.Logf("[AUDIT] @ticket IoT command dispatched: cmd=%s type=%s sn=%s operator=%s topic=%s", commandID, cmdType, sn, operator, topic)，完成对应的业务处理

	pkg.Logf("[AUDIT] @ticket IoT command dispatched: cmd=%s type=%s sn=%s operator=%s topic=%s", commandID, cmdType, sn, operator, topic)
	// 返回指令回执给调用方，结束当前函数

	return map[string]interface{}{
		// 执行 "command_id": commandID,，完成对应的业务处理

		"command_id": commandID,
		// 执行 "device_id": d.ID,，完成对应的业务处理

		"device_id":  d.ID,
		// 执行 "cmd_type": cmdType,，完成对应的业务处理

		"cmd_type":   cmdType,
		// 执行 "timeout_sec": timeoutSec,，完成对应的业务处理

		"timeout_sec": timeoutSec,
		// 执行 "max_retry": maxRetry,，完成对应的业务处理

		"max_retry":  maxRetry,
		// 执行 "status": rec.Status,，完成对应的业务处理

		"status":     rec.Status,
		// 执行 "published": pubErr == nil,，完成对应的业务处理

		"published":  pubErr == nil,
	// 当前代码块结束

	}, pkg.CodeOK
// 当前代码块结束

}

// markCommandSent 标记指令为已下发待响应
// 定义函数 markCommandSent：封装一段可复用的处理逻辑（参数：rec *model.DeviceCommand）

func markCommandSent(rec *model.DeviceCommand) {
	// 执行 pkg.DB.Model(&model.DeviceCommand{}).Where("id = ?", rec.ID).Update("status", model.CmdSent)，完成对应的业务处理

	pkg.DB.Model(&model.DeviceCommand{}).Where("id = ?", rec.ID).Update("status", model.CmdSent)
	// 为 rec.Status 赋值：model.CmdSent

	rec.Status = model.CmdSent
// 当前代码块结束

}

// scheduleCommandTimeout 启动指令超时监控：超时未响应则按策略重试或判定超时（方案 5.2）
// 定义函数 scheduleCommandTimeout：封装一段可复用的处理逻辑（参数：rec *model.DeviceCommand）

func scheduleCommandTimeout(rec *model.DeviceCommand) {
	// 加锁登记定时器

	cmdTimersMu.Lock()
	// 延迟注册解锁

	defer cmdTimersMu.Unlock()
	// 定义并初始化局部变量 t，取值为 time.AfterFunc(time.Duration(rec.TimeoutSec)*time.Second, func() { onCommandTimeout(rec) })

	t := time.AfterFunc(time.Duration(rec.TimeoutSec)*time.Second, func() { onCommandTimeout(rec) })
	// 为 cmdTimers[rec.CommandID] 赋值：&cmdTimer{timer: t}

	cmdTimers[rec.CommandID] = &cmdTimer{timer: t}
// 当前代码块结束

}

// cancelCommandTimeout 取消指令超时监控（收到设备响应时调用）
// 定义函数 cancelCommandTimeout：封装一段可复用的处理逻辑（参数：commandID string）

func cancelCommandTimeout(commandID string) {
	// 加锁读取并移除定时器

	cmdTimersMu.Lock()
	// 延迟注册解锁

	defer cmdTimersMu.Unlock()
	// 定义并初始化局部变量 ct, ok，取值为 cmdTimers[commandID]

	ct, ok := cmdTimers[commandID]
	// 条件判断：ok 时执行下面逻辑

	if ok {
		// 执行 ct.timer.Stop()，完成对应的业务处理

		ct.timer.Stop()
		// 删除 map 条目

		delete(cmdTimers, commandID)
	// 当前代码块结束

	}
// 当前代码块结束

}

// onCommandTimeout 超时回调：可重试则重发指令，否则判定超时并告警（方案 5.2 超时处理流程）
// 定义函数 onCommandTimeout：封装一段可复用的处理逻辑（参数：rec *model.DeviceCommand）

func onCommandTimeout(rec *model.DeviceCommand) {
	// 条件判断：pkg.DB == nil 时执行下面逻辑

	if pkg.DB == nil {
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 声明变量 cur model.DeviceCommand，用于保存运行过程中的数据

	var cur model.DeviceCommand
	// 条件判断：err := pkg.DB.Where("command_id = ?", rec.CommandID).First(&cur).Error; err != nil 时执行下面逻辑

	if err := pkg.DB.Where("command_id = ?", rec.CommandID).First(&cur).Error; err != nil {
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 已有终态（成功/失败/超时）则直接返回
	// 条件判断：cur.Status != model.CmdSent && cur.Status != model.CmdPending 时执行下面逻辑

	if cur.Status != model.CmdSent && cur.Status != model.CmdPending {
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 条件判断：cur.RetryCount < cur.MaxRetry 时执行下面逻辑（可重试）

	if cur.RetryCount < cur.MaxRetry {
		// 重发指令：RetryCount+1 后重新发布并重新计时 @ticket IoT-MQTT设备接入
		// 执行 pkg.DB.Model(&model.DeviceCommand{}).Where("id = ?", cur.ID).Updates(map[string]interface{}{"retry_count": gorm.Expr("retry_count + 1"), "status": model.CmdPending})，完成对应的业务处理

		pkg.DB.Model(&model.DeviceCommand{}).Where("id = ?", cur.ID).Updates(map[string]interface{}{"retry_count": gorm.Expr("retry_count + 1"), "status": model.CmdPending})
		// 定义并初始化局部变量 topic，取值为 fmt.Sprintf("%scommand/%d/%s", mqttTopicPlatformPrefix, cur.DeviceID, cur.CmdType)

		topic := fmt.Sprintf("%scommand/%d/%s", mqttTopicPlatformPrefix, cur.DeviceID, cur.CmdType)
		// 条件判断：err := pkg.MqttPublish(topic, 1, []byte(cur.Payload)); err == nil 时执行下面逻辑

		if err := pkg.MqttPublish(topic, 1, []byte(cur.Payload)); err == nil {
			// 重发成功：标记已下发并重新调度超时
			// 定义并初始化局部变量 next，取值为 cur（副本）用于重新计时

			next := cur
			// 执行 markCommandSent(&next)，完成对应的业务处理

			markCommandSent(&next)
			// 执行 scheduleCommandTimeout(&next)，完成对应的业务处理

			scheduleCommandTimeout(&next)
			// 执行 pkg.Logf("[INFO] @ticket IoT command retried: cmd=%s retry=%d", cur.CommandID, cur.RetryCount+1)，完成对应的业务处理

			pkg.Logf("[INFO] @ticket IoT command retried: cmd=%s retry=%d", cur.CommandID, cur.RetryCount+1)
			// 返回，结束当前函数

			return
		// 当前代码块结束

		}
	// 当前代码块结束

	}
	// 重试耗尽或不可重试：判定超时并触发告警
	// 定义并初始化局部变量 nowT，取值为 time.Now()

	nowT := time.Now()
	// 执行 pkg.DB.Model(&model.DeviceCommand{}).Where("id = ?", cur.ID).Updates(map[string]interface{}{"status": model.CmdTimeout, "responded_at": &nowT})，完成对应的业务处理

	pkg.DB.Model(&model.DeviceCommand{}).Where("id = ?", cur.ID).Updates(map[string]interface{}{"status": model.CmdTimeout, "responded_at": &nowT})
	// 执行 pkg.Logf("[WARN] @ticket IoT command timeout: cmd=%s type=%s device=%d", cur.CommandID, cur.CmdType, cur.DeviceID)，完成对应的业务处理

	pkg.Logf("[WARN] @ticket IoT command timeout: cmd=%s type=%s device=%d", cur.CommandID, cur.CmdType, cur.DeviceID)
	// 超时联动工单01 场景通知（开门失败属于操作失败→管理员 / 设备异常→网点管理员）
	// 执行 sceneDeviceAlert(cur.DeviceID, fmt.Sprintf("platform/command/%d/%s", cur.DeviceID, cur.CmdType), "指令超时未响应: "+cur.CommandID)，完成对应的业务处理

	sceneDeviceAlert(cur.DeviceID, fmt.Sprintf("platform/command/%d/%s", cur.DeviceID, cur.CmdType), "指令超时未响应: "+cur.CommandID)
// 当前代码块结束

}

// HandleCommandResponse 处理设备响应（device/response/{id}/{type}）：更新指令终态并联动订单业务（方案 5.1 步骤 12~13）
// 定义函数 HandleCommandResponse：封装一段可复用的处理逻辑（参数：deviceID uint, cmdType, payload string）

func HandleCommandResponse(deviceID uint, cmdType, payload string) {
	// 条件判断：pkg.DB == nil 时执行下面逻辑

	if pkg.DB == nil {
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 定义并初始化局部变量 rp，取值为 map[string]interface{}{} 并解析报文

	rp := map[string]interface{}{}
	// 条件判断：err := json.Unmarshal([]byte(payload), &rp); err != nil 时执行下面逻辑

	if err := json.Unmarshal([]byte(payload), &rp); err != nil {
		// 执行 pkg.Logf("[WARN] @ticket IoT response payload invalid: device_id=%d err=%v", deviceID, err)，完成对应的业务处理

		pkg.Logf("[WARN] @ticket IoT response payload invalid: device_id=%d err=%v", deviceID, err)
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 定义并初始化局部变量 commandID, _，取值为 rp["command_id"].(string)

	commandID, _ := rp["command_id"].(string)
	// 定义并初始化局部变量 result, _，取值为 rp["result"].(string)

	result, _ := rp["result"].(string)
	// 定义并初始化局部变量 message, _，取值为 rp["message"].(string)

	message, _ := rp["message"].(string)
	// 条件判断：commandID == "" 时执行下面逻辑

	if commandID == "" {
		// 执行 pkg.Logf("[WARN] @ticket IoT response missing command_id: device_id=%d", deviceID)，完成对应的业务处理

		pkg.Logf("[WARN] @ticket IoT response missing command_id: device_id=%d", deviceID)
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 取消超时监控（收到响应即停止计时）@ticket IoT-MQTT设备接入
	// 执行 cancelCommandTimeout(commandID)，完成对应的业务处理

	cancelCommandTimeout(commandID)
	// 声明变量 cur model.DeviceCommand，用于保存运行过程中的数据

	var cur model.DeviceCommand
	// 条件判断：err := pkg.DB.Where("command_id = ?", commandID).First(&cur).Error; err != nil 时执行下面逻辑

	if err := pkg.DB.Where("command_id = ?", commandID).First(&cur).Error; err != nil {
		// 未知指令响应留存事件即可（可能为重复响应或跨实例指令）
		// 执行 saveDeviceEvent(deviceID, "", "device/response/"+strconv.FormatUint(uint64(deviceID), 10)+"/"+cmdType, "response", "warn", payload)，完成对应的业务处理

		saveDeviceEvent(deviceID, "", "device/response/"+strconv.FormatUint(uint64(deviceID), 10)+"/"+cmdType, "response", "warn", payload)
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 更新指令终态 @ticket IoT-MQTT设备接入
	// 定义并初始化局部变量 nowT，取值为 time.Now()

	nowT := time.Now()
	// 定义并初始化局部变量 finalStatus，取值为 model.CmdSuccess

	finalStatus := model.CmdSuccess
	// 条件判断：result != "success" 时执行下面逻辑

	if result != "success" {
		// 为 finalStatus 赋值：model.CmdFailed

		finalStatus = model.CmdFailed
	// 当前代码块结束

	}
	// 执行 pkg.DB.Model(&model.DeviceCommand{}).Where("id = ?", cur.ID).Updates(map[string]interface{}{"status": finalStatus, "response": payload, "responded_at": &nowT})，完成对应的业务处理

	pkg.DB.Model(&model.DeviceCommand{}).Where("id = ?", cur.ID).Updates(map[string]interface{}{"status": finalStatus, "response": payload, "responded_at": &nowT})
	// 响应报文留存为事件
	// 执行 saveDeviceEvent(deviceID, cur.DeviceSN, "device/response/"+strconv.FormatUint(uint64(deviceID), 10)+"/"+cmdType, "response", "info", payload)，完成对应的业务处理

	saveDeviceEvent(deviceID, cur.DeviceSN, "device/response/"+strconv.FormatUint(uint64(deviceID), 10)+"/"+cmdType, "response", "info", payload)
	// 开门响应联动订单/柜门状态（方案 5.1 步骤 13：更新订单状态与柜格状态）@ticket IoT-MQTT设备接入
	// 条件判断：cmdType == "open_door" && result == "success" && cur.OrderNo != "" 时执行下面逻辑

	if cmdType == "open_door" && result == "success" && cur.OrderNo != "" {
		// 执行 onOpenDoorSuccess(cur.OrderNo)，完成对应的业务处理

		onOpenDoorSuccess(cur.OrderNo)
	// 当前代码块结束

	}
	// 失败响应联动场景通知
	// 条件判断：result != "success" 时执行下面逻辑

	if result != "success" {
		// 执行 sceneDeviceAlert(deviceID, "device/response/"+strconv.FormatUint(uint64(deviceID), 10)+"/"+cmdType, "指令执行失败: "+message)，完成对应的业务处理

		sceneDeviceAlert(deviceID, "device/response/"+strconv.FormatUint(uint64(deviceID), 10)+"/"+cmdType, "指令执行失败: "+message)
	// 当前代码块结束

	}
// 当前代码块结束

}

// onOpenDoorSuccess 开门成功后置处理：记录审计（订单状态更新仍由订单流程接口负责，避免双写冲突）
// 定义函数 onOpenDoorSuccess：封装一段可复用的处理逻辑（参数：orderNo string）

func onOpenDoorSuccess(orderNo string) {
	// 执行 pkg.Logf("[AUDIT] @ticket IoT open_door success: order_no=%s", orderNo)，完成对应的业务处理

	pkg.Logf("[AUDIT] @ticket IoT open_door success: order_no=%s", orderNo)
// 当前代码块结束

}

// DispatchConfig 下发价格配置（platform/config/{id}/price，方案 3.3 配置下发）
// 定义函数 DispatchConfig：封装一段可复用的处理逻辑（参数：sn string, priceConfig map[string]interface{}）(map[string]interface{}, int)

func DispatchConfig(sn string, priceConfig map[string]interface{}) (map[string]interface{}, int) {
	// 校验价格配置非空
	// 条件判断：priceConfig == nil || len(priceConfig) == 0 时执行下面逻辑

	if priceConfig == nil || len(priceConfig) == 0 {
		// 返回 nil, pkg.ErrBad 给调用方，结束当前函数

		return nil, pkg.ErrBad
	// 当前代码块结束

	}
	// 复用通用指令下发通道：cmdType=config，extra 携带 price_config
	// 定义并初始化局部变量 configID，取值为 fmt.Sprintf("CFG%s", time.Now().Format("20060102150405"))

	configID := fmt.Sprintf("CFG%s", time.Now().Format("20060102150405"))
	// 返回 DispatchCommand(sn, "config", "", "", "platform", map[string]interface{}{"config_id": configID, "price_config": priceConfig}) 给调用方，结束当前函数

	return DispatchCommand(sn, "config", "", "", "platform", map[string]interface{}{"config_id": configID, "price_config": priceConfig})
// 当前代码块结束

}

// CreateOta 创建并下发 OTA 升级任务（platform/ota/{id}/upgrade，方案 3.3 OTA）
// 定义函数 CreateOta：封装一段可复用的处理逻辑（参数：sn, version, url, md5 string, size int64, force bool）(map[string]interface{}, int)

func CreateOta(sn, version, url, md5 string, size int64, force bool) (map[string]interface{}, int) {
	// 参数校验：版本与下载地址必填
	// 条件判断：sn == "" || version == "" || url == "" 时执行下面逻辑

	if sn == "" || version == "" || url == "" {
		// 返回 nil, pkg.ErrBad 给调用方，结束当前函数

		return nil, pkg.ErrBad
	// 当前代码块结束

	}
	// 定义并初始化局部变量 d, err，取值为 deviceBySN(sn)

	d, err := deviceBySN(sn)
	// 条件判断：err != nil 时执行下面逻辑

	if err != nil {
		// 返回 nil, pkg.ErrDeviceNF 给调用方，结束当前函数

		return nil, pkg.ErrDeviceNF
	// 当前代码块结束

	}
	// 生成 OTA 任务号（方案 3.3 示例 OTA20240101001）
	// 定义并初始化局部变量 otaID，取值为 fmt.Sprintf("OTA%s%s", time.Now().Format("20060102150405"), strings.ToUpper(randomHex(3)))

	otaID := fmt.Sprintf("OTA%s%s", time.Now().Format("20060102150405"), strings.ToUpper(randomHex(3)))
	// 任务落库 @ticket IoT-MQTT设备接入
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
		// 执行 Status: model.OtaIssued,，完成对应的业务处理

		Status:          model.OtaIssued,
		// 执行 CreatedAt: time.Now(), UpdatedAt: time.Now(),，完成对应的业务处理

		CreatedAt:       time.Now(), UpdatedAt: time.Now(),
	// 当前代码块结束

	}
	// 条件判断：err := pkg.DB.Create(task).Error; err != nil 时执行下面逻辑

	if err := pkg.DB.Create(task).Error; err != nil {
		// 执行 pkg.Logf("[WARN] @ticket IoT ota save failed: %v", err)，完成对应的业务处理

		pkg.Logf("[WARN] @ticket IoT ota save failed: %v", err)
		// 返回 nil, pkg.ErrInternal 给调用方，结束当前函数

		return nil, pkg.ErrInternal
	// 当前代码块结束

	}
	// 组装 OTA 下行报文（方案 3.3 platform/ota/{id}/upgrade 报文）
	// 定义并初始化局部变量 payload，取值为 buildCommandPayload(otaID, d.ID, "", "", 300, map[string]interface{}{"ota_id": otaID, "firmware_version": version, "download_url": url, "md5": md5, "file_size": size, "force_upgrade": force})

	payload := buildCommandPayload(otaID, d.ID, "", "", 300, map[string]interface{}{"ota_id": otaID, "firmware_version": version, "download_url": url, "md5": md5, "file_size": size, "force_upgrade": force})
	// 发布 OTA 通知（QoS1；300s 超时不重试，等待设备主动拉取）
	// 定义并初始化局部变量 pubErr，取值为 pkg.MqttPublish(fmt.Sprintf("%sota/%d/upgrade", mqttTopicPlatformPrefix, d.ID), 1, []byte(payload))

	pubErr := pkg.MqttPublish(fmt.Sprintf("%sota/%d/upgrade", mqttTopicPlatformPrefix, d.ID), 1, []byte(payload))
	// 执行 pkg.Logf("[AUDIT] @ticket IoT ota dispatched: ota=%s sn=%s version=%s published=%v", otaID, sn, version, pubErr == nil)，完成对应的业务处理

	pkg.Logf("[AUDIT] @ticket IoT ota dispatched: ota=%s sn=%s version=%s published=%v", otaID, sn, version, pubErr == nil)
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
	// 当前代码块结束

	}, pkg.CodeOK
// 当前代码块结束

}

// otaProgressPayload OTA 进度上报报文（方案 3.3 device/ota/{id}/progress）
// 定义结构体 otaProgressPayload：用于承载一组相关的字段数据

type otaProgressPayload struct {
	// 结构体字段 OtaID，类型为 string，附带标签 `json:"ota_id"`

	OtaID     string `json:"ota_id"`
	// 结构体字段 Status，类型为 string，附带标签 `json:"status"`；downloading/upgrading/success/failed

	Status    string `json:"status"`
	// 结构体字段 Progress，类型为 int，附带标签 `json:"progress"`

	Progress  int    `json:"progress"`
	// 结构体字段 Message，类型为 string，附带标签 `json:"message"`

	Message   string `json:"message"`
// 当前代码块结束

}

// HandleOtaProgress 处理 OTA 进度上报：更新任务状态与进度，成功后同步设备固件版本
// 定义函数 HandleOtaProgress：封装一段可复用的处理逻辑（参数：deviceID uint, payload string）

func HandleOtaProgress(deviceID uint, payload string) {
	// 条件判断：pkg.DB == nil 时执行下面逻辑

	if pkg.DB == nil {
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 解析进度报文
	// 定义并初始化局部变量 op otaProgressPayload，取值为零值

	var op otaProgressPayload
	// 条件判断：err := json.Unmarshal([]byte(payload), &op); err != nil 时执行下面逻辑

	if err := json.Unmarshal([]byte(payload), &op); err != nil {
		// 执行 pkg.Logf("[WARN] @ticket IoT ota progress invalid: device_id=%d err=%v", deviceID, err)，完成对应的业务处理

		pkg.Logf("[WARN] @ticket IoT ota progress invalid: device_id=%d err=%v", deviceID, err)
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 条件判断：op.OtaID == "" 时执行下面逻辑

	if op.OtaID == "" {
		// 执行 pkg.Logf("[WARN] @ticket IoT ota progress missing ota_id: device_id=%d", deviceID)，完成对应的业务处理

		pkg.Logf("[WARN] @ticket IoT ota progress missing ota_id: device_id=%d", deviceID)
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 状态字符串映射到 Ota* 常量
	// 定义并初始化局部变量 st，取值为 mapOtaStatus(op.Status)

	st := mapOtaStatus(op.Status)
	// 更新任务状态/进度/消息 @ticket IoT-MQTT设备接入
	// 执行 pkg.DB.Model(&model.OtaTask{}).Where("ota_id = ?", op.OtaID).Updates(...)，完成对应的业务处理

	pkg.DB.Model(&model.OtaTask{}).Where("ota_id = ?", op.OtaID).Updates(map[string]interface{}{
		// 执行 "status": st,，完成对应的业务处理

		"status":   st,
		// 执行 "progress": op.Progress,，完成对应的业务处理

		"progress": op.Progress,
		// 执行 "message": op.Message,，完成对应的业务处理

		"message":  op.Message,
		// 执行 "updated_at": time.Now(),，完成对应的业务处理

		"updated_at": time.Now(),
	// 当前代码块结束

	})
	// 升级成功：同步设备固件版本（方案 6.5 测试场景：固件版本更新）
	// 条件判断：st == model.OtaDone 时执行下面逻辑

	if st == model.OtaDone {
		// 声明变量 t model.OtaTask，用于保存运行过程中的数据

		var t model.OtaTask
		// 条件判断：err := pkg.DB.Where("ota_id = ?", op.OtaID).First(&t).Error; err == nil 时执行下面逻辑

		if err := pkg.DB.Where("ota_id = ?", op.OtaID).First(&t).Error; err == nil {
			// 执行 pkg.DB.Model(&model.IoTDevice{}).Where("id = ?", t.DeviceID).Updates(map[string]interface{}{"firmware_version": t.FirmwareVersion, "updated_at": time.Now()})，完成对应的业务处理

			pkg.DB.Model(&model.IoTDevice{}).Where("id = ?", t.DeviceID).Updates(map[string]interface{}{"firmware_version": t.FirmwareVersion, "updated_at": time.Now()})
		// 当前代码块结束

		}
	// 当前代码块结束

	}
	// 进度留存为事件
	// 执行 saveDeviceEvent(deviceID, "", "device/ota/"+strconv.FormatUint(uint64(deviceID), 10)+"/progress", "ota", "info", payload)，完成对应的业务处理

	saveDeviceEvent(deviceID, "", "device/ota/"+strconv.FormatUint(uint64(deviceID), 10)+"/progress", "ota", "info", payload)
// 当前代码块结束

}

// mapOtaStatus OTA 进度状态字符串 → 常量映射
// 定义函数 mapOtaStatus：封装一段可复用的处理逻辑（参数：s string）int

func mapOtaStatus(s string) int {
	// 多分支判断：根据 s 的值选择分支

	switch s {
	// 当取值匹配 downloading 时进入该分支

	case "downloading":
		// 返回 model.OtaDownloading 给调用方，结束当前函数

		return model.OtaDownloading
	// 当取值匹配 upgrading 时进入该分支

	case "upgrading":
		// 返回 model.OtaUpgrading 给调用方，结束当前函数

		return model.OtaUpgrading
	// 当取值匹配 success 时进入该分支

	case "success":
		// 返回 model.OtaDone 给调用方，结束当前函数

		return model.OtaDone
	// 当取值匹配 failed 时进入该分支

	case "failed":
		// 返回 model.OtaFailed 给调用方，结束当前函数

		return model.OtaFailed
// 当前代码块结束

	}
	// 兜底：保持已下发状态
	// 返回 model.OtaIssued 给调用方，结束当前函数

	return model.OtaIssued
// 当前代码块结束

}

// lookupAuthOrderNo 按临时开柜 Token 反查订单号（供订单流程联动 MQTT 开门）
// 定义函数 lookupAuthOrderNo：封装一段可复用的处理逻辑（参数：token string）(string, error)

func lookupAuthOrderNo(token string) (string, error) {
	// 条件判断：pkg.DB == nil 时执行下面逻辑

	if pkg.DB == nil {
		// 返回 "", fmt.Errorf("database unavailable") 给调用方，结束当前函数

		return "", fmt.Errorf("database unavailable")
	// 当前代码块结束

	}
	// 声明变量 a model.TempOpenAuth，用于保存运行过程中的数据

	var a model.TempOpenAuth
	// 条件判断：err := pkg.DB.Where("token = ?", token).First(&a).Error; err != nil 时执行下面逻辑

	if err := pkg.DB.Where("token = ?", token).First(&a).Error; err != nil {
		// 返回 "", err 给调用方，结束当前函数

		return "", err
	// 当前代码块结束

	}
	// 返回 a.OrderNo, nil 给调用方，结束当前函数

	return a.OrderNo, nil
// 当前代码块结束

}

// DispatchOpenDoorForOrder 订单流程联动：按订单定位设备并下发开门指令（供工单05 取件/临时开柜调用）
// 定义函数 DispatchOpenDoorForOrder：封装一段可复用的处理逻辑（参数：orderNo, operator string）error

func DispatchOpenDoorForOrder(orderNo, operator string) error {
	// 条件判断：pkg.DB == nil 时执行下面逻辑

	if pkg.DB == nil {
		// 返回 fmt.Errorf("database unavailable") 给调用方，结束当前函数

		return fmt.Errorf("database unavailable")
	// 当前代码块结束

	}
	// 条件判断：!pkg.MqttReady() 时执行下面逻辑（Broker 未启用时静默跳过，不影响既有订单流程）

	if !pkg.MqttReady() {
		// 返回 nil 给调用方，结束当前函数

		return nil
	// 当前代码块结束

	}
	// 声明变量 o model.DepositOrder，用于保存运行过程中的数据

	var o model.DepositOrder
	// 条件判断：err := pkg.DB.Where("order_no = ?", orderNo).First(&o).Error; err != nil 时执行下面逻辑

	if err := pkg.DB.Where("order_no = ?", orderNo).First(&o).Error; err != nil {
		// 返回 err 给调用方，结束当前函数

		return err
	// 当前代码块结束

	}
	// 声明变量 door model.LockerDoor，用于保存运行过程中的数据

	var door model.LockerDoor
	// 条件判断：err := pkg.DB.First(&door, o.DoorID).Error; err != nil 时执行下面逻辑

	if err := pkg.DB.First(&door, o.DoorID).Error; err != nil {
		// 返回 err 给调用方，结束当前函数

		return err
	// 当前代码块结束

	}
	// 按柜体定位在线设备（设备绑定 cabinet_id）@ticket IoT-MQTT设备接入
	// 声明变量 d model.IoTDevice，用于保存运行过程中的数据

	var d model.IoTDevice
	// 条件判断：err := pkg.DB.Where("cabinet_id = ? AND status = ?", door.CabinetID, model.IoTDevOnline).First(&d).Error; err != nil 时执行下面逻辑

	if err := pkg.DB.Where("cabinet_id = ? AND status = ?", door.CabinetID, model.IoTDevOnline).First(&d).Error; err != nil {
		// 设备未绑定或离线：记录告警日志但不阻断既有流程
		// 执行 pkg.Logf("[WARN] @ticket IoT open door skip (device offline/unbound): order=%s door=%s", orderNo, door.Code)，完成对应的业务处理

		pkg.Logf("[WARN] @ticket IoT open door skip (device offline/unbound): order=%s door=%s", orderNo, door.Code)
		// 返回 nil 给调用方，结束当前函数

		return nil
	// 当前代码块结束

	}
	// 下发开门指令（失败仅记录日志，不阻断既有订单流程）
	// 定义并初始化局部变量 _, code，取值为 DispatchCommand(d.DeviceSN, "open_door", door.Code, orderNo, operator, nil)

	_, code := DispatchCommand(d.DeviceSN, "open_door", door.Code, orderNo, operator, nil)
	// 条件判断：code != pkg.CodeOK 时执行下面逻辑

	if code != pkg.CodeOK {
		// 执行 pkg.Logf("[WARN] @ticket IoT open door dispatch failed: order=%s code=%d", orderNo, code)，完成对应的业务处理

		pkg.Logf("[WARN] @ticket IoT open door dispatch failed: order=%s code=%d", orderNo, code)
		// 返回 fmt.Errorf("dispatch failed") 给调用方，结束当前函数

		return fmt.Errorf("dispatch failed")
	// 当前代码块结束

	}
	// 返回 nil 给调用方，结束当前函数

	return nil
// 当前代码块结束

}

// RouteMessage device/# 上行消息统一入口（paho 默认处理器调用）：按 Topic 分段路由 @ticket IoT-MQTT设备接入
// 定义函数 RouteMessage：封装一段可复用的处理逻辑（参数：topic string, payload []byte）

func RouteMessage(topic string, payload []byte) {
	// 条件判断：!strings.HasPrefix(topic, mqttTopicDevicePrefix) 时执行下面逻辑

	if !strings.HasPrefix(topic, mqttTopicDevicePrefix) {
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 按 / 切分 Topic：device/{kind}/{device_id}/{action}（方案 3.1 命名规范）
	// 定义并初始化局部变量 parts，取值为 strings.Split(topic, "/")

	parts := strings.Split(topic, "/")
	// 条件判断：len(parts) < 4 时执行下面逻辑

	if len(parts) < 4 {
		// 执行 log.Printf("[WARN] @ticket IoT topic malformed: %s", topic)，完成对应的业务处理

		log.Printf("[WARN] @ticket IoT topic malformed: %s", topic)
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 提取消息类型与设备 ID
	// 定义并初始化局部变量 kind，取值为 parts[1]

	kind := parts[1]
	// 定义并初始化局部变量 idStr，取值为 parts[2]

	idStr := parts[2]
	// 定义并初始化局部变量 action，取值为 parts[3]

	action := parts[3]
	// 解析设备 ID
	// 定义并初始化局部变量 deviceID, err，取值为 strconv.ParseUint(idStr, 10, 64)

	deviceID, err := strconv.ParseUint(idStr, 10, 64)
	// 条件判断：err != nil 时执行下面逻辑

	if err != nil {
		// 执行 log.Printf("[WARN] @ticket IoT device id invalid: %s topic=%s", idStr, topic)，完成对应的业务处理

		log.Printf("[WARN] @ticket IoT device id invalid: %s topic=%s", idStr, topic)
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 定义并初始化局部变量 body，取值为 string(payload)

	body := string(payload)
	// 多分支判断：根据 kind 的值选择分支

	switch kind {
	// 当取值匹配 status 时进入该分支
	case "status":
		// 多分支判断：根据 action 的值选择分支

		switch action {
		// 当取值匹配 online 时进入该分支
		case "online":
			// 执行 HandleConnected(uint(deviceID), extractJSONString(body, "ip_address"))，完成对应的业务处理

			HandleConnected(uint(deviceID), extractJSONString(body, "ip_address"))
		// 当取值匹配 offline 时进入该分支
		case "offline":
			// 执行 HandleDisconnected(uint(deviceID))，完成对应的业务处理

			HandleDisconnected(uint(deviceID))
		// 当取值匹配 heartbeat 时进入该分支
		case "heartbeat":
			// 执行 HandleHeartbeat(uint(deviceID), body)，完成对应的业务处理

			HandleHeartbeat(uint(deviceID), body)
		// 当取值匹配 locker 时进入该分支
		case "locker":
			// 执行 HandleLockerReport(uint(deviceID), body)，完成对应的业务处理

			HandleLockerReport(uint(deviceID), body)
		// 当前代码块结束

		}
	// 当取值匹配 log 时进入该分支
	case "log":
		// 执行 HandleDeviceLog(uint(deviceID), topic, body, action == "error")，完成对应的业务处理

		HandleDeviceLog(uint(deviceID), topic, body, action == "error")
	// 当取值匹配 alert 时进入该分支
	case "alert":
		// 告警级别取 payload.alert_level（方案 3.2 fault 示例），缺省按 warn
		// 执行 HandleDeviceAlert(uint(deviceID), topic, body, extractJSONString(body, "alert_level"))，完成对应的业务处理

		lvl := extractJSONString(body, "alert_level")
		// 条件判断：lvl == "" 时执行下面逻辑

		if lvl == "" {
			// 为 lvl 赋值："warn"

			lvl = "warn"
		// 当前代码块结束

		}
		// 执行 HandleDeviceAlert(uint(deviceID), topic, body, lvl)，完成对应的业务处理

		HandleDeviceAlert(uint(deviceID), topic, body, lvl)
	// 当取值匹配 response 时进入该分支
	case "response":
		// 执行 HandleCommandResponse(uint(deviceID), action, body)，完成对应的业务处理

		HandleCommandResponse(uint(deviceID), action, body)
	// 当取值匹配 ota 时进入该分支
	case "ota":
		// 条件判断：action == "progress" 时执行下面逻辑

		if action == "progress" {
			// 执行 HandleOtaProgress(uint(deviceID), body)，完成对应的业务处理

			HandleOtaProgress(uint(deviceID), body)
		// 当前代码块结束

		}
	// 当前代码块结束

	}
// 当前代码块结束

}

// extractJSONString 从 JSON 报文中提取顶层字符串字段（轻量解析，避免引入 schema 依赖）
// 定义函数 extractJSONString：封装一段可复用的处理逻辑（参数：body, key string）string

func extractJSONString(body, key string) string {
	// 定义并初始化局部变量 m，取值为 map[string]interface{}{}

	m := map[string]interface{}{}
	// 条件判断：err := json.Unmarshal([]byte(body), &m); err != nil 时执行下面逻辑

	if err := json.Unmarshal([]byte(body), &m); err != nil {
		// 返回 "" 给调用方，结束当前函数

		return ""
	// 当前代码块结束

	}
	// 定义并初始化局部变量 v, ok，取值为 m[key]

	v, ok := m[key]
	// 条件判断：!ok 时执行下面逻辑

	if !ok {
		// 返回 "" 给调用方，结束当前函数

		return ""
	// 当前代码块结束

	}
	// 定义并初始化局部变量 s, ok，取值为 v.(string)

	s, ok := v.(string)
	// 条件判断：!ok 时执行下面逻辑

	if !ok {
		// 返回 "" 给调用方，结束当前函数

		return ""
	// 当前代码块结束

	}
	// 返回 s 给调用方，结束当前函数

	return s
// 当前代码块结束

}

// DeviceStatusList 设备状态列表（管理端查询；IP/密钥脱敏由 model 标签控制）@ticket IoT-MQTT设备接入
// 定义函数 DeviceStatusList：封装一段可复用的处理逻辑（参数：siteID uint）([]model.IoTDevice, int)

func DeviceStatusList(siteID uint) ([]model.IoTDevice, int) {
	// 条件判断：pkg.DB == nil 时执行下面逻辑

	if pkg.DB == nil {
		// 返回 nil, pkg.ErrInternal 给调用方，结束当前函数

		return nil, pkg.ErrInternal
	// 当前代码块结束

	}
	// 定义并初始化局部变量 out，取值为 []model.IoTDevice{}

	out := []model.IoTDevice{}
	// 定义并初始化局部变量 q，取值为 pkg.DB.Order("id")

	q := pkg.DB.Order("id")
	// 条件判断：siteID > 0 时执行下面逻辑

	if siteID > 0 {
		// 执行 q = q.Where("site_id = ?", siteID)，完成对应的业务处理

		q = q.Where("site_id = ?", siteID)
	// 当前代码块结束

	}
	// 条件判断：err := q.Find(&out).Error; err != nil 时执行下面逻辑

	if err := q.Find(&out).Error; err != nil {
		// 返回 nil, pkg.ErrInternal 给调用方，结束当前函数

		return nil, pkg.ErrInternal
	// 当前代码块结束

	}
	// 返回 out, pkg.CodeOK 给调用方，结束当前函数

	return out, pkg.CodeOK
// 当前代码块结束

}

// CommandRecords 指令记录分页查询
// 定义函数 CommandRecords：封装一段可复用的处理逻辑（参数：deviceID uint, page, pageSize int）(map[string]interface{}, int)

func CommandRecords(deviceID uint, page, pageSize int) (map[string]interface{}, int) {
	// 条件判断：pkg.DB == nil 时执行下面逻辑

	if pkg.DB == nil {
		// 返回 nil, pkg.ErrInternal 给调用方，结束当前函数

		return nil, pkg.ErrInternal
	// 当前代码块结束

	}
	// 分页参数兜底（对齐既有接口分页口径）
	// 条件判断：page <= 0 时执行下面逻辑

	if page <= 0 {
		// 为 page 赋值：1

		page = 1
	// 当前代码块结束

	}
	// 条件判断：pageSize <= 0 || pageSize > 100 时执行下面逻辑

	if pageSize <= 0 || pageSize > 100 {
		// 为 pageSize 赋值：20

		pageSize = 20
	// 当前代码块结束

	}
	// 定义并初始化局部变量 q，取值为 pkg.DB.Model(&model.DeviceCommand{})

	q := pkg.DB.Model(&model.DeviceCommand{})
	// 条件判断：deviceID > 0 时执行下面逻辑

	if deviceID > 0 {
		// 执行 q = q.Where("device_id = ?", deviceID)，完成对应的业务处理

		q = q.Where("device_id = ?", deviceID)
	// 当前代码块结束

	}
	// 声明变量 total int64，用于保存运行过程中的数据

	var total int64
	// 执行 q.Count(&total)，完成对应的业务处理

	q.Count(&total)
	// 定义并初始化局部变量 rows，取值为 []model.DeviceCommand{}

	rows := []model.DeviceCommand{}
	// 条件判断：err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&rows).Error; err != nil 时执行下面逻辑

	if err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&rows).Error; err != nil {
		// 返回 nil, pkg.ErrInternal 给调用方，结束当前函数

		return nil, pkg.ErrInternal
	// 当前代码块结束

	}
	// 返回分页结果给调用方，结束当前函数

	return map[string]interface{}{
		// 执行 "total": total,，完成对应的业务处理

		"total": total,
		// 执行 "page": page,，完成对应的业务处理

		"page":  page,
		// 执行 "page_size": pageSize,，完成对应的业务处理

		"page_size": pageSize,
		// 执行 "records": rows,，完成对应的业务处理

		"records": rows,
	// 当前代码块结束

	}, pkg.CodeOK
// 当前代码块结束

}

// EventRecords 设备上行事件分页查询
// 定义函数 EventRecords：封装一段可复用的处理逻辑（参数：deviceID uint, kind string, page, pageSize int）(map[string]interface{}, int)

func EventRecords(deviceID uint, kind string, page, pageSize int) (map[string]interface{}, int) {
	// 条件判断：pkg.DB == nil 时执行下面逻辑

	if pkg.DB == nil {
		// 返回 nil, pkg.ErrInternal 给调用方，结束当前函数

		return nil, pkg.ErrInternal
	// 当前代码块结束

	}
	// 分页参数兜底
	// 条件判断：page <= 0 时执行下面逻辑

	if page <= 0 {
		// 为 page 赋值：1

		page = 1
	// 当前代码块结束

	}
	// 条件判断：pageSize <= 0 || pageSize > 100 时执行下面逻辑

	if pageSize <= 0 || pageSize > 100 {
		// 为 pageSize 赋值：20

		pageSize = 20
	// 当前代码块结束

	}
	// 定义并初始化局部变量 q，取值为 pkg.DB.Model(&model.DeviceEvent{})

	q := pkg.DB.Model(&model.DeviceEvent{})
	// 条件判断：deviceID > 0 时执行下面逻辑

	if deviceID > 0 {
		// 执行 q = q.Where("device_id = ?", deviceID)，完成对应的业务处理

		q = q.Where("device_id = ?", deviceID)
	// 当前代码块结束

	}
	// 条件判断：kind != "" 时执行下面逻辑

	if kind != "" {
		// 执行 q = q.Where("kind = ?", kind)，完成对应的业务处理

		q = q.Where("kind = ?", kind)
	// 当前代码块结束

	}
	// 声明变量 total int64，用于保存运行过程中的数据

	var total int64
	// 执行 q.Count(&total)，完成对应的业务处理

	q.Count(&total)
	// 定义并初始化局部变量 rows，取值为 []model.DeviceEvent{}

	rows := []model.DeviceEvent{}
	// 条件判断：err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&rows).Error; err != nil 时执行下面逻辑

	if err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&rows).Error; err != nil {
		// 返回 nil, pkg.ErrInternal 给调用方，结束当前函数

		return nil, pkg.ErrInternal
	// 当前代码块结束

	}
	// 返回分页结果给调用方，结束当前函数

	return map[string]interface{}{
		// 执行 "total": total,，完成对应的业务处理

		"total": total,
		// 执行 "page": page,，完成对应的业务处理

		"page":  page,
		// 执行 "page_size": pageSize,，完成对应的业务处理

		"page_size": pageSize,
		// 执行 "records": rows,，完成对应的业务处理

		"records": rows,
	// 当前代码块结束

	}, pkg.CodeOK
// 当前代码块结束

}

// OtaRecords OTA 任务列表查询
// 定义函数 OtaRecords：封装一段可复用的处理逻辑（参数：deviceID uint, page, pageSize int）(map[string]interface{}, int)

func OtaRecords(deviceID uint, page, pageSize int) (map[string]interface{}, int) {
	// 条件判断：pkg.DB == nil 时执行下面逻辑

	if pkg.DB == nil {
		// 返回 nil, pkg.ErrInternal 给调用方，结束当前函数

		return nil, pkg.ErrInternal
	// 当前代码块结束

	}
	// 分页参数兜底
	// 条件判断：page <= 0 时执行下面逻辑

	if page <= 0 {
		// 为 page 赋值：1

		page = 1
	// 当前代码块结束

	}
	// 条件判断：pageSize <= 0 || pageSize > 100 时执行下面逻辑

	if pageSize <= 0 || pageSize > 100 {
		// 为 pageSize 赋值：20

		pageSize = 20
	// 当前代码块结束

	}
	// 定义并初始化局部变量 q，取值为 pkg.DB.Model(&model.OtaTask{})

	q := pkg.DB.Model(&model.OtaTask{})
	// 条件判断：deviceID > 0 时执行下面逻辑

	if deviceID > 0 {
		// 执行 q = q.Where("device_id = ?", deviceID)，完成对应的业务处理

		q = q.Where("device_id = ?", deviceID)
	// 当前代码块结束

	}
	// 声明变量 total int64，用于保存运行过程中的数据

	var total int64
	// 执行 q.Count(&total)，完成对应的业务处理

	q.Count(&total)
	// 定义并初始化局部变量 rows，取值为 []model.OtaTask{}

	rows := []model.OtaTask{}
	// 条件判断：err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&rows).Error; err != nil 时执行下面逻辑

	if err := q.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&rows).Error; err != nil {
		// 返回 nil, pkg.ErrInternal 给调用方，结束当前函数

		return nil, pkg.ErrInternal
	// 当前代码块结束

	}
	// 返回分页结果给调用方，结束当前函数

	return map[string]interface{}{
		// 执行 "total": total,，完成对应的业务处理

		"total": total,
		// 执行 "page": page,，完成对应的业务处理

		"page":  page,
		// 执行 "page_size": pageSize,，完成对应的业务处理

		"page_size": pageSize,
		// 执行 "records": rows,，完成对应的业务处理

		"records": rows,
	// 当前代码块结束

	}, pkg.CodeOK
// 当前代码块结束

}
