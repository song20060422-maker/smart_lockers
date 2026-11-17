// @ticket 04-物品寄存：收费计算单元测试
package service

import (
	"testing"
	"time"

	"smart_lockers/internal/model"
)

func rule() model.PricingRule {
	return model.PricingRule{
		Size: "small", FirstHours: 1, FirstFee: 300, ExtraHourFee: 200,
		DailyCap: 0, Deposit: 100, Enabled: true, Priority: 1,
	}
}

// 首段时长内收首费 + 押金 @ticket 04
func TestCalcFeeWithinFirstHours(t *testing.T) {
	start := time.Unix(0, 0)
	got := CalcFee(CalcFeeInput{Rule: rule(), StartTime: start, EndTime: start.Add(30 * time.Minute)})
	if got.BaseFee != 300 || got.Deposit != 100 || got.Payable != 400 {
		t.Fatalf("unexpected fee: %+v", got)
	}
}

// 超出首段按小时向上取整累加 @ticket 04
func TestCalcFeeExtraHours(t *testing.T) {
	start := time.Unix(0, 0)
	got := CalcFee(CalcFeeInput{Rule: rule(), StartTime: start, EndTime: start.Add(150 * time.Minute)})
	// 2.5h -> 超出 1.5h -> ceil=2 -> 300 + 400 = 700
	if got.BaseFee != 700 || got.Payable != 800 {
		t.Fatalf("unexpected fee: %+v", got)
	}
}

// 优惠抵扣不超过基础费用且不为负 @ticket 04
func TestCalcFeeCouponCap(t *testing.T) {
	start := time.Unix(0, 0)
	got := CalcFee(CalcFeeInput{Rule: rule(), StartTime: start, EndTime: start.Add(30 * time.Minute), Coupon: 500})
	if got.Discount != 300 || got.Payable != 100 {
		t.Fatalf("unexpected coupon cap: %+v", got)
	}
	neg := CalcFee(CalcFeeInput{Rule: rule(), StartTime: start, EndTime: start.Add(30 * time.Minute), Coupon: -50})
	if neg.Discount != 0 {
		t.Fatalf("negative coupon should clamp to 0: %+v", neg)
	}
}

// 单日封顶 @ticket 04
func TestCalcFeeDailyCap(t *testing.T) {
	r := rule()
	r.DailyCap = 400
	start := time.Unix(0, 0)
	got := CalcFee(CalcFeeInput{Rule: r, StartTime: start, EndTime: start.Add(2 * time.Hour)})
	if got.BaseFee != 400 {
		t.Fatalf("daily cap not applied: %+v", got)
	}
}

// 结束时间早于开始只返回押金 @ticket 04
func TestCalcFeeInvalidRange(t *testing.T) {
	start := time.Unix(0, 0)
	got := CalcFee(CalcFeeInput{Rule: rule(), StartTime: start, EndTime: start.Add(-time.Hour)})
	if got.BaseFee != 0 || got.Payable != 100 {
		t.Fatalf("unexpected invalid range fee: %+v", got)
	}
}

// 按尺寸匹配最高优先级启用规则 @ticket 04
func TestMatchRule(t *testing.T) {
	rules := []model.PricingRule{
		{Size: "small", Priority: 1, Enabled: true},
		{Size: "small", Priority: 5, Enabled: true},
		{Size: "small", Priority: 9, Enabled: false},
		{Size: "large", Priority: 10, Enabled: true},
	}
	got, ok := MatchRule(rules, "small")
	if !ok || got.Priority != 5 {
		t.Fatalf("expected priority 5, got %+v ok=%v", got, ok)
	}
	if _, ok := MatchRule(rules, "medium"); ok {
		t.Fatal("medium should not match")
	}
}
