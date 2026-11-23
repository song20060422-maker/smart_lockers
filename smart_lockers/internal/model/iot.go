// @ticket IoT-MQTT设备接入：设备/令牌/指令/事件/OTA 实体定义
// 字段与《智能快递柜IoT设备MQTT接入技术方案》v1.0 对齐：设备生命周期、24h Token、
// command_id 指令幂等、上行事件（状态/日志/告警）、OTA 升级任务。
package model

import "time"

// IoT 设备生命周期状态常量 @ticket IoT-MQTT设备接入（对齐方案 2.1 设备生命周期）
// 开始定义一组常量的 const 声明块

const (
	// 为 IoTDevFactory 赋值：0 // 出厂

	IoTDevFactory = 0 // 出厂
	// 为 IoTDevRegistered 赋值：1 // 已注册（录入系统）

	IoTDevRegistered = 1 // 已注册（录入系统）
	// 为 IoTDevPending 赋值：2 // 待激活（分配到站点）

	IoTDevPending = 2 // 待激活（分配到站点）
	// 为 IoTDevOnline 赋值：3 // 在线（MQTT 连接成功）

	IoTDevOnline = 3 // 在线（MQTT 连接成功）
	// 为 IoTDevOffline 赋值：4 // 离线（MQTT 断开）

	IoTDevOffline = 4 // 离线（MQTT 断开）
	// 为 IoTDevFault 赋值：5 // 异常（连续上报错误）

	IoTDevFault = 5 // 异常（连续上报错误）
	// 为 IoTDevMaintain 赋值：6 // 维护中（人工维护）

	IoTDevMaintain = 6 // 维护中（人工维护）
	// 为 IoTDevDisabled 赋值：7 // 已停用（管理员手动）

	IoTDevDisabled = 7 // 已停用（管理员手动）
	// 为 IoTDevScrapped 赋值：8 // 已报废（设备淘汰）

	IoTDevScrapped = 8 // 已报废（设备淘汰）
// 结束参数/表达式列表

)

// 指令状态常量 @ticket IoT-MQTT设备接入（对齐方案 5.2 指令超时处理）
// 开始定义一组常量的 const 声明块

const (
	// 为 CmdPending 赋值：0 // 待下发

	CmdPending = 0 // 待下发
	// 为 CmdSent 赋值：1 // 已下发待响应

	CmdSent = 1 // 已下发待响应
	// 为 CmdSuccess 赋值：2 // 设备响应成功

	CmdSuccess = 2 // 设备响应成功
	// 为 CmdFailed 赋值：3 // 设备响应失败

	CmdFailed = 3 // 设备响应失败
	// 为 CmdTimeout 赋值：4 // 指令超时（含重试耗尽）

	CmdTimeout = 4 // 指令超时（含重试耗尽）
// 结束参数/表达式列表

)

// OTA 任务状态常量 @ticket IoT-MQTT设备接入（对齐方案 3.3 OTA 升级）
// 开始定义一组常量的 const 声明块

const (
	// 为 OtaIssued 赋值：0 // 升级通知已下发

	OtaIssued = 0 // 升级通知已下发
	// 为 OtaDownloading 赋值：1 // 设备下载中

	OtaDownloading = 1 // 设备下载中
	// 为 OtaUpgrading 赋值：2 // 设备升级中

	OtaUpgrading = 2 // 设备升级中
	// 为 OtaDone 赋值：3 // 升级成功

	OtaDone = 3 // 升级成功
	// 为 OtaFailed 赋值：4 // 升级失败

	OtaFailed = 4 // 升级失败
// 结束参数/表达式列表

)

// IoTDevice 智能快递柜设备主数据 @ticket IoT-MQTT设备接入
// 定义结构体 IoTDevice：用于承载一组相关的字段数据

