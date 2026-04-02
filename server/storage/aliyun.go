package storage

import (
	"fmt"
	"log"
)

// AliyunStorage 阿里云 OSS 存储服务
// 注意：当前版本尚未完整实现阿里云 OSS 功能
type AliyunStorage struct {
	config StorageConfig // 存储配置
}

// newAliyunStorage 创建阿里云存储服务实例
func newAliyunStorage(config StorageConfig) *AliyunStorage {
	return &AliyunStorage{config: config}
}

// Name 返回存储类型名称
func (s *AliyunStorage) Name() string {
	return "aliyun"
}

// Upload 上传文件到阿里云 OSS（尚未实现）
// key: 文件存储键
// data: 文件数据
// contentType: 内容类型
// 返回: 访问 URL 和错误
func (s *AliyunStorage) Upload(key string, data []byte, contentType string) (string, error) {
	if s.config.Aliyun.AccessKeyId == "" || s.config.Aliyun.Bucket == "" {
		log.Printf("[aliyun storage] config incomplete, accessKeyId: %s, bucket: %s",
			s.config.Aliyun.AccessKeyId, s.config.Aliyun.Bucket)
		return "", fmt.Errorf("aliyun oss config incomplete")
	}

	// TODO: 实现阿里云 OSS 上传功能
	return "", fmt.Errorf("aliyun oss storage not implemented yet")
}

// Delete 删除阿里云 OSS 文件（暂未实现）
// key: 文件存储键
// 返回: 错误
func (s *AliyunStorage) Delete(key string) error {
	return nil
}

// GetURL 获取阿里云 OSS 文件访问 URL
// key: 文件存储键
// 返回: 访问 URL
func (s *AliyunStorage) GetURL(key string) string {
	if s.config.Aliyun.Domain == "" {
		return ""
	}
	return fmt.Sprintf("https://%s.%s/%s", s.config.Aliyun.Bucket, s.config.Aliyun.Endpoint, key)
}

// TestConnection 测试阿里云 OSS 连接（暂未实现）
// 返回: 错误
func (s *AliyunStorage) TestConnection() error {
	if s.config.Aliyun.AccessKeyId == "" || s.config.Aliyun.Bucket == "" {
		return fmt.Errorf("aliyun oss config incomplete")
	}
	return nil
}

// NewAliyunStorageWithCheck 创建阿里云存储并检查配置
func NewAliyunStorageWithCheck(config StorageConfig) (Storage, error) {
	storage := newAliyunStorage(config)
	if err := storage.TestConnection(); err != nil {
		log.Printf("[aliyun storage] connection check failed: %v", err)
		return nil, err
	}
	return storage, nil
}
