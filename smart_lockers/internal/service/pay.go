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
func GetBalance(userID uint) (map[string]interface{}, int) {
	bal := int64(0)
	if pkg.DB != nil {
		var a model.Account
		if err := pkg.DB.Where("user_id = ?", userID).First(&a).Error; err == nil {
			bal = a.Balance
		}
	}
	return map[string]interface{}{"user_id": userID, "balance": bal, "currency": "CNY"}, pkg.CodeOK
}

// BalanceLogs 余额明细查询 @ticket 02-支付宝支付
func BalanceLogs(userID uint, page, pageSize int) (map[string]interface{}, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	var list []model.BalanceLog
	if pkg.DB != nil {
		pkg.DB.Where("user_id = ?", userID).Order("id desc").
			Offset((page - 1) * pageSize).Limit(pageSize).Find(&list)
	}
	return map[string]interface{}{"list": list, "page": page, "page_size": pageSize}, pkg.CodeOK
}

// CreateRecharge 创建充值单 @ticket 02-支付宝支付
func CreateRecharge(userID uint, amount int64) (map[string]interface{}, int) {
	if amount <= 0 {
		return nil, pkg.ErrBad
	}
	chargeNo := fmt.Sprintf("RC%s%06d", time.Now().Format("20060102150405"), time.Now().Nanosecond()/1000)
	if pkg.DB != nil {
		ro := model.RechargeOrder{ChargeNo: chargeNo, UserID: userID, Amount: amount, Status: 1}
		if err := pkg.DB.Create(&ro).Error; err != nil {
			return nil, pkg.ErrInternal
		}
	}
	// 生成支付宝沙箱支付链接（SDK 本地签名，无需出站调用）@ticket 02-支付宝支付
	out := map[string]interface{}{"charge_no": chargeNo, "amount": amount, "status": 1}
	if pkg.AlipayReady() {
		cfg := config.Get().Alipay
		payURL, err := pkg.CreateAlipayPagePayURL(cfg, chargeNo, amount, "账户充值")
		if err != nil {
			pkg.Logf("[WARN] create alipay pay url: %v", err)
		} else {
			out["pay_url"] = payURL
			out["pay_channel"] = map[bool]string{true: "alipay_sandbox", false: "alipay_production"}[cfg.Sandbox]
		}
	}
	return out, pkg.CodeOK
}

// RechargeRecords 充值记录查询 @ticket 02-支付宝支付
func RechargeRecords(userID uint, page, pageSize int) (map[string]interface{}, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	var list []model.RechargeOrder
	if pkg.DB != nil {
		pkg.DB.Where("user_id = ?", userID).Order("id desc").
			Offset((page - 1) * pageSize).Limit(pageSize).Find(&list)
	}
	return map[string]interface{}{"list": list, "page": page, "page_size": pageSize}, pkg.CodeOK
}

// RechargeCallbackInput 支付回调入参 @ticket 02-支付宝支付
type RechargeCallbackInput struct {
	ChargeNo       string
	TradeNo        string
	Amount         int64
	Signature      string
	IdempotencyKey string
}

// verifyAlipaySign 回调签名校验（本地占位：生产用支付宝公钥验签）@ticket 02-支付宝支付
func verifyAlipaySign(in RechargeCallbackInput) bool {
	// TODO @ticket 02：接入支付宝 SDK 的 VerifySign；本地开发放行但记录 trade_no
	return in.ChargeNo != "" && in.TradeNo != ""
}

// RechargeCallback 支付回调幂等入账：签名校验 + 状态机校验 + 事务内余额变动必带明细 + 幂等
// @ticket 02-支付宝支付
func RechargeCallback(in RechargeCallbackInput) (map[string]interface{}, int) {
	if !verifyAlipaySign(in) {
		return nil, pkg.ErrPayCallback
	}
	if in.IdempotencyKey == "" {
		return nil, pkg.ErrBad
	}
	// 幂等：同一 biz_no 已有明细直接返回 @ticket 02-支付宝支付
	var dup model.BalanceLog
	if pkg.DB != nil && pkg.DB.Where("biz_type = ? AND biz_no = ?", "recharge", in.IdempotencyKey).First(&dup).Error == nil {
		return map[string]interface{}{"charge_no": in.ChargeNo, "idempotent": true}, pkg.CodeOK
	}
	err := pkg.Tx(func(tx *gorm.DB) error {
		var ro model.RechargeOrder
		if err := tx.Where("charge_no = ?", in.ChargeNo).First(&ro).Error; err != nil {
			return err
		}
		if ro.Status == 2 { // 已支付，重复回调幂等返回
			return nil
		}
		if ro.Status != 1 || ro.Amount != in.Amount {
			return fmt.Errorf("invalid recharge state or amount")
		}
		// 账户余额：无则创建 @ticket 02-支付宝支付
		var acc model.Account
		if err := tx.Where("user_id = ?", ro.UserID).First(&acc).Error; err != nil {
			acc = model.Account{UserID: ro.UserID, Balance: 0}
			if err := tx.Create(&acc).Error; err != nil {
				return err
			}
		}
		newBal := acc.Balance + in.Amount
		if err := tx.Model(&model.Account{}).Where("id = ?", acc.ID).Update("balance", newBal).Error; err != nil {
			return err
		}
		// 任意余额变动必有明细 @ticket 02-支付宝支付
		log := model.BalanceLog{
			UserID: ro.UserID, Change: in.Amount, BalanceAfter: newBal,
			BizType: "recharge", BizNo: in.IdempotencyKey,
		}
		if err := tx.Create(&log).Error; err != nil {
			return err
		}
		now := time.Now()
		if err := tx.Model(&model.RechargeOrder{}).Where("id = ?", ro.ID).
			Updates(map[string]interface{}{"status": 2, "trade_no": in.TradeNo, "paid_at": now}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		pkg.Logf("[ERROR] recharge callback: %v", err)
		return nil, pkg.ErrInternal
	}
	return map[string]interface{}{"charge_no": in.ChargeNo, "idempotent": false}, pkg.CodeOK
}

// PayFromBalance 余额支付扣款（供寄存/补缴复用）：事务 + 明细 + 幂等 @ticket 02-支付宝支付
func PayFromBalance(tx *gorm.DB, userID uint, amount int64, bizType, bizNo string) error {
	var acc model.Account
	if err := tx.Where("user_id = ?", userID).First(&acc).Error; err != nil {
		return fmt.Errorf("account not found: %w", err)
	}
	if acc.Balance < amount {
		return fmt.Errorf("insufficient balance")
	}
	newBal := acc.Balance - amount
	if err := tx.Model(&model.Account{}).Where("id = ?", acc.ID).Update("balance", newBal).Error; err != nil {
		return err
	}
	log := model.BalanceLog{
		UserID: userID, Change: -amount, BalanceAfter: newBal, BizType: bizType, BizNo: bizNo,
	}
	return tx.Create(&log).Error
}
