// @ticket B端-真实对接：GORM 数据库仓储实现（内存仓储的可选真实后端），使运营端操作真正落库
// 设计：New() 仍为纯内存实现（保留单测契约）；NewWithDB() 启用 GORM 后端。
//
//	公开方法固定使用 s.db；事务内由 Tx 设置 s.tx，NoLock 方法经 activeDB() 使用事务句柄。
package store

import (
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"smart_lockers_admin/internal/model"
)

// idemDBRow 幂等键持久化行 @ticket B端-幂等
// 定义结构体 idemDBRow：承载已消费过的业务幂等键
type idemDBRow struct {
	// 结构体字段 Key，主键，列名 key，保存形如 recon:<key> 的幂等键
	Key string `gorm:"primaryKey;column:key"`
}

// TableName 返回幂等表表名 @ticket B端-幂等
// 定义方法 TableName：显式指定 GORM 使用的表名
func (idemDBRow) TableName() string { return "admin_idempotencies" }

// OpenDB 依据 DSN 建立 MySQL 连接并返回 GORM 句柄 @ticket B端-真实对接
// 定义函数 OpenDB：封装 GORM 的 mysql 驱动打开动作
func OpenDB(dsn string) (*gorm.DB, error) {
	// 调用 gorm.Open，以 mysql.Open(dsn) 作为驱动
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// 条件判断：连接出错时
	if err != nil {
		// 返回 nil 句柄与错误给调用方
		return nil, err
	}
	// 返回成功建立的数据库句柄
	return db, nil
}

// OpenDBWithPool 依据 DSN 建立连接并配置连接池 @ticket B端-真实对接
// 定义函数 OpenDBWithPool：在 OpenDB 基础上设置最大连接数、空闲数与生命周期
func OpenDBWithPool(dsn string, maxOpen, maxIdle, lifeSec int) (*gorm.DB, error) {
	// 先建立基础连接
	db, err := OpenDB(dsn)
	// 条件判断：连接失败直接返回
	if err != nil {
		// 返回错误
		return nil, err
	}
	// 取得底层 sql.DB 以配置连接池
	sqlDB, err := db.DB()
	// 条件判断：取底层连接失败
	if err != nil {
		// 返回错误
		return nil, err
	}
	// 设置最大打开连接数（非正值则保持默认）
	if maxOpen > 0 {
		// 应用最大连接数
		sqlDB.SetMaxOpenConns(maxOpen)
	}
	// 设置最大空闲连接数（非正值则保持默认）
	if maxIdle > 0 {
		// 应用最大空闲连接数
		sqlDB.SetMaxIdleConns(maxIdle)
	}
	// 设置连接最大生命周期（非正值则保持默认）
	if lifeSec > 0 {
		// 应用连接生命周期
		sqlDB.SetConnMaxLifetime(time.Duration(lifeSec) * time.Second)
	}
	// 返回配置完成的句柄
	return db, nil
}

// NewWithDB 构造带真实 GORM 后端的 Store，自动建表并写入种子数据 @ticket B端-真实对接
// 定义函数 NewWithDB：把数据库句柄注入仓储并完成初始化
func NewWithDB(db *gorm.DB) (*Store, error) {
	// 构造 Store 实例，内存 map 初始化为空，db 指向真实连接
	s := &Store{
		// 初始化 sites 映射
		sites: map[uint]*model.Site{},
		// 初始化 devices 映射
		devices: map[uint]*model.Device{},
		// 初始化 recons 映射
		recons: map[uint]*model.ReconcileRecord{},
		// 初始化 alarms 映射
		alarms: map[uint]*model.Alarm{},
		// 初始化 idem 映射
		idem: map[string]bool{},
		// 初始化自增编号
		nextID: 1,
		// 注入真实数据库句柄
		db: db,
	}
	// 自动迁移 B 端四张业务表、IoT 指令/OTA 两张表与一张幂等表 @ticket B端-IoT中台
	if err := db.AutoMigrate(&model.Site{}, &model.Device{}, &model.ReconcileRecord{}, &model.Alarm{}, &model.DeviceCommand{}, &model.OtaTask{}, &idemDBRow{}); err != nil {
		// 迁移失败时返回错误
		return nil, err
	}
	// 写入演示种子数据（仅当对应表为空）
	if err := s.seedDB(); err != nil {
		// 种子写入失败时返回错误
		return nil, err
	}
	// 打印初始化成功日志
	log.Printf("[INFO] store initialized with GORM backend: sites/devices/reconcile_records/alarms/admin_idempotencies")
	// 返回构造好的仓储实例
	return s, nil
}

// activeDB 返回当前有效句柄：事务内返回事务句柄，否则返回普通句柄 @ticket B端-事务
// 定义方法 activeDB：为 NoLock 系列在事务内提供同一事务句柄
func (s *Store) activeDB() *gorm.DB {
	// 条件判断：当前处于事务中
	if s.tx != nil {
		// 返回事务句柄
		return s.tx
	}
	// 返回普通数据库句柄
	return s.db
}

