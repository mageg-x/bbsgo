package handlers

import (
	"bbsgo/config"
	"bbsgo/database"
	"bbsgo/models"
	"bbsgo/services"
	"bbsgo/utils"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// SendCodeRequest 发送验证码请求结构
type SendCodeRequest struct {
	Email string `json:"email"` // 邮箱地址
	Type  string `json:"type"`  // 验证码类型
}

// RegisterWithCodeRequest 邮箱注册请求结构
type RegisterWithCodeRequest struct {
	Username        string `json:"username"`        // 用户名
	Nickname        string `json:"nickname"`        // 昵称
	Email           string `json:"email"`           // 邮箱
	Password        string `json:"password"`       // 密码
	ConfirmPassword string `json:"confirm_password"` // 确认密码
	Code            string `json:"code"`            // 验证码
}

// SendVerificationCode 发送邮箱验证码处理器
// 生成6位随机验证码并发送到指定邮箱
func SendVerificationCode(w http.ResponseWriter, r *http.Request) {
	// 解析请求体
	var req SendCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("send verification code: failed to decode request body, error: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	// 验证邮箱
	if req.Email == "" {
		log.Printf("send verification code: email is empty")
		utils.Error(w, 400, "邮箱不能为空")
		return
	}

	// 检查邮件服务是否启用
	emailEnabled := config.GetConfigBool("email_enabled", false)
	if !emailEnabled {
		log.Printf("send verification code: email service is disabled")
		utils.Error(w, 500, "验证码发送失败，请重试")
		return
	}

	// 检查邮箱是否已被注册
	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		log.Printf("send verification code: email already registered, email: %s", req.Email)
		utils.Error(w, 400, "该邮箱已被注册")
		return
	}

	// 生成6位随机验证码
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	expiresAt := time.Now().Add(5 * time.Minute) // 5分钟过期

	// 创建验证码记录
	verificationCode := models.VerificationCode{
		Email:     req.Email,
		Code:      code,
		Type:      "register",
		ExpiresAt: expiresAt,
	}

	if err := database.DB.Create(&verificationCode).Error; err != nil {
		log.Printf("send verification code: failed to save verification code, email: %s, error: %v", req.Email, err)
		utils.Error(w, 500, "发送验证码失败")
		return
	}

	// 发送邮件
	if err := services.SendVerificationCode(req.Email, code); err != nil {
		log.Printf("send verification code: failed to send email, email: %s, error: %v", req.Email, err)
		utils.Error(w, 500, "发送验证码失败")
		return
	}

	log.Printf("send verification code: verification code sent successfully, email: %s", req.Email)
	utils.Success(w, map[string]string{
		"message": "验证码已发送",
	})
}

// RegisterWithCode 邮箱注册处理器
// 使用邮箱验证码进行注册
func RegisterWithCode(w http.ResponseWriter, r *http.Request) {
	// 解析请求体
	var req RegisterWithCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("register with code: failed to decode request body, error: %v", err)
		utils.Error(w, 400, "无效的请求参数")
		return
	}

	// 验证必填字段
	if req.Username == "" || req.Nickname == "" || req.Email == "" || req.Password == "" || req.ConfirmPassword == "" {
		log.Printf("register with code: incomplete registration info, username: %s, nickname: %s, email: %s", req.Username, req.Nickname, req.Email)
		utils.Error(w, 400, "请填写完整信息")
		return
	}

	// 验证两次密码
	if req.Password != req.ConfirmPassword {
		log.Printf("register with code: password mismatch, username: %s", req.Username)
		utils.Error(w, 400, "两次密码输入不一致")
		return
	}

	// 检查邮件服务是否启用
	emailEnabled := config.GetConfigBool("email_enabled", false)
	if emailEnabled {
		// 验证验证码
		if req.Code == "" {
			log.Printf("register with code: verification code is empty, username: %s", req.Username)
			utils.Error(w, 400, "请输入验证码")
			return
		}

		// 查询有效的验证码
		var verificationCode models.VerificationCode
		result := database.DB.Where("email = ? AND code = ? AND type = ? AND expires_at > ?",
			req.Email, req.Code, "register", time.Now()).First(&verificationCode)
		if result.Error != nil {
			log.Printf("register with code: invalid or expired verification code, username: %s, email: %s, code: %s", req.Username, req.Email, req.Code)
			utils.Error(w, 400, "验证码无效或已过期")
			return
		}

		// 删除已使用的验证码
		if err := database.DB.Unscoped().Delete(&verificationCode).Error; err != nil {
			log.Printf("register with code: failed to delete used verification code, error: %v", err)
		}
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := database.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		log.Printf("register with code: username already exists, username: %s", req.Username)
		utils.Error(w, 400, "用户名已存在")
		return
	}

	// 检查邮箱是否已被注册
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		log.Printf("register with code: email already registered, username: %s, email: %s", req.Username, req.Email)
		utils.Error(w, 400, "邮箱已被注册")
		return
	}

	// 密码加密
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Printf("register with code: failed to hash password, username: %s, error: %v", req.Username, err)
		utils.Error(w, 500, "密码加密失败")
		return
	}

	// 创建用户
	user := models.User{
		Username:     req.Username,
		Nickname:     req.Nickname,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         0,  // 普通用户
		Credits:      0,  // 初始积分
		Level:        1,  // 初始等级
	}

	if err := database.DB.Create(&user).Error; err != nil {
		log.Printf("register with code: failed to create user, username: %s, email: %s, error: %v", req.Username, req.Email, err)
		utils.Error(w, 500, "注册失败")
		return
	}

	// 生成 JWT Token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		log.Printf("register with code: failed to generate token, userID: %d, username: %s, error: %v", user.ID, user.Username, err)
		utils.Error(w, 500, "生成令牌失败")
		return
	}

	log.Printf("register with code: user registered successfully, userID: %d, username: %s, email: %s", user.ID, user.Username, user.Email)
	utils.Success(w, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}
