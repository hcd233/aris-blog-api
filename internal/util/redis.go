package util

import (
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient *redis.Client
	redisOnce   sync.Once
)

// GetRedis 获取Redis客户端
//
//	return *redis.Client
//	author system
//	update 2025-01-19 12:00:00
func GetRedis() *redis.Client {
	redisOnce.Do(func() {
		// 这里应该从配置文件读取Redis配置，暂时使用默认配置
		redisClient = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})
	})
	return redisClient
}