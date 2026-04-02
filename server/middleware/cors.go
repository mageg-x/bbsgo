package middleware

import (
	"log"
	"net/http"
)

// CORS 跨域资源共享中间件
// 为所有响应添加 CORS 头，支持跨域请求
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 获取请求来源，如果为空则允许所有来源
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}

		// 设置 CORS 响应头
		w.Header().Set("Access-Control-Allow-Origin", origin)                                                           // 允许的来源
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")                               // 允许的 HTTP 方法
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept, Origin") // 允许的请求头
		w.Header().Set("Access-Control-Allow-Credentials", "true")                                                      // 是否允许携带凭证（cookies）
		w.Header().Set("Access-Control-Max-Age", "86400")                                                               // 预检请求缓存时间（24小时）

		// 处理 OPTIONS 预检请求，直接返回 200
		if r.Method == "OPTIONS" {
			log.Printf("cors middleware: handling preflight request for path: %s", r.URL.Path)
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
