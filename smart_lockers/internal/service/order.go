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
// 定义结构体 OrderVO：用于承载一组相关的字段数据

type OrderVO struct {
	// 结构体字段 OrderNo，类型为 string，附带标签 `json:"order_no"`

	OrderNo   string    `json:"order_no"`
	// 结构体字段 SiteID，类型为 uint，附带标签 `json:"site_id"`

	SiteID    uint      `json:"site_id"`
	// 结构体字段 DoorID，类型为 uint，附带标签 `json:"door_id"`

	DoorID    uint      `json:"door_id"`
	// 结构体字段 StartTime，类型为 time.Time，附带标签 `json:"start_time"`

	StartTime time.Time `json:"start_time"`
	// 结构体字段 EndTime，类型为 time.Time，附带标签 `json:"end_time"`

	EndTime   time.Time `json:"end_time"`
	// 结构体字段 Amount，类型为 int64，附带标签 `json:"amount"`

	Amount    int64     `json:"amount"`
	// 结构体字段 Deposit，类型为 int64，附带标签 `json:"deposit"`

	Deposit   int64     `json:"deposit"`
	// 结构体字段 Status，类型为 int，附带标签 `json:"status"`

	Status    int       `json:"status"`
// 当前代码块结束

}

// CursorList 订单游标分页（当前/历史）@ticket 05-订单查看与取出
// status: 1=进行中(当前) 2/3=历史；cursor 为上一页最后一条 order_no
// 定义函数 CursorList：封装一段可复用的处理逻辑（参数：userID uint, status int, cursor string, limit int）

func CursorList(userID uint, status int, cursor string, limit int) (map[string]interface{}, int) {
	// 条件判断：limit <= 0 时执行下面逻辑

	if limit <= 0 {
		// 为 limit 赋值：20

		limit = 20
	// 当前代码块结束

	}
	// 定义并初始化局部变量 list，取值为 make([]OrderVO, 0)

	list := make([]OrderVO, 0)
	// 条件判断：pkg.DB == nil 时执行下面逻辑

	if pkg.DB == nil {
		// 返回 map[string]interface{}{"list": list, "next_cursor": "", "has_more": false}, pkg.CodeOK 给调用方，结束当前函数

		return map[string]interface{}{"list": list, "next_cursor": "", "has_more": false}, pkg.CodeOK
	// 当前代码块结束

	}
	// 定义并初始化局部变量 q，取值为 pkg.DB.Model(&model.DepositOrder{}).Where("user_id = ?", userID)

	q := pkg.DB.Model(&model.DepositOrder{}).Where("user_id = ?", userID)
	// 条件判断：status == 1 时执行下面逻辑

	if status == 1 {
		// 为 q 赋值：q.Where("status = ?", 1)

		q = q.Where("status = ?", 1)
	// 上述条件不成立时，执行 else 分支逻辑

	} else {
		// 为 q 赋值：q.Where("status IN ?", []int{2, 3})

		q = q.Where("status IN ?", []int{2, 3})
	// 当前代码块结束

	}
	// 条件判断：cursor != "" 时执行下面逻辑

	if cursor != "" {
		// 为 q 赋值：q.Where("order_no < ?", cursor)

		q = q.Where("order_no < ?", cursor)
	// 当前代码块结束

	}
	// 声明变量 orders []model.DepositOrder，用于保存运行过程中的数据

	var orders []model.DepositOrder
	// 条件判断：err := q.Order("order_no desc").Limit(limit + 1).Find(&orders).Error; err != nil 时执行下面逻辑

	if err := q.Order("order_no desc").Limit(limit + 1).Find(&orders).Error; err != nil {
		// 返回 nil, pkg.ErrInternal 给调用方，结束当前函数

		return nil, pkg.ErrInternal
	// 当前代码块结束

	}
	// 定义并初始化局部变量 hasMore，取值为 len(orders) > limit

	hasMore := len(orders) > limit
	// 条件判断：hasMore 时执行下面逻辑

	if hasMore {
		// 为 orders 赋值：orders[:limit]

		orders = orders[:limit]
	// 当前代码块结束

	}
	// 循环处理：for _, o := range orders，逐轮迭代

	for _, o := range orders {
		// 为 list 赋值：append(list, OrderVO{o.OrderNo, o.SiteID, o.DoorID, o.StartTime, o.EndTime, o.Amount, o.Deposit, o.Status})

		list = append(list, OrderVO{o.OrderNo, o.SiteID, o.DoorID, o.StartTime, o.EndTime, o.Amount, o.Deposit, o.Status})
	// 当前代码块结束

	}
	// 定义并初始化局部变量 next，取值为 ""

	next := ""
	// 条件判断：hasMore && len(list) > 0 时执行下面逻辑

	if hasMore && len(list) > 0 {
		// 为 next 赋值：list[len(list)-1].OrderNo

		next = list[len(list)-1].OrderNo
	// 当前代码块结束

	}
	// 返回 map[string]interface{}{"list": list, "next_cursor": next, "has_more": hasMore}, pkg.CodeOK 给调用方，结束当前函数

	return map[string]interface{}{"list": list, "next_cursor": next, "has_more": hasMore}, pkg.CodeOK
// 当前代码块结束

}

