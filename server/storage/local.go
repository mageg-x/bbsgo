package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LocalStorage 本地存储服务
// 将文件存储在本地文件系统
type LocalStorage struct {
	config StorageConfig // 存储配置
}

// newLocalStorage 创建本地存储服务实例
func newLocalStorage(config StorageConfig) *LocalStorage {
	return &LocalStorage{config: config}
}

// NewLocalStorageWithCheck 创建本地存储并检查配置
func NewLocalStorageWithCheck(config StorageConfig) (Storage, error) {
	storage := newLocalStorage(config)
	return storage, nil
}

// Name 返回存储类型名称
func (s *LocalStorage) Name() string {
	return "local"
}

// Upload 上传文件到本地
// key: 文件存储路径（相对于基础路径）
// data: 文件数据
// contentType: 内容类型
// 返回: 访问 URL 和错误
func (s *LocalStorage) Upload(key string, data []byte, contentType string) (string, error) {
	// 获取基础存储路径
	basePath := s.config.Local.Path
	if basePath == "" {
		basePath = "./uploads"
	}

	// 确保目录存在
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return "", fmt.Errorf("failed to create base directory: %v", err)
	}

	// 拼接完整文件路径
	filePath := filepath.Join(basePath, key)

	// 确保子目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create subdirectory: %v", err)
	}

	// 写入文件
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %v", err)
	}

	return s.GetURL(key), nil
}

// Delete 删除本地文件
// key: 文件存储路径
// 返回: 错误
func (s *LocalStorage) Delete(key string) error {
	basePath := s.config.Local.Path
	if basePath == "" {
		basePath = "./uploads"
	}

	filePath := filepath.Join(basePath, key)
	return os.Remove(filePath)
}

// GetURL 获取本地文件访问 URL
// key: 文件存储路径
// 返回: 访问 URL
func (s *LocalStorage) GetURL(key string) string {
	baseURL := s.config.Local.BaseURL
	if baseURL == "" {
		return "/uploads/" + key
	}

	// 确保 baseURL 以 / 结尾
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	// URL 编码 key
	encodedKey := url.PathEscape(key)
	return baseURL + encodedKey
}

// SaveUploadedFile 保存上传的文件
// file: 上传的文件句柄
// key: 存储键
// storage: 存储服务实例
// 返回: 访问 URL 和错误
func SaveUploadedFile(file *multipart.FileHeader, key string, storage Storage) (string, error) {
	// 打开文件
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer src.Close()

	// 读取文件内容
	data, err := io.ReadAll(src)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	// 获取内容类型
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// 上传文件
	return storage.Upload(key, data, contentType)
}

// GenerateFileKey 生成文件存储键
// dir: 存储目录
// filename: 原文件名
// 返回: 存储键
func GenerateFileKey(dir string, filename string) string {
	ext := filepath.Ext(filename)
	now := time.Now()
	key := fmt.Sprintf("%s/%d%s%s", dir, now.Unix(), fmt.Sprintf("%x", now.Nanosecond()), ext)
	return key
}
