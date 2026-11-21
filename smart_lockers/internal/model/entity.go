// @ticket 00-project：全局数据表实体定义（各实体注释标注来源工单）
package model

import "time"

// ============ 工单03 首页 ============

// City 城市 @ticket 03-首页
// 定义结构体 City：用于承载一组相关的字段数据

type City struct {
	// 结构体字段 ID，类型为 uint，附带标签 `gorm:"primaryKey" json:"id"`

	ID       uint   `gorm:"primaryKey" json:"id"`
	// 结构体字段 Name，类型为 string，附带标签 `gorm:"size:64;not null" json:"name"`

	Name     string `gorm:"size:64;not null" json:"name"`
	// 结构体字段 Pinyin，类型为 string，附带标签 `gorm:"size:128" json:"pinyin"`

	Pinyin   string `gorm:"size:128" json:"pinyin"`
	// 结构体字段 Initial，类型为 string，附带标签 `gorm:"size:2" json:"initial"`

	Initial  string `gorm:"size:2" json:"initial"`
	// 结构体字段 IsHot，类型为 bool，附带标签 `json:"is_hot"`

	IsHot    bool   `json:"is_hot"`
	// 结构体字段 Sort，类型为 int，附带标签 `json:"sort"`

	Sort     int    `json:"sort"`
	// 结构体字段 Location，类型为 string，附带标签 `gorm:"size:128" json:"location"`

	Location string `gorm:"size:128" json:"location"`
// 当前代码块结束

}

// LockerSite 寄存点 @ticket 03-首页
// 定义结构体 LockerSite：用于承载一组相关的字段数据

type LockerSite struct {
	// 结构体字段 ID，类型为 uint，附带标签 `gorm:"primaryKey" json:"id"`

	ID        uint    `gorm:"primaryKey" json:"id"`
	// 结构体字段 CityID，类型为 uint，附带标签 `gorm:"index" json:"city_id"`

	CityID    uint    `gorm:"index" json:"city_id"`
	// 结构体字段 Name，类型为 string，附带标签 `gorm:"size:128;not null" json:"name"`

	Name      string  `gorm:"size:128;not null" json:"name"`
	// 结构体字段 Address，类型为 string，附带标签 `gorm:"size:256" json:"address"`

	Address   string  `gorm:"size:256" json:"address"`
	// 结构体字段 Lat，类型为 float64，附带标签 `gorm:"index" json:"lat"`

	Lat       float64 `gorm:"index" json:"lat"`
	// 结构体字段 Lng，类型为 float64，附带标签 `json:"lng"`

	Lng       float64 `json:"lng"`
	// 结构体字段 OpenTime，类型为 string，附带标签 `gorm:"size:16" json:"open_time"`

	OpenTime  string  `gorm:"size:16" json:"open_time"`
	// 结构体字段 CloseTime，类型为 string，附带标签 `gorm:"size:16" json:"close_time"`

	CloseTime string  `gorm:"size:16" json:"close_time"`
	// 结构体字段 Status，类型为 int，附带标签 `gorm:"default:1" json:"status"`

	Status    int     `gorm:"default:1" json:"status"`
// 当前代码块结束

}

// Guide 寄存指南 @ticket 03-首页
// 定义结构体 Guide：用于承载一组相关的字段数据

type Guide struct {
	// 结构体字段 ID，类型为 uint，附带标签 `gorm:"primaryKey" json:"id"`

	ID        uint   `gorm:"primaryKey" json:"id"`
	// 结构体字段 Version，类型为 string，附带标签 `gorm:"size:32;uniqueIndex" json:"version"`

	Version   string `gorm:"size:32;uniqueIndex" json:"version"`
	// 结构体字段 Lang，类型为 string，附带标签 `gorm:"size:16" json:"lang"`

	Lang      string `gorm:"size:16" json:"lang"`
	// 结构体字段 Content，类型为 string，附带标签 `gorm:"type:text" json:"content"`

	Content   string `gorm:"type:text" json:"content"`
	// 结构体字段 IsCurrent，类型为 bool，附带标签 `json:"is_current"`

	IsCurrent bool   `json:"is_current"`
// 当前代码块结束

}

// ============ 工单04 物品寄存 ============

