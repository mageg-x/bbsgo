package middleware

import (
	"bbsgo/database"
	"bbsgo/errors"
	"bbsgo/models"
	"bbsgo/utils"
	"context"
	"log"
	"net/http"
	"strings"
)

// contextKey 上下文键类型，用于在 context 中存储和获取值
type contextKey string

// UserContextKey 用户上下文的键名
const UserContextKey = contextKey("user")

// Auth JWT 认证中间件
// 验证请求头中的 Bearer Token，解析并验证 JWT
// 验证成功后，将用户信息存入 context 供后续处理器使用
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// OPTIONS 请求直接放行，让 CORS 预检请求通过
		if r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}

		// 从请求头获取 Authorization 字段
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Printf("auth middleware: authorization header is empty, path: %s, method: %s", r.URL.Path, r.Method)
			errors.ErrorWithStatus(w, http.StatusUnauthorized, errors.CodeUnauthorized, "")
			return
		}

		// 提取 Bearer Token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 解析并验证 Token
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			log.Printf("auth middleware: parse token failed, path: %s, method: %s, error: %v", r.URL.Path, r.Method, err)
			errors.ErrorWithStatus(w, http.StatusUnauthorized, errors.CodeUnauthorized, "")
			return
		}

		// 将用户信息存入 context
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminAuth 管理员认证中间件
// 在 Auth 中间件之后使用，验证用户是否为管理员
// 管理员权限：Role >= 2
func AdminAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 从 context 获取用户ID
		userID, ok := GetUserIDFromContext(r.Context())
		if !ok {
			log.Printf("admin auth middleware: failed to get user id from context, path: %s, method: %s", r.URL.Path, r.Method)
			errors.ErrorWithStatus(w, http.StatusUnauthorized, errors.CodeUnauthorized, "")
			return
		}

		// 查询用户信息
		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			log.Printf("admin auth middleware: user not found, userID: %d, path: %s, error: %v", userID, r.URL.Path, err)
			errors.ErrorWithStatus(w, http.StatusUnauthorized, errors.CodeUserNotFound, "")
			return
		}

		// 检查用户角色权限
		if user.Role < 2 {
			log.Printf("admin auth middleware: insufficient permissions, userID: %d, role: %d, path: %s", userID, user.Role, r.URL.Path)
			errors.ErrorWithStatus(w, http.StatusForbidden, errors.CodeNoPermission, "")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// GetUserIDFromContext 从 context 中获取当前用户ID
// ctx: 请求上下文
// 返回: 用户ID和是否获取成功
func GetUserIDFromContext(ctx context.Context) (uint, bool) {
	claims, ok := ctx.Value(UserContextKey).(*utils.Claims)
	if !ok {
		return 0, false
	}
	return claims.UserID, true
}

// GetAdminIDFromContext 从 context 中获取当前管理员用户ID
// 验证用户存在且角色为管理员
// ctx: 请求上下文
// 返回: 管理员用户ID和是否获取成功
func GetAdminIDFromContext(ctx context.Context) (uint, bool) {
	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return 0, false
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		log.Printf("get admin id from context: user not found, userID: %d, error: %v", userID, err)
		return 0, false
	}

	// 检查是否为管理员
	if user.Role < 2 {
		log.Printf("get admin id from context: user is not admin, userID: %d, role: %d", userID, user.Role)
		return 0, false
	}

	return userID, true
}
