// @ticket B端-数据模型：网点/设备/对账/告警/账号实体（字段对齐C端语义，不臆造指标）
package model

import "time"

// AdminRole 三类角色 @ticket B端-RBAC
// 定义类型 AdminRole，其底层类型为 string

type AdminRole string

// 开始定义一组常量的 const 声明块

const (
	// 为 RoleNormalUser AdminRole 赋值："user"       // 普通用户：只读

	RoleNormalUser AdminRole = "user"       // 普通用户：只读
	// 为 RoleSiteAdmin  AdminRole 赋值："site_admin" // 网点管理员：本网点读写

	RoleSiteAdmin  AdminRole = "site_admin" // 网点管理员：本网点读写
	// 为 RolePlatformOp AdminRole 赋值："operator"   // 平台运营：全局读写

	RolePlatformOp AdminRole = "operator"   // 平台运营：全局读写
// 结束参数/表达式列表

)

// AdminUser 后台账号 @ticket B端-RBAC
// 定义结构体 AdminUser：用于承载一组相关的字段数据

type AdminUser struct {
	// 结构体字段 ID，类型为 uint，附带标签 `json:"id"`

	ID       uint      `json:"id"`
	// 结构体字段 Username，类型为 string，附带标签 `json:"username"`

	Username string    `json:"username"`
	// 结构体字段 Password，类型为 string，附带标签 `json:"-"`

	Password string    `json:"-"` // 脱敏：永不输出
	// 结构体字段 Role，类型为 AdminRole，附带标签 `json:"role"`

	Role     AdminRole `json:"role"`
	// 结构体字段 SiteID，类型为 uint，附带标签 `json:"site_id"`

	SiteID   uint      `json:"site_id"`
	// 结构体字段 Status，类型为 int，附带标签 `json:"status"`

	Status   int       `json:"status"`
// 当前代码块结束

}

// Site 网点 @ticket B端-网点管理
// 定义结构体 Site：用于承载一组相关的字段数据

type Site struct {
	// 结构体字段 ID，类型为 uint，附带标签 `json:"id"`

	ID           uint      `json:"id"`
	// 结构体字段 Code，类型为 string，附带标签 `json:"code"`

	Code         string    `json:"code"`
	// 结构体字段 Name，类型为 string，附带标签 `json:"name"`

	Name         string    `json:"name"`
	// 结构体字段 City，类型为 string，附带标签 `json:"city"`

	City         string    `json:"city"`
	// 结构体字段 Address，类型为 string，附带标签 `json:"address"`

	Address      string    `json:"address"`
	// 结构体字段 ManagerPhone，类型为 string，附带标签 `json:"manager_phone"`

	ManagerPhone string    `json:"manager_phone"` // 展示前脱敏
	// 结构体字段 Status，类型为 int，附带标签 `json:"status"`

	Status       int       `json:"status"`        // 1启用 0停用
	// 结构体字段 DoorTotal，类型为 int，附带标签 `json:"door_total"`

	DoorTotal    int       `json:"door_total"`
	// 结构体字段 DoorFault，类型为 int，附带标签 `json:"door_fault"`

	DoorFault    int       `json:"door_fault"`
	// 结构体字段 CreatedAt，类型为 time.Time，附带标签 `json:"created_at"`

	CreatedAt    time.Time `json:"created_at"`
	// 结构体字段 UpdatedAt，类型为 time.Time，附带标签 `json:"updated_at"`

	UpdatedAt    time.Time `json:"updated_at"`
// 当前代码块结束

}

// Device 设备（柜组/柜门监控）@ticket B端-设备监控
// 定义结构体 Device：用于承载一组相关的字段数据

