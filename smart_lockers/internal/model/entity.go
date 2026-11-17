// @ticket 00-project：全局数据表实体定义（各实体注释标注来源工单）
package model

import "time"

// ============ 工单03 首页 ============

// City 城市 @ticket 03-首页
type City struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `gorm:"size:64;not null" json:"name"`
	Pinyin   string `gorm:"size:128" json:"pinyin"`
	Initial  string `gorm:"size:2" json:"initial"`
	IsHot    bool   `json:"is_hot"`
	Sort     int    `json:"sort"`
	Location string `gorm:"size:128" json:"location"`
}

// LockerSite 寄存点 @ticket 03-首页
type LockerSite struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	CityID    uint    `gorm:"index" json:"city_id"`
	Name      string  `gorm:"size:128;not null" json:"name"`
	Address   string  `gorm:"size:256" json:"address"`
	Lat       float64 `gorm:"index" json:"lat"`
	Lng       float64 `json:"lng"`
	OpenTime  string  `gorm:"size:16" json:"open_time"`
	CloseTime string  `gorm:"size:16" json:"close_time"`
	Status    int     `gorm:"default:1" json:"status"`
}

// Guide 寄存指南 @ticket 03-首页
type Guide struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Version   string `gorm:"size:32;uniqueIndex" json:"version"`
	Lang      string `gorm:"size:16" json:"lang"`
	Content   string `gorm:"type:text" json:"content"`
	IsCurrent bool   `json:"is_current"`
}

// ============ 工单04 物品寄存 ============

// LockerCabinet 柜组 @ticket 04-物品寄存
type LockerCabinet struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	SiteID    uint   `gorm:"index" json:"site_id"`
	Code      string `gorm:"size:32" json:"code"`
	Size      string `gorm:"size:16" json:"size"` // small/medium/large
	Total     int    `json:"total"`
	Available int    `json:"available"`
	Status    int    `gorm:"default:1" json:"status"`
}

// LockerDoor 柜门 @ticket 04-物品寄存
type LockerDoor struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CabinetID uint      `gorm:"index" json:"cabinet_id"`
	Code      string    `gorm:"size:32" json:"code"`
	Size      string    `gorm:"size:16" json:"size"`
	Status    int       `gorm:"default:0" json:"status"` // 0空闲 1占用 2故障
	Version   int       `gorm:"default:0" json:"-"`      // 乐观锁
	UpdatedAt time.Time `json:"updated_at"`
}

// PricingRule 收费规则 @ticket 04-物品寄存
type PricingRule struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Size         string `gorm:"size:16;index" json:"size"`
	FirstHours   int    `json:"first_hours"`
	FirstFee     int64  `json:"first_fee"` // 分
	ExtraHourFee int64  `json:"extra_hour_fee"`
	DailyCap     int64  `json:"daily_cap"`
	Deposit      int64  `json:"deposit"` // 押金（分）
	Priority     int    `json:"priority"`
	Enabled      bool   `json:"enabled"`
}

// DepositOrder 寄存订单 @ticket 04-物品寄存 / 05-订单查看与取出
type DepositOrder struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	OrderNo        string    `gorm:"size:64;uniqueIndex" json:"order_no"`
	UserID         uint      `gorm:"index" json:"user_id"`
	DoorID         uint      `gorm:"index" json:"door_id"`
	SiteID         uint      `json:"site_id"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	Amount         int64     `json:"amount"`
	Deposit        int64     `json:"deposit"`
	Status         int       `json:"status"` // 1进行中 2已完成 3已取消
	IdempotencyKey string    `gorm:"size:128;uniqueIndex" json:"-"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ============ 工单05 订单查看与取出 ============

// TempOpenAuth 临时开柜授权 @ticket 05-订单查看与取出
type TempOpenAuth struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	OrderNo   string    `gorm:"size:64;index" json:"order_no"`
	Token     string    `gorm:"size:128;uniqueIndex" json:"token"`
	MaxTimes  int       `json:"max_times"`
	UsedTimes int       `json:"used_times"`
	ExpireAt  time.Time `json:"expire_at"`
	Used      bool      `json:"used"`
	CreatedAt time.Time `json:"created_at"`
}

