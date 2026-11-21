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
// 定义结构体 CabinetQuery：用于承载一组相关的字段数据

type CabinetQuery struct {
	// 结构体字段 ID，类型为 uint，附带标签 `json:"id"`

	ID        uint   `json:"id"`
	// 结构体字段 SiteID，类型为 uint，附带标签 `json:"site_id"`

	SiteID    uint   `json:"site_id"`
	// 结构体字段 Code，类型为 string，附带标签 `json:"code"`

	Code      string `json:"code"`
	// 结构体字段 Size，类型为 string，附带标签 `json:"size"`

	Size      string `json:"size"`
	// 结构体字段 Total，类型为 int，附带标签 `json:"total"`

	Total     int    `json:"total"`
	// 结构体字段 Available，类型为 int，附带标签 `json:"available"`

	Available int    `json:"available"`
// 当前代码块结束

}

// QueryCabinets 柜子查询 @ticket 04-物品寄存 P95≤200ms
// 定义函数 QueryCabinets：封装一段可复用的处理逻辑（参数：siteID uint, size string）

func QueryCabinets(siteID uint, size string) ([]CabinetQuery, int) {
	// 条件判断：pkg.DB == nil 时执行下面逻辑

	if pkg.DB == nil {
		// 返回 []CabinetQuery{}, pkg.CodeOK 给调用方，结束当前函数

		return []CabinetQuery{}, pkg.CodeOK
	// 当前代码块结束

	}
	// 定义并初始化局部变量 q，取值为 pkg.DB.Model(&model.LockerCabinet{}).Where("status = ?", 1)

	q := pkg.DB.Model(&model.LockerCabinet{}).Where("status = ?", 1)
	// 条件判断：siteID > 0 时执行下面逻辑

	if siteID > 0 {
		// 为 q 赋值：q.Where("site_id = ?", siteID)

		q = q.Where("site_id = ?", siteID)
	// 当前代码块结束

	}
	// 条件判断：size != "" 时执行下面逻辑

	if size != "" {
		// 为 q 赋值：q.Where("size = ?", size)

		q = q.Where("size = ?", size)
	// 当前代码块结束

	}
	// 声明变量 cab []model.LockerCabinet，用于保存运行过程中的数据

	var cab []model.LockerCabinet
	// 条件判断：err := q.Find(&cab).Error; err != nil 时执行下面逻辑

	if err := q.Find(&cab).Error; err != nil {
		// 返回 nil, pkg.ErrInternal 给调用方，结束当前函数

		return nil, pkg.ErrInternal
	// 当前代码块结束

	}
	// 定义并初始化局部变量 res，取值为 make([]CabinetQuery, 0, len(cab))

	res := make([]CabinetQuery, 0, len(cab))
	// 循环处理：for _, c := range cab，逐轮迭代

	for _, c := range cab {
		// 为 res 赋值：append(res, CabinetQuery{c.ID, c.SiteID, c.Code, c.Size, c.Total, c.Available})

		res = append(res, CabinetQuery{c.ID, c.SiteID, c.Code, c.Size, c.Total, c.Available})
	// 当前代码块结束

	}
	// 返回 res, pkg.CodeOK 给调用方，结束当前函数

	return res, pkg.CodeOK
// 当前代码块结束

}

// QuoteFee 收费试算：读取规则 -> 计算 @ticket 04-物品寄存
// 定义函数 QuoteFee：封装一段可复用的处理逻辑（参数：size string, hours float64, coupon int64）

