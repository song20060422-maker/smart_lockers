// @ticket 04-物品寄存：收费计算单元测试
package service

import (
	"testing"
	"time"

	"smart_lockers/internal/model"
)

// 定义函数 rule：封装一段可复用的处理逻辑（参数：）

func rule() model.PricingRule {
	// 返回 model.PricingRule{ 给调用方，结束当前函数

	return model.PricingRule{
		// 执行 Size: "small", FirstHours: 1, FirstFee: 300, ExtraHourFee: 200,，完成对应的业务处理

		Size: "small", FirstHours: 1, FirstFee: 300, ExtraHourFee: 200,
		// 执行 DailyCap: 0, Deposit: 100, Enabled: true, Priority: 1,，完成对应的业务处理

		DailyCap: 0, Deposit: 100, Enabled: true, Priority: 1,
	// 当前代码块结束

	}
// 当前代码块结束

}

// 首段时长内收首费 + 押金 @ticket 04
// 定义函数 TestCalcFeeWithinFirstHours：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestCalcFeeWithinFirstHours(t *testing.T) {
	// 定义并初始化局部变量 start，取值为 time.Unix(0, 0)

	start := time.Unix(0, 0)
	// 定义并初始化局部变量 got，取值为 CalcFee(CalcFeeInput{Rule: rule(), StartTime: start, EndTime: start.Add(30 * time.Minute)})

	got := CalcFee(CalcFeeInput{Rule: rule(), StartTime: start, EndTime: start.Add(30 * time.Minute)})
	// 条件判断：got.BaseFee != 300 || got.Deposit != 100 || got.Payable != 400 时执行下面逻辑

	if got.BaseFee != 300 || got.Deposit != 100 || got.Payable != 400 {
		// 执行 t.Fatalf("unexpected fee: %+v", got)，完成对应的业务处理

		t.Fatalf("unexpected fee: %+v", got)
	// 当前代码块结束

	}
// 当前代码块结束

}

// 超出首段按小时向上取整累加 @ticket 04
// 定义函数 TestCalcFeeExtraHours：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestCalcFeeExtraHours(t *testing.T) {
	// 定义并初始化局部变量 start，取值为 time.Unix(0, 0)

	start := time.Unix(0, 0)
	// 定义并初始化局部变量 got，取值为 CalcFee(CalcFeeInput{Rule: rule(), StartTime: start, EndTime: start.Add(150 * time.Minute)})

	got := CalcFee(CalcFeeInput{Rule: rule(), StartTime: start, EndTime: start.Add(150 * time.Minute)})
	// 2.5h -> 超出 1.5h -> ceil=2 -> 300 + 400 = 700
	// 条件判断：got.BaseFee != 700 || got.Payable != 800 时执行下面逻辑

	if got.BaseFee != 700 || got.Payable != 800 {
		// 执行 t.Fatalf("unexpected fee: %+v", got)，完成对应的业务处理

		t.Fatalf("unexpected fee: %+v", got)
	// 当前代码块结束

	}
// 当前代码块结束

}

// 优惠抵扣不超过基础费用且不为负 @ticket 04
// 定义函数 TestCalcFeeCouponCap：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestCalcFeeCouponCap(t *testing.T) {
	// 定义并初始化局部变量 start，取值为 time.Unix(0, 0)

	start := time.Unix(0, 0)
	// 定义并初始化局部变量 got，取值为 CalcFee(CalcFeeInput{Rule: rule(), StartTime: start, EndTime: start.Add(30 * time.Minute), Coupon: 500})

	got := CalcFee(CalcFeeInput{Rule: rule(), StartTime: start, EndTime: start.Add(30 * time.Minute), Coupon: 500})
	// 条件判断：got.Discount != 300 || got.Payable != 100 时执行下面逻辑

	if got.Discount != 300 || got.Payable != 100 {
		// 执行 t.Fatalf("unexpected coupon cap: %+v", got)，完成对应的业务处理

		t.Fatalf("unexpected coupon cap: %+v", got)
	// 当前代码块结束

	}
	// 定义并初始化局部变量 neg，取值为 CalcFee(CalcFeeInput{Rule: rule(), StartTime: start, EndTime: start.Add(30 * time.Minute), Coupon: -50})

	neg := CalcFee(CalcFeeInput{Rule: rule(), StartTime: start, EndTime: start.Add(30 * time.Minute), Coupon: -50})
	// 条件判断：neg.Discount != 0 时执行下面逻辑

	if neg.Discount != 0 {
		// 执行 t.Fatalf("negative coupon should clamp to 0: %+v", neg)，完成对应的业务处理

		t.Fatalf("negative coupon should clamp to 0: %+v", neg)
	// 当前代码块结束

	}
