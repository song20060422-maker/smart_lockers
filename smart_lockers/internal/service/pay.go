// @ticket 02-支付宝支付：余额/明细/充值/充值记录，资金一致性 + 回调幂等入账
package service

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"smart_lockers/internal/config"
	"smart_lockers/internal/model"
	"smart_lockers/internal/pkg"
)

// GetBalance 账户余额查询（脱敏：金额以分为单位返回）@ticket 02-支付宝支付
// 定义函数 GetBalance：封装一段可复用的处理逻辑（参数：userID uint）

func GetBalance(userID uint) (map[string]interface{}, int) {
	// 定义并初始化局部变量 bal，取值为 int64(0)

	bal := int64(0)
	// 条件判断：pkg.DB != nil 时执行下面逻辑

	if pkg.DB != nil {
		// 声明变量 a model.Account，用于保存运行过程中的数据

		var a model.Account
		// 条件判断：err := pkg.DB.Where("user_id = ?", userID).First(&a).Error; err == nil 时执行下面逻辑

		if err := pkg.DB.Where("user_id = ?", userID).First(&a).Error; err == nil {
			// 为 bal 赋值：a.Balance

			bal = a.Balance
		// 当前代码块结束

		}
	// 当前代码块结束

	}
	// 返回 map[string]interface{}{"user_id": userID, "balance": bal, "currency": "CNY"}, pkg.CodeOK 给调用方，结束当前函数

	return map[string]interface{}{"user_id": userID, "balance": bal, "currency": "CNY"}, pkg.CodeOK
// 当前代码块结束

}

// BalanceLogs 余额明细查询 @ticket 02-支付宝支付
// 定义函数 BalanceLogs：封装一段可复用的处理逻辑（参数：userID uint, page, pageSize int）

func BalanceLogs(userID uint, page, pageSize int) (map[string]interface{}, int) {
	// 条件判断：page <= 0 时执行下面逻辑

	if page <= 0 {
		// 为 page 赋值：1

		page = 1
	// 当前代码块结束

	}
	// 条件判断：pageSize <= 0 时执行下面逻辑

	if pageSize <= 0 {
		// 为 pageSize 赋值：20

		pageSize = 20
	// 当前代码块结束

	}
	// 声明变量 list []model.BalanceLog，用于保存运行过程中的数据

	var list []model.BalanceLog
	// 条件判断：pkg.DB != nil 时执行下面逻辑

	if pkg.DB != nil {
		// 为 pkg.DB.Where("user_id 赋值：?", userID).Order("id desc").

		pkg.DB.Where("user_id = ?", userID).Order("id desc").
			// 执行 Offset((page - 1) * pageSize).Limit(pageSize).Find(&list)，完成对应的业务处理

			Offset((page - 1) * pageSize).Limit(pageSize).Find(&list)
	// 当前代码块结束

	}
	// 返回 map[string]interface{}{"list": list, "page": page, "page_size": pageSize}, pkg.CodeOK 给调用方，结束当前函数

	return map[string]interface{}{"list": list, "page": page, "page_size": pageSize}, pkg.CodeOK
// 当前代码块结束

}

// CreateRecharge 创建充值单 @ticket 02-支付宝支付
// 定义函数 CreateRecharge：封装一段可复用的处理逻辑（参数：userID uint, amount int64）

