// @ticket B端-IoT中台：指令/OTA 仓储读写（内存模式与 GORM 模式双实现，复用 Store 事务框架）
package store

import (
	"sort"

	"gorm.io/gorm"

	"smart_lockers_admin/internal/model"
)

// ---- 指令记录 ----

// Commands 指令记录列表（内存模式；siteID 经 DeviceID 关联过滤由 service 层完成，这里按 deviceID=0 返回全部）
// 定义方法 Commands：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) Commands(deviceID uint) []model.DeviceCommand {
	// GORM 模式：改由数据库读取指令记录

	if s.db != nil {
		// 调用数据库实现并返回

		return commandsDB(s.db, deviceID)
	// 当前代码块结束

	}
	// 执行 s.mu.RLock(); defer s.mu.RUnlock()，完成对应的业务处理

	s.mu.RLock()
	defer s.mu.RUnlock()
	// 返回 s.CommandsNoLock(deviceID) 给调用方，结束当前函数

	return s.CommandsNoLock(deviceID)
// 当前代码块结束

}

// CommandsNoLock 指令记录列表（须持锁版本）
// 定义方法 CommandsNoLock：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) CommandsNoLock(deviceID uint) []model.DeviceCommand {
	// 定义并初始化局部变量 out，取值为 []model.DeviceCommand{}

	out := []model.DeviceCommand{}
	// 循环处理：for _, v := range s.cmds，逐轮迭代

	for _, v := range s.cmds {
		// 条件判断：deviceID == 0 || v.DeviceID == deviceID 时执行下面逻辑

		if deviceID == 0 || v.DeviceID == deviceID {
			// 追加记录副本
			// 执行 out = append(out, *v)，完成对应的业务处理

			out = append(out, *v)
		// 当前代码块结束

		}
	// 当前代码块结束

	}
	// 按 ID 倒序排序（最新在前）
	// 执行 sort.Slice(out, func(i, j int) bool { return out[i].ID > out[j].ID })，完成对应的业务处理

	sort.Slice(out, func(i, j int) bool { return out[i].ID > out[j].ID })
	// 返回 out 给调用方，结束当前函数

	return out
// 当前代码块结束

}

// AddCommand 写入指令记录
// 定义方法 AddCommand：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) AddCommand(v *model.DeviceCommand) {
	// GORM 模式：改由数据库写入

	if s.db != nil {
		// 调用数据库写入实现

		addCommandDB(s.db, v)
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 执行 s.mu.Lock(); defer s.mu.Unlock()，完成对应的业务处理

	s.mu.Lock()
	defer s.mu.Unlock()
	// 内存模式：分配自增 ID 并写入 map
	// 执行 s.AddCommandNoLock(v)，完成对应的业务处理

	s.AddCommandNoLock(v)
// 当前代码块结束

}

// AddCommandNoLock 写入指令记录（须持锁版本）
// 定义方法 AddCommandNoLock：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) AddCommandNoLock(v *model.DeviceCommand) {
	// GORM 模式：在当前（事务）句柄上写入

	if s.db != nil {
		// 调用数据库写入实现

		addCommandDB(s.activeDB(), v)
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 分配自增 ID
	// 为 v.ID 赋值：s.allocID()

	v.ID = s.allocID()
	// 写入 map
	// 执行 s.cmds[v.ID] = v，完成对应的业务处理

	s.cmds[v.ID] = v
// 当前代码块结束

}

// UpdateCommand 更新指令记录
// 定义方法 UpdateCommand：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) UpdateCommand(id uint, mutate func(*model.DeviceCommand) bool) bool {
	// GORM 模式：改由数据库更新

	if s.db != nil {
		// 调用数据库更新实现并返回结果

		return updateCommandDB(s.db, id, mutate)
	// 当前代码块结束

	}
	// 执行 s.mu.Lock(); defer s.mu.Unlock()，完成对应的业务处理

	s.mu.Lock()
	defer s.mu.Unlock()
	// 内存模式更新
	// 定义并初始化局部变量 v, ok，取值为 s.cmds[id]

	v, ok := s.cmds[id]
	// 条件判断：!ok || !mutate(v) 时执行下面逻辑

	if !ok || !mutate(v) {
		// 返回 false 给调用方，结束当前函数

		return false
	// 当前代码块结束

	}
	// 返回 true 给调用方，结束当前函数

	return true
