// @ticket B端-IoT中台：指令下发/OTA 升级实体（对齐《智能快递柜IoT设备MQTT接入技术方案》与运营中台原型 devices/ota 视图）
package model

import "time"

// DeviceCommand B端指令记录（platform/command/{id}/… 下发 + device/response/{id}/… 响应）
// 定义结构体 DeviceCommand：用于承载一组相关的字段数据

type DeviceCommand struct {
	// 结构体字段 ID，类型为 uint，附带标签 `gorm:"primaryKey" json:"id"`

	ID          uint       `gorm:"primaryKey" json:"id"`
	// 结构体字段 CommandID，类型为 string，附带标签 `gorm:"size:64;uniqueIndex" json:"command_id"`

	CommandID   string     `gorm:"size:64;uniqueIndex" json:"command_id"`
	// 结构体字段 DeviceID，类型为 uint，附带标签 `gorm:"index" json:"device_id"`

	DeviceID    uint       `gorm:"index" json:"device_id"`
	// 结构体字段 DeviceSN，类型为 string，附带标签 `gorm:"size:64" json:"device_sn"`

	DeviceSN    string     `gorm:"size:64" json:"device_sn"`
	// 结构体字段 CmdType，类型为 string，附带标签 `gorm:"size:32" json:"cmd_type"`；open_door/close_door/reboot/config

	CmdType     string     `gorm:"size:32" json:"cmd_type"`
	// 结构体字段 LockerNo，类型为 string，附带标签 `gorm:"size:32" json:"locker_no"`

	LockerNo    string     `gorm:"size:32" json:"locker_no"`
	// 结构体字段 Payload，类型为 string，附带标签 `gorm:"type:text" json:"payload"`

	Payload     string     `gorm:"type:text" json:"payload"`
	// 结构体字段 Status，类型为 int，附带标签 `json:"status"`；0已下发 1成功 2失败 3超时

	Status      int        `json:"status"`
	// 结构体字段 Response，类型为 string，附带标签 `gorm:"type:text" json:"response"`

	Response    string     `gorm:"type:text" json:"response"`
	// 结构体字段 Operator，类型为 string，附带标签 `gorm:"size:64" json:"operator"`

	Operator    string     `gorm:"size:64" json:"operator"`
	// 结构体字段 CreatedAt，类型为 time.Time，附带标签 `json:"created_at"`

	CreatedAt   time.Time  `json:"created_at"`
	// 结构体字段 RespondedAt，类型为 *time.Time，附带标签 `json:"responded_at"`

	RespondedAt *time.Time `json:"responded_at"`
// 当前代码块结束

}

// OtaTask B端 OTA 升级任务（platform/ota/{id}/upgrade 下发 + device/ota/{id}/progress 进度）
// 定义结构体 OtaTask：用于承载一组相关的字段数据

type OtaTask struct {
	// 结构体字段 ID，类型为 uint，附带标签 `gorm:"primaryKey" json:"id"`

	ID              uint      `gorm:"primaryKey" json:"id"`
	// 结构体字段 OtaID，类型为 string，附带标签 `gorm:"size:64;uniqueIndex" json:"ota_id"`

	OtaID           string    `gorm:"size:64;uniqueIndex" json:"ota_id"`
	// 结构体字段 DeviceID，类型为 uint，附带标签 `gorm:"index" json:"device_id"`

	DeviceID        uint      `gorm:"index" json:"device_id"`
	// 结构体字段 DeviceSN，类型为 string，附带标签 `gorm:"size:64" json:"device_sn"`

	DeviceSN        string    `gorm:"size:64" json:"device_sn"`
	// 结构体字段 FirmwareVersion，类型为 string，附带标签 `gorm:"size:32" json:"firmware_version"`

	FirmwareVersion string    `gorm:"size:32" json:"firmware_version"`
	// 结构体字段 DownloadURL，类型为 string，附带标签 `gorm:"size:512" json:"download_url"`

	DownloadURL     string    `gorm:"size:512" json:"download_url"`
	// 结构体字段 MD5，类型为 string，附带标签 `gorm:"size:64" json:"md5"`

	MD5             string    `gorm:"size:64" json:"md5"`
	// 结构体字段 FileSize，类型为 int64，附带标签 `json:"file_size"`

	FileSize        int64     `json:"file_size"`
	// 结构体字段 ForceUpgrade，类型为 bool，附带标签 `json:"force_upgrade"`

	ForceUpgrade    bool      `json:"force_upgrade"`
	// 结构体字段 Status，类型为 int，附带标签 `json:"status"`；0已下发 1下载中 2升级中 3成功 4失败

	Status          int       `json:"status"`
	// 结构体字段 Progress，类型为 int，附带标签 `json:"progress"`

	Progress        int       `json:"progress"`
	// 结构体字段 Message，类型为 string，附带标签 `gorm:"size:255" json:"message"`

	Message         string    `gorm:"size:255" json:"message"`
	// 结构体字段 Operator，类型为 string，附带标签 `gorm:"size:64" json:"operator"`

	Operator        string    `gorm:"size:64" json:"operator"`
	// 结构体字段 CreatedAt，类型为 time.Time，附带标签 `json:"created_at"`

	CreatedAt       time.Time `json:"created_at"`
	// 结构体字段 UpdatedAt，类型为 time.Time，附带标签 `json:"updated_at"`

	UpdatedAt       time.Time `json:"updated_at"`
// 当前代码块结束

}