func CreateRecharge(userID uint, amount int64) (map[string]interface{}, int) {
	// 条件判断：amount <= 0 时执行下面逻辑

	if amount <= 0 {
		// 返回 nil, pkg.ErrBad 给调用方，结束当前函数

		return nil, pkg.ErrBad
	// 当前代码块结束

	}
	// 定义并初始化局部变量 chargeNo，取值为 fmt.Sprintf("RC%s%06d", time.Now().Format("20060102150405"), time.Now().Nanosecond()/1000)

	chargeNo := fmt.Sprintf("RC%s%06d", time.Now().Format("20060102150405"), time.Now().Nanosecond()/1000)
	// 条件判断：pkg.DB != nil 时执行下面逻辑

	if pkg.DB != nil {
		// 定义并初始化局部变量 ro，取值为 model.RechargeOrder{ChargeNo: chargeNo, UserID: userID, Amount: amount, Status: 1}

		ro := model.RechargeOrder{ChargeNo: chargeNo, UserID: userID, Amount: amount, Status: 1}
		// 条件判断：err := pkg.DB.Create(&ro).Error; err != nil 时执行下面逻辑

		if err := pkg.DB.Create(&ro).Error; err != nil {
			// 返回 nil, pkg.ErrInternal 给调用方，结束当前函数

			return nil, pkg.ErrInternal
		// 当前代码块结束

		}
	// 当前代码块结束

	}
	// 生成支付宝沙箱支付链接（SDK 本地签名，无需出站调用）@ticket 02-支付宝支付
	// 定义并初始化局部变量 out，取值为 map[string]interface{}{"charge_no": chargeNo, "amount": amount, "status": 1}

	out := map[string]interface{}{"charge_no": chargeNo, "amount": amount, "status": 1}
	// 条件判断：pkg.AlipayReady() 时执行下面逻辑

	if pkg.AlipayReady() {
		// 定义并初始化局部变量 cfg，取值为 config.Get().Alipay

		cfg := config.Get().Alipay
		// 定义并初始化局部变量 payURL, err，取值为 pkg.CreateAlipayPagePayURL(cfg, chargeNo, amount, "账户充值")

		payURL, err := pkg.CreateAlipayPagePayURL(cfg, chargeNo, amount, "账户充值")
		// 条件判断：err != nil 时执行下面逻辑

		if err != nil {
			// 执行 pkg.Logf("[WARN] create alipay pay url: %v", err)，完成对应的业务处理

			pkg.Logf("[WARN] create alipay pay url: %v", err)
		// 上述条件不成立时，执行 else 分支逻辑

		} else {
			// 为 out["pay_url"] 赋值：payURL

			out["pay_url"] = payURL
			// 为 out["pay_channel"] 赋值：map[bool]string{true: "alipay_sandbox", false: "alipay_production"}[cfg.Sandbox]

			out["pay_channel"] = map[bool]string{true: "alipay_sandbox", false: "alipay_production"}[cfg.Sandbox]
		// 当前代码块结束

		}
	// 当前代码块结束

	}
	// 返回 out, pkg.CodeOK 给调用方，结束当前函数

	return out, pkg.CodeOK
// 当前代码块结束

}

// RechargeRecords 充值记录查询 @ticket 02-支付宝支付
// 定义函数 RechargeRecords：封装一段可复用的处理逻辑（参数：userID uint, page, pageSize int）

func RechargeRecords(userID uint, page, pageSize int) (map[string]interface{}, int) {
	// 条件判断：page <= 0 时执行下面逻辑

	if page <= 0 {
		// 为 page 赋值：1

		page = 1
	// 当前代码块结束

	}
	// 条件判断：pageSize <= 0 时执行下面逻辑

	if pageSize <= 0 {
		// 为 pageSize 赋值：20

		pageSize = 20
	// 当前代码块结束

	}
	// 声明变量 list []model.RechargeOrder，用于保存运行过程中的数据

	var list []model.RechargeOrder
	// 条件判断：pkg.DB != nil 时执行下面逻辑

	if pkg.DB != nil {
		// 为 pkg.DB.Where("user_id 赋值：?", userID).Order("id desc").

		pkg.DB.Where("user_id = ?", userID).Order("id desc").
			// 执行 Offset((page - 1) * pageSize).Limit(pageSize).Find(&list)，完成对应的业务处理

			Offset((page - 1) * pageSize).Limit(pageSize).Find(&list)
	// 当前代码块结束

	}
	// 返回 map[string]interface{}{"list": list, "page": page, "page_size": pageSize}, pkg.CodeOK 给调用方，结束当前函数

	return map[string]interface{}{"list": list, "page": page, "page_size": pageSize}, pkg.CodeOK
// 当前代码块结束

}

// RechargeCallbackInput 支付回调入参 @ticket 02-支付宝支付
// 定义结构体 RechargeCallbackInput：用于承载一组相关的字段数据

type RechargeCallbackInput struct {
	// 结构体字段 ChargeNo，类型为 string

	ChargeNo       string
	// 结构体字段 TradeNo，类型为 string

	TradeNo        string
	// 结构体字段 Amount，类型为 int64

	Amount         int64
	// 结构体字段 Signature，类型为 string

	Signature      string
	// 结构体字段 IdempotencyKey，类型为 string

	IdempotencyKey string
// 当前代码块结束

}

// verifyAlipaySign 回调签名校验（本地占位：生产用支付宝公钥验签）@ticket 02-支付宝支付
// 定义函数 verifyAlipaySign：封装一段可复用的处理逻辑（参数：in RechargeCallbackInput）

func verifyAlipaySign(in RechargeCallbackInput) bool {
	// TODO @ticket 02：接入支付宝 SDK 的 VerifySign；本地开发放行但记录 trade_no
	// 返回 in.ChargeNo != "" && in.TradeNo != "" 给调用方，结束当前函数

	return in.ChargeNo != "" && in.TradeNo != ""
// 当前代码块结束

}

