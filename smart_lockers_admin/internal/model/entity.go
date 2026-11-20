// @ticket B端-数据模型：网点/设备/对账/告警/账号实体（字段对齐C端语义，不臆造指标）
package model

import "time"

// AdminRole 三类角色 @ticket B端-RBAC
type AdminRole string

const (
	RoleNormalUser AdminRole = "user"       // 普通用户：只读
	RoleSiteAdmin  AdminRole = "site_admin" // 网点管理员：本网点读写
	RolePlatformOp AdminRole = "operator"   // 平台运营：全局读写
)

// AdminUser 后台账号 @ticket B端-RBAC
type AdminUser struct {
	ID       uint      `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"-"` // 脱敏：永不输出
	Role     AdminRole `json:"role"`
	SiteID   uint      `json:"site_id"`
	Status   int       `json:"status"`
}

// Site 网点 @ticket B端-网点管理
type Site struct {
	ID           uint      `json:"id"`
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	City         string    `json:"city"`
	Address      string    `json:"address"`
	ManagerPhone string    `json:"manager_phone"` // 展示前脱敏
	Status       int       `json:"status"`        // 1启用 0停用
	DoorTotal    int       `json:"door_total"`
	DoorFault    int       `json:"door_fault"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Device 设备（柜组/柜门监控）@ticket B端-设备监控
type Device struct {
	ID            uint      `json:"id"`
	SiteID        uint      `json:"site_id"`
	Code          string    `json:"code"`
	Type          string    `json:"type"`   // cabinet / door
	Status        int       `json:"status"` // 0在线 1离线 2故障
	LastHeartbeat time.Time `json:"last_heartbeat"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ReconcileRecord 充值对账记录 @ticket B端-充值对账
type ReconcileRecord struct {
	ID             uint      `json:"id"`
	BatchNo        string    `json:"batch_no"`
	SiteID         uint      `json:"site_id"`
	ExpectedAmount int64     `json:"expected_amount"` // 分
	ActualAmount   int64     `json:"actual_amount"`
	DiffAmount     int64     `json:"diff_amount"`
	Status         int       `json:"status"` // 1平 2不平
	IdempotencyKey string    `json:"-"`
	OperatorID     uint      `json:"operator_id"`
	CreatedAt      time.Time `json:"created_at"`
}

// Alarm 告警 @ticket B端-告警处理
type Alarm struct {
	ID             uint       `json:"id"`
	SiteID         uint       `json:"site_id"`
	DeviceID       uint       `json:"device_id"`
	Level          string     `json:"level"` // info/warn/critical
	Content        string     `json:"content"`
	Status         int        `json:"status"` // 0待处理 1已处理 2已忽略
	Handler        string     `json:"handler"`
	HandleNote     string     `json:"handle_note"`
	IdempotencyKey string     `json:"-"`
	CreatedAt      time.Time  `json:"created_at"`
	HandledAt      *time.Time `json:"handled_at"`
}