// OrderDetail 订单详情（寄存地点 + 位置聚合）@ticket 05-订单查看与取出
// 定义函数 OrderDetail：封装一段可复用的处理逻辑（参数：userID uint, orderNo string）

func OrderDetail(userID uint, orderNo string) (map[string]interface{}, int) {
	// 条件判断：pkg.DB == nil 时执行下面逻辑

	if pkg.DB == nil {
		// 返回 nil, pkg.ErrOrderNF 给调用方，结束当前函数

		return nil, pkg.ErrOrderNF
	// 当前代码块结束

	}
	// 声明变量 o model.DepositOrder，用于保存运行过程中的数据

	var o model.DepositOrder
	// 条件判断：err := pkg.DB.Where("order_no = ? AND user_id = ?", orderNo, userID).First(&o).Error; err != nil 时执行下面逻辑

	if err := pkg.DB.Where("order_no = ? AND user_id = ?", orderNo, userID).First(&o).Error; err != nil {
		// 返回 nil, pkg.ErrOrderNF 给调用方，结束当前函数

		return nil, pkg.ErrOrderNF
	// 当前代码块结束

	}
	// 定义并初始化局部变量 detail，取值为 map[string]interface{}{

	detail := map[string]interface{}{
		// 执行 "order_no": o.OrderNo, "site_id": o.SiteID, "door_id": o.DoorID,，完成对应的业务处理

		"order_no": o.OrderNo, "site_id": o.SiteID, "door_id": o.DoorID,
		// 执行 "start_time": o.StartTime, "end_time": o.EndTime,，完成对应的业务处理

		"start_time": o.StartTime, "end_time": o.EndTime,
		// 执行 "amount": o.Amount, "deposit": o.Deposit, "status": o.Status,，完成对应的业务处理

		"amount": o.Amount, "deposit": o.Deposit, "status": o.Status,
	// 当前代码块结束

	}
	// 聚合寄存地点与位置 @ticket 05
	// 条件判断：o.SiteID > 0 时执行下面逻辑

	if o.SiteID > 0 {
		// 声明变量 site model.LockerSite，用于保存运行过程中的数据

		var site model.LockerSite
		// 条件判断：pkg.DB.First(&site, o.SiteID).Error == nil 时执行下面逻辑

		if pkg.DB.First(&site, o.SiteID).Error == nil {
			// 为 detail["site"] 赋值：map[string]interface{}{

			detail["site"] = map[string]interface{}{
				// 执行 "name": site.Name, "address": site.Address, "lat": site.Lat, "lng": site.Lng,，完成对应的业务处理

				"name": site.Name, "address": site.Address, "lat": site.Lat, "lng": site.Lng,
			// 当前代码块结束

			}
		// 当前代码块结束

		}
	// 当前代码块结束

	}
	// 返回 detail, pkg.CodeOK 给调用方，结束当前函数

	return detail, pkg.CodeOK
// 当前代码块结束

}

// OngoingSummary 进行中订单聚合 @ticket 05-订单查看与取出
// 定义函数 OngoingSummary：封装一段可复用的处理逻辑（参数：userID uint）