// LockerCabinet 柜组 @ticket 04-物品寄存
// 定义结构体 LockerCabinet：用于承载一组相关的字段数据

type LockerCabinet struct {
	// 结构体字段 ID，类型为 uint，附带标签 `gorm:"primaryKey" json:"id"`

	ID        uint   `gorm:"primaryKey" json:"id"`
	// 结构体字段 SiteID，类型为 uint，附带标签 `gorm:"index" json:"site_id"`

	SiteID    uint   `gorm:"index" json:"site_id"`
	// 结构体字段 Code，类型为 string，附带标签 `gorm:"size:32" json:"code"`

	Code      string `gorm:"size:32" json:"code"`
	// 结构体字段 Size，类型为 string，附带标签 `gorm:"size:16" json:"size"`

	Size      string `gorm:"size:16" json:"size"` // small/medium/large
	// 结构体字段 Total，类型为 int，附带标签 `json:"total"`

	Total     int    `json:"total"`
	// 结构体字段 Available，类型为 int，附带标签 `json:"available"`

	Available int    `json:"available"`
	// 结构体字段 Status，类型为 int，附带标签 `gorm:"default:1" json:"status"`

	Status    int    `gorm:"default:1" json:"status"`
// 当前代码块结束

}

// LockerDoor 柜门 @ticket 04-物品寄存
// 定义结构体 LockerDoor：用于承载一组相关的字段数据

type LockerDoor struct {
	// 结构体字段 ID，类型为 uint，附带标签 `gorm:"primaryKey" json:"id"`

	ID        uint      `gorm:"primaryKey" json:"id"`
	// 结构体字段 CabinetID，类型为 uint，附带标签 `gorm:"index" json:"cabinet_id"`

	CabinetID uint      `gorm:"index" json:"cabinet_id"`
	// 结构体字段 Code，类型为 string，附带标签 `gorm:"size:32" json:"code"`

	Code      string    `gorm:"size:32" json:"code"`
	// 结构体字段 Size，类型为 string，附带标签 `gorm:"size:16" json:"size"`

	Size      string    `gorm:"size:16" json:"size"`
	// 结构体字段 Status，类型为 int，附带标签 `gorm:"default:0" json:"status"`

	Status    int       `gorm:"default:0" json:"status"` // 0空闲 1占用 2故障
	// 结构体字段 Version，类型为 int，附带标签 `gorm:"default:0" json:"-"`

	Version   int       `gorm:"default:0" json:"-"`      // 乐观锁
	// 结构体字段 UpdatedAt，类型为 time.Time，附带标签 `json:"updated_at"`

	UpdatedAt time.Time `json:"updated_at"`
// 当前代码块结束

}

// PricingRule 收费规则 @ticket 04-物品寄存
// 定义结构体 PricingRule：用于承载一组相关的字段数据

type PricingRule struct {
	// 结构体字段 ID，类型为 uint，附带标签 `gorm:"primaryKey" json:"id"`

	ID           uint   `gorm:"primaryKey" json:"id"`
	// 结构体字段 Size，类型为 string，附带标签 `gorm:"size:16;index" json:"size"`

	Size         string `gorm:"size:16;index" json:"size"`
	// 结构体字段 FirstHours，类型为 int，附带标签 `json:"first_hours"`

	FirstHours   int    `json:"first_hours"`
	// 结构体字段 FirstFee，类型为 int64，附带标签 `json:"first_fee"`

	FirstFee     int64  `json:"first_fee"` // 分
	// 结构体字段 ExtraHourFee，类型为 int64，附带标签 `json:"extra_hour_fee"`

	ExtraHourFee int64  `json:"extra_hour_fee"`
	// 结构体字段 DailyCap，类型为 int64，附带标签 `json:"daily_cap"`

	DailyCap     int64  `json:"daily_cap"`
	// 结构体字段 Deposit，类型为 int64，附带标签 `json:"deposit"`

	Deposit      int64  `json:"deposit"` // 押金（分）
	// 结构体字段 Priority，类型为 int，附带标签 `json:"priority"`

	Priority     int    `json:"priority"`
	// 结构体字段 Enabled，类型为 bool，附带标签 `json:"enabled"`

	Enabled      bool   `json:"enabled"`
// 当前代码块结束

}

