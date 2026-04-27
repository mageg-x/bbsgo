package database

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局数据库连接实例
// 整个应用程序共享同一个数据库连接
var DB *gorm.DB

// getDBPath 获取平台特定的数据库路径
func getDBPath() string {
	var dir string

	// 优先使用环境变量指定的路径
	if envDir := os.Getenv("BBSGO_DATA_DIR"); envDir != "" {
		dir = envDir
	} else {
		switch runtime.GOOS {
		case "windows":
			dir = os.Getenv("APPDATA")
			if dir == "" {
				dir = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming")
			}
			dir = filepath.Join(dir, "bbsgo")
		case "darwin":
			dir = filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "bbsgo")
		case "linux":
			dir = filepath.Join(os.Getenv("HOME"), ".bbsgo")
		default:
			dir = "."
		}
	}

	// 确保目录存在
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("failed to create database directory: %v", err)
	}

	return filepath.Join(dir, "bbsgo.db")
}

// InitDB 初始化数据库连接
// 设置数据库连接参数：连接池大小、WAL 模式等
func InitDB() {
	var err error

	dbPath := getDBPath()
	log.Printf("database path: %s", dbPath)

	// 打开 SQLite 数据库文件
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
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
