package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword 使用 bcrypt 对密码进行哈希加密
// password: 明文密码
// 返回: 加密后的哈希字符串和错误信息
func HashPassword(password string) (string, error) {
	// GenerateFromPassword 使用默认 cost 生成哈希
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("hash password failed: %v", err)
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword 验证密码是否与哈希匹配
// password: 明文密码
// hash: 加密后的哈希字符串
// 返回: 布尔值，表示是否匹配
func CheckPassword(password, hash string) bool {
	// CompareHashAndPassword 比较明文密码与哈希是否匹配
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Printf("password check failed: %v", err)
		return false
	}
	return true
}
