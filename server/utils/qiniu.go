package utils

import (
	"bbsgo/config"
	"log"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

// GetQiniuUploadToken 获取七牛云上传令牌
// 返回: 上传令牌字符串
// 注意: 该函数已被 services/qiniu.go 中的更完整实现取代
func GetQiniuUploadToken() string {
	// 从配置获取七牛云凭证
	accessKey := config.GetConfig("qiniu_access_key")
	secretKey := config.GetConfig("qiniu_secret_key")
	bucket := config.GetConfig("qiniu_bucket")

	// 检查配置是否完整
	if accessKey == "" || secretKey == "" || bucket == "" {
		log.Printf("qiniu config incomplete: accessKey=%s, bucket=%s", accessKey, bucket)
		return ""
	}

	// 创建 MAC 认证
	mac := qbox.NewMac(accessKey, secretKey)

	// 创建上传策略，指定存储空间
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}

	// 生成上传令牌
	return putPolicy.UploadToken(mac)
}