// DepositOrder 寄存订单 @ticket 04-物品寄存 / 05-订单查看与取出
// 定义结构体 DepositOrder：用于承载一组相关的字段数据

type DepositOrder struct {
	// 结构体字段 ID，类型为 uint，附带标签 `gorm:"primaryKey" json:"id"`

	ID             uint      `gorm:"primaryKey" json:"id"`
	// 结构体字段 OrderNo，类型为 string，附带标签 `gorm:"size:64;uniqueIndex" json:"order_no"`

	OrderNo        string    `gorm:"size:64;uniqueIndex" json:"order_no"`
	// 结构体字段 UserID，类型为 uint，附带标签 `gorm:"index" json:"user_id"`

	UserID         uint      `gorm:"index" json:"user_id"`
	// 结构体字段 DoorID，类型为 uint，附带标签 `gorm:"index" json:"door_id"`

	DoorID         uint      `gorm:"index" json:"door_id"`
	// 结构体字段 SiteID，类型为 uint，附带标签 `json:"site_id"`

	SiteID         uint      `json:"site_id"`
	// 结构体字段 StartTime，类型为 time.Time，附带标签 `json:"start_time"`

	StartTime      time.Time `json:"start_time"`
	// 结构体字段 EndTime，类型为 time.Time，附带标签 `json:"end_time"`

	EndTime        time.Time `json:"end_time"`
	// 结构体字段 Amount，类型为 int64，附带标签 `json:"amount"`

	Amount         int64     `json:"amount"`
	// 结构体字段 Deposit，类型为 int64，附带标签 `json:"deposit"`

	Deposit        int64     `json:"deposit"`
	// 结构体字段 Status，类型为 int，附带标签 `json:"status"`

	Status         int       `json:"status"` // 1进行中 2已完成 3已取消
	// 结构体字段 IdempotencyKey，类型为 string，附带标签 `gorm:"size:128;uniqueIndex" json:"-"`

	IdempotencyKey string    `gorm:"size:128;uniqueIndex" json:"-"`
	// 结构体字段 CreatedAt，类型为 time.Time，附带标签 `json:"created_at"`

	CreatedAt      time.Time `json:"created_at"`
	// 结构体字段 UpdatedAt，类型为 time.Time，附带标签 `json:"updated_at"`

	UpdatedAt      time.Time `json:"updated_at"`
// 当前代码块结束

}

// ============ 工单05 订单查看与取出 ============

// TempOpenAuth 临时开柜授权 @ticket 05-订单查看与取出
// 定义结构体 TempOpenAuth：用于承载一组相关的字段数据

type TempOpenAuth struct {
	// 结构体字段 ID，类型为 uint，附带标签 `gorm:"primaryKey" json:"id"`

	ID        uint      `gorm:"primaryKey" json:"id"`
	// 结构体字段 OrderNo，类型为 string，附带标签 `gorm:"size:64;index" json:"order_no"`

	OrderNo   string    `gorm:"size:64;index" json:"order_no"`
	// 结构体字段 Token，类型为 string，附带标签 `gorm:"size:128;uniqueIndex" json:"token"`

	Token     string    `gorm:"size:128;uniqueIndex" json:"token"`
	// 结构体字段 MaxTimes，类型为 int，附带标签 `json:"max_times"`

	MaxTimes  int       `json:"max_times"`
	// 结构体字段 UsedTimes，类型为 int，附带标签 `json:"used_times"`

	UsedTimes int       `json:"used_times"`
	// 结构体字段 ExpireAt，类型为 time.Time，附带标签 `json:"expire_at"`

	ExpireAt  time.Time `json:"expire_at"`
	// 结构体字段 Used，类型为 bool，附带标签 `json:"used"`

	Used      bool      `json:"used"`
	// 结构体字段 CreatedAt，类型为 time.Time，附带标签 `json:"created_at"`

	CreatedAt time.Time `json:"created_at"`
// 当前代码块结束

}

// MakeupRecord 超时补缴记录 @ticket 05-订单查看与取出
// 定义结构体 MakeupRecord：用于承载一组相关的字段数据

