package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"time"

	"bbsgo/config"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

// QiniuService 七牛云服务结构
// 提供文件上传、删除等功能
type QiniuService struct {
	accessKey string // AccessKey
	secretKey string // SecretKey
	bucket    string // 存储空间名称
	domain    string // CDN 域名
	mac       *qbox.Mac
	uploadMgr *storage.FormUploader // 上传管理器
	bucketMgr *storage.BucketManager // 存储空间管理器
}

// NewQiniuService 创建七牛云服务实例
// 从配置中读取七牛云相关参数
func NewQiniuService() *QiniuService {
	accessKey := config.GetConfig("qiniu_access_key")
	secretKey := config.GetConfig("qiniu_secret_key")
	bucket := config.GetConfig("qiniu_bucket")
	domain := config.GetConfig("qiniu_domain")

	// 创建 MAC 认证
	mac := qbox.NewMac(accessKey, secretKey)

	// 配置
	cfg := &storage.Config{
		Zone:          &storage.ZoneHuanan, // 华南区
		UseHTTPS:      false,
		UseCdnDomains: false,
	}

	return &QiniuService{
		accessKey: accessKey,
		secretKey: secretKey,
		bucket:    bucket,
		domain:    domain,
		mac:       mac,
		uploadMgr: storage.NewFormUploader(cfg),
		bucketMgr: storage.NewBucketManager(mac, cfg),
	}
}

// UploadFile 上传文件
// key: 文件存储 key
// data: 文件数据
// size: 文件大小
// 返回: 文件访问 URL 和错误信息
func (s *QiniuService) UploadFile(key string, data io.Reader, size int64) (string, error) {
	log.Printf("[qiniu] starting upload, key: %s, size: %d bytes", key, size)

	// 检查配置是否完整
	if s.accessKey == "" {
		log.Printf("[qiniu] error: access key is empty")
		return "", fmt.Errorf("qiniu access_key not configured")
	}
	if s.secretKey == "" {
		log.Printf("[qiniu] error: secret key is empty")
		return "", fmt.Errorf("qiniu secret_key not configured")
	}
	if s.bucket == "" {
		log.Printf("[qiniu] error: bucket is empty")
		return "", fmt.Errorf("qiniu bucket not configured")
	}
	if s.domain == "" {
		log.Printf("[qiniu] error: domain is empty")
		return "", fmt.Errorf("qiniu domain not configured")
	}

	log.Printf("[qiniu] config loaded, bucket: %s, domain: %s", s.bucket, s.domain)

	// 创建上传策略
	putPolicy := storage.PutPolicy{
		Scope: s.bucket,
	}

	// 生成上传令牌
	upToken := putPolicy.UploadToken(s.mac)
	log.Printf("[qiniu] upload token generated")

	// 上传结果
	ret := storage.PutRet{}

	// 执行上传
	err := s.uploadMgr.Put(context.Background(), &ret, upToken, key, data, size, nil)
	if err != nil {
		log.Printf("[qiniu] error: upload failed: %v", err)
		return "", fmt.Errorf("failed to upload file: %v", err)
	}

	log.Printf("[qiniu] upload success, key: %s, hash: %s", ret.Key, ret.Hash)

	// 拼接访问 URL
	publicURL := fmt.Sprintf("https://%s/%s", s.domain, ret.Key)
	log.Printf("[qiniu] public url: %s", publicURL)

	return publicURL, nil
}

// UploadLocalFile 上传本地文件
// key: 文件存储 key
// localPath: 本地文件路径
// 返回: 文件访问 URL 和错误信息
func (s *QiniuService) UploadLocalFile(key, localPath string) (string, error) {
	log.Printf("[qiniu] starting local file upload, key: %s, localPath: %s", key, localPath)

	// 检查配置
	if s.accessKey == "" || s.secretKey == "" || s.bucket == "" || s.domain == "" {
		return "", fmt.Errorf("qiniu config incomplete")
	}

	// 创建上传策略
	putPolicy := storage.PutPolicy{
		Scope: s.bucket,
	}
	upToken := putPolicy.UploadToken(s.mac)

	// 上传结果
	ret := storage.PutRet{}

	// 执行本地上传
	err := s.uploadMgr.PutFile(context.Background(), &ret, upToken, key, localPath, nil)
	if err != nil {
		log.Printf("[qiniu] error: local upload failed: %v", err)
		return "", fmt.Errorf("failed to upload local file: %v", err)
	}

	log.Printf("[qiniu] local upload success, key: %s, hash: %s", ret.Key, ret.Hash)

	publicURL := fmt.Sprintf("https://%s/%s", s.domain, ret.Key)
	return publicURL, nil
}

