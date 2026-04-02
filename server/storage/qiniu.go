package storage

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

// QiniuStorage 七牛云存储服务
type QiniuStorage struct {
	config  StorageConfig // 存储配置
	mac     *qbox.Mac     // MAC 认证
	bucket  string        // 存储空间
	domain  string        // CDN 域名
	upToken string        // 上传令牌
}

// newQiniuStorage 创建七牛云存储服务实例
func newQiniuStorage(config StorageConfig) *QiniuStorage {
	return &QiniuStorage{config: config}
}

// Name 返回存储类型名称
func (s *QiniuStorage) Name() string {
	return "qiniu"
}

// getAuth 获取认证信息并生成上传令牌
// 返回: 错误
func (s *QiniuStorage) getAuth() error {
	// 检查配置
	if s.config.Qiniu.AccessKey == "" || s.config.Qiniu.SecretKey == "" {
		return fmt.Errorf("qiniu config incomplete")
	}

	// 创建 MAC
	s.mac = qbox.NewMac(s.config.Qiniu.AccessKey, s.config.Qiniu.SecretKey)
	s.bucket = s.config.Qiniu.Bucket
	s.domain = s.config.Qiniu.Domain

	// 创建上传策略
	putPolicy := storage.PutPolicy{
		Scope: s.bucket,
	}
	s.upToken = putPolicy.UploadToken(s.mac)

	return nil
}

// Upload 上传文件到七牛云
// key: 文件存储键
// data: 文件数据
// contentType: 内容类型
// 返回: 访问 URL 和错误
func (s *QiniuStorage) Upload(key string, data []byte, contentType string) (string, error) {
	// 获取认证
	if err := s.getAuth(); err != nil {
		log.Printf("[qiniu storage] auth error: %v", err)
		return "", err
	}

	// 配置
	cfg := storage.Config{
		Zone:          &storage.ZoneHuanan,
		UseCdnDomains: false,
		UseHTTPS:      false,
	}

	// 创建 FormUploader
	formUploader := storage.NewFormUploader(&cfg)

	// 上传结果
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{}

	// 执行上传
	dataReader := bytes.NewReader(data)
	err := formUploader.Put(context.Background(), &ret, s.upToken, key, dataReader, int64(len(data)), &putExtra)
	if err != nil {
		return "", fmt.Errorf("failed to upload to qiniu: %v", err)
	}

	return s.GetURL(ret.Key), nil
}

// Delete 删除七牛云文件（暂未实现）
// key: 文件存储键
// 返回: 错误
func (s *QiniuStorage) Delete(key string) error {
	// 七牛云 SDK 的删除功能暂未集成
	return nil
}

// GetURL 获取七牛云文件访问 URL
// key: 文件存储键
// 返回: 访问 URL
func (s *QiniuStorage) GetURL(key string) string {
	if s.domain == "" {
		return ""
	}
	return fmt.Sprintf("https://%s/%s", s.domain, key)
}

// TestConnection 测试七牛云连接
// 返回: 错误
func (s *QiniuStorage) TestConnection() error {
	if err := s.getAuth(); err != nil {
		return err
	}

	if s.bucket == "" || s.domain == "" {
		return fmt.Errorf("qiniu bucket or domain not configured")
	}

	return nil
}

// NewQiniuStorageWithCheck 创建七牛云存储并检查配置
func NewQiniuStorageWithCheck(config StorageConfig) (Storage, error) {
	storage := newQiniuStorage(config)
	if err := storage.TestConnection(); err != nil {
		log.Printf("[qiniu storage] connection check failed: %v", err)
		return nil, err
	}
	return storage, nil
}
