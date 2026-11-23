// @ticket 07-数字人Coze：数字人账号密码真实查库校验——bcrypt 比对 + 签发 JWT，与用户端/运营端账号体系完全隔离
// 定义包 service：为 handler 提供可复用的业务能力
package service

import (
	"golang.org/x/crypto/bcrypt"

	"smart_lockers/internal/middleware"
	"smart_lockers/internal/model"
	"smart_lockers/internal/pkg"
)

// DigitalHumanLogin 数字人账号密码登录：查库校验 + bcrypt 比对 + 签发 token @ticket 07-数字人Coze
// 定义函数 DigitalHumanLogin：返回登录数据与业务码；账号不存在/禁用/密码错误统一返回 ErrDHLogin，避免账号枚举
func DigitalHumanLogin(username, password string) (map[string]interface{}, int) {
	// 条件判断：数据库未就绪时返回内部错误，不降级为前端演示校验
	if pkg.DB == nil {
		// 返回 nil, pkg.ErrInternal 给调用方
		return nil, pkg.ErrInternal
	}
	// 定义并初始化局部变量 acc，承载查库结果
	var acc model.DigitalHumanAccount
	// 按账号精确查询，未命中或查询异常均按登录失败处理
	if err := pkg.DB.Where("username = ?", username).First(&acc).Error; err != nil {
		// 返回 nil, pkg.ErrDHLogin 给调用方
		return nil, pkg.ErrDHLogin
	}
	// 条件判断：账号被禁用时按登录失败处理
	if acc.Status != 1 {
		// 返回 nil, pkg.ErrDHLogin 给调用方
		return nil, pkg.ErrDHLogin
	}
	// bcrypt 比对数据库哈希与用户输入的明文密码
	if bcrypt.CompareHashAndPassword([]byte(acc.PasswordHash), []byte(password)) != nil {
		// 返回 nil, pkg.ErrDHLogin 给调用方
		return nil, pkg.ErrDHLogin
	}
	// 复用现有 JWT 中间件签发令牌，uid 使用独立账号表主键，tokenVer 取账号版本号
	token, err := middleware.GenerateToken(acc.ID, acc.TokenVer)
	// 条件判断：签发失败返回内部错误
	if err != nil {
		// 返回 nil, pkg.ErrInternal 给调用方
		return nil, pkg.ErrInternal
	}
	// 返回固定字段契约：token/username/role，前端仅依赖这些字段
	return map[string]interface{}{
		// 令牌字段
		"token": token,
		// 账号字段
		"username": acc.Username,
		// 内部角色标记字段
		"role": acc.Role,
	}, pkg.CodeOK
}

// SeedDigitalHumanAccount 幂等写入数字人种子账号：分角色补种三业务账号 + 默认运维账号 @ticket 07
// 定义函数 SeedDigitalHumanAccount：保证不同账号登录进入不同角色，密码统一为 bcrypt("123456")
func SeedDigitalHumanAccount() {
	// 条件判断：数据库未就绪直接返回
	if pkg.DB == nil {
		// 结束当前函数
		return
	}
	// bcrypt("123456") 固定哈希，三角色与运维账号统一明文密码便于演示
	const pwdHash = "$2b$10$MiQj0bjdmXuT6XFuyVjPYuGC8QpFffCiixzlEKv7O4Y7kIMIaeL6S"
	// 定义种子账号列表：账号名即登录名，role 决定登录后进入的业务角色 @ticket 07
	seeds := []model.DigitalHumanAccount{
		// 默认运维账号（旧数据兼容保留）
		{Username: "admin", PasswordHash: pwdHash, Role: "operator", Status: 1, TokenVer: 0},
		// 寄件引导员账号
		{Username: "sender", PasswordHash: pwdHash, Role: "sender", Status: 1, TokenVer: 0},
		// 取件指引员账号
		{Username: "receiver", PasswordHash: pwdHash, Role: "receiver", Status: 1, TokenVer: 0},
		// 智能客服专员账号
		{Username: "service", PasswordHash: pwdHash, Role: "service", Status: 1, TokenVer: 0},
	}
	// 逐账号按用户名幂等补种：已存在则跳过，保证老库升级后也能补齐三角色
	for _, seed := range seeds {
		// 定义并初始化局部变量 count，统计同名账号数量
		var count int64
		// 按账号名查询，失败则跳过该账号
		if err := pkg.DB.Model(&model.DigitalHumanAccount{}).Where("username = ?", seed.Username).Count(&count).Error; err != nil {
			// 进入下一个账号
			continue
		}
		// 条件判断：同名账号已存在时跳过，避免覆盖已有密码或角色
		if count > 0 {
			// 进入下一个账号
			continue
		}
		// 不存在则写入种子账号
		_ = pkg.DB.Create(&seed).Error
	}
}