// 当前代码块结束

}

// CommandByID 按 CommandID（业务指令号）查询指令
// 定义方法 CommandByID：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) CommandByID(commandID string) (model.DeviceCommand, bool) {
	// GORM 模式：改由数据库查询

	if s.db != nil {
		// 调用数据库查询实现

		return commandByCommandIDDB(s.db, commandID)
	// 当前代码块结束

	}
	// 执行 s.mu.RLock(); defer s.mu.RUnlock()，完成对应的业务处理

	s.mu.RLock()
	defer s.mu.RUnlock()
	// 循环处理：for _, v := range s.cmds，逐轮迭代

	for _, v := range s.cmds {
		// 条件判断：v.CommandID == commandID 时执行下面逻辑

		if v.CommandID == commandID {
			// 返回 *v, true 给调用方，结束当前函数

			return *v, true
		// 当前代码块结束

		}
	// 当前代码块结束

	}
	// 返回 model.DeviceCommand{}, false 给调用方，结束当前函数

	return model.DeviceCommand{}, false
// 当前代码块结束

}

// ---- OTA 任务 ----

// OtaTasks OTA 任务列表
// 定义方法 OtaTasks：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) OtaTasks(deviceID uint) []model.OtaTask {
	// GORM 模式：改由数据库读取

	if s.db != nil {
		// 调用数据库实现并返回

		return otasDB(s.db, deviceID)
	// 当前代码块结束

	}
	// 执行 s.mu.RLock(); defer s.mu.RUnlock()，完成对应的业务处理

	s.mu.RLock()
	defer s.mu.RUnlock()
	// 返回 s.OtaTasksNoLock(deviceID) 给调用方，结束当前函数

	return s.OtaTasksNoLock(deviceID)
// 当前代码块结束

}

// OtaTasksNoLock OTA 任务列表（须持锁版本）
// 定义方法 OtaTasksNoLock：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) OtaTasksNoLock(deviceID uint) []model.OtaTask {
	// 定义并初始化局部变量 out，取值为 []model.OtaTask{}

	out := []model.OtaTask{}
	// 循环处理：for _, v := range s.otas，逐轮迭代

	for _, v := range s.otas {
		// 条件判断：deviceID == 0 || v.DeviceID == deviceID 时执行下面逻辑

		if deviceID == 0 || v.DeviceID == deviceID {
			// 追加任务副本
			// 执行 out = append(out, *v)，完成对应的业务处理

			out = append(out, *v)
		// 当前代码块结束

		}
	// 当前代码块结束

	}
	// 按 ID 倒序排序
	// 执行 sort.Slice(out, func(i, j int) bool { return out[i].ID > out[j].ID })，完成对应的业务处理

	sort.Slice(out, func(i, j int) bool { return out[i].ID > out[j].ID })
	// 返回 out 给调用方，结束当前函数

	return out
// 当前代码块结束

}

// AddOta 写入 OTA 任务
// 定义方法 AddOta：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) AddOta(v *model.OtaTask) {
	// GORM 模式：改由数据库写入

	if s.db != nil {
		// 调用数据库写入实现

		addOtaDB(s.db, v)
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 执行 s.mu.Lock(); defer s.mu.Unlock()，完成对应的业务处理

	s.mu.Lock()
	defer s.mu.Unlock()
	// 执行 s.AddOtaNoLock(v)，完成对应的业务处理

	s.AddOtaNoLock(v)
// 当前代码块结束

}

// AddOtaNoLock 写入 OTA 任务（须持锁版本）
// 定义方法 AddOtaNoLock：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) AddOtaNoLock(v *model.OtaTask) {
	// GORM 模式：在当前（事务）句柄上写入

	if s.db != nil {
		// 调用数据库写入实现

		addOtaDB(s.activeDB(), v)
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 分配自增 ID
	// 为 v.ID 赋值：s.allocID()

	v.ID = s.allocID()
	// 写入 map
	// 执行 s.otas[v.ID] = v，完成对应的业务处理

	s.otas[v.ID] = v
// 当前代码块结束

}