func QuoteFee(size string, hours float64, coupon int64) (FeeResult, int) {
	// 条件判断：pkg.DB == nil 时执行下面逻辑

	if pkg.DB == nil {
		// 返回 FeeResult{}, pkg.ErrBillingRule 给调用方，结束当前函数

		return FeeResult{}, pkg.ErrBillingRule
	// 当前代码块结束

	}
	// 声明变量 rules []model.PricingRule，用于保存运行过程中的数据

	var rules []model.PricingRule
	// 条件判断：err := pkg.DB.Where("enabled = ?", true).Find(&rules).Error; err != nil 时执行下面逻辑

	if err := pkg.DB.Where("enabled = ?", true).Find(&rules).Error; err != nil {
		// 返回 FeeResult{}, pkg.ErrInternal 给调用方，结束当前函数

		return FeeResult{}, pkg.ErrInternal
	// 当前代码块结束

	}
	// 定义并初始化局部变量 rule, ok，取值为 MatchRule(rules, size)

	rule, ok := MatchRule(rules, size)
	// 条件判断：!ok 时执行下面逻辑

	if !ok {
		// 返回 FeeResult{}, pkg.ErrBillingRule 给调用方，结束当前函数

		return FeeResult{}, pkg.ErrBillingRule
	// 当前代码块结束

	}
	// 定义并初始化局部变量 now，取值为 time.Now()

	now := time.Now()
	// 定义并初始化局部变量 in，取值为 CalcFeeInput{

	in := CalcFeeInput{
		// 执行 Rule: rule, StartTime: now, EndTime: now.Add(time.Duration(hours * float64(time.Hour))),，完成对应的业务处理

		Rule: rule, StartTime: now, EndTime: now.Add(time.Duration(hours * float64(time.Hour))),
		// 执行 Coupon: coupon,，完成对应的业务处理

		Coupon: coupon,
	// 当前代码块结束

	}
	// 返回 CalcFee(in), pkg.CodeOK 给调用方，结束当前函数

	return CalcFee(in), pkg.CodeOK
// 当前代码块结束

}

// CreateDepositInput 落单入参 @ticket 04-物品寄存
// 定义结构体 CreateDepositInput：用于承载一组相关的字段数据

type CreateDepositInput struct {
	// 结构体字段 UserID，类型为 uint

	UserID         uint
	// 结构体字段 DoorID，类型为 uint

	DoorID         uint
	// 结构体字段 SiteID，类型为 uint

	SiteID         uint
	// 结构体字段 IdempotencyKey，类型为 string

	IdempotencyKey string
	// 结构体字段 Hours，类型为 float64

	Hours          float64
	// 结构体字段 Coupon，类型为 int64

	Coupon         int64
// 当前代码块结束

}

// CreateDeposit 支付完成落单：占用柜门 + 生成订单 + 支付流水，事务 + 幂等 @ticket 04-物品寄存
// 定义函数 CreateDeposit：封装一段可复用的处理逻辑（参数：in CreateDepositInput）