// MakeupRecord 超时补缴记录 @ticket 05-订单查看与取出
type MakeupRecord struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	OrderNo        string    `gorm:"size:64;index" json:"order_no"`
	Amount         int64     `json:"amount"`
	Status         int       `json:"status"`
	IdempotencyKey string    `gorm:"size:128;uniqueIndex" json:"-"`
	CreatedAt      time.Time `json:"created_at"`
}

// ============ 工单02 支付宝支付 ============

// Account 账户余额 @ticket 02-支付宝支付
type Account struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"uniqueIndex" json:"user_id"`
	Balance   int64     `json:"balance"` // 分，以 DB 为准
	UpdatedAt time.Time `json:"updated_at"`
}

// BalanceLog 余额明细 @ticket 02-支付宝支付
type BalanceLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"index" json:"user_id"`
	Change       int64     `json:"change"` // 正为入，负为出
	BalanceAfter int64     `json:"balance_after"`
	BizType      string    `gorm:"size:32" json:"biz_type"` // recharge/pay/makeup/refund
	BizNo        string    `gorm:"size:64;index" json:"biz_no"`
	CreatedAt    time.Time `json:"created_at"`
}

// RechargeOrder 充值订单 @ticket 02-支付宝支付
type RechargeOrder struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	ChargeNo  string     `gorm:"size:64;uniqueIndex" json:"charge_no"`
	UserID    uint       `gorm:"index" json:"user_id"`
	Amount    int64      `json:"amount"`
	Status    int        `json:"status"` // 1待支付 2已支付 3失败
	TradeNo   string     `gorm:"size:128" json:"trade_no"`
	CreatedAt time.Time  `json:"created_at"`
	PaidAt    *time.Time `json:"paid_at"`
}

// ============ 工单06 个人中心 ============

// UserProfile 用户资料 @ticket 06-个人中心
type UserProfile struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Phone     string    `gorm:"size:32;uniqueIndex" json:"phone"`
	Nickname  string    `gorm:"size:64" json:"nickname"`
	Avatar    string    `gorm:"size:256" json:"avatar"`
	Role      string    `gorm:"size:32" json:"role"`
	Status    int       `gorm:"default:1" json:"status"`
	TokenVer  int       `gorm:"default:0" json:"-"` // 退出登录时自增，使旧 token 失效
	CreatedAt time.Time `json:"created_at"`
}

// Coupon 优惠卡券 @ticket 06-个人中心
type Coupon struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Code      string    `gorm:"size:64;uniqueIndex" json:"code"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Discount  int64     `json:"discount"`
	Scope     string    `gorm:"size:64" json:"scope"`
	ExpireAt  time.Time `json:"expire_at"`
	Used      bool      `json:"used"`
	CreatedAt time.Time `json:"created_at"`
}

// MyService 我的服务 @ticket 06-个人中心
type MyService struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	Code    string `gorm:"size:64;uniqueIndex" json:"code"`
	Name    string `gorm:"size:64" json:"name"`
	URL     string `gorm:"size:256" json:"url"`
	Perm    string `gorm:"size:64" json:"perm"`
	Roles   string `gorm:"size:128" json:"roles"` // 逗号分隔，可动态配置
	Enabled bool   `json:"enabled"`
}

// ============ 工单01 消息通知 ============

// NotifyTemplate 通知模板 @ticket 01-消息通知
type NotifyTemplate struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Code      string `gorm:"size:64;uniqueIndex" json:"code"`
	Channel   string `gorm:"size:32" json:"channel"`
	Content   string `gorm:"type:text" json:"content"`
	Variables string `gorm:"type:text" json:"variables"`
	Enabled   bool   `json:"enabled"`
}

// NotifyMessage 通知消息 @ticket 01-消息通知
type NotifyMessage struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	TplCode      string    `gorm:"size:64;index" json:"tpl_code"`
	ReceiverType string    `gorm:"size:32" json:"receiver_type"` // admin/site_admin/user
	ReceiverID   uint      `gorm:"index" json:"receiver_id"`
	Channel      string    `gorm:"size:32" json:"channel"`
	Content      string    `gorm:"type:text" json:"content"`
	Status       int       `json:"status"` // 0待发 1成功 2失败 3重试中
	RetryCount   int       `json:"retry_count"`
	DedupeKey    string    `gorm:"size:128;index" json:"dedupe_key"`
	CreatedAt    time.Time `json:"created_at"`
}

// ============ 工单08 性能分析（无需持久化表，指标通过 expvar/Prometheus 暴露）============