func OngoingSummary(userID uint) (map[string]interface{}, int) {
	// 声明变量 count int64，用于保存运行过程中的数据

	var count int64
	// 声明变量 totalAmount int64，用于保存运行过程中的数据

	var totalAmount int64
	// 条件判断：pkg.DB != nil 时执行下面逻辑

	if pkg.DB != nil {
		// 为 pkg.DB.Model(&model.DepositOrder{}).Where("user_id 赋值：? AND status = ?", userID, 1).Count(&count)

		pkg.DB.Model(&model.DepositOrder{}).Where("user_id = ? AND status = ?", userID, 1).Count(&count)
		// 为 pkg.DB.Model(&model.DepositOrder{}).Where("user_id 赋值：? AND status = ?", userID, 1).

		pkg.DB.Model(&model.DepositOrder{}).Where("user_id = ? AND status = ?", userID, 1).
			// 执行 Select("COALESCE(SUM(amount),0)").Scan(&totalAmount)，完成对应的业务处理

			Select("COALESCE(SUM(amount),0)").Scan(&totalAmount)
	// 当前代码块结束

	}
	// 返回 map[string]interface{}{"ongoing_count": count, "total_amount": totalAmount}, pkg.CodeOK 给调用方，结束当前函数

	return map[string]interface{}{"ongoing_count": count, "total_amount": totalAmount}, pkg.CodeOK
// 当前代码块结束

}

// randomToken 生成开柜令牌 @ticket 05-订单查看与取出
// 定义函数 randomToken：封装一段可复用的处理逻辑（参数：）

func randomToken() string {
	// 定义并初始化局部变量 b，取值为 make([]byte, 16)

	b := make([]byte, 16)
	// 为 _, _ 赋值：rand.Read(b)

	_, _ = rand.Read(b)
	// 返回 hex.EncodeToString(b) 给调用方，结束当前函数

	return hex.EncodeToString(b)
// 当前代码块结束

}

// CreateOpenAuth 临时开柜：限时限次 + 审计 + 状态机校验 @ticket 05-订单查看与取出
// 定义函数 CreateOpenAuth：封装一段可复用的处理逻辑（参数：userID uint, orderNo string, maxTimes int, ttlSec int）

func CreateOpenAuth(userID uint, orderNo string, maxTimes int, ttlSec int) (map[string]interface{}, int) {
	// 条件判断：pkg.DB == nil 时执行下面逻辑

	if pkg.DB == nil {
		// 返回 nil, pkg.ErrOrderNF 给调用方，结束当前函数

		return nil, pkg.ErrOrderNF
	// 当前代码块结束

	}
	// 声明变量 o model.DepositOrder，用于保存运行过程中的数据

	var o model.DepositOrder
	// 条件判断：err := pkg.DB.Where("order_no = ? AND user_id = ?", orderNo, userID).First(&o).Error; err != nil 时执行下面逻辑

	if err := pkg.DB.Where("order_no = ? AND user_id = ?", orderNo, userID).First(&o).Error; err != nil {
		// 返回 nil, pkg.ErrOrderNF 给调用方，结束当前函数

		return nil, pkg.ErrOrderNF
	// 当前代码块结束

	}
	// 条件判断：o.Status != 1 时执行下面逻辑

	if o.Status != 1 {
		// 返回 nil, pkg.ErrOrderState 给调用方，结束当前函数

		return nil, pkg.ErrOrderState
	// 当前代码块结束

	}
	// 条件判断：maxTimes <= 0 时执行下面逻辑

	if maxTimes <= 0 {
		// 为 maxTimes 赋值：1

		maxTimes = 1
	// 当前代码块结束

	}
	// 条件判断：ttlSec <= 0 时执行下面逻辑

	if ttlSec <= 0 {
		// 为 ttlSec 赋值：300

		ttlSec = 300
	// 当前代码块结束

	}
	// 定义并初始化局部变量 token，取值为 randomToken()

	token := randomToken()
	// 定义并初始化局部变量 auth，取值为 model.TempOpenAuth{

	auth := model.TempOpenAuth{
		// 执行 OrderNo: orderNo, Token: token, MaxTimes: maxTimes, ExpireAt: time.Now().Add(time.Duration(ttlSec) * time.Second),，完成对应的业务处理

		OrderNo: orderNo, Token: token, MaxTimes: maxTimes, ExpireAt: time.Now().Add(time.Duration(ttlSec) * time.Second),
	// 当前代码块结束

	}
	// 条件判断：err := pkg.DB.Create(&auth).Error; err != nil 时执行下面逻辑

	if err := pkg.DB.Create(&auth).Error; err != nil {
		// 返回 nil, pkg.ErrInternal 给调用方，结束当前函数

		return nil, pkg.ErrInternal
	// 当前代码块结束

	}
	// 执行 pkg.Logf("[AUDIT] @ticket 05 temp open auth issued: order=%s uid=%d token=%s", orderNo, userID, token[:8])，完成对应的业务处理

	pkg.Logf("[AUDIT] @ticket 05 temp open auth issued: order=%s uid=%d token=%s", orderNo, userID, token[:8])
	// 返回 map[string]interface{}{"token": token, "expire_at": auth.ExpireAt, "max_times": maxTimes}, pkg.CodeOK 给调用方，结束当前函数

	return map[string]interface{}{"token": token, "expire_at": auth.ExpireAt, "max_times": maxTimes}, pkg.CodeOK
// 当前代码块结束

}

