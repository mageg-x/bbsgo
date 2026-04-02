package handlers

import (
	"bbsgo/storage"
	"bbsgo/utils"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

// UploadFile 文件上传处理器
// 支持图片上传到配置的存储服务（本地/七牛云/阿里云/腾讯云）
// 最大文件大小：50MB
func UploadFile(w http.ResponseWriter, r *http.Request) {
	log.Printf("[upload] upload handler started")

	// 解析 multipart 表单，最大50MB
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		log.Printf("[upload] failed to parse multipart form, error: %v", err)
		utils.Error(w, 400, "文件大小超过限制(最大50MB)")
		return
	}

	// 获取上传文件
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("[upload] failed to get form file, error: %v", err)
		utils.Error(w, 400, "获取文件失败")
		return
	}
	defer file.Close()

	log.Printf("[upload] received file: %s, size: %d bytes", header.Filename, header.Size)

	// 验证文件大小
	if header.Size > 50*1024*1024 {
		log.Printf("[upload] file too large: %d bytes", header.Size)
		utils.Error(w, 400, "文件大小超过限制(最大50MB)")
		return
	}

	// 获取文件扩展名并验证
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
		".svg":  true,
		".bmp":  true,
	}

	if !allowedExts[ext] {
		log.Printf("[upload] unsupported file type: %s", ext)
		utils.Error(w, 400, "不支持的文件类型")
		return
	}

	// 读取文件内容
	fileData, err := io.ReadAll(file)
	if err != nil {
		log.Printf("[upload] failed to read file data, error: %v", err)
		utils.Error(w, 500, "读取文件失败")
		return
	}
	log.Printf("[upload] file read success, size: %d bytes", len(fileData))

	// 获取存储服务实例
	storageSvc, err := storage.GetStorage()
	if err != nil {
		log.Printf("[upload] failed to get storage service, error: %v", err)
		utils.Error(w, 500, "存储服务不可用")
		return
	}

	// 获取上传目录参数
	dir := r.URL.Query().Get("dir")
	if dir == "" {
		dir = "uploads"
	}

	// 生成存储文件key
	key := storage.GenerateFileKey(dir, header.Filename)
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// 上传到存储服务
	url, err := storageSvc.Upload(key, fileData, contentType)
	if err != nil {
		log.Printf("[upload] failed to upload to storage, error: %v", err)
		utils.Error(w, 500, "文件上传失败")
		return
	}

	log.Printf("[upload] upload success, url: %s", url)
	log.Printf("[upload] upload handler finished")

	utils.Success(w, map[string]string{
		"url": url,
	})
}