// seedDB 写入演示种子数据（与内存 seed 完全一致，表非空则跳过）@ticket B端-真实对接
// 定义方法 seedDB：逐表按主键 FirstOrCreate 写入种子
func (s *Store) seedDB() error {
	// 取当前时间作为种子时间基准
	t := now()
	// 定义站点种子切片
	sites := []*model.Site{
		// 站点 1001：与 config 中 site_admin.site_id 对齐
		&model.Site{ID: 1001, Code: "SITE-1001", Name: "北京中关村网点", City: "北京", Address: "演示地址1", ManagerPhone: "13800138001", Status: 1, DoorTotal: 24, DoorFault: 0, CreatedAt: t, UpdatedAt: t},
		// 站点 1002
		&model.Site{ID: 1002, Code: "SITE-1002", Name: "上海陆家嘴网点", City: "上海", Address: "演示地址2", ManagerPhone: "13800138002", Status: 1, DoorTotal: 24, DoorFault: 1, CreatedAt: t, UpdatedAt: t},
	}
	// 循环写入站点种子
	for _, v := range sites {
		// 按主键存在则跳过、不存在则创建
		if err := s.db.Where("id = ?", v.ID).FirstOrCreate(v).Error; err != nil {
			// 出错时返回错误
			return err
		}
	}
	// 定义设备种子切片
	devices := []*model.Device{
		// 设备 2001：柜组
		&model.Device{ID: 2001, SiteID: 1001, Code: "CAB-001", Type: "cabinet", Status: 0, LastHeartbeat: t, UpdatedAt: t},
		// 设备 2002：柜门（故障）
		&model.Device{ID: 2002, SiteID: 1001, Code: "DOOR-001", Type: "door", Status: 2, LastHeartbeat: t, UpdatedAt: t},
	}
	// 循环写入设备种子
	for _, v := range devices {
		// 按主键存在则跳过、不存在则创建
		if err := s.db.Where("id = ?", v.ID).FirstOrCreate(v).Error; err != nil {
			// 出错时返回错误
			return err
		}
	}
	// 定义对账种子（单条，ID 3001）
	recon := &model.ReconcileRecord{ID: 3001, BatchNo: "RC-DEMO-001", SiteID: 1001, ExpectedAmount: 12800, ActualAmount: 12800, Status: 1, IdempotencyKey: "seed-recon-3001", CreatedAt: t}
	// 按主键写入对账种子
	if err := s.db.Where("id = ?", recon.ID).FirstOrCreate(recon).Error; err != nil {
		// 出错时返回错误
		return err
	}
	// 定义告警种子（单条，ID 4001）
	alarm := &model.Alarm{ID: 4001, SiteID: 1001, DeviceID: 2002, Level: "warn", Content: "演示：柜门离线告警", Status: 0, IdempotencyKey: "seed-alarm-4001", CreatedAt: t}
	// 按主键写入告警种子
	if err := s.db.Where("id = ?", alarm.ID).FirstOrCreate(alarm).Error; err != nil {
		// 出错时返回错误
		return err
	}
	// 全部成功返回 nil
	return nil
}

// sitesDB 从数据库读取全部网点 @ticket B端-真实对接
// 定义函数 sitesDB：按 id 升序返回网点列表
func sitesDB(db *gorm.DB) []model.Site {
	// 定义结果切片
	var out []model.Site
	// 按 id 升序查询全部
	db.Order("id").Find(&out)
	// 返回结果
	return out
}

// siteDB 按主键读取网点 @ticket B端-真实对接
// 定义函数 siteDB：查询单个网点并返回是否存在
func siteDB(db *gorm.DB, id uint) (model.Site, bool) {
	// 定义目标变量
	var v model.Site
	// 按主键查询
	if db.First(&v, id).Error != nil {
		// 未找到返回空值与 false
		return model.Site{}, false
	}
	// 找到返回网点与 true
	return v, true
}

// addSiteDB 写入网点 @ticket B端-真实对接
// 定义函数 addSiteDB：新增一条网点记录
func addSiteDB(db *gorm.DB, v *model.Site) {
	// 执行新增，主键由数据库自增回填
	db.Create(v)
}

// updateSiteDB 更新网点 @ticket B端-真实对接
// 定义函数 updateSiteDB：读出后执行 mutate，再保存回库
func updateSiteDB(db *gorm.DB, id uint, mutate func(*model.Site) bool) bool {
	// 定义目标变量
	var v model.Site
	// 按主键查询，未找到返回 false
	if db.First(&v, id).Error != nil {
		// 未找到
		return false
	}
	// 执行业务变更闭包，返回 false 表示无需更新
	if !mutate(&v) {
		// 无需更新
		return false
	}
	// 刷新更新时间
	v.UpdatedAt = now()
	// 保存回库并返回结果
	return db.Save(&v).Error == nil
}

