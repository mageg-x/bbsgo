package storage

import (
	"log"
	"sync"

	"bbsgo/database"
	"bbsgo/models"
)

var (
	storageInstance Storage
	storageOnce     sync.Once
	storageInitErr  error
)

func GetStorage() (Storage, error) {
	storageOnce.Do(func() {
		storageInstance, storageInitErr = NewStorageFromConfig()
	})

	if storageInitErr != nil {
		log.Printf("[storage] failed to get storage: %v", storageInitErr)
	}

	return storageInstance, storageInitErr
}

func NewStorageFromConfig() (Storage, error) {
	config := GetStorageConfigFromDB()

	activeStorage := config.ActiveStorage
	if activeStorage == "" {
		activeStorage = "local"
	}

	log.Printf("[storage] using storage: %s", activeStorage)

	var storage Storage
	var err error

	switch activeStorage {
	case "qiniu":
		storage, err = NewQiniuStorageWithCheck(config)
	case "local":
		storage, err = NewLocalStorageWithCheck(config)
	case "aliyun":
		storage, err = NewAliyunStorageWithCheck(config)
	case "tencent":
		storage, err = NewTencentStorageWithCheck(config)
	default:
		storage, err = NewLocalStorageWithCheck(config)
	}

	if err != nil {
		log.Printf("[storage] selected storage error: %v, falling back to local", err)
		return NewLocalStorageWithCheck(config)
	}

	return storage, nil
}

func ReloadStorage() {
	storageOnce = sync.Once{}
	storageInstance = nil
	storageInitErr = nil
}

type StorageConfig struct {
	ActiveStorage string

	Qiniu struct {
		AccessKey string
		SecretKey string
		Bucket    string
		Domain    string
	}

	Local struct {
		Path    string
		BaseURL string
	}

	Aliyun struct {
		AccessKeyId     string
		AccessKeySecret string
		Bucket          string
		Endpoint        string
		Domain          string
	}

	Tencent struct {
		SecretId  string
		SecretKey string
		Bucket    string
		Region    string
		Domain    string
	}
}

func GetStorageConfigFromDB() StorageConfig {
	config := StorageConfig{}

	config.ActiveStorage = getConfigValueFromDB("active_storage")

	config.Qiniu.AccessKey = getConfigValueFromDB("qiniu_access_key")
	config.Qiniu.SecretKey = getConfigValueFromDB("qiniu_secret_key")
	config.Qiniu.Bucket = getConfigValueFromDB("qiniu_bucket")
	config.Qiniu.Domain = getConfigValueFromDB("qiniu_domain")

	config.Local.Path = getConfigValueFromDB("local_storage_path")
	config.Local.BaseURL = getConfigValueFromDB("local_storage_base_url")

	config.Aliyun.AccessKeyId = getConfigValueFromDB("aliyun_access_key_id")
	config.Aliyun.AccessKeySecret = getConfigValueFromDB("aliyun_access_key_secret")
	config.Aliyun.Bucket = getConfigValueFromDB("aliyun_bucket")
	config.Aliyun.Endpoint = getConfigValueFromDB("aliyun_endpoint")
	config.Aliyun.Domain = getConfigValueFromDB("aliyun_domain")

	config.Tencent.SecretId = getConfigValueFromDB("tencent_secret_id")
	config.Tencent.SecretKey = getConfigValueFromDB("tencent_secret_key")
	config.Tencent.Bucket = getConfigValueFromDB("tencent_bucket")
	config.Tencent.Region = getConfigValueFromDB("tencent_region")
	config.Tencent.Domain = getConfigValueFromDB("tencent_domain")

	return config
}

func getConfigValueFromDB(key string) string {
	var siteConfig models.SiteConfig
	result := database.DB.Where("key = ?", key).First(&siteConfig)
	if result.Error != nil {
		return ""
	}
	return siteConfig.Value
}

func GetSiteConfigValue(key string) string {
	return getConfigValueFromDB(key)
}
