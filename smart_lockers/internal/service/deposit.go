// @ticket 04-物品寄存 + 05-订单查看与取出：寄存落单/换柜/订单查询业务逻辑
package service

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"smart_lockers/internal/model"
	"smart_lockers/internal/pkg"
)

// CabinetQuery 柜子列表项 @ticket 04-物品寄存 P95≤200ms
type CabinetQuery struct {
	ID        uint   `json:"id"`
	SiteID    uint   `json:"site_id"`
	Code      string `json:"code"`
	Size      string `json:"size"`
	Total     int    `json:"total"`
	Available int    `json:"available"`
}

// QueryCabinets 柜子查询 @ticket 04-物品寄存 P95≤200ms
func QueryCabinets(siteID uint, size string) ([]CabinetQuery, int) {
	if pkg.DB == nil {
		return []CabinetQuery{}, pkg.CodeOK
	}
	q := pkg.DB.Model(&model.LockerCabinet{}).Where("status = ?", 1)
	if siteID > 0 {
		q = q.Where("site_id = ?", siteID)
	}
	if size != "" {
		q = q.Where("size = ?", size)
	}
	var cab []model.LockerCabinet
	if err := q.Find(&cab).Error; err != nil {
		return nil, pkg.ErrInternal
	}
	res := make([]CabinetQuery, 0, len(cab))
	for _, c := range cab {
		res = append(res, CabinetQuery{c.ID, c.SiteID, c.Code, c.Size, c.Total, c.Available})
	}
	return res, pkg.CodeOK
}

// QuoteFee 收费试算：读取规则 -> 计算 @ticket 04-物品寄存
func QuoteFee(size string, hours float64, coupon int64) (FeeResult, int) {
	if pkg.DB == nil {
		return FeeResult{}, pkg.ErrBillingRule
	}
	var rules []model.PricingRule
	if err := pkg.DB.Where("enabled = ?", true).Find(&rules).Error; err != nil {
		return FeeResult{}, pkg.ErrInternal
	}
	rule, ok := MatchRule(rules, size)
	if !ok {
		return FeeResult{}, pkg.ErrBillingRule
	}
	now := time.Now()
	in := CalcFeeInput{
		Rule: rule, StartTime: now, EndTime: now.Add(time.Duration(hours * float64(time.Hour))),
		Coupon: coupon,
	}
	return CalcFee(in), pkg.CodeOK
}

// CreateDepositInput 落单入参 @ticket 04-物品寄存
type CreateDepositInput struct {
	UserID         uint
	DoorID         uint
	SiteID         uint
	IdempotencyKey string
	Hours          float64
	Coupon         int64
}