// devicesDB 按网点读取设备列表，siteID 为 0 时返回全部 @ticket B端-真实对接
// 定义函数 devicesDB：按条件查询设备
func devicesDB(db *gorm.DB, siteID uint) []model.Device {
	// 定义结果切片
	var out []model.Device
	// 构造查询
	q := db.Order("id")
	// 条件判断：指定了网点
	if siteID != 0 {
		// 追加网点过滤
		q = q.Where("site_id = ?", siteID)
	}
	// 执行查询
	q.Find(&out)
	// 返回结果
	return out
}

// deviceDB 按主键读取设备 @ticket B端-真实对接
// 定义函数 deviceDB：查询单个设备并返回是否存在
func deviceDB(db *gorm.DB, id uint) (model.Device, bool) {
	// 定义目标变量
	var v model.Device
	// 按主键查询
	if db.First(&v, id).Error != nil {
		// 未找到返回空值与 false
		return model.Device{}, false
	}
	// 找到返回设备与 true
	return v, true
}

// updateDeviceDB 更新设备 @ticket B端-真实对接
// 定义函数 updateDeviceDB：读出后执行 mutate，再保存回库
func updateDeviceDB(db *gorm.DB, id uint, mutate func(*model.Device) bool) bool {
	// 定义目标变量
	var v model.Device
	// 按主键查询，未找到返回 false
	if db.First(&v, id).Error != nil {
		// 未找到
		return false
	}
	// 执行业务变更闭包，返回 false 表示无需更新
	if !mutate(&v) {
		// 无需更新
		return false
	}
	// 刷新更新时间
	v.UpdatedAt = now()
	// 保存回库并返回结果
	return db.Save(&v).Error == nil
}

// reconsDB 按网点读取对账记录，siteID 为 0 时返回全部 @ticket B端-真实对接
// 定义函数 reconsDB：按条件查询对账记录
func reconsDB(db *gorm.DB, siteID uint) []model.ReconcileRecord {
	// 定义结果切片
	var out []model.ReconcileRecord
	// 构造查询
	q := db.Order("id")
	// 条件判断：指定了网点
	if siteID != 0 {
		// 追加网点过滤
		q = q.Where("site_id = ?", siteID)
	}
	// 执行查询
	q.Find(&out)
	// 返回结果
	return out
}

// addReconDB 写入对账记录 @ticket B端-真实对接
// 定义函数 addReconDB：新增一条对账记录
func addReconDB(db *gorm.DB, v *model.ReconcileRecord) {
	// 执行新增，主键由数据库自增回填
	db.Create(v)
}

// alarmsDB 按网点读取告警列表，siteID 为 0 时返回全部 @ticket B端-真实对接
// 定义函数 alarmsDB：按条件查询告警
func alarmsDB(db *gorm.DB, siteID uint) []model.Alarm {
	// 定义结果切片
	var out []model.Alarm
	// 构造查询
	q := db.Order("id")
	// 条件判断：指定了网点
	if siteID != 0 {
		// 追加网点过滤
		q = q.Where("site_id = ?", siteID)
	}
	// 执行查询
	q.Find(&out)
	// 返回结果
	return out
}

// alarmDB 按主键读取告警 @ticket B端-真实对接
// 定义函数 alarmDB：查询单个告警并返回是否存在
func alarmDB(db *gorm.DB, id uint) (model.Alarm, bool) {
	// 定义目标变量
	var v model.Alarm
	// 按主键查询
	if db.First(&v, id).Error != nil {
		// 未找到返回空值与 false
		return model.Alarm{}, false
	}
	// 找到返回告警与 true
	return v, true
}

// updateAlarmDB 更新告警 @ticket B端-真实对接
// 定义函数 updateAlarmDB：读出后执行 mutate，再保存回库
func updateAlarmDB(db *gorm.DB, id uint, mutate func(*model.Alarm) bool) bool {
	// 定义目标变量
	var v model.Alarm
	// 按主键查询，未找到返回 false
	if db.First(&v, id).Error != nil {
		// 未找到
		return false
	}
	// 执行业务变更闭包，返回 false 表示无需更新
	if !mutate(&v) {
		// 无需更新
		return false
	}
	// 保存回库并返回结果
	return db.Save(&v).Error == nil
}

// checkIdemDB 数据库幂等校验：INSERT IGNORE 成功表示首次消费 @ticket B端-幂等
// 定义函数 checkIdemDB：用唯一键冲突保证幂等
func checkIdemDB(db *gorm.DB, key string) bool {
	// 条件判断：空键不参与幂等
	if key == "" {
		// 返回 false
		return false
	}
	// 以 INSERT IGNORE 写入幂等表，靠主键冲突去重
	res := db.Exec("INSERT IGNORE INTO admin_idempotencies(`key`) VALUES(?)", key)
	// 条件判断：执行出错时按未消费处理
	if res.Error != nil {
		// 返回 false
		return false
	}
	// 影响行数为 1 表示首次消费，0 表示重复
	return res.RowsAffected > 0
}