type MakeupRecord struct {
	// 结构体字段 ID，类型为 uint，附带标签 `gorm:"primaryKey" json:"id"`

	ID             uint      `gorm:"primaryKey" json:"id"`
	// 结构体字段 OrderNo，类型为 string，附带标签 `gorm:"size:64;index" json:"order_no"`

	OrderNo        string    `gorm:"size:64;index" json:"order_no"`
	// 结构体字段 Amount，类型为 int64，附带标签 `json:"amount"`

	Amount         int64     `json:"amount"`
	// 结构体字段 Status，类型为 int，附带标签 `json:"status"`

	Status         int       `json:"status"`
	// 结构体字段 IdempotencyKey，类型为 string，附带标签 `gorm:"size:128;uniqueIndex" json:"-"`

	IdempotencyKey string    `gorm:"size:128;uniqueIndex" json:"-"`
	// 结构体字段 CreatedAt，类型为 time.Time，附带标签 `json:"created_at"`

	CreatedAt      time.Time `json:"created_at"`
// 当前代码块结束

}

// ============ 工单02 支付宝支付 ============

// Account 账户余额 @ticket 02-支付宝支付
// 定义结构体 Account：用于承载一组相关的字段数据

type Account struct {
	// 结构体字段 ID，类型为 uint，附带标签 `gorm:"primaryKey" json:"id"`

	ID        uint      `gorm:"primaryKey" json:"id"`
	// 结构体字段 UserID，类型为 uint，附带标签 `gorm:"uniqueIndex" json:"user_id"`

	UserID    uint      `gorm:"uniqueIndex" json:"user_id"`
	// 结构体字段 Balance，类型为 int64，附带标签 `json:"balance"`

	Balance   int64     `json:"balance"` // 分，以 DB 为准
	// 结构体字段 UpdatedAt，类型为 time.Time，附带标签 `json:"updated_at"`

	UpdatedAt time.Time `json:"updated_at"`
// 当前代码块结束

}

// BalanceLog 余额明细 @ticket 02-支付宝支付
// 定义结构体 BalanceLog：用于承载一组相关的字段数据

type BalanceLog struct {
	// 结构体字段 ID，类型为 uint，附带标签 `gorm:"primaryKey" json:"id"`

	ID           uint      `gorm:"primaryKey" json:"id"`
	// 结构体字段 UserID，类型为 uint，附带标签 `gorm:"index" json:"user_id"`

	UserID       uint      `gorm:"index" json:"user_id"`
	// 结构体字段 Change，类型为 int64，附带标签 `json:"change"`

	Change       int64     `json:"change"` // 正为入，负为出
	// 结构体字段 BalanceAfter，类型为 int64，附带标签 `json:"balance_after"`

	BalanceAfter int64     `json:"balance_after"`
	// 结构体字段 BizType，类型为 string，附带标签 `gorm:"size:32" json:"biz_type"`

	BizType      string    `gorm:"size:32" json:"biz_type"` // recharge/pay/makeup/refund
	// 结构体字段 BizNo，类型为 string，附带标签 `gorm:"size:64;index" json:"biz_no"`

	BizNo        string    `gorm:"size:64;index" json:"biz_no"`
	// 结构体字段 CreatedAt，类型为 time.Time，附带标签 `json:"created_at"`

	CreatedAt    time.Time `json:"created_at"`
// 当前代码块结束

}

// RechargeOrder 充值订单 @ticket 02-支付宝支付
// 定义结构体 RechargeOrder：用于承载一组相关的字段数据

type RechargeOrder struct {
	// 结构体字段 ID，类型为 uint，附带标签 `gorm:"primaryKey" json:"id"`

	ID        uint       `gorm:"primaryKey" json:"id"`
	// 结构体字段 ChargeNo，类型为 string，附带标签 `gorm:"size:64;uniqueIndex" json:"charge_no"`

	ChargeNo  string     `gorm:"size:64;uniqueIndex" json:"charge_no"`
	// 结构体字段 UserID，类型为 uint，附带标签 `gorm:"index" json:"user_id"`

	UserID    uint       `gorm:"index" json:"user_id"`
	// 结构体字段 Amount，类型为 int64，附带标签 `json:"amount"`

	Amount    int64      `json:"amount"`
	// 结构体字段 Status，类型为 int，附带标签 `json:"status"`

	Status    int        `json:"status"` // 1待支付 2已支付 3失败
	// 结构体字段 TradeNo，类型为 string，附带标签 `gorm:"size:128" json:"trade_no"`

	TradeNo   string     `gorm:"size:128" json:"trade_no"`
	// 结构体字段 CreatedAt，类型为 time.Time，附带标签 `json:"created_at"`

	CreatedAt time.Time  `json:"created_at"`
	// 结构体字段 PaidAt，类型为 *time.Time，附带标签 `json:"paid_at"`

	PaidAt    *time.Time `json:"paid_at"`
// 当前代码块结束

}

