package main

import (
	"bbsgo/cache"
	"bbsgo/config"
	"bbsgo/database"
	"bbsgo/fileserver"
	"bbsgo/routes"
	"bbsgo/seed"
	"bbsgo/utils"
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Server 应用实例
type Server struct {
	Addr   string
	Srv    *http.Server
	ShutCh chan struct{}
}

// ParseFlags 解析命令行参数
func ParseFlags() (ip string, port int) {
	var bindIP = flag.String("a", "", "IP address to bind (default: all interfaces)")
	var bindPort = flag.Int("p", 8080, "Port number to listen on (default: 8080)")
	flag.Parse()
	return *bindIP, *bindPort
}

// Init 初始化应用（数据库、缓存、配置等）
func Init() {
	database.InitDB()
	database.AutoMigrate()
	cache.Init()
	config.InitConfigCache()
	seed.Init()

	// 启动匿名清理定时任务
	utils.StartAnonymityCleanupTask()
}

// SetupRouter 设置路由
func SetupRouter() *mux.Router {
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

	return r
}

// NewServer 创建服务器实例
func NewServer(addr string) *Server {
	r := SetupRouter()
	srv := &http.Server{Addr: addr, Handler: r}
	return &Server{
		Addr: addr,
		Srv:  srv,
	}
}

// Start 启动服务器
func (s *Server) Start() error {
	log.Printf("server starting on %s...", s.Addr)
	return s.Srv.ListenAndServe()
}

// Shutdown 关闭服务器
func (s *Server) Shutdown(ctx context.Context) error {
	return s.Srv.Shutdown(ctx)
}