type Device struct {
	// 结构体字段 ID，类型为 uint，附带标签 `json:"id"`

	ID            uint      `json:"id"`
	// 结构体字段 SiteID，类型为 uint，附带标签 `json:"site_id"`

	SiteID        uint      `json:"site_id"`
	// 结构体字段 Code，类型为 string，附带标签 `json:"code"`

	Code          string    `json:"code"`
	// 结构体字段 Type，类型为 string，附带标签 `json:"type"`

	Type          string    `json:"type"`   // cabinet / door
	// 结构体字段 Status，类型为 int，附带标签 `json:"status"`

	Status        int       `json:"status"` // 0在线 1离线 2故障
	// 结构体字段 LastHeartbeat，类型为 time.Time，附带标签 `json:"last_heartbeat"`

	LastHeartbeat time.Time `json:"last_heartbeat"`
	// 结构体字段 UpdatedAt，类型为 time.Time，附带标签 `json:"updated_at"`

	UpdatedAt     time.Time `json:"updated_at"`
// 当前代码块结束

}

// ReconcileRecord 充值对账记录 @ticket B端-充值对账
// 定义结构体 ReconcileRecord：用于承载一组相关的字段数据

type ReconcileRecord struct {
	// 结构体字段 ID，类型为 uint，附带标签 `json:"id"`

	ID             uint      `json:"id"`
	// 结构体字段 BatchNo，类型为 string，附带标签 `json:"batch_no"`

	BatchNo        string    `json:"batch_no"`
	// 结构体字段 SiteID，类型为 uint，附带标签 `json:"site_id"`

	SiteID         uint      `json:"site_id"`
	// 结构体字段 ExpectedAmount，类型为 int64，附带标签 `json:"expected_amount"`

	ExpectedAmount int64     `json:"expected_amount"` // 分
	// 结构体字段 ActualAmount，类型为 int64，附带标签 `json:"actual_amount"`

	ActualAmount   int64     `json:"actual_amount"`
	// 结构体字段 DiffAmount，类型为 int64，附带标签 `json:"diff_amount"`

	DiffAmount     int64     `json:"diff_amount"`
	// 结构体字段 Status，类型为 int，附带标签 `json:"status"`

	Status         int       `json:"status"` // 1平 2不平
	// 修复偏差6：输出幂等键字段，前端 ReconcileRecord 类型声明含 idempotency_key @ticket B端-对账

	IdempotencyKey string    `json:"idempotency_key"`
	// 结构体字段 OperatorID，类型为 uint，附带标签 `json:"operator_id"`

	OperatorID     uint      `json:"operator_id"`
	// 结构体字段 CreatedAt，类型为 time.Time，附带标签 `json:"created_at"`

	CreatedAt      time.Time `json:"created_at"`
// 当前代码块结束

}

// Alarm 告警 @ticket B端-告警处理
// 定义结构体 Alarm：用于承载一组相关的字段数据

type Alarm struct {
	// 结构体字段 ID，类型为 uint，附带标签 `json:"id"`

	ID             uint       `json:"id"`
	// 结构体字段 SiteID，类型为 uint，附带标签 `json:"site_id"`

	SiteID         uint       `json:"site_id"`
	// 结构体字段 DeviceID，类型为 uint，附带标签 `json:"device_id"`

	DeviceID       uint       `json:"device_id"`
	// 结构体字段 Level，类型为 string，附带标签 `json:"level"`

	Level          string     `json:"level"` // info/warn/critical
	// 结构体字段 Content，类型为 string，附带标签 `json:"content"`

	Content        string     `json:"content"`
	// 结构体字段 Status，类型为 int，附带标签 `json:"status"`

	Status         int        `json:"status"` // 0待处理 1已处理 2已忽略
	// 结构体字段 Handler，类型为 string，附带标签 `json:"handler"`

	Handler        string     `json:"handler"`
	// 结构体字段 HandleNote，类型为 string，附带标签 `json:"handle_note"`

	HandleNote     string     `json:"handle_note"`
	// 修复偏差6：输出幂等键字段，前端 Alarm 类型声明含 idempotency_key @ticket B端-告警

	IdempotencyKey string     `json:"idempotency_key"`
	// 结构体字段 CreatedAt，类型为 time.Time，附带标签 `json:"created_at"`

	CreatedAt      time.Time  `json:"created_at"`
	// 结构体字段 HandledAt，类型为 *time.Time，附带标签 `json:"handled_at"`

	HandledAt      *time.Time `json:"handled_at"`
// 当前代码块结束

}