// ============ 工单06 个人中心 ============

// UserProfile 用户资料 @ticket 06-个人中心
// 定义结构体 UserProfile：用于承载一组相关的字段数据

type UserProfile struct {
	// 结构体字段 ID，类型为 uint，附带标签 `gorm:"primaryKey" json:"id"`

	ID        uint      `gorm:"primaryKey" json:"id"`
	// 结构体字段 Phone，类型为 string，附带标签 `gorm:"size:32;uniqueIndex" json:"phone"`

	Phone     string    `gorm:"size:32;uniqueIndex" json:"phone"`
	// 结构体字段 Nickname，类型为 string，附带标签 `gorm:"size:64" json:"nickname"`

	Nickname  string    `gorm:"size:64" json:"nickname"`
	// 结构体字段 Avatar，类型为 string，附带标签 `gorm:"size:256" json:"avatar"`

	Avatar    string    `gorm:"size:256" json:"avatar"`
	// 结构体字段 Role，类型为 string，附带标签 `gorm:"size:32" json:"role"`

	Role      string    `gorm:"size:32" json:"role"`
	// 结构体字段 Status，类型为 int，附带标签 `gorm:"default:1" json:"status"`

	Status    int       `gorm:"default:1" json:"status"`
	// 结构体字段 TokenVer，类型为 int，附带标签 `gorm:"default:0" json:"-"`

	TokenVer  int       `gorm:"default:0" json:"-"` // 退出登录时自增，使旧 token 失效
	// 结构体字段 CreatedAt，类型为 time.Time，附带标签 `json:"created_at"`

	CreatedAt time.Time `json:"created_at"`
// 当前代码块结束

}

// Coupon 优惠卡券 @ticket 06-个人中心
// 定义结构体 Coupon：用于承载一组相关的字段数据

type Coupon struct {
	// 结构体字段 ID，类型为 uint，附带标签 `gorm:"primaryKey" json:"id"`

	ID        uint      `gorm:"primaryKey" json:"id"`
	// 结构体字段 Code，类型为 string，附带标签 `gorm:"size:64;uniqueIndex" json:"code"`

	Code      string    `gorm:"size:64;uniqueIndex" json:"code"`
	// 结构体字段 UserID，类型为 uint，附带标签 `gorm:"index" json:"user_id"`

	UserID    uint      `gorm:"index" json:"user_id"`
	// 结构体字段 Discount，类型为 int64，附带标签 `json:"discount"`

	Discount  int64     `json:"discount"`
	// 结构体字段 Scope，类型为 string，附带标签 `gorm:"size:64" json:"scope"`

	Scope     string    `gorm:"size:64" json:"scope"`
	// 结构体字段 ExpireAt，类型为 time.Time，附带标签 `json:"expire_at"`

	ExpireAt  time.Time `json:"expire_at"`
	// 结构体字段 Used，类型为 bool，附带标签 `json:"used"`

	Used      bool      `json:"used"`
	// 结构体字段 CreatedAt，类型为 time.Time，附带标签 `json:"created_at"`

	CreatedAt time.Time `json:"created_at"`
// 当前代码块结束

}

// MyService 我的服务 @ticket 06-个人中心
// 定义结构体 MyService：用于承载一组相关的字段数据