// DeleteFile 删除文件
// key: 文件存储 key
// 返回: 错误信息
func (s *QiniuService) DeleteFile(key string) error {
	err := s.bucketMgr.Delete(s.bucket, key)
	if err != nil {
		log.Printf("[qiniu] error: delete file failed, key: %s, error: %v", key, err)
		return fmt.Errorf("failed to delete file: %v", err)
	}
	log.Printf("[qiniu] file deleted, key: %s", key)
	return nil
}

// GetFileURL 获取文件访问 URL
// key: 文件存储 key
// 返回: 文件访问 URL
func (s *QiniuService) GetFileURL(key string) string {
	return fmt.Sprintf("https://%s/%s", s.domain, key)
}

// UploadImage 上传图片（简化接口）
// data: 图片数据
// size: 图片大小
// filename: 原文件名
// 返回: 图片访问 URL 和错误信息
func UploadImage(data io.Reader, size int64, filename string) (string, error) {
	service := NewQiniuService()

	ext := filepath.Ext(filename)
	key := fmt.Sprintf("images/%d%s%s", time.Now().Unix(), fmt.Sprintf("%x", time.Now().Nanosecond()), ext)

	return service.UploadFile(key, data, size)
}

// UploadLocalImage 上传本地图片（简化接口）
// localPath: 本地图片路径
// filename: 原文件名
// 返回: 图片访问 URL 和错误信息
func UploadLocalImage(localPath, filename string) (string, error) {
	service := NewQiniuService()

	ext := filepath.Ext(filename)
	key := fmt.Sprintf("images/%d%s%s", time.Now().Unix(), fmt.Sprintf("%x", time.Now().Nanosecond()), ext)

	return service.UploadLocalFile(key, localPath)
}

// UploadToQiniu 上传文件到七牛云（简化接口）
// fileData: 文件数据
// fileName: 原文件名
// dir: 存储目录
// 返回: 文件访问 URL 和错误信息
func UploadToQiniu(fileData []byte, fileName string, dir string) (string, error) {
	log.Printf("[qiniu] starting upload, fileName: %s, size: %d bytes, dir: %s", fileName, len(fileData), dir)

	// 获取配置
	accessKey := config.GetConfig("qiniu_access_key")
	secretKey := config.GetConfig("qiniu_secret_key")
	bucket := config.GetConfig("qiniu_bucket")
	domain := config.GetConfig("qiniu_domain")

	log.Printf("[qiniu] config loaded, bucket: %s, domain: %s", bucket, domain)

	// 检查配置
	if accessKey == "" {
		log.Printf("[qiniu] error: access key is empty")
		return "", fmt.Errorf("qiniu access_key not configured")
	}
	if secretKey == "" {
		log.Printf("[qiniu] error: secret key is empty")
		return "", fmt.Errorf("qiniu secret_key not configured")
	}
	if bucket == "" {
		log.Printf("[qiniu] error: bucket is empty")
		return "", fmt.Errorf("qiniu bucket not configured")
	}
	if domain == "" {
		log.Printf("[qiniu] error: domain is empty")
		return "", fmt.Errorf("qiniu domain not configured")
	}

	// 创建 MAC
	mac := qbox.NewMac(accessKey, secretKey)

	// 创建上传策略
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	upToken := putPolicy.UploadToken(mac)

	// 创建 FormUploader
	cfg := storage.Config{
		Zone:          &storage.ZoneHuanan,
		UseCdnDomains: false,
		UseHTTPS:      false,
	}
	formUploader := storage.NewFormUploader(&cfg)

	// 生成存储 key
	ext := filepath.Ext(fileName)
	if dir == "" {
		dir = "files"
	}
	key := fmt.Sprintf("%s/%d%s%s", dir, time.Now().Unix(), fmt.Sprintf("%x", time.Now().Nanosecond()), ext)
	log.Printf("[qiniu] storage key: %s", key)

	// 执行上传
	data := bytes.NewReader(fileData)
	ret := storage.PutRet{}
	err := formUploader.Put(context.Background(), &ret, upToken, key, data, int64(len(fileData)), nil)
	if err != nil {
		log.Printf("[qiniu] error: upload failed: %v", err)
		return "", fmt.Errorf("upload failed: %v", err)
	}

	log.Printf("[qiniu] upload success, key: %s, hash: %s", ret.Key, ret.Hash)

	// 拼接访问 URL
	publicURL := fmt.Sprintf("https://%s/%s", domain, ret.Key)
	log.Printf("[qiniu] public url: %s", publicURL)

	return publicURL, nil
}