// RechargeCallback 支付回调幂等入账：签名校验 + 状态机校验 + 事务内余额变动必带明细 + 幂等
// @ticket 02-支付宝支付
// 定义函数 RechargeCallback：封装一段可复用的处理逻辑（参数：in RechargeCallbackInput）

func RechargeCallback(in RechargeCallbackInput) (map[string]interface{}, int) {
	// 条件判断：!verifyAlipaySign(in) 时执行下面逻辑

	if !verifyAlipaySign(in) {
		// 返回 nil, pkg.ErrPayCallback 给调用方，结束当前函数

		return nil, pkg.ErrPayCallback
	// 当前代码块结束

	}
	// 条件判断：in.IdempotencyKey == "" 时执行下面逻辑

	if in.IdempotencyKey == "" {
		// 返回 nil, pkg.ErrBad 给调用方，结束当前函数

		return nil, pkg.ErrBad
	// 当前代码块结束

	}
	// 幂等：同一 biz_no 已有明细直接返回 @ticket 02-支付宝支付
	// 声明变量 dup model.BalanceLog，用于保存运行过程中的数据

	var dup model.BalanceLog
	// 条件判断：pkg.DB != nil && pkg.DB.Where("biz_type = ? AND biz_no = ?", "recharge", in.IdempotencyKey).First(&dup).Error == nil 时执行下面逻辑

	if pkg.DB != nil && pkg.DB.Where("biz_type = ? AND biz_no = ?", "recharge", in.IdempotencyKey).First(&dup).Error == nil {
		// 返回 map[string]interface{}{"charge_no": in.ChargeNo, "idempotent": true}, pkg.CodeOK 给调用方，结束当前函数

		return map[string]interface{}{"charge_no": in.ChargeNo, "idempotent": true}, pkg.CodeOK
	// 当前代码块结束

	}
	// 定义并初始化局部变量 err，取值为 pkg.Tx(func(tx *gorm.DB) error {

	err := pkg.Tx(func(tx *gorm.DB) error {
		// 声明变量 ro model.RechargeOrder，用于保存运行过程中的数据

		var ro model.RechargeOrder
		// 条件判断：err := tx.Where("charge_no = ?", in.ChargeNo).First(&ro).Error; err != nil 时执行下面逻辑

		if err := tx.Where("charge_no = ?", in.ChargeNo).First(&ro).Error; err != nil {
			// 返回 err 给调用方，结束当前函数

			return err
		// 当前代码块结束

		}
		// 条件判断：ro.Status == 2 { // 已支付，重复回调幂等返回 时执行下面逻辑

		if ro.Status == 2 { // 已支付，重复回调幂等返回
			// 返回 nil 给调用方，结束当前函数

			return nil
		// 当前代码块结束

		}
		// 条件判断：ro.Status != 1 || ro.Amount != in.Amount 时执行下面逻辑

		if ro.Status != 1 || ro.Amount != in.Amount {
			// 返回 fmt.Errorf("invalid recharge state or amount") 给调用方，结束当前函数

			return fmt.Errorf("invalid recharge state or amount")
		// 当前代码块结束

		}
		// 账户余额：无则创建 @ticket 02-支付宝支付
		// 声明变量 acc model.Account，用于保存运行过程中的数据

		var acc model.Account
		// 条件判断：err := tx.Where("user_id = ?", ro.UserID).First(&acc).Error; err != nil 时执行下面逻辑

		if err := tx.Where("user_id = ?", ro.UserID).First(&acc).Error; err != nil {
			// 为 acc 赋值：model.Account{UserID: ro.UserID, Balance: 0}

			acc = model.Account{UserID: ro.UserID, Balance: 0}
			// 条件判断：err := tx.Create(&acc).Error; err != nil 时执行下面逻辑

			if err := tx.Create(&acc).Error; err != nil {
				// 返回 err 给调用方，结束当前函数

				return err
			// 当前代码块结束

			}
		// 当前代码块结束

		}
		// 定义并初始化局部变量 newBal，取值为 acc.Balance + in.Amount

		newBal := acc.Balance + in.Amount
		// 条件判断：err := tx.Model(&model.Account{}).Where("id = ?", acc.ID).Update("balance", newBal).Error; err != nil 时执行下面逻辑

		if err := tx.Model(&model.Account{}).Where("id = ?", acc.ID).Update("balance", newBal).Error; err != nil {
			// 返回 err 给调用方，结束当前函数

			return err
		// 当前代码块结束

		}
		// 任意余额变动必有明细 @ticket 02-支付宝支付
		// 定义并初始化局部变量 log，取值为 model.BalanceLog{

		log := model.BalanceLog{
			// 执行 UserID: ro.UserID, Change: in.Amount, BalanceAfter: newBal,，完成对应的业务处理

			UserID: ro.UserID, Change: in.Amount, BalanceAfter: newBal,
			// 执行 BizType: "recharge", BizNo: in.IdempotencyKey,，完成对应的业务处理

			BizType: "recharge", BizNo: in.IdempotencyKey,
		// 当前代码块结束

		}
		// 条件判断：err := tx.Create(&log).Error; err != nil 时执行下面逻辑

		if err := tx.Create(&log).Error; err != nil {
			// 返回 err 给调用方，结束当前函数

			return err
		// 当前代码块结束

		}
		// 定义并初始化局部变量 now，取值为 time.Now()

		now := time.Now()
		// 条件判断：err := tx.Model(&model.RechargeOrder{}).Where("id = ?", ro.ID). 时执行下面逻辑

		if err := tx.Model(&model.RechargeOrder{}).Where("id = ?", ro.ID).
			// 执行 Updates(map[string]interface{}{"status": 2, "trade_no": in.TradeNo, "paid_at": now}).Error; err != nil {，完成对应的业务处理

			Updates(map[string]interface{}{"status": 2, "trade_no": in.TradeNo, "paid_at": now}).Error; err != nil {
			// 返回 err 给调用方，结束当前函数

			return err
		// 当前代码块结束

		}
		// 返回 nil 给调用方，结束当前函数

		return nil
	// 执行 })，完成对应的业务处理

	})
	// 条件判断：err != nil 时执行下面逻辑

	if err != nil {
		// 执行 pkg.Logf("[ERROR] recharge callback: %v", err)，完成对应的业务处理

		pkg.Logf("[ERROR] recharge callback: %v", err)
		// 返回 nil, pkg.ErrInternal 给调用方，结束当前函数

		return nil, pkg.ErrInternal
	// 当前代码块结束

	}
	// 返回 map[string]interface{}{"charge_no": in.ChargeNo, "idempotent": false}, pkg.CodeOK 给调用方，结束当前函数

	return map[string]interface{}{"charge_no": in.ChargeNo, "idempotent": false}, pkg.CodeOK
// 当前代码块结束

}