type IoTDevice struct {
	// 结构体字段 ID，类型为 uint，附带标签 `gorm:"primaryKey" json:"id"`；MQTT Topic 中使用该 ID

	ID               uint       `gorm:"primaryKey" json:"id"`
	// 结构体字段 DeviceSN，类型为 string，附带标签 `gorm:"size:64;uniqueIndex" json:"device_sn"`；设备序列号（如 SL20240101001）

	DeviceSN         string     `gorm:"size:64;uniqueIndex" json:"device_sn"`
	// 结构体字段 DeviceSecret，类型为 string，附带标签 `gorm:"size:128" json:"-"`；出厂密钥，脱敏永不输出

	DeviceSecret     string     `gorm:"size:128" json:"-"`
	// 结构体字段 Model，类型为 string，附带标签 `gorm:"size:64" json:"model"`；设备型号

	Model            string     `gorm:"size:64" json:"model"`
	// 结构体字段 HardwareVersion，类型为 string，附带标签 `gorm:"size:32" json:"hardware_version"`；硬件版本

	HardwareVersion  string     `gorm:"size:32" json:"hardware_version"`
	// 结构体字段 FirmwareVersion，类型为 string，附带标签 `gorm:"size:32" json:"firmware_version"`；固件版本（OTA 后更新）

	FirmwareVersion  string     `gorm:"size:32" json:"firmware_version"`
	// 结构体字段 FactoryDate，类型为 time.Time，附带标签 `json:"factory_date"`；出厂日期

	FactoryDate      time.Time  `json:"factory_date"`
	// 结构体字段 SiteID，类型为 uint，附带标签 `gorm:"index" json:"site_id"`；分配站点

	SiteID           uint       `gorm:"index" json:"site_id"`
	// 结构体字段 CabinetID，类型为 uint，附带标签 `json:"cabinet_id"`；绑定柜体（用于柜格状态上报落库）

	CabinetID        uint       `json:"cabinet_id"`
	// 结构体字段 Status，类型为 int，附带标签 `gorm:"default:0" json:"status"`；生命周期状态（IoTDev* 常量）

	Status           int        `gorm:"default:0" json:"status"`
	// 结构体字段 IPAddress，类型为 string，附带标签 `gorm:"size:64" json:"ip_address"`；最近连接 IP

	IPAddress        string     `gorm:"size:64" json:"ip_address"`
	// 结构体字段 LastOnlineAt，类型为 *time.Time，附带标签 `json:"last_online_at"`；最近上线时间

	LastOnlineAt     *time.Time `json:"last_online_at"`
	// 结构体字段 LastHeartbeat，类型为 *time.Time，附带标签 `json:"last_heartbeat"`；最近心跳时间

	LastHeartbeat    *time.Time `json:"last_heartbeat"`
	// 结构体字段 CreatedAt，类型为 time.Time，附带标签 `json:"created_at"`

	CreatedAt        time.Time  `json:"created_at"`
	// 结构体字段 UpdatedAt，类型为 time.Time，附带标签 `json:"updated_at"`

	UpdatedAt        time.Time  `json:"updated_at"`
// 当前代码块结束

}

// DeviceToken 设备 MQTT 连接令牌（24h 有效）@ticket IoT-MQTT设备接入
// 定义结构体 DeviceToken：用于承载一组相关的字段数据

type DeviceToken struct {
	// 结构体字段 ID，类型为 uint，附带标签 `gorm:"primaryKey" json:"id"`

	ID        uint      `gorm:"primaryKey" json:"id"`
	// 结构体字段 DeviceID，类型为 uint，附带标签 `gorm:"index" json:"device_id"`

	DeviceID  uint      `gorm:"index" json:"device_id"`
	// 结构体字段 Token，类型为 string，附带标签 `gorm:"size:128;uniqueIndex" json:"token"`；MQTT password 字段使用的 Token

	Token     string    `gorm:"size:128;uniqueIndex" json:"token"`
	// 结构体字段 ExpireAt，类型为 time.Time，附带标签 `json:"expire_at"`；过期时间（默认 24h）

	ExpireAt  time.Time `json:"expire_at"`
	// 结构体字段 Revoked，类型为 bool，附带标签 `gorm:"default:false" json:"revoked"`；是否已吊销

	Revoked   bool      `gorm:"default:false" json:"revoked"`
	// 结构体字段 CreatedAt，类型为 time.Time，附带标签 `json:"created_at"`

	CreatedAt time.Time `json:"created_at"`
// 当前代码块结束

}

