package cache

import (
	"bbsgo/config"
	"log"
	"time"

	"github.com/dgraph-io/ristretto"
)

// Cache 全局缓存实例
// 使用 Ristretto 内存缓存库
var Cache *ristretto.Cache

// Init 初始化缓存实例
// 从配置中读取缓存参数并创建 Ristretto 缓存
func Init() {
	var err error

	// 创建 Ristretto 缓存实例
	Cache, err = ristretto.NewCache(&ristretto.Config{
		// 最大计数器数量，用于统计和 eviction
		NumCounters: int64(config.GetConfigInt("cache_num_counters", 10000)),
		// 最大缓存成本（字节）
		MaxCost: int64(config.GetConfigInt("cache_max_cost", 10000000)),
		// 每个 get 请求的缓冲区项目数
		BufferItems: 64,
	})
	if err != nil {
		log.Fatalf("failed to initialize cache: %v", err)
	}
	log.Println("cache initialized successfully")
}

// Set 设置缓存值
// key: 缓存键
// value: 缓存值
// ttl: 过期时间
func Set(key string, value interface{}, ttl time.Duration) {
	if Cache == nil {
		return
	}
	// SetWithTTL 设置缓存项，第三个参数 1 表示成本
	Cache.SetWithTTL(key, value, 1, ttl)
}

// Get 获取缓存值
// key: 缓存键
// 返回: 缓存值和是否存在
func Get(key string) (interface{}, bool) {
	if Cache == nil {
		log.Printf("cache.Get: [DEBUG] Cache is nil, key=%s", key)
		return nil, false
	}
	val, ok := Cache.Get(key)
	log.Printf("cache.Get: [DEBUG] key=%s, found=%v", key, ok)
	return val, ok
}

// Delete 删除缓存
// key: 缓存键
func Delete(key string) {
	if Cache == nil {
		return
	}
	Cache.Del(key)
}

// DeletePattern 删除匹配指定模式的缓存（暂未实现）
// pattern: 通配符模式
func DeletePattern(pattern string) {
	// Ristretto 不支持模式匹配删除，需要自行维护键列表
	log.Printf("delete pattern not implemented, pattern: %s", pattern)
}