func CreateDeposit(in CreateDepositInput) (map[string]interface{}, int) {
	// 条件判断：in.IdempotencyKey == "" 时执行下面逻辑

	if in.IdempotencyKey == "" {
		// 返回 nil, pkg.ErrBad 给调用方，结束当前函数

		return nil, pkg.ErrBad
	// 当前代码块结束

	}
	// 幂等前置校验：同一 key 直接返回已有结果 @ticket 04-物品寄存
	// 声明变量 exists model.DepositOrder，用于保存运行过程中的数据

	var exists model.DepositOrder
	// 条件判断：pkg.DB != nil 时执行下面逻辑

	if pkg.DB != nil {
		// 条件判断：err := pkg.DB.Where("idempotency_key = ?", in.IdempotencyKey).First(&exists).Error; err == nil 时执行下面逻辑

		if err := pkg.DB.Where("idempotency_key = ?", in.IdempotencyKey).First(&exists).Error; err == nil {
			// 返回 map[string]interface{}{"order_no": exists.OrderNo, "idempotent": true}, pkg.CodeOK 给调用方，结束当前函数

			return map[string]interface{}{"order_no": exists.OrderNo, "idempotent": true}, pkg.CodeOK
		// 当前代码块结束

		}
	// 当前代码块结束

	}
	// 声明变量 orderNo string，用于保存运行过程中的数据

	var orderNo string
	// 声明变量 payable int64，用于保存运行过程中的数据

	var payable int64
	// 定义并初始化局部变量 err，取值为 pkg.Tx(func(tx *gorm.DB) error {

	err := pkg.Tx(func(tx *gorm.DB) error {
		// 乐观锁占用柜门：仅空闲(status=0)可占用 @ticket 04-物品寄存
		// 声明变量 door model.LockerDoor，用于保存运行过程中的数据

		var door model.LockerDoor
		// 条件判断：err := tx.First(&door, in.DoorID).Error; err != nil 时执行下面逻辑

		if err := tx.First(&door, in.DoorID).Error; err != nil {
			// 返回 fmt.Errorf("door not found: %w", err) 给调用方，结束当前函数

			return fmt.Errorf("door not found: %w", err)
		// 当前代码块结束

		}
		// 条件判断：door.Status != 0 时执行下面逻辑

		if door.Status != 0 {
			// 返回 fmt.Errorf("door unavailable: %w", gorm.ErrRecordNotFound) 给调用方，结束当前函数

			return fmt.Errorf("door unavailable: %w", gorm.ErrRecordNotFound)
		// 当前代码块结束

		}
		// 定义并初始化局部变量 res，取值为 tx.Model(&model.LockerDoor{}).

		res := tx.Model(&model.LockerDoor{}).
			// 为 Where("id 赋值：? AND status = 0 AND version = ?", in.DoorID, door.Version).

			Where("id = ? AND status = 0 AND version = ?", in.DoorID, door.Version).
			// 执行 Updates(map[string]interface{}{"status": 1, "version": door.Version + 1, "updated_at": time.Now()})，完成对应的业务处理

			Updates(map[string]interface{}{"status": 1, "version": door.Version + 1, "updated_at": time.Now()})
		// 条件判断：res.Error != nil || res.RowsAffected == 0 时执行下面逻辑

		if res.Error != nil || res.RowsAffected == 0 {
			// 返回 fmt.Errorf("door occupy conflict") 给调用方，结束当前函数

			return fmt.Errorf("door occupy conflict")
		// 当前代码块结束

		}

		// 读取收费规则并计费 @ticket 04-物品寄存
		// 声明变量 rules []model.PricingRule，用于保存运行过程中的数据

		var rules []model.PricingRule
		// 条件判断：err := tx.Where("enabled = ?", true).Find(&rules).Error; err != nil 时执行下面逻辑

		if err := tx.Where("enabled = ?", true).Find(&rules).Error; err != nil {
			// 返回 err 给调用方，结束当前函数

			return err
		// 当前代码块结束

		}
		// 定义并初始化局部变量 rule, ok，取值为 MatchRule(rules, door.Size)

		rule, ok := MatchRule(rules, door.Size)
		// 条件判断：!ok 时执行下面逻辑

		if !ok {
			// 返回 fmt.Errorf("billing rule not matched") 给调用方，结束当前函数

			return fmt.Errorf("billing rule not matched")
		// 当前代码块结束

		}
		// 定义并初始化局部变量 now，取值为 time.Now()

		now := time.Now()
		// 定义并初始化局部变量 fee，取值为 CalcFee(CalcFeeInput{

		fee := CalcFee(CalcFeeInput{
			// 执行 Rule: rule, StartTime: now, EndTime: now.Add(time.Duration(in.Hours * float64(time.Hour))),，完成对应的业务处理

			Rule: rule, StartTime: now, EndTime: now.Add(time.Duration(in.Hours * float64(time.Hour))),
			// 执行 Coupon: in.Coupon,，完成对应的业务处理

			Coupon: in.Coupon,
		// 执行 })，完成对应的业务处理

		})
		// 为 payable 赋值：fee.Payable

		payable = fee.Payable

		// 为 orderNo 赋值：genOrderNo()

		orderNo = genOrderNo()
		// 定义并初始化局部变量 order，取值为 model.DepositOrder{

		order := model.DepositOrder{
			// 执行 OrderNo: orderNo, UserID: in.UserID, DoorID: in.DoorID, SiteID: in.SiteID,，完成对应的业务处理

			OrderNo: orderNo, UserID: in.UserID, DoorID: in.DoorID, SiteID: in.SiteID,
			// 执行 StartTime: now, EndTime: now.Add(time.Duration(in.Hours * float64(time.Hour))),，完成对应的业务处理

			StartTime: now, EndTime: now.Add(time.Duration(in.Hours * float64(time.Hour))),
			// 执行 Amount: fee.BaseFee - fee.Discount, Deposit: fee.Deposit, Status: 1,，完成对应的业务处理

			Amount: fee.BaseFee - fee.Discount, Deposit: fee.Deposit, Status: 1,
			// 执行 IdempotencyKey: in.IdempotencyKey,，完成对应的业务处理

			IdempotencyKey: in.IdempotencyKey,
		// 当前代码块结束

		}
		// 条件判断：err := tx.Create(&order).Error; err != nil 时执行下面逻辑

		if err := tx.Create(&order).Error; err != nil {
			// 唯一索引冲突视为幂等命中 @ticket 04-物品寄存
			// 返回 fmt.Errorf("create order: %w", err) 给调用方，结束当前函数

			return fmt.Errorf("create order: %w", err)
		// 当前代码块结束

		}
		// 返回 nil 给调用方，结束当前函数

		return nil
	// 执行 })，完成对应的业务处理

	})
	// 条件判断：err != nil 时执行下面逻辑

	if err != nil {
		// 唯一键冲突 -> 幂等返回 @ticket 04-物品寄存
		// 声明变量 existing model.DepositOrder，用于保存运行过程中的数据

		var existing model.DepositOrder
		// 条件判断：pkg.DB != nil && pkg.DB.Where("idempotency_key = ?", in.IdempotencyKey).First(&existing).Error == nil 时执行下面逻辑

		if pkg.DB != nil && pkg.DB.Where("idempotency_key = ?", in.IdempotencyKey).First(&existing).Error == nil {
			// 返回 map[string]interface{}{"order_no": existing.OrderNo, "idempotent": true}, pkg.CodeOK 给调用方，结束当前函数

			return map[string]interface{}{"order_no": existing.OrderNo, "idempotent": true}, pkg.CodeOK
		// 当前代码块结束

		}
		// 执行 pkg.Logf("[ERROR] create deposit: %v", err)，完成对应的业务处理

		pkg.Logf("[ERROR] create deposit: %v", err)
		// 返回 nil, pkg.ErrLock 给调用方，结束当前函数

		return nil, pkg.ErrLock
	// 当前代码块结束

	}
	// 返回 map[string]interface{}{"order_no": orderNo, "payable": payable, "idempotent": false}, pkg.CodeOK 给调用方，结束当前函数

	return map[string]interface{}{"order_no": orderNo, "payable": payable, "idempotent": false}, pkg.CodeOK
// 当前代码块结束

}

