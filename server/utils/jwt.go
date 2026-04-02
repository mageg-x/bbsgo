package utils

import (
	"bbsgo/config"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT 载荷结构
// 包含用户的基本身份信息
type Claims struct {
	UserID               uint   `json:"user_id"`  // 用户ID
	Username             string `json:"username"` // 用户名
	jwt.RegisteredClaims        // JWT 标准声明（过期时间、签发时间等）
}

// GenerateToken 生成 JWT 令牌
// userID: 用户ID
// username: 用户名
// 返回: 生成的令牌字符串和错误信息
func GenerateToken(userID uint, username string) (string, error) {
	// 从配置获取 JWT 密钥，如果未配置则使用默认值
	secret := config.GetConfig("jwt_secret")
	if secret == "" {
		secret = "change-this-secret-in-production"
		log.Printf("warning: using default jwt secret, please configure in production")
	}

	// 获取 token 过期天数配置，默认为 7 天
	expireDays := config.GetConfigInt("jwt_expire_days", 7)

	// 构建 JWT Claims
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 0, expireDays)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                           // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                           // 生效时间
		},
	}

	// 使用 HS256 算法签名生成 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken 解析并验证 JWT 令牌
// tokenString: 令牌字符串
// 返回: 解析后的 Claims 和错误信息
func ParseToken(tokenString string) (*Claims, error) {
	// 从配置获取 JWT 密钥
	secret := config.GetConfig("jwt_secret")
	if secret == "" {
		secret = "change-this-secret-in-production"
	}

	// 解析 token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	// 解析失败，记录日志
	if err != nil {
		log.Printf("parse token failed, token: %s, error: %v", tokenString, err)
		return nil, err
	}

	// 验证 token 有效性并提取 claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	// token 无效
	log.Printf("token is invalid, token: %s", tokenString)
	return nil, errors.New("invalid token")
}
