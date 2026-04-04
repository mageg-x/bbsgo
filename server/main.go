package main

import (
	"bbsgo/cache"
	"bbsgo/config"
	"bbsgo/database"
	"bbsgo/fileserver"
	"bbsgo/routes"
	"bbsgo/seed"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	database.InitDB()
	database.AutoMigrate()
	cache.Init()
	config.InitConfigCache()
	seed.Init()

	r := mux.NewRouter()

	// API 路由 - 最高优先级
	apiRouter := r.PathPrefix("/api/v1").Subrouter()
	routes.SetupAPIRoutes(apiRouter)

	// 管理后台 - /console 下的所有请求都由 admin 处理（包括 assets）
	consoleRouter := r.PathPrefix("/console").Subrouter()
	// /console 重定向
	consoleRouter.HandleFunc("", func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, "/console/", http.StatusFound)
	})
	// /console/ 下的所有路径都由 admin 的 SPA 处理
	consoleRouter.PathPrefix("/").Handler(http.HandlerFunc(fileserver.ServeAdmin))

	// 上传文件 - 优先处理
	r.HandleFunc("/uploads/{file:.*}", func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		fp := vars["file"]
		fullPath := "./uploads/" + fp
		log.Printf("[static] serving file: %s", fullPath)
		http.ServeFile(w, req, fullPath)
	})

	// 主站 - 所有其他路径（SPA）
	r.PathPrefix("/").Handler(http.HandlerFunc(fileserver.ServeSite))

	log.Printf("server starting on :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
