package handlers

import (
	"bbsgo/database"
	"bbsgo/errors"
	"bbsgo/models"
	"bbsgo/utils"
	"encoding/json"
	"log"
	"net/http"
)

// RegisterRequest 注册请求结构
type RegisterRequest struct {
	Username string `json:"username"` // 用户名
	Email    string `json:"email"`    // 邮箱
	Password string `json:"password"` // 密码
}

// LoginRequest 登录请求结构
type LoginRequest struct {
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
}

// Register 用户注册处理器
// 处理用户注册请求，创建新用户并返回 JWT Token
func Register(w http.ResponseWriter, r *http.Request) {
	// 检查是否允许注册
	if !utils.GetConfigBool("allow_register", true) {
		log.Printf("register: registration disabled")
		errors.Error(w, errors.CodeRegisterDisabled, "")
		return
	}

	// 解析请求体
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("register: failed to decode request body, error: %v", err)
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	// 验证必填字段
	if req.Username == "" || req.Email == "" || req.Password == "" {
		log.Printf("register: incomplete registration info, username: %s, email: %s", req.Username, req.Email)
		errors.Error(w, errors.CodeIncompleteInfo, "")
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := database.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		log.Printf("register: username already exists, username: %s", req.Username)
		errors.Error(w, errors.CodeUsernameExists, "")
		return
	}

	// 检查邮箱是否已被注册
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		log.Printf("register: email already registered, email: %s", req.Email)
		errors.Error(w, errors.CodeEmailExists, "")
		return
	}

	// 密码加密
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Printf("register: failed to hash password, error: %v", err)
		errors.Error(w, errors.CodePasswordHashFailed, "")
		return
	}

	// 创建用户
	user := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         0, // 普通用户
		Credits:      0, // 初始积分
		Level:        1, // 初始等级
	}

	if err := database.DB.Create(&user).Error; err != nil {
		log.Printf("register: failed to create user, username: %s, email: %s, error: %v", req.Username, req.Email, err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	// 生成 JWT Token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		log.Printf("register: failed to generate token, userID: %d, error: %v", user.ID, err)
		errors.Error(w, errors.CodeTokenGenerateFailed, "")
		return
	}

	log.Printf("register: user registered successfully, userID: %d, username: %s", user.ID, req.Username)
	errors.Success(w, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

// Login 用户登录处理器
// 验证用户名密码，返回 JWT Token
func Login(w http.ResponseWriter, r *http.Request) {
	// 解析请求体
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("login: failed to decode request body, error: %v", err)
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	// 查询用户
	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		log.Printf("login: user not found, username: %s", req.Username)
		errors.Error(w, errors.CodeUsernameOrPassword, "")
		return
	}

	// 验证密码
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		log.Printf("login: password mismatch, username: %s", req.Username)
		errors.Error(w, errors.CodeUsernameOrPassword, "")
		return
	}

	// 生成 JWT Token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		log.Printf("login: failed to generate token, userID: %d, error: %v", user.ID, err)
		errors.Error(w, errors.CodeTokenGenerateFailed, "")
		return
	}

	log.Printf("login: user logged in successfully, userID: %d, username: %s", user.ID, req.Username)
	errors.Success(w, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}
