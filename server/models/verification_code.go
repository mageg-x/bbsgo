package models

import "time"

// VerificationCode 验证码模型
// 存储用户注册时发送的邮箱验证码
type VerificationCode struct {
	ID        uint      `gorm:"primaryKey" json:"id"`                    // 验证码唯一标识
	Email     string    `gorm:"index;not null" json:"email"`             // 邮箱地址
	Code      string    `gorm:"not null" json:"code"`                    // 验证码
	Type      string    `gorm:"not null;default:'register'" json:"type"` // 验证码类型：register=注册
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`              // 过期时间
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`        // 创建时间
}

// TableName 指定表名为 verification_codes（默认情况下 GORM 会转换为 verification_codes）
func (VerificationCode) TableName() string {
	return "verification_codes"
}
