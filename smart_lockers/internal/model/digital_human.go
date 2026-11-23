// @ticket 07-数字人Coze：数字人独立登录账号表——与 C 端 user_profiles、B 端 admin 完全隔离，账号密码入库查库校验
// 定义包 model：全局数据表实体定义
package model

import "time"

// DigitalHumanAccount 数字人独立登录账号 @ticket 07-数字人Coze
// 定义结构体 DigitalHumanAccount：承载数字人登录账号的数据库记录，独立于用户资料表与运营后台账号
type DigitalHumanAccount struct {
	// 结构体字段 ID，主键自增
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`
	// 结构体字段 Username，登录账号，唯一索引防重复
	Username string `gorm:"size:64;not null;uniqueIndex:uk_dh_username" json:"username"`
	// 结构体字段 PasswordHash，bcrypt 密码哈希，json:"-" 禁止对外输出
	PasswordHash string `gorm:"size:128;not null" json:"-"`
	// 结构体字段 Role，登录后直接进入的业务角色 sender/receiver/service（账号即角色，无需选择页）@ticket 07
	Role string `gorm:"size:32;default:operator" json:"role"`
	// 结构体字段 Status，1 启用 0 禁用
	Status int `gorm:"default:1" json:"status"`
	// 结构体字段 TokenVer，令牌版本号，退出登录自增使旧 token 失效
	TokenVer int `gorm:"default:0" json:"token_ver"`
	// 结构体字段 CreatedAt，创建时间
	CreatedAt time.Time `json:"-"`
}

// TableName 指定独立表名，避免与 user_profiles/admin 混淆
// 定义函数 TableName：返回真实数据表名
func (DigitalHumanAccount) TableName() string {
	// 返回表名 digital_human_accounts
	return "digital_human_accounts"
}
