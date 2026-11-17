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

var DB *gorm.DB

func InitDB(dsn string, maxOpen, maxIdle, lifeSec int) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Printf("[WARN] database connect failed: %v (服务将以无DB模式启动，相关接口返回系统错误)", err)
		DB = nil
		return
	}
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxOpenConns(maxOpen)
		sqlDB.SetMaxIdleConns(maxIdle)
		sqlDB.SetConnMaxLifetime(time.Duration(lifeSec) * time.Second)
	}
	DB = db
	log.Println("[INFO] database connected")
}

// Tx 在事务内执行，任一环节返回 error 则整体回滚
func Tx(fn func(tx *gorm.DB) error) error {
	if DB == nil {
		return fmt.Errorf("database unavailable")
	}
	return DB.Transaction(fn)
}