// UpdateOta 更新 OTA 任务
// 定义方法 UpdateOta：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) UpdateOta(id uint, mutate func(*model.OtaTask) bool) bool {
	// GORM 模式：改由数据库更新

	if s.db != nil {
		// 调用数据库更新实现并返回结果

		return updateOtaDB(s.db, id, mutate)
	// 当前代码块结束

	}
	// 执行 s.mu.Lock(); defer s.mu.Unlock()，完成对应的业务处理

	s.mu.Lock()
	defer s.mu.Unlock()
	// 内存模式更新
	// 定义并初始化局部变量 v, ok，取值为 s.otas[id]

	v, ok := s.otas[id]
	// 条件判断：!ok || !mutate(v) 时执行下面逻辑

	if !ok || !mutate(v) {
		// 返回 false 给调用方，结束当前函数

		return false
	// 当前代码块结束

	}
	// 返回 true 给调用方，结束当前函数

	return true
// 当前代码块结束

}

// ---- 告警写入（IoT 告警联动使用） ----

// AddAlarm 写入告警记录 @ticket B端-IoT中台
// 定义方法 AddAlarm：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) AddAlarm(v *model.Alarm) {
	// GORM 模式：改由数据库写入

	if s.db != nil {
		// 调用数据库写入实现

		addAlarmDB(s.db, v)
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 执行 s.mu.Lock(); defer s.mu.Unlock()，完成对应的业务处理

	s.mu.Lock()
	defer s.mu.Unlock()
	// 执行 s.AddAlarmNoLock(v)，完成对应的业务处理

	s.AddAlarmNoLock(v)
// 当前代码块结束

}

// AddAlarmNoLock 写入告警记录（须持锁版本）
// 定义方法 AddAlarmNoLock：接收者 s *Store，实现该类型对应的业务行为

func (s *Store) AddAlarmNoLock(v *model.Alarm) {
	// GORM 模式：在当前（事务）句柄上写入

	if s.db != nil {
		// 调用数据库写入实现

		addAlarmDB(s.activeDB(), v)
		// 返回，结束当前函数

		return
	// 当前代码块结束

	}
	// 分配自增 ID
	// 为 v.ID 赋值：s.allocID()

	v.ID = s.allocID()
	// 写入 map
	// 执行 s.alarms[v.ID] = v，完成对应的业务处理

	s.alarms[v.ID] = v
// 当前代码块结束

}

// addAlarmDB 写入告警记录（GORM 实现）
// 定义函数 addAlarmDB：封装一段可复用的处理逻辑（参数：db *gorm.DB, v *model.Alarm）

func addAlarmDB(db *gorm.DB, v *model.Alarm) {
	// 忽略该调用的返回值，仅执行其副作用

	_ = db.Create(v).Error
// 当前代码块结束

}

// ---- GORM 实现 ----

// commandsDB 读取指令记录列表
// 定义函数 commandsDB：封装一段可复用的处理逻辑（参数：db *gorm.DB, deviceID uint）[]model.DeviceCommand

func commandsDB(db *gorm.DB, deviceID uint) []model.DeviceCommand {
	// 定义并初始化局部变量 out，取值为 []model.DeviceCommand{}

	out := []model.DeviceCommand{}
	// 定义并初始化局部变量 q，取值为 db.Order("id DESC")

	q := db.Order("id DESC")
	// 条件判断：deviceID > 0 时执行下面逻辑

	if deviceID > 0 {
		// 执行 q = q.Where("device_id = ?", deviceID)，完成对应的业务处理

		q = q.Where("device_id = ?", deviceID)
	// 当前代码块结束

	}
	// 忽略该调用的返回值，仅执行其副作用

	_ = q.Find(&out).Error
	// 返回 out 给调用方，结束当前函数

	return out
// 当前代码块结束

}

// addCommandDB 写入指令记录
// 定义函数 addCommandDB：封装一段可复用的处理逻辑（参数：db *gorm.DB, v *model.DeviceCommand）

func addCommandDB(db *gorm.DB, v *model.DeviceCommand) {
	// 忽略该调用的返回值，仅执行其副作用

	_ = db.Create(v).Error
// 当前代码块结束

}