// DeviceCommand 平台下发指令记录 @ticket IoT-MQTT设备接入（command_id 幂等 + 超时重试）
// 定义结构体 DeviceCommand：用于承载一组相关的字段数据

type DeviceCommand struct {
	// 结构体字段 ID，类型为 uint，附带标签 `gorm:"primaryKey" json:"id"`

	ID          uint       `gorm:"primaryKey" json:"id"`
	// 结构体字段 CommandID，类型为 string，附带标签 `gorm:"size:64;uniqueIndex" json:"command_id"`；全局唯一指令号（设备端幂等去重依据）

	CommandID   string     `gorm:"size:64;uniqueIndex" json:"command_id"`
	// 结构体字段 DeviceID，类型为 uint，附带标签 `gorm:"index" json:"device_id"`

	DeviceID    uint       `gorm:"index" json:"device_id"`
	// 结构体字段 DeviceSN，类型为 string，附带标签 `gorm:"size:64" json:"device_sn"`

	DeviceSN    string     `gorm:"size:64" json:"device_sn"`
	// 结构体字段 CmdType，类型为 string，附带标签 `gorm:"size:32" json:"cmd_type"`；open_door/close_door/reboot/config

	CmdType     string     `gorm:"size:32" json:"cmd_type"`
	// 结构体字段 LockerNo，类型为 string，附带标签 `gorm:"size:32" json:"locker_no"`；柜格编号（开门类指令）

	LockerNo    string     `gorm:"size:32" json:"locker_no"`
	// 结构体字段 OrderNo，类型为 string，附带标签 `gorm:"size:64" json:"order_no"`；关联订单号

	OrderNo     string     `gorm:"size:64" json:"order_no"`
	// 结构体字段 Payload，类型为 string，附带标签 `gorm:"type:text" json:"payload"`；下发 JSON 报文

	Payload     string     `gorm:"type:text" json:"payload"`
	// 结构体字段 TimeoutSec，类型为 int，附带标签 `json:"timeout_sec"`；超时时间（开门/关门 10s、重启 30s、配置 60s）

	TimeoutSec  int        `json:"timeout_sec"`
	// 结构体字段 MaxRetry，类型为 int，附带标签 `json:"max_retry"`；最大重试次数（开门/关门 1、配置 2、重启/OTA 0）

	MaxRetry    int        `json:"max_retry"`
	// 结构体字段 RetryCount，类型为 int，附带标签 `json:"retry_count"`；已重试次数

	RetryCount  int        `json:"retry_count"`
	// 结构体字段 Status，类型为 int，附带标签 `json:"status"`；Cmd* 常量

	Status      int        `json:"status"`
	// 结构体字段 Response，类型为 string，附带标签 `gorm:"type:text" json:"response"`；设备响应报文

	Response    string     `gorm:"type:text" json:"response"`
	// 结构体字段 Operator，类型为 string，附带标签 `gorm:"size:64" json:"operator"`；触发人（用户/运维）

	Operator    string     `gorm:"size:64" json:"operator"`
	// 结构体字段 CreatedAt，类型为 time.Time，附带标签 `json:"created_at"`

	CreatedAt   time.Time  `json:"created_at"`
	// 结构体字段 RespondedAt，类型为 *time.Time，附带标签 `json:"responded_at"`；响应时间

	RespondedAt *time.Time `json:"responded_at"`
// 当前代码块结束

}

// DeviceEvent 设备上行事件（状态/日志/告警统一留存）@ticket IoT-MQTT设备接入
// 定义结构体 DeviceEvent：用于承载一组相关的字段数据

