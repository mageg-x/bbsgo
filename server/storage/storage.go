package storage

// Storage 存储接口
// 定义文件存储的标准操作
type Storage interface {
	// Upload 上传文件
	// key: 文件存储键
	// data: 文件数据
	// contentType: 内容类型
	// 返回: 访问 URL 和错误
	Upload(key string, data []byte, contentType string) (string, error)

	// Delete 删除文件
	// key: 文件存储键
	// 返回: 错误
	Delete(key string) error

	// GetURL 获取文件访问 URL
	// key: 文件存储键
	// 返回: 访问 URL
	GetURL(key string) string

	// Name 获取存储类型名称
	// 返回: 存储类型名称
	Name() string
}

// Config 存储配置结构
// 包含各种存储类型的配置参数
type Config struct {
	StorageType string `json:"storage_type"` // 存储类型：local/qiniu/aliyun/tencent

	// 七牛云配置
	Qiniu struct {
		AccessKey string `json:"access_key"` // AccessKey
		SecretKey string `json:"secret_key"` // SecretKey
		Bucket    string `json:"bucket"`     // 存储空间
		Domain    string `json:"domain"`     // CDN 域名
	} `json:"qiniu"`

	// 本地存储配置
	Local struct {
		Path    string `json:"path"`     // 存储路径
		BaseURL string `json:"base_url"` // 访问基础 URL
	} `json:"local"`

	// 阿里云 OSS 配置
	Aliyun struct {
		AccessKeyId     string `json:"access_key_id"`     // AccessKey ID
		AccessKeySecret string `json:"access_key_secret"` // AccessKey Secret
		Bucket          string `json:"bucket"`            // 存储空间
		Endpoint        string `json:"endpoint"`          // 区域节点
		Domain          string `json:"domain"`            // CDN 域名
	} `json:"aliyun"`

	// 腾讯云 COS 配置
	Tencent struct {
		SecretId  string `json:"secret_id"`  // SecretId
		SecretKey string `json:"secret_key"` // SecretKey
		Bucket    string `json:"bucket"`     // 存储桶
		Region    string `json:"region"`     // 区域
		Domain    string `json:"domain"`     // CDN 域名
	} `json:"tencent"`
}