type MyService struct {
	// 结构体字段 ID，类型为 uint，附带标签 `gorm:"primaryKey" json:"id"`

	ID      uint   `gorm:"primaryKey" json:"id"`
	// 结构体字段 Code，类型为 string，附带标签 `gorm:"size:64;uniqueIndex" json:"code"`

	Code    string `gorm:"size:64;uniqueIndex" json:"code"`
	// 结构体字段 Name，类型为 string，附带标签 `gorm:"size:64" json:"name"`

	Name    string `gorm:"size:64" json:"name"`
	// 结构体字段 URL，类型为 string，附带标签 `gorm:"size:256" json:"url"`

	URL     string `gorm:"size:256" json:"url"`
	// 结构体字段 Perm，类型为 string，附带标签 `gorm:"size:64" json:"perm"`

	Perm    string `gorm:"size:64" json:"perm"`
	// 结构体字段 Roles，类型为 string，附带标签 `gorm:"size:128" json:"roles"`

	Roles   string `gorm:"size:128" json:"roles"` // 逗号分隔，可动态配置
	// 结构体字段 Enabled，类型为 bool，附带标签 `json:"enabled"`

	Enabled bool   `json:"enabled"`
// 当前代码块结束

}

// ============ 工单01 消息通知 ============

// NotifyTemplate 通知模板 @ticket 01-消息通知
// 定义结构体 NotifyTemplate：用于承载一组相关的字段数据

type NotifyTemplate struct {
	// 结构体字段 ID，类型为 uint，附带标签 `gorm:"primaryKey" json:"id"`

	ID        uint   `gorm:"primaryKey" json:"id"`
	// 结构体字段 Code，类型为 string，附带标签 `gorm:"size:64;uniqueIndex" json:"code"`

	Code      string `gorm:"size:64;uniqueIndex" json:"code"`
	// 结构体字段 Channel，类型为 string，附带标签 `gorm:"size:32" json:"channel"`

	Channel   string `gorm:"size:32" json:"channel"`
	// 结构体字段 Content，类型为 string，附带标签 `gorm:"type:text" json:"content"`

	Content   string `gorm:"type:text" json:"content"`
	// 结构体字段 Variables，类型为 string，附带标签 `gorm:"type:text" json:"variables"`

	Variables string `gorm:"type:text" json:"variables"`
	// 结构体字段 Enabled，类型为 bool，附带标签 `json:"enabled"`

	Enabled   bool   `json:"enabled"`
// 当前代码块结束

}

// NotifyMessage 通知消息 @ticket 01-消息通知
// 定义结构体 NotifyMessage：用于承载一组相关的字段数据

type NotifyMessage struct {
	// 结构体字段 ID，类型为 uint，附带标签 `gorm:"primaryKey" json:"id"`

	ID           uint      `gorm:"primaryKey" json:"id"`
	// 结构体字段 TplCode，类型为 string，附带标签 `gorm:"size:64;index" json:"tpl_code"`

	TplCode      string    `gorm:"size:64;index" json:"tpl_code"`
	// 结构体字段 ReceiverType，类型为 string，附带标签 `gorm:"size:32" json:"receiver_type"`

	ReceiverType string    `gorm:"size:32" json:"receiver_type"` // admin/site_admin/user
	// 结构体字段 ReceiverID，类型为 uint，附带标签 `gorm:"index" json:"receiver_id"`

	ReceiverID   uint      `gorm:"index" json:"receiver_id"`
	// 结构体字段 Channel，类型为 string，附带标签 `gorm:"size:32" json:"channel"`

	Channel      string    `gorm:"size:32" json:"channel"`
	// 结构体字段 Content，类型为 string，附带标签 `gorm:"type:text" json:"content"`

	Content      string    `gorm:"type:text" json:"content"`
	// 结构体字段 Status，类型为 int，附带标签 `json:"status"`

	Status       int       `json:"status"` // 0待发 1成功 2失败 3重试中
	// 结构体字段 RetryCount，类型为 int，附带标签 `json:"retry_count"`

	RetryCount   int       `json:"retry_count"`
	// 结构体字段 DedupeKey，类型为 string，附带标签 `gorm:"size:128;index" json:"dedupe_key"`

	DedupeKey    string    `gorm:"size:128;index" json:"dedupe_key"`
	// 结构体字段 CreatedAt，类型为 time.Time，附带标签 `json:"created_at"`

	CreatedAt    time.Time `json:"created_at"`
// 当前代码块结束

}

// ============ 工单08 性能分析（无需持久化表，指标通过 expvar/Prometheus 暴露）============