type DeviceEvent struct {
	// 结构体字段 ID，类型为 uint，附带标签 `gorm:"primaryKey" json:"id"`

	ID        uint      `gorm:"primaryKey" json:"id"`
	// 结构体字段 DeviceID，类型为 uint，附带标签 `gorm:"index" json:"device_id"`

	DeviceID  uint      `gorm:"index" json:"device_id"`
	// 结构体字段 DeviceSN，类型为 string，附带标签 `gorm:"size:64" json:"device_sn"`

	DeviceSN  string    `gorm:"size:64" json:"device_sn"`
	// 结构体字段 Topic，类型为 string，附带标签 `gorm:"size:128" json:"topic"`；上行 Topic 全文

	Topic     string    `gorm:"size:128" json:"topic"`
	// 结构体字段 Kind，类型为 string，附带标签 `gorm:"size:16;index" json:"kind"`；status/log/alert/response/ota

	Kind      string    `gorm:"size:16;index" json:"kind"`
	// 结构体字段 Level，类型为 string，附带标签 `gorm:"size:16" json:"level"`；info/warn/high/critical（告警时按方案分级）

	Level     string    `gorm:"size:16" json:"level"`
	// 结构体字段 Content，类型为 string，附带标签 `gorm:"type:text" json:"content"`；原始 payload JSON

	Content   string    `gorm:"type:text" json:"content"`
	// 结构体字段 CreatedAt，类型为 time.Time，附带标签 `gorm:"index" json:"created_at"`

	CreatedAt time.Time `gorm:"index" json:"created_at"`
// 当前代码块结束

}

// OtaTask OTA 固件升级任务 @ticket IoT-MQTT设备接入
// 定义结构体 OtaTask：用于承载一组相关的字段数据

type OtaTask struct {
	// 结构体字段 ID，类型为 uint，附带标签 `gorm:"primaryKey" json:"id"`

	ID              uint      `gorm:"primaryKey" json:"id"`
	// 结构体字段 OtaID，类型为 string，附带标签 `gorm:"size:64;uniqueIndex" json:"ota_id"`；任务号（如 OTA20240101001）

	OtaID           string    `gorm:"size:64;uniqueIndex" json:"ota_id"`
	// 结构体字段 DeviceID，类型为 uint，附带标签 `gorm:"index" json:"device_id"`

	DeviceID        uint      `gorm:"index" json:"device_id"`
	// 结构体字段 DeviceSN，类型为 string，附带标签 `gorm:"size:64" json:"device_sn"`

	DeviceSN        string    `gorm:"size:64" json:"device_sn"`
	// 结构体字段 FirmwareVersion，类型为 string，附带标签 `gorm:"size:32" json:"firmware_version"`；目标固件版本

	FirmwareVersion string    `gorm:"size:32" json:"firmware_version"`
	// 结构体字段 DownloadURL，类型为 string，附带标签 `gorm:"size:512" json:"download_url"`；固件下载地址

	DownloadURL     string    `gorm:"size:512" json:"download_url"`
	// 结构体字段 MD5，类型为 string，附带标签 `gorm:"size:64" json:"md5"`；固件 MD5 校验

	MD5             string    `gorm:"size:64" json:"md5"`
	// 结构体字段 FileSize，类型为 int64，附带标签 `json:"file_size"`；固件大小（字节）

	FileSize        int64     `json:"file_size"`
	// 结构体字段 ForceUpgrade，类型为 bool，附带标签 `json:"force_upgrade"`；是否强制升级

	ForceUpgrade    bool      `json:"force_upgrade"`
	// 结构体字段 Status，类型为 int，附带标签 `json:"status"`；Ota* 常量

	Status          int       `json:"status"`
	// 结构体字段 Progress，类型为 int，附带标签 `json:"progress"`；进度百分比 0~100

	Progress        int       `json:"progress"`
	// 结构体字段 Message，类型为 string，附带标签 `gorm:"size:255" json:"message"`；设备上报进度消息

	Message         string    `gorm:"size:255" json:"message"`
	// 结构体字段 CreatedAt，类型为 time.Time，附带标签 `json:"created_at"`

	CreatedAt       time.Time `json:"created_at"`
	// 结构体字段 UpdatedAt，类型为 time.Time，附带标签 `json:"updated_at"`

	UpdatedAt       time.Time `json:"updated_at"`
// 当前代码块结束

}
