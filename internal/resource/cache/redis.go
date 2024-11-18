package cache

import (
	"context"
	"fmt"

	"github.com/hcd233/Aris-blog/internal/config"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
)

var rdb *redis.Client

func GetCacheInstance() *redis.Client {
	return rdb
}

func InitCache() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
		Password: config.RedisPassword,
		DB:       0,
	})

	_ = lo.Must1(rdb.Ping(context.Background()).Result())
}