// updateCommandDB 更新指令记录
// 定义函数 updateCommandDB：封装一段可复用的处理逻辑（参数：db *gorm.DB, id uint, mutate func(*model.DeviceCommand) bool）bool

func updateCommandDB(db *gorm.DB, id uint, mutate func(*model.DeviceCommand) bool) bool {
	// 声明变量 v model.DeviceCommand，用于保存运行过程中的数据

	var v model.DeviceCommand
	// 条件判断：err := db.First(&v, id).Error; err != nil 时执行下面逻辑

	if err := db.First(&v, id).Error; err != nil {
		// 返回 false 给调用方，结束当前函数

		return false
	// 当前代码块结束

	}
	// 条件判断：!mutate(&v) 时执行下面逻辑

	if !mutate(&v) {
		// 返回 false 给调用方，结束当前函数

		return false
	// 当前代码块结束

	}
	// 返回 db.Save(&v).Error == nil 给调用方，结束当前函数

	return db.Save(&v).Error == nil
// 当前代码块结束

}

// commandByCommandIDDB 按业务指令号查询
// 定义函数 commandByCommandIDDB：封装一段可复用的处理逻辑（参数：db *gorm.DB, commandID string）(model.DeviceCommand, bool)

func commandByCommandIDDB(db *gorm.DB, commandID string) (model.DeviceCommand, bool) {
	// 声明变量 v model.DeviceCommand，用于保存运行过程中的数据

	var v model.DeviceCommand
	// 条件判断：err := db.Where("command_id = ?", commandID).First(&v).Error; err != nil 时执行下面逻辑

	if err := db.Where("command_id = ?", commandID).First(&v).Error; err != nil {
		// 返回 model.DeviceCommand{}, false 给调用方，结束当前函数

		return model.DeviceCommand{}, false
	// 当前代码块结束

	}
	// 返回 v, true 给调用方，结束当前函数

	return v, true
// 当前代码块结束

}

// otasDB 读取 OTA 任务列表
// 定义函数 otasDB：封装一段可复用的处理逻辑（参数：db *gorm.DB, deviceID uint）[]model.OtaTask

func otasDB(db *gorm.DB, deviceID uint) []model.OtaTask {
	// 定义并初始化局部变量 out，取值为 []model.OtaTask{}

	out := []model.OtaTask{}
	// 定义并初始化局部变量 q，取值为 db.Order("id DESC")

	q := db.Order("id DESC")
	// 条件判断：deviceID > 0 时执行下面逻辑

	if deviceID > 0 {
		// 执行 q = q.Where("device_id = ?", deviceID)，完成对应的业务处理

		q = q.Where("device_id = ?", deviceID)
	// 当前代码块结束

	}
	// 忽略该调用的返回值，仅执行其副作用

	_ = q.Find(&out).Error
	// 返回 out 给调用方，结束当前函数

	return out
// 当前代码块结束

}

// addOtaDB 写入 OTA 任务
// 定义函数 addOtaDB：封装一段可复用的处理逻辑（参数：db *gorm.DB, v *model.OtaTask）

func addOtaDB(db *gorm.DB, v *model.OtaTask) {
	// 忽略该调用的返回值，仅执行其副作用

	_ = db.Create(v).Error
// 当前代码块结束

}

// updateOtaDB 更新 OTA 任务
// 定义函数 updateOtaDB：封装一段可复用的处理逻辑（参数：db *gorm.DB, id uint, mutate func(*model.OtaTask) bool）bool

func updateOtaDB(db *gorm.DB, id uint, mutate func(*model.OtaTask) bool) bool {
	// 声明变量 v model.OtaTask，用于保存运行过程中的数据

	var v model.OtaTask
	// 条件判断：err := db.First(&v, id).Error; err != nil 时执行下面逻辑

	if err := db.First(&v, id).Error; err != nil {
		// 返回 false 给调用方，结束当前函数

		return false
	// 当前代码块结束

	}
	// 条件判断：!mutate(&v) 时执行下面逻辑

	if !mutate(&v) {
		// 返回 false 给调用方，结束当前函数

		return false
	// 当前代码块结束

	}
	// 返回 db.Save(&v).Error == nil 给调用方，结束当前函数

	return db.Save(&v).Error == nil
// 当前代码块结束

}