// UseOpenAuth 使用临时开柜令牌：限时限次校验 + 审计 @ticket 05-订单查看与取出
// 定义函数 UseOpenAuth：封装一段可复用的处理逻辑（参数：token string）

func UseOpenAuth(token string) (map[string]interface{}, int) {
	// 条件判断：pkg.DB == nil 时执行下面逻辑

	if pkg.DB == nil {
		// 返回 nil, pkg.ErrOpenAuthExp 给调用方，结束当前函数

		return nil, pkg.ErrOpenAuthExp
	// 当前代码块结束

	}
	// 定义并初始化局部变量 err，取值为 pkg.Tx(func(tx *gorm.DB) error {

	err := pkg.Tx(func(tx *gorm.DB) error {
		// 声明变量 a model.TempOpenAuth，用于保存运行过程中的数据

		var a model.TempOpenAuth
		// 条件判断：err := tx.Where("token = ?", token).First(&a).Error; err != nil 时执行下面逻辑

		if err := tx.Where("token = ?", token).First(&a).Error; err != nil {
			// 返回 err 给调用方，结束当前函数

			return err
		// 当前代码块结束

		}
		// 条件判断：a.Used || time.Now().After(a.ExpireAt) || a.UsedTimes >= a.MaxTimes 时执行下面逻辑

		if a.Used || time.Now().After(a.ExpireAt) || a.UsedTimes >= a.MaxTimes {
			// 返回 fmt.Errorf("auth expired or exhausted") 给调用方，结束当前函数

			return fmt.Errorf("auth expired or exhausted")
		// 当前代码块结束

		}
		// 定义并初始化局部变量 res，取值为 tx.Model(&model.TempOpenAuth{}).

		res := tx.Model(&model.TempOpenAuth{}).
			// 为 Where("id 赋值：? AND used_times < max_times AND used = 0", a.ID).

			Where("id = ? AND used_times < max_times AND used = 0", a.ID).
			// 执行 UpdateColumn("used_times", gorm.Expr("used_times + 1"))，完成对应的业务处理

			UpdateColumn("used_times", gorm.Expr("used_times + 1"))
		// 条件判断：res.Error != nil || res.RowsAffected == 0 时执行下面逻辑

		if res.Error != nil || res.RowsAffected == 0 {
			// 返回 fmt.Errorf("auth concurrent conflict") 给调用方，结束当前函数

			return fmt.Errorf("auth concurrent conflict")
		// 当前代码块结束

		}
		// 返回 nil 给调用方，结束当前函数

		return nil
	// 执行 })，完成对应的业务处理

	})
	// 条件判断：err != nil 时执行下面逻辑

	if err != nil {
		// 返回 nil, pkg.ErrOpenAuthExp 给调用方，结束当前函数

		return nil, pkg.ErrOpenAuthExp
	// 当前代码块结束

	}
	// 执行 pkg.Logf("[AUDIT] @ticket 05 temp open auth used: token=%s", token[:8])，完成对应的业务处理

	pkg.Logf("[AUDIT] @ticket 05 temp open auth used: token=%s", token[:8])
	// 返回 map[string]interface{}{"token": token, "ok": true}, pkg.CodeOK 给调用方，结束当前函数

	return map[string]interface{}{"token": token, "ok": true}, pkg.CodeOK
// 当前代码块结束

}

