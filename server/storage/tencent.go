package storage

import (
	"fmt"
	"log"
)

// TencentStorage 腾讯云 COS 存储服务
// 注意：当前版本尚未完整实现腾讯云 COS 功能
type TencentStorage struct {
	config StorageConfig // 存储配置
}

// newTencentStorage 创建腾讯云存储服务实例
func newTencentStorage(config StorageConfig) *TencentStorage {
	return &TencentStorage{config: config}
}

// Name 返回存储类型名称
func (s *TencentStorage) Name() string {
	return "tencent"
}

// Upload 上传文件到腾讯云 COS（尚未实现）
// key: 文件存储键
// data: 文件数据
// contentType: 内容类型
// 返回: 访问 URL 和错误
func (s *TencentStorage) Upload(key string, data []byte, contentType string) (string, error) {
	if s.config.Tencent.SecretId == "" || s.config.Tencent.Bucket == "" {
		log.Printf("[tencent storage] config incomplete, secretId: %s, bucket: %s",
			s.config.Tencent.SecretId, s.config.Tencent.Bucket)
		return "", fmt.Errorf("tencent cos config incomplete")
	}

	// TODO: 实现腾讯云 COS 上传功能
	return "", fmt.Errorf("tencent cos storage not implemented yet")
}

// Delete 删除腾讯云 COS 文件（暂未实现）
// key: 文件存储键
// 返回: 错误
func (s *TencentStorage) Delete(key string) error {
	return nil
}

// GetURL 获取腾讯云 COS 文件访问 URL
// key: 文件存储键
// 返回: 访问 URL
func (s *TencentStorage) GetURL(key string) string {
	if s.config.Tencent.Domain == "" {
		return ""
	}
	return fmt.Sprintf("https://%s/%s", s.config.Tencent.Domain, key)
}

// TestConnection 测试腾讯云 COS 连接（暂未实现）
// 返回: 错误
func (s *TencentStorage) TestConnection() error {
	if s.config.Tencent.SecretId == "" || s.config.Tencent.Bucket == "" {
		return fmt.Errorf("tencent cos config incomplete")
	}
	return nil
}

// NewTencentStorageWithCheck 创建腾讯云存储并检查配置
func NewTencentStorageWithCheck(config StorageConfig) (Storage, error) {
	storage := newTencentStorage(config)
	if err := storage.TestConnection(); err != nil {
		log.Printf("[tencent storage] connection check failed: %v", err)
		return nil, err
	}
	return storage, nil
}
