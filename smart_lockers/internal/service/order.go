// @ticket 05-订单查看与取出：游标分页、详情聚合、临时开柜、超时补交、进行中聚合
package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"gorm.io/gorm"

	"smart_lockers/internal/model"
	"smart_lockers/internal/pkg"
)

// OrderVO 订单视图（脱敏聚合）@ticket 05-订单查看与取出
type OrderVO struct {
	OrderNo   string    `json:"order_no"`
	SiteID    uint      `json:"site_id"`
	DoorID    uint      `json:"door_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Amount    int64     `json:"amount"`
	Deposit   int64     `json:"deposit"`
	Status    int       `json:"status"`
}

// CursorList 订单游标分页（当前/历史）@ticket 05-订单查看与取出
// status: 1=进行中(当前) 2/3=历史；cursor 为上一页最后一条 order_no
func CursorList(userID uint, status int, cursor string, limit int) (map[string]interface{}, int) {
	if limit <= 0 {
		limit = 20
	}
	list := make([]OrderVO, 0)
	if pkg.DB == nil {
		return map[string]interface{}{"list": list, "next_cursor": "", "has_more": false}, pkg.CodeOK
	}
	q := pkg.DB.Model(&model.DepositOrder{}).Where("user_id = ?", userID)
	if status == 1 {
		q = q.Where("status = ?", 1)
	} else {
		q = q.Where("status IN ?", []int{2, 3})
	}
	if cursor != "" {
		q = q.Where("order_no < ?", cursor)
	}
	var orders []model.DepositOrder
	if err := q.Order("order_no desc").Limit(limit + 1).Find(&orders).Error; err != nil {
		return nil, pkg.ErrInternal
	}
	hasMore := len(orders) > limit
	if hasMore {
		orders = orders[:limit]
	}
	for _, o := range orders {
		list = append(list, OrderVO{o.OrderNo, o.SiteID, o.DoorID, o.StartTime, o.EndTime, o.Amount, o.Deposit, o.Status})
	}
	next := ""
	if hasMore && len(list) > 0 {
		next = list[len(list)-1].OrderNo
	}
	return map[string]interface{}{"list": list, "next_cursor": next, "has_more": hasMore}, pkg.CodeOK
}

// OrderDetail 订单详情（寄存地点 + 位置聚合）@ticket 05-订单查看与取出
func OrderDetail(userID uint, orderNo string) (map[string]interface{}, int) {
	if pkg.DB == nil {
		return nil, pkg.ErrOrderNF
	}
	var o model.DepositOrder
	if err := pkg.DB.Where("order_no = ? AND user_id = ?", orderNo, userID).First(&o).Error; err != nil {
		return nil, pkg.ErrOrderNF
	}
	detail := map[string]interface{}{
		"order_no": o.OrderNo, "site_id": o.SiteID, "door_id": o.DoorID,
		"start_time": o.StartTime, "end_time": o.EndTime,
		"amount": o.Amount, "deposit": o.Deposit, "status": o.Status,
	}
	// 聚合寄存地点与位置 @ticket 05
	if o.SiteID > 0 {
		var site model.LockerSite
		if pkg.DB.First(&site, o.SiteID).Error == nil {
			detail["site"] = map[string]interface{}{
				"name": site.Name, "address": site.Address, "lat": site.Lat, "lng": site.Lng,
			}
		}
	}
	return detail, pkg.CodeOK
}

// OngoingSummary 进行中订单聚合 @ticket 05-订单查看与取出
func OngoingSummary(userID uint) (map[string]interface{}, int) {
	var count int64
	var totalAmount int64
	if pkg.DB != nil {
		pkg.DB.Model(&model.DepositOrder{}).Where("user_id = ? AND status = ?", userID, 1).Count(&count)
		pkg.DB.Model(&model.DepositOrder{}).Where("user_id = ? AND status = ?", userID, 1).
			Select("COALESCE(SUM(amount),0)").Scan(&totalAmount)
	}
	return map[string]interface{}{"ongoing_count": count, "total_amount": totalAmount}, pkg.CodeOK
}

// randomToken 生成开柜令牌 @ticket 05-订单查看与取出
func randomToken() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// CreateOpenAuth 临时开柜：限时限次 + 审计 + 状态机校验 @ticket 05-订单查看与取出
func CreateOpenAuth(userID uint, orderNo string, maxTimes int, ttlSec int) (map[string]interface{}, int) {
	if pkg.DB == nil {
		return nil, pkg.ErrOrderNF
	}
	var o model.DepositOrder
	if err := pkg.DB.Where("order_no = ? AND user_id = ?", orderNo, userID).First(&o).Error; err != nil {
		return nil, pkg.ErrOrderNF
	}
	if o.Status != 1 {
		return nil, pkg.ErrOrderState
	}
	if maxTimes <= 0 {
		maxTimes = 1
	}
	if ttlSec <= 0 {
		ttlSec = 300
	}
	token := randomToken()
	auth := model.TempOpenAuth{
		OrderNo: orderNo, Token: token, MaxTimes: maxTimes, ExpireAt: time.Now().Add(time.Duration(ttlSec) * time.Second),
	}
	if err := pkg.DB.Create(&auth).Error; err != nil {
		return nil, pkg.ErrInternal
	}
	pkg.Logf("[AUDIT] @ticket 05 temp open auth issued: order=%s uid=%d token=%s", orderNo, userID, token[:8])
	return map[string]interface{}{"token": token, "expire_at": auth.ExpireAt, "max_times": maxTimes}, pkg.CodeOK
}

// UseOpenAuth 使用临时开柜令牌：限时限次校验 + 审计 @ticket 05-订单查看与取出
func UseOpenAuth(token string) (map[string]interface{}, int) {
	if pkg.DB == nil {
		return nil, pkg.ErrOpenAuthExp
	}
	err := pkg.Tx(func(tx *gorm.DB) error {
		var a model.TempOpenAuth
		if err := tx.Where("token = ?", token).First(&a).Error; err != nil {
			return err
		}
		if a.Used || time.Now().After(a.ExpireAt) || a.UsedTimes >= a.MaxTimes {
			return fmt.Errorf("auth expired or exhausted")
		}
		res := tx.Model(&model.TempOpenAuth{}).
			Where("id = ? AND used_times < max_times AND used = 0", a.ID).
			UpdateColumn("used_times", gorm.Expr("used_times + 1"))
		if res.Error != nil || res.RowsAffected == 0 {
			return fmt.Errorf("auth concurrent conflict")
		}
		return nil
	})
	if err != nil {
		return nil, pkg.ErrOpenAuthExp
	}
	pkg.Logf("[AUDIT] @ticket 05 temp open auth used: token=%s", token[:8])
	return map[string]interface{}{"token": token, "ok": true}, pkg.CodeOK
}

// MakeupPayment 超时补交：幂等 + 扣款 + 解锁/续期，事务 @ticket 05-订单查看与取出
func MakeupPayment(userID uint, orderNo string, amount int64, idempotencyKey string) (map[string]interface{}, int) {
	if idempotencyKey == "" || amount <= 0 {
		return nil, pkg.ErrBad
	}
	var dup model.MakeupRecord
	if pkg.DB != nil && pkg.DB.Where("idempotency_key = ?", idempotencyKey).First(&dup).Error == nil {
		return map[string]interface{}{"order_no": orderNo, "idempotent": true}, pkg.CodeOK
	}
	err := pkg.Tx(func(tx *gorm.DB) error {
		var o model.DepositOrder
		if err := tx.Where("order_no = ? AND user_id = ?", orderNo, userID).First(&o).Error; err != nil {
			return err
		}
		if o.Status != 1 {
			return fmt.Errorf("order state not allowed for makeup")
		}
		// 余额扣款 + 明细 @ticket 05 + 02
		if err := PayFromBalance(tx, userID, amount, "makeup", idempotencyKey); err != nil {
			return err
		}
		// 补缴记录（幂等唯一键）@ticket 05
		mr := model.MakeupRecord{OrderNo: orderNo, Amount: amount, Status: 1, IdempotencyKey: idempotencyKey}
		if err := tx.Create(&mr).Error; err != nil {
			return err
		}
		// 续期：延后结束时间 @ticket 05
		if err := tx.Model(&model.DepositOrder{}).Where("order_no = ?", orderNo).
			Update("end_time", o.EndTime.Add(24*time.Hour)).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		var existing model.MakeupRecord
		if pkg.DB != nil && pkg.DB.Where("idempotency_key = ?", idempotencyKey).First(&existing).Error == nil {
			return map[string]interface{}{"order_no": orderNo, "idempotent": true}, pkg.CodeOK
		}
		pkg.Logf("[ERROR] makeup payment: %v", err)
		if err.Error() == "insufficient balance" {
			return nil, pkg.ErrBalanceLow
		}
		return nil, pkg.ErrMakeup
	}
	return map[string]interface{}{"order_no": orderNo, "idempotent": false}, pkg.CodeOK
}

// CompleteOrder 订单完成并释放柜门 @ticket 05-订单查看与取出
func CompleteOrder(userID uint, orderNo string) (map[string]interface{}, int) {
	err := pkg.Tx(func(tx *gorm.DB) error {
		var o model.DepositOrder
		if err := tx.Where("order_no = ? AND user_id = ?", orderNo, userID).First(&o).Error; err != nil {
			return err
		}
		if o.Status != 1 {
			return fmt.Errorf("order state not allowed")
		}
		// 释放柜门 @ticket 05
		if err := tx.Model(&model.LockerDoor{}).Where("id = ?", o.DoorID).
			Updates(map[string]interface{}{"status": 0, "version": gorm.Expr("version + 1")}).Error; err != nil {
			return err
		}
		return tx.Model(&model.DepositOrder{}).Where("order_no = ?", orderNo).
			Update("status", 2).Error
	})
	if err != nil {
		return nil, pkg.ErrOrderState
	}
	return map[string]interface{}{"order_no": orderNo, "status": 2}, pkg.CodeOK
}
