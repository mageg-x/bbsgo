package handlers

import (
	"bbsgo/errors"
	"bbsgo/storage"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

// CheckFileExists 检查文件是否已存在（用于秒传）
// 通过文件内容hash生成相同key，如果文件存在则直接返回URL
func CheckFileExists(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")
	contentHash := r.URL.Query().Get("content_hash") // 文件内容MD5
	if filename == "" {
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}
	if contentHash == "" {
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}

	// 获取存储服务实例
	storageSvc, err := storage.GetStorage()
	if err != nil {
		log.Printf("[upload/exists] failed to get storage service, error: %v", err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	// 根据文件扩展名确定目录
	ext := strings.ToLower(filepath.Ext(filename))
	imageExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true, ".svg": true, ".bmp": true}
	videoExts := map[string]bool{".mp4": true, ".webm": true, ".ogg": true, ".mov": true, ".mkv": true, ".avi": true}

	dir := ""
	if imageExts[ext] {
		dir = "images"
	} else if videoExts[ext] {
		dir = "videos"
	}

	// 生成文件key（使用content_hash确保相同内容生成相同key）
	key := storage.GenerateFileKeyWithHash(dir, filename, contentHash)

	// 检查文件是否存在
	if storageSvc.Exists(key) {
		url := storageSvc.GetURL(key)
		log.Printf("[upload/exists] file exists, key: %s, url: %s", key, url)
		errors.Success(w, map[string]interface{}{
			"exists": true,
			"url":    url,
			"key":    key,
		})
		return
	}

	log.Printf("[upload/exists] file not exists, key: %s", key)
	errors.Success(w, map[string]interface{}{
		"exists": false,
		"key":    key,
	})
}

// UploadFile 文件上传处理器
// 支持图片和视频上传到配置的存储服务（本地/七牛云/阿里云/腾讯云）
// 最大文件大小：图片50MB，视频500MB
func UploadFile(w http.ResponseWriter, r *http.Request) {
	log.Printf("[upload] upload handler started")

	// 解析 multipart 表单，最大500MB
	if err := r.ParseMultipartForm(500 << 20); err != nil {
		log.Printf("[upload] failed to parse multipart form, error: %v", err)
		errors.Error(w, errors.CodeFileSizeExceeded, "")
		return
	}

	// 获取上传文件
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("[upload] failed to get form file, error: %v", err)
		errors.Error(w, errors.CodeInvalidParams, "")
		return
	}
	defer file.Close()

	log.Printf("[upload] received file: %s, size: %d bytes", header.Filename, header.Size)

	// 获取文件扩展名并验证
	ext := strings.ToLower(filepath.Ext(header.Filename))

	// 支持的文件格式
	imageExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
		".svg":  true,
		".bmp":  true,
	}

	videoExts := map[string]bool{
		".mp4":  true,
		".webm": true,
		".ogg":  true,
		".mov":  true,
		".mkv":  true,
		".avi":  true,
	}

	// 验证文件类型
	isImage := imageExts[ext]
	isVideo := videoExts[ext]

	if !isImage && !isVideo {
		log.Printf("[upload] unsupported file type: %s", ext)
		errors.Error(w, errors.CodeFileTypeUnsupported, "")
		return
	}

	// 验证文件大小（图片50MB，视频500MB）
	if isImage && header.Size > 50*1024*1024 {
		log.Printf("[upload] image too large: %d bytes", header.Size)
		errors.Error(w, errors.CodeImageSizeExceeded, "")
		return
	}

	if isVideo && header.Size > 500*1024*1024 {
		log.Printf("[upload] video too large: %d bytes", header.Size)
		errors.Error(w, errors.CodeFileSizeExceeded, "")
		return
	}

	// 读取文件内容
	fileData, err := io.ReadAll(file)
	if err != nil {
		log.Printf("[upload] failed to read file data, error: %v", err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}
	log.Printf("[upload] file read success, size: %d bytes", len(fileData))

	// 获取存储服务实例
	storageSvc, err := storage.GetStorage()
	if err != nil {
		log.Printf("[upload] failed to get storage service, error: %v", err)
		errors.Error(w, errors.CodeServerInternal, "")
		return
	}

	// 获取上传目录参数
	dir := r.URL.Query().Get("dir")
	if dir == "" {
		if isImage {
			dir = "images"
		} else if isVideo {
			dir = "videos"
		}
	}

	// 获取content_hash参数（可选，用于秒传）
	contentHash := r.URL.Query().Get("content_hash")

	// 生成存储文件key
	// 如果提供了content_hash，使用GenerateFileKeyWithHash确保与秒传检查一致
	var key string
	if contentHash != "" {
		key = storage.GenerateFileKeyWithHash(dir, header.Filename, contentHash)
		log.Printf("[upload] using content hash for key, hash: %s", contentHash[:16])
	} else {
		key = storage.GenerateFileKey(dir, header.Filename)
		log.Printf("[upload] using filename hash for key")
	}

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		if isImage {
			contentType = "image/jpeg"
		} else if isVideo {
			contentType = "video/mp4"
		} else {
			contentType = "application/octet-stream"
		}
	}

	// 上传到存储服务
	url, err := storageSvc.Upload(key, fileData, contentType)
	if err != nil {
		log.Printf("[upload] failed to upload to storage, error: %v", err)
		errors.Error(w, errors.CodeUploadFailed, "")
		return
	}

	log.Printf("[upload] upload success, url: %s", url)
	log.Printf("[upload] upload handler finished")

	errors.Success(w, map[string]string{
		"url": url,
	})
}