// ChangeCabinet 异常换柜：释放旧柜占用 -> 重分配新柜 -> 重算费用，事务保护 @ticket 04-物品寄存
// 定义函数 ChangeCabinet：封装一段可复用的处理逻辑（参数：orderNo string, newDoorID uint）

func ChangeCabinet(orderNo string, newDoorID uint) (map[string]interface{}, int) {
	// 定义并初始化局部变量 err，取值为 pkg.Tx(func(tx *gorm.DB) error {

	err := pkg.Tx(func(tx *gorm.DB) error {
		// 声明变量 order model.DepositOrder，用于保存运行过程中的数据

		var order model.DepositOrder
		// 条件判断：err := tx.Where("order_no = ?", orderNo).First(&order).Error; err != nil 时执行下面逻辑

		if err := tx.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
			// 返回 err 给调用方，结束当前函数

			return err
		// 当前代码块结束

		}
		// 条件判断：order.Status != 1 时执行下面逻辑

		if order.Status != 1 {
			// 返回 fmt.Errorf("order state not allowed") 给调用方，结束当前函数

			return fmt.Errorf("order state not allowed")
		// 当前代码块结束

		}
		// 释放旧柜 @ticket 04-物品寄存
		// 条件判断：err := tx.Model(&model.LockerDoor{}).Where("id = ?", order.DoorID). 时执行下面逻辑

		if err := tx.Model(&model.LockerDoor{}).Where("id = ?", order.DoorID).
			// 执行 Updates(map[string]interface{}{"status": 0, "version": gorm.Expr("version + 1")}).Error; err != nil {，完成对应的业务处理

			Updates(map[string]interface{}{"status": 0, "version": gorm.Expr("version + 1")}).Error; err != nil {
			// 返回 err 给调用方，结束当前函数

			return err
		// 当前代码块结束

		}
		// 占用新柜（乐观锁）@ticket 04-物品寄存
		// 声明变量 door model.LockerDoor，用于保存运行过程中的数据

		var door model.LockerDoor
		// 条件判断：err := tx.First(&door, newDoorID).Error; err != nil 时执行下面逻辑

		if err := tx.First(&door, newDoorID).Error; err != nil {
			// 返回 err 给调用方，结束当前函数

			return err
		// 当前代码块结束

		}
		// 定义并初始化局部变量 res，取值为 tx.Model(&model.LockerDoor{}).

		res := tx.Model(&model.LockerDoor{}).
			// 为 Where("id 赋值：? AND status = 0 AND version = ?", newDoorID, door.Version).

			Where("id = ? AND status = 0 AND version = ?", newDoorID, door.Version).
			// 执行 Updates(map[string]interface{}{"status": 1, "version": door.Version + 1})，完成对应的业务处理

			Updates(map[string]interface{}{"status": 1, "version": door.Version + 1})
		// 条件判断：res.Error != nil || res.RowsAffected == 0 时执行下面逻辑

		if res.Error != nil || res.RowsAffected == 0 {
			// 返回 fmt.Errorf("new door unavailable") 给调用方，结束当前函数

			return fmt.Errorf("new door unavailable")
		// 当前代码块结束

		}
		// 更新订单柜门（重算费用：尺寸变化才需重算）@ticket 04-物品寄存
		// 条件判断：err := tx.Model(&model.DepositOrder{}).Where("order_no = ?", orderNo). 时执行下面逻辑

		if err := tx.Model(&model.DepositOrder{}).Where("order_no = ?", orderNo).
			// 执行 Update("door_id", newDoorID).Error; err != nil {，完成对应的业务处理

			Update("door_id", newDoorID).Error; err != nil {
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
		// 执行 pkg.Logf("[ERROR] change cabinet: %v", err)，完成对应的业务处理

		pkg.Logf("[ERROR] change cabinet: %v", err)
		// 返回 nil, pkg.ErrLock 给调用方，结束当前函数

		return nil, pkg.ErrLock
	// 当前代码块结束

	}
	// 返回 map[string]interface{}{"order_no": orderNo, "door_id": newDoorID}, pkg.CodeOK 给调用方，结束当前函数

	return map[string]interface{}{"order_no": orderNo, "door_id": newDoorID}, pkg.CodeOK
// 当前代码块结束

}

// genOrderNo 生成订单号 @ticket 04-物品寄存
// 定义函数 genOrderNo：封装一段可复用的处理逻辑（参数：）

func genOrderNo() string {
	// 返回 fmt.Sprintf("DP%s%06d", time.Now().Format("20060102150405"), time.Now().Nanosecond()/1000) 给调用方，结束当前函数

	return fmt.Sprintf("DP%s%06d", time.Now().Format("20060102150405"), time.Now().Nanosecond()/1000)
// 当前代码块结束

}