// MakeupPayment 超时补交：幂等 + 扣款 + 解锁/续期，事务 @ticket 05-订单查看与取出
// 定义函数 MakeupPayment：封装一段可复用的处理逻辑（参数：userID uint, orderNo string, amount int64, idempotencyKey string）

func MakeupPayment(userID uint, orderNo string, amount int64, idempotencyKey string) (map[string]interface{}, int) {
	// 条件判断：idempotencyKey == "" || amount <= 0 时执行下面逻辑

	if idempotencyKey == "" || amount <= 0 {
		// 返回 nil, pkg.ErrBad 给调用方，结束当前函数

		return nil, pkg.ErrBad
	// 当前代码块结束

	}
	// 声明变量 dup model.MakeupRecord，用于保存运行过程中的数据

	var dup model.MakeupRecord
	// 条件判断：pkg.DB != nil && pkg.DB.Where("idempotency_key = ?", idempotencyKey).First(&dup).Error == nil 时执行下面逻辑

	if pkg.DB != nil && pkg.DB.Where("idempotency_key = ?", idempotencyKey).First(&dup).Error == nil {
		// 返回 map[string]interface{}{"order_no": orderNo, "idempotent": true}, pkg.CodeOK 给调用方，结束当前函数

		return map[string]interface{}{"order_no": orderNo, "idempotent": true}, pkg.CodeOK
	// 当前代码块结束

	}
	// 定义并初始化局部变量 err，取值为 pkg.Tx(func(tx *gorm.DB) error {

	err := pkg.Tx(func(tx *gorm.DB) error {
		// 声明变量 o model.DepositOrder，用于保存运行过程中的数据

		var o model.DepositOrder
		// 条件判断：err := tx.Where("order_no = ? AND user_id = ?", orderNo, userID).First(&o).Error; err != nil 时执行下面逻辑

		if err := tx.Where("order_no = ? AND user_id = ?", orderNo, userID).First(&o).Error; err != nil {
			// 返回 err 给调用方，结束当前函数

			return err
		// 当前代码块结束

		}
		// 条件判断：o.Status != 1 时执行下面逻辑

		if o.Status != 1 {
			// 返回 fmt.Errorf("order state not allowed for makeup") 给调用方，结束当前函数

			return fmt.Errorf("order state not allowed for makeup")
		// 当前代码块结束

		}
		// 余额扣款 + 明细 @ticket 05 + 02
		// 条件判断：err := PayFromBalance(tx, userID, amount, "makeup", idempotencyKey); err != nil 时执行下面逻辑

		if err := PayFromBalance(tx, userID, amount, "makeup", idempotencyKey); err != nil {
			// 返回 err 给调用方，结束当前函数

			return err
		// 当前代码块结束

		}
		// 补缴记录（幂等唯一键）@ticket 05
		// 定义并初始化局部变量 mr，取值为 model.MakeupRecord{OrderNo: orderNo, Amount: amount, Status: 1, IdempotencyKey: idempotencyKey}

		mr := model.MakeupRecord{OrderNo: orderNo, Amount: amount, Status: 1, IdempotencyKey: idempotencyKey}
		// 条件判断：err := tx.Create(&mr).Error; err != nil 时执行下面逻辑

		if err := tx.Create(&mr).Error; err != nil {
			// 返回 err 给调用方，结束当前函数

			return err
		// 当前代码块结束

		}
		// 续期：延后结束时间 @ticket 05
		// 条件判断：err := tx.Model(&model.DepositOrder{}).Where("order_no = ?", orderNo). 时执行下面逻辑

		if err := tx.Model(&model.DepositOrder{}).Where("order_no = ?", orderNo).
			// 执行 Update("end_time", o.EndTime.Add(24*time.Hour)).Error; err != nil {，完成对应的业务处理

			Update("end_time", o.EndTime.Add(24*time.Hour)).Error; err != nil {
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
		// 声明变量 existing model.MakeupRecord，用于保存运行过程中的数据

		var existing model.MakeupRecord
		// 条件判断：pkg.DB != nil && pkg.DB.Where("idempotency_key = ?", idempotencyKey).First(&existing).Error == nil 时执行下面逻辑

		if pkg.DB != nil && pkg.DB.Where("idempotency_key = ?", idempotencyKey).First(&existing).Error == nil {
			// 返回 map[string]interface{}{"order_no": orderNo, "idempotent": true}, pkg.CodeOK 给调用方，结束当前函数

			return map[string]interface{}{"order_no": orderNo, "idempotent": true}, pkg.CodeOK
		// 当前代码块结束

		}
		// 执行 pkg.Logf("[ERROR] makeup payment: %v", err)，完成对应的业务处理

		pkg.Logf("[ERROR] makeup payment: %v", err)
		// 条件判断：err.Error() == "insufficient balance" 时执行下面逻辑

		if err.Error() == "insufficient balance" {
			// 返回 nil, pkg.ErrBalanceLow 给调用方，结束当前函数

			return nil, pkg.ErrBalanceLow
		// 当前代码块结束

		}
		// 返回 nil, pkg.ErrMakeup 给调用方，结束当前函数

		return nil, pkg.ErrMakeup
	// 当前代码块结束

	}
	// 返回 map[string]interface{}{"order_no": orderNo, "idempotent": false}, pkg.CodeOK 给调用方，结束当前函数

	return map[string]interface{}{"order_no": orderNo, "idempotent": false}, pkg.CodeOK
// 当前代码块结束

}

// CompleteOrder 订单完成并释放柜门 @ticket 05-订单查看与取出
// 定义函数 CompleteOrder：封装一段可复用的处理逻辑（参数：userID uint, orderNo string）

func CompleteOrder(userID uint, orderNo string) (map[string]interface{}, int) {
	// 定义并初始化局部变量 err，取值为 pkg.Tx(func(tx *gorm.DB) error {

	err := pkg.Tx(func(tx *gorm.DB) error {
		// 声明变量 o model.DepositOrder，用于保存运行过程中的数据

		var o model.DepositOrder
		// 条件判断：err := tx.Where("order_no = ? AND user_id = ?", orderNo, userID).First(&o).Error; err != nil 时执行下面逻辑

		if err := tx.Where("order_no = ? AND user_id = ?", orderNo, userID).First(&o).Error; err != nil {
			// 返回 err 给调用方，结束当前函数

			return err
		// 当前代码块结束

		}
		// 条件判断：o.Status != 1 时执行下面逻辑

		if o.Status != 1 {
			// 返回 fmt.Errorf("order state not allowed") 给调用方，结束当前函数

			return fmt.Errorf("order state not allowed")
		// 当前代码块结束

		}
		// 释放柜门 @ticket 05
		// 条件判断：err := tx.Model(&model.LockerDoor{}).Where("id = ?", o.DoorID). 时执行下面逻辑

		if err := tx.Model(&model.LockerDoor{}).Where("id = ?", o.DoorID).
			// 执行 Updates(map[string]interface{}{"status": 0, "version": gorm.Expr("version + 1")}).Error; err != nil {，完成对应的业务处理

			Updates(map[string]interface{}{"status": 0, "version": gorm.Expr("version + 1")}).Error; err != nil {
			// 返回 err 给调用方，结束当前函数

			return err
		// 当前代码块结束

		}
		// 返回 tx.Model(&model.DepositOrder{}).Where("order_no = ?", orderNo). 给调用方，结束当前函数

		return tx.Model(&model.DepositOrder{}).Where("order_no = ?", orderNo).
			// 执行 Update("status", 2).Error，完成对应的业务处理

			Update("status", 2).Error
	// 执行 })，完成对应的业务处理

	})
	// 条件判断：err != nil 时执行下面逻辑

	if err != nil {
		// 返回 nil, pkg.ErrOrderState 给调用方，结束当前函数

		return nil, pkg.ErrOrderState
	// 当前代码块结束

	}
	// 返回 map[string]interface{}{"order_no": orderNo, "status": 2}, pkg.CodeOK 给调用方，结束当前函数

	return map[string]interface{}{"order_no": orderNo, "status": 2}, pkg.CodeOK
// 当前代码块结束

}