// 当前代码块结束

}

// 单日封顶 @ticket 04
// 定义函数 TestCalcFeeDailyCap：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestCalcFeeDailyCap(t *testing.T) {
	// 定义并初始化局部变量 r，取值为 rule()

	r := rule()
	// 为 r.DailyCap 赋值：400

	r.DailyCap = 400
	// 定义并初始化局部变量 start，取值为 time.Unix(0, 0)

	start := time.Unix(0, 0)
	// 定义并初始化局部变量 got，取值为 CalcFee(CalcFeeInput{Rule: r, StartTime: start, EndTime: start.Add(2 * time.Hour)})

	got := CalcFee(CalcFeeInput{Rule: r, StartTime: start, EndTime: start.Add(2 * time.Hour)})
	// 条件判断：got.BaseFee != 400 时执行下面逻辑

	if got.BaseFee != 400 {
		// 执行 t.Fatalf("daily cap not applied: %+v", got)，完成对应的业务处理

		t.Fatalf("daily cap not applied: %+v", got)
	// 当前代码块结束

	}
// 当前代码块结束

}

// 结束时间早于开始只返回押金 @ticket 04
// 定义函数 TestCalcFeeInvalidRange：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestCalcFeeInvalidRange(t *testing.T) {
	// 定义并初始化局部变量 start，取值为 time.Unix(0, 0)

	start := time.Unix(0, 0)
	// 定义并初始化局部变量 got，取值为 CalcFee(CalcFeeInput{Rule: rule(), StartTime: start, EndTime: start.Add(-time.Hour)})

	got := CalcFee(CalcFeeInput{Rule: rule(), StartTime: start, EndTime: start.Add(-time.Hour)})
	// 条件判断：got.BaseFee != 0 || got.Payable != 100 时执行下面逻辑

	if got.BaseFee != 0 || got.Payable != 100 {
		// 执行 t.Fatalf("unexpected invalid range fee: %+v", got)，完成对应的业务处理

		t.Fatalf("unexpected invalid range fee: %+v", got)
	// 当前代码块结束

	}
// 当前代码块结束

}

// 按尺寸匹配最高优先级启用规则 @ticket 04
// 定义函数 TestMatchRule：封装一段可复用的处理逻辑（参数：t *testing.T）

func TestMatchRule(t *testing.T) {
	// 定义并初始化局部变量 rules，取值为 []model.PricingRule{

	rules := []model.PricingRule{
		// 执行 {Size: "small", Priority: 1, Enabled: true},，完成对应的业务处理

		{Size: "small", Priority: 1, Enabled: true},
		// 执行 {Size: "small", Priority: 5, Enabled: true},，完成对应的业务处理

		{Size: "small", Priority: 5, Enabled: true},
		// 执行 {Size: "small", Priority: 9, Enabled: false},，完成对应的业务处理

		{Size: "small", Priority: 9, Enabled: false},
		// 执行 {Size: "large", Priority: 10, Enabled: true},，完成对应的业务处理

		{Size: "large", Priority: 10, Enabled: true},
	// 当前代码块结束

	}
	// 定义并初始化局部变量 got, ok，取值为 MatchRule(rules, "small")

	got, ok := MatchRule(rules, "small")
	// 条件判断：!ok || got.Priority != 5 时执行下面逻辑

	if !ok || got.Priority != 5 {
		// 执行 t.Fatalf("expected priority 5, got %+v ok=%v", got, ok)，完成对应的业务处理

		t.Fatalf("expected priority 5, got %+v ok=%v", got, ok)
	// 当前代码块结束

	}
	// 条件判断：_, ok := MatchRule(rules, "medium"); ok 时执行下面逻辑

	if _, ok := MatchRule(rules, "medium"); ok {
		// 执行 t.Fatal("medium should not match")，完成对应的业务处理

		t.Fatal("medium should not match")
	// 当前代码块结束

	}
// 当前代码块结束

}
