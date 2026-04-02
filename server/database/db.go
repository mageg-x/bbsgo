package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局数据库连接实例
// 整个应用程序共享同一个数据库连接
var DB *gorm.DB

// InitDB 初始化数据库连接
// 设置数据库连接参数：连接池大小、WAL 模式等
func InitDB() {
	var err error

	// 打开 SQLite 数据库文件
	DB, err = gorm.Open(sqlite.Open("bbsgo.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 启用 SQL 日志
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// 获取底层 sql.DB 实例以设置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("failed to get database instance: %v", err)
	}

	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(10)
	// 设置最大打开连接数
	sqlDB.SetMaxOpenConns(100)

	// 启用 WAL 模式，提升并发读写性能
	if err := DB.Exec("PRAGMA journal_mode=WAL").Error; err != nil {
		log.Printf("failed to set WAL mode: %v", err)
	}

	log.Println("database connected successfully")
}