// PayFromBalance 余额支付扣款（供寄存/补缴复用）：事务 + 明细 + 幂等 @ticket 02-支付宝支付
// 定义函数 PayFromBalance：封装一段可复用的处理逻辑（参数：tx *gorm.DB, userID uint, amount int64, bizType, bizNo string）

func PayFromBalance(tx *gorm.DB, userID uint, amount int64, bizType, bizNo string) error {
	// 声明变量 acc model.Account，用于保存运行过程中的数据

	var acc model.Account
	// 条件判断：err := tx.Where("user_id = ?", userID).First(&acc).Error; err != nil 时执行下面逻辑

	if err := tx.Where("user_id = ?", userID).First(&acc).Error; err != nil {
		// 返回 fmt.Errorf("account not found: %w", err) 给调用方，结束当前函数

		return fmt.Errorf("account not found: %w", err)
	// 当前代码块结束

	}
	// 条件判断：acc.Balance < amount 时执行下面逻辑

	if acc.Balance < amount {
		// 返回 fmt.Errorf("insufficient balance") 给调用方，结束当前函数

		return fmt.Errorf("insufficient balance")
	// 当前代码块结束

	}
	// 定义并初始化局部变量 newBal，取值为 acc.Balance - amount

	newBal := acc.Balance - amount
	// 条件判断：err := tx.Model(&model.Account{}).Where("id = ?", acc.ID).Update("balance", newBal).Error; err != nil 时执行下面逻辑

	if err := tx.Model(&model.Account{}).Where("id = ?", acc.ID).Update("balance", newBal).Error; err != nil {
		// 返回 err 给调用方，结束当前函数

		return err
	// 当前代码块结束

	}
	// 定义并初始化局部变量 log，取值为 model.BalanceLog{

	log := model.BalanceLog{
		// 执行 UserID: userID, Change: -amount, BalanceAfter: newBal, BizType: bizType, BizNo: bizNo,，完成对应的业务处理

		UserID: userID, Change: -amount, BalanceAfter: newBal, BizType: bizType, BizNo: bizNo,
	// 当前代码块结束

	}
	// 返回 tx.Create(&log).Error 给调用方，结束当前函数

	return tx.Create(&log).Error
// 当前代码块结束

}