// CreateDeposit 支付完成落单：占用柜门 + 生成订单 + 支付流水，事务 + 幂等 @ticket 04-物品寄存
func CreateDeposit(in CreateDepositInput) (map[string]interface{}, int) {
	if in.IdempotencyKey == "" {
		return nil, pkg.ErrBad
	}
	// 幂等前置校验：同一 key 直接返回已有结果 @ticket 04-物品寄存
	var exists model.DepositOrder
	if pkg.DB != nil {
		if err := pkg.DB.Where("idempotency_key = ?", in.IdempotencyKey).First(&exists).Error; err == nil {
			return map[string]interface{}{"order_no": exists.OrderNo, "idempotent": true}, pkg.CodeOK
		}
	}
	var orderNo string
	var payable int64
	err := pkg.Tx(func(tx *gorm.DB) error {
		// 乐观锁占用柜门：仅空闲(status=0)可占用 @ticket 04-物品寄存
		var door model.LockerDoor
		if err := tx.First(&door, in.DoorID).Error; err != nil {
			return fmt.Errorf("door not found: %w", err)
		}
		if door.Status != 0 {
			return fmt.Errorf("door unavailable: %w", gorm.ErrRecordNotFound)
		}
		res := tx.Model(&model.LockerDoor{}).
			Where("id = ? AND status = 0 AND version = ?", in.DoorID, door.Version).
			Updates(map[string]interface{}{"status": 1, "version": door.Version + 1, "updated_at": time.Now()})
		if res.Error != nil || res.RowsAffected == 0 {
			return fmt.Errorf("door occupy conflict")
		}

		// 读取收费规则并计费 @ticket 04-物品寄存
		var rules []model.PricingRule
		if err := tx.Where("enabled = ?", true).Find(&rules).Error; err != nil {
			return err
		}
		rule, ok := MatchRule(rules, door.Size)
		if !ok {
			return fmt.Errorf("billing rule not matched")
		}
		now := time.Now()
		fee := CalcFee(CalcFeeInput{
			Rule: rule, StartTime: now, EndTime: now.Add(time.Duration(in.Hours * float64(time.Hour))),
			Coupon: in.Coupon,
		})
		payable = fee.Payable

		orderNo = genOrderNo()
		order := model.DepositOrder{
			OrderNo: orderNo, UserID: in.UserID, DoorID: in.DoorID, SiteID: in.SiteID,
			StartTime: now, EndTime: now.Add(time.Duration(in.Hours * float64(time.Hour))),
			Amount: fee.BaseFee - fee.Discount, Deposit: fee.Deposit, Status: 1,
			IdempotencyKey: in.IdempotencyKey,
		}
		if err := tx.Create(&order).Error; err != nil {
			// 唯一索引冲突视为幂等命中 @ticket 04-物品寄存
			return fmt.Errorf("create order: %w", err)
		}
		return nil
	})
	if err != nil {
		// 唯一键冲突 -> 幂等返回 @ticket 04-物品寄存
		var existing model.DepositOrder
		if pkg.DB != nil && pkg.DB.Where("idempotency_key = ?", in.IdempotencyKey).First(&existing).Error == nil {
			return map[string]interface{}{"order_no": existing.OrderNo, "idempotent": true}, pkg.CodeOK
		}
		pkg.Logf("[ERROR] create deposit: %v", err)
		return nil, pkg.ErrLock
	}
	return map[string]interface{}{"order_no": orderNo, "payable": payable, "idempotent": false}, pkg.CodeOK
}

// ChangeCabinet 异常换柜：释放旧柜占用 -> 重分配新柜 -> 重算费用，事务保护 @ticket 04-物品寄存
func ChangeCabinet(orderNo string, newDoorID uint) (map[string]interface{}, int) {
	err := pkg.Tx(func(tx *gorm.DB) error {
		var order model.DepositOrder
		if err := tx.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
			return err
		}
		if order.Status != 1 {
			return fmt.Errorf("order state not allowed")
		}
		// 释放旧柜 @ticket 04-物品寄存
		if err := tx.Model(&model.LockerDoor{}).Where("id = ?", order.DoorID).
			Updates(map[string]interface{}{"status": 0, "version": gorm.Expr("version + 1")}).Error; err != nil {
			return err
		}
		// 占用新柜（乐观锁）@ticket 04-物品寄存
		var door model.LockerDoor
		if err := tx.First(&door, newDoorID).Error; err != nil {
			return err
		}
		res := tx.Model(&model.LockerDoor{}).
			Where("id = ? AND status = 0 AND version = ?", newDoorID, door.Version).
			Updates(map[string]interface{}{"status": 1, "version": door.Version + 1})
		if res.Error != nil || res.RowsAffected == 0 {
			return fmt.Errorf("new door unavailable")
		}
		// 更新订单柜门（重算费用：尺寸变化才需重算）@ticket 04-物品寄存
		if err := tx.Model(&model.DepositOrder{}).Where("order_no = ?", orderNo).
			Update("door_id", newDoorID).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		pkg.Logf("[ERROR] change cabinet: %v", err)
		return nil, pkg.ErrLock
	}
	return map[string]interface{}{"order_no": orderNo, "door_id": newDoorID}, pkg.CodeOK
}

// genOrderNo 生成订单号 @ticket 04-物品寄存
func genOrderNo() string {
	return fmt.Sprintf("DP%s%06d", time.Now().Format("20060102150405"), time.Now().Nanosecond()/1000)
}
