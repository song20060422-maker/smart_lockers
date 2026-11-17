// @ticket 04-物品寄存：分时/多维收费计算（纯函数，便于单测）
package service

import (
	"math"
	"time"

	"smart_lockers/internal/model"
)

// CalcFeeInput 计费入参 @ticket 04-物品寄存
type CalcFeeInput struct {
	Rule      model.PricingRule
	StartTime time.Time
	EndTime   time.Time
	Coupon    int64 // 优惠券抵扣（分）
}

// FeeResult 计费结果 @ticket 04-物品寄存
type FeeResult struct {
	BaseFee  int64 `json:"base_fee"` // 时长费用
	Discount int64 `json:"discount"` // 实际抵扣
	Deposit  int64 `json:"deposit"`  // 押金
	Payable  int64 `json:"payable"`  // 应付（含押金）
}

// CalcFee 按规则计算费用：首段时长内收 first_fee，超出按小时累加，单日封顶，叠加优惠(不超基础费用)与押金
// @ticket 04-物品寄存 P95≤250ms
func CalcFee(in CalcFeeInput) FeeResult {
	r := in.Rule
	if in.EndTime.Before(in.StartTime) || in.EndTime.Equal(in.StartTime) {
		return FeeResult{Deposit: r.Deposit, Payable: r.Deposit}
	}
	hours := in.EndTime.Sub(in.StartTime).Hours()
	var base int64
	if hours <= float64(r.FirstHours) {
		base = r.FirstFee
	} else {
		base = r.FirstFee
		extra := hours - float64(r.FirstHours)
		extraHours := int64(math.Ceil(extra)) // 不足一小时按一小时
		base += extraHours * r.ExtraHourFee
	}
	if r.DailyCap > 0 {
		days := int64(math.Ceil(hours / 24))
		cap := r.DailyCap * days
		if base > cap {
			base = cap
		}
	}
	// 优惠上限：不超过基础费用，且不为负 @ticket 04-物品寄存
	discount := in.Coupon
	if discount < 0 {
		discount = 0
	}
	if discount > base {
		discount = base
	}
	payable := base - discount + r.Deposit
	return FeeResult{BaseFee: base, Discount: discount, Deposit: r.Deposit, Payable: payable}
}

// MatchRule 按尺寸匹配启用且优先级最高的收费规则 @ticket 04-物品寄存
func MatchRule(rules []model.PricingRule, size string) (model.PricingRule, bool) {
	var best model.PricingRule
	found := false
	for _, r := range rules {
		if r.Enabled && r.Size == size {
			if !found || r.Priority > best.Priority {
				best = r
				found = true
			}
		}
	}
	return best, found
}
