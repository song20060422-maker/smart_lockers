// @ticket 04-物品寄存：分时/多维收费计算（纯函数，便于单测）
package service

import (
	"math"
	"time"

	"smart_lockers/internal/model"
)

// CalcFeeInput 计费入参 @ticket 04-物品寄存
// 定义结构体 CalcFeeInput：用于承载一组相关的字段数据

type CalcFeeInput struct {
	// 结构体字段 Rule，类型为 model.PricingRule

	Rule      model.PricingRule
	// 结构体字段 StartTime，类型为 time.Time

	StartTime time.Time
	// 结构体字段 EndTime，类型为 time.Time

	EndTime   time.Time
	// 结构体字段 Coupon，类型为 int64

	Coupon    int64 // 优惠券抵扣（分）
// 当前代码块结束

}

// FeeResult 计费结果 @ticket 04-物品寄存
// 定义结构体 FeeResult：用于承载一组相关的字段数据

type FeeResult struct {
	// 结构体字段 BaseFee，类型为 int64，附带标签 `json:"base_fee"`

	BaseFee  int64 `json:"base_fee"` // 时长费用
	// 结构体字段 Discount，类型为 int64，附带标签 `json:"discount"`

	Discount int64 `json:"discount"` // 实际抵扣
	// 结构体字段 Deposit，类型为 int64，附带标签 `json:"deposit"`

	Deposit  int64 `json:"deposit"`  // 押金
	// 结构体字段 Payable，类型为 int64，附带标签 `json:"payable"`

	Payable  int64 `json:"payable"`  // 应付（含押金）
// 当前代码块结束

}

// CalcFee 按规则计算费用：首段时长内收 first_fee，超出按小时累加，单日封顶，叠加优惠(不超基础费用)与押金
// @ticket 04-物品寄存 P95≤250ms
// 定义函数 CalcFee：封装一段可复用的处理逻辑（参数：in CalcFeeInput）

func CalcFee(in CalcFeeInput) FeeResult {
	// 定义并初始化局部变量 r，取值为 in.Rule

	r := in.Rule
	// 条件判断：in.EndTime.Before(in.StartTime) || in.EndTime.Equal(in.StartTime) 时执行下面逻辑

	if in.EndTime.Before(in.StartTime) || in.EndTime.Equal(in.StartTime) {
		// 返回 FeeResult{Deposit: r.Deposit, Payable: r.Deposit} 给调用方，结束当前函数

		return FeeResult{Deposit: r.Deposit, Payable: r.Deposit}
	// 当前代码块结束

	}
	// 定义并初始化局部变量 hours，取值为 in.EndTime.Sub(in.StartTime).Hours()

	hours := in.EndTime.Sub(in.StartTime).Hours()
	// 声明变量 base int64，用于保存运行过程中的数据

	var base int64
	// 条件判断：hours <= float64(r.FirstHours) 时执行下面逻辑

	if hours <= float64(r.FirstHours) {
		// 为 base 赋值：r.FirstFee

		base = r.FirstFee
	// 上述条件不成立时，执行 else 分支逻辑

	} else {
		// 为 base 赋值：r.FirstFee

		base = r.FirstFee
		// 定义并初始化局部变量 extra，取值为 hours - float64(r.FirstHours)

		extra := hours - float64(r.FirstHours)
		// 定义并初始化局部变量 extraHours，取值为 int64(math.Ceil(extra)) // 不足一小时按一小时

		extraHours := int64(math.Ceil(extra)) // 不足一小时按一小时
		// 执行 base += extraHours * r.ExtraHourFee，完成对应的业务处理

		base += extraHours * r.ExtraHourFee
	// 当前代码块结束

	}
	// 条件判断：r.DailyCap > 0 时执行下面逻辑

	if r.DailyCap > 0 {
		// 定义并初始化局部变量 days，取值为 int64(math.Ceil(hours / 24))

		days := int64(math.Ceil(hours / 24))
		// 定义并初始化局部变量 cap，取值为 r.DailyCap * days

		cap := r.DailyCap * days
		// 条件判断：base > cap 时执行下面逻辑

		if base > cap {
			// 为 base 赋值：cap

			base = cap
		// 当前代码块结束

		}
	// 当前代码块结束

	}
	// 优惠上限：不超过基础费用，且不为负 @ticket 04-物品寄存
	// 定义并初始化局部变量 discount，取值为 in.Coupon

	discount := in.Coupon
	// 条件判断：discount < 0 时执行下面逻辑

	if discount < 0 {
		// 为 discount 赋值：0

		discount = 0
	// 当前代码块结束

	}
	// 条件判断：discount > base 时执行下面逻辑

	if discount > base {
		// 为 discount 赋值：base

		discount = base
	// 当前代码块结束

	}
	// 定义并初始化局部变量 payable，取值为 base - discount + r.Deposit

	payable := base - discount + r.Deposit
	// 返回 FeeResult{BaseFee: base, Discount: discount, Deposit: r.Deposit, Payable: payable} 给调用方，结束当前函数

	return FeeResult{BaseFee: base, Discount: discount, Deposit: r.Deposit, Payable: payable}
// 当前代码块结束

}

// MatchRule 按尺寸匹配启用且优先级最高的收费规则 @ticket 04-物品寄存
// 定义函数 MatchRule：封装一段可复用的处理逻辑（参数：rules []model.PricingRule, size string）

func MatchRule(rules []model.PricingRule, size string) (model.PricingRule, bool) {
	// 声明变量 best model.PricingRule，用于保存运行过程中的数据

	var best model.PricingRule
	// 定义并初始化局部变量 found，取值为 false

	found := false
	// 循环处理：for _, r := range rules，逐轮迭代

	for _, r := range rules {
		// 条件判断：r.Enabled && r.Size == size 时执行下面逻辑

		if r.Enabled && r.Size == size {
			// 条件判断：!found || r.Priority > best.Priority 时执行下面逻辑

			if !found || r.Priority > best.Priority {
				// 为 best 赋值：r

				best = r
				// 为 found 赋值：true

				found = true
			// 当前代码块结束

			}
		// 当前代码块结束

		}
	// 当前代码块结束

	}
	// 返回 best, found 给调用方，结束当前函数

	return best, found
// 当前代码块结束

}
