package main

import (
	"bbsgo/cache"
	"bbsgo/config"
	"bbsgo/database"
	"bbsgo/routes"
	"bbsgo/seed"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

// main 程序入口
// 初始化数据库、缓存、路由并启动 HTTP 服务器
func main() {
	// 初始化数据库连接
	database.InitDB()

	// 执行数据库迁移
	database.AutoMigrate()

	// 初始化缓存
	cache.Init()

	// 初始化配置缓存
	config.InitConfigCache()

	// 初始化默认数据（版块、配置、标签、管理员、勋章等）
	seed.Init()

	// 配置路由
	r := routes.SetupRoutes()

	// 配置静态文件服务（处理本地存储的文件访问）
	configureStaticFiles(r)

	// 启动 HTTP 服务器
	log.Printf("server starting on :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}

// configureStaticFiles 配置静态文件服务
// 处理本地存储上传文件的访问
func configureStaticFiles(r *mux.Router) {
	staticURL := "/uploads"

	// 使用相对路径，直接返回文件内容
	r.HandleFunc(staticURL+"/{file:.*}", func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		fp := vars["file"]
		fullPath := filepath.Join(".", "uploads", fp)
		log.Printf("[static] serving file: %s", fullPath)
		http.ServeFile(w, req, fullPath)
	})

	log.Printf("[static] route registered: %s/*", staticURL)
}
