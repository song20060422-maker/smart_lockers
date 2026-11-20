// @ticket 00-project：数据库/Redis 初始化；连接失败仅告警不 panic，保证本地 go run 可启动
package pkg

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 声明变量 DB *gorm.DB，用于保存运行过程中的数据

var DB *gorm.DB

// 定义函数 InitDB：封装一段可复用的处理逻辑（参数：dsn string, maxOpen, maxIdle, lifeSec int）

func InitDB(dsn string, maxOpen, maxIdle, lifeSec int) {
	// 定义并初始化局部变量 db, err，取值为 gorm.Open(mysql.Open(dsn), &gorm.Config{

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 执行 Logger: logger.Default.LogMode(logger.Warn),，完成对应的业务处理

		Logger: logger.Default.LogMode(logger.Warn),
	// 执行 })，完成对应的业务处理

	})
	// 条件判断：err != nil 时执行下面逻辑

	if err != nil {
		// 执行 log.Printf("[WARN] database connect failed: %v (服务将以无DB模式启动，相关接口返回系统错误)", err)，完成对应的业务处理

		log.Printf("[WARN] database connect failed: %v (服务将以无DB模式启动，相关接口返回系统错误)", err)
		// 为 DB 赋值：nil

		DB = nil
		// 返回  给调用方，结束当前函数

		return
	// 当前代码块结束

	}
	// 定义并初始化局部变量 sqlDB, err，取值为 db.DB()

	sqlDB, err := db.DB()
	// 条件判断：err == nil 时执行下面逻辑

	if err == nil {
		// 执行 sqlDB.SetMaxOpenConns(maxOpen)，完成对应的业务处理

		sqlDB.SetMaxOpenConns(maxOpen)
		// 执行 sqlDB.SetMaxIdleConns(maxIdle)，完成对应的业务处理

		sqlDB.SetMaxIdleConns(maxIdle)
		// 执行 sqlDB.SetConnMaxLifetime(time.Duration(lifeSec) * time.Second)，完成对应的业务处理

		sqlDB.SetConnMaxLifetime(time.Duration(lifeSec) * time.Second)
	// 当前代码块结束

	}
	// 为 DB 赋值：db

	DB = db
	// 执行 log.Println("[INFO] database connected")，完成对应的业务处理

	log.Println("[INFO] database connected")
// 当前代码块结束

}

// Tx 在事务内执行，任一环节返回 error 则整体回滚
// 定义函数 Tx：封装一段可复用的处理逻辑（参数：fn func(tx *gorm.DB）

func Tx(fn func(tx *gorm.DB) error) error {
	// 条件判断：DB == nil 时执行下面逻辑

	if DB == nil {
		// 返回 fmt.Errorf("database unavailable") 给调用方，结束当前函数

		return fmt.Errorf("database unavailable")
	// 当前代码块结束

	}
	// 返回 DB.Transaction(fn) 给调用方，结束当前函数

	return DB.Transaction(fn)
// 当前代码块结束

}
