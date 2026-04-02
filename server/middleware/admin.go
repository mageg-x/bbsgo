package middleware

import (
	"bbsgo/database"
	"bbsgo/models"
	"bbsgo/utils"
	"log"
	"net/http"
)

// Admin 管理员权限中间件（简化版）
// 验证当前用户是否为管理员（Role >= 2）
// 注意：需要在 Auth 中间件之后使用
func Admin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 从 context 获取用户ID
		userID, ok := GetUserIDFromContext(r.Context())
		if !ok {
			log.Printf("admin middleware: failed to get user id from context")
			utils.Error(w, 401, "未认证")
			return
		}

		// 查询用户信息
		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			log.Printf("admin middleware: user not found, userID: %d, error: %v", userID, err)
			utils.Error(w, 401, "用户不存在")
			return
		}

		// 检查用户角色权限
		if user.Role < 2 {
			log.Printf("admin middleware: insufficient permissions, userID: %d, role: %d", userID, user.Role)
			utils.Error(w, 403, "需要管理员权限")
			return
		}

		next.ServeHTTP(w, r)
	})
}
