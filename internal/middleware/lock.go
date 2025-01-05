package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hcd233/Aris-blog/internal/logger"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/cache"
	"github.com/hcd233/Aris-blog/internal/util"
	"go.uber.org/zap"
)

// RedisLockMiddleware Redis锁中间件
//
//	param serviceName string
//	param key string
//	param expire time.Duration
//	return gin.HandlerFunc
//	author centonhuang
//	update 2025-01-05 15:06:51
func RedisLockMiddleware(serviceName, key string, expire time.Duration) gin.HandlerFunc {
	redis := cache.GetRedisClient()

	return func(c *gin.Context) {
		ctx := context.Background()

		value := c.MustGet(key)

		lockKey := fmt.Sprintf("%s:%s:%v", serviceName, key, value)
		lockValue := uuid.New().String()

		success, err := redis.SetNX(ctx, lockKey, lockValue, expire).Result()
		if err != nil {
			logger.Logger.Error("[RedisLockMiddleware] failed to get lock", zap.Error(err))
			util.SendHTTPResponse(c, nil, protocol.ErrInternalError)
			c.Abort()
			return
		}

		if !success {
			lockValue, err = redis.Get(ctx, lockKey).Result()
			if err != nil {
				logger.Logger.Error("[RedisLockMiddleware] failed to get lock info", zap.String("lockKey", lockKey), zap.Error(err))
				util.SendHTTPResponse(c, nil, protocol.ErrInternalError)
				c.Abort()
				return
			}
			logger.Logger.Info("[RedisLockMiddleware] resource is locked", zap.String("lockKey", lockKey), zap.String("lockValue", lockValue))
			util.SendHTTPResponse(c, nil, protocol.ErrTooManyRequests)
			c.Abort()
			return
		}

		c.Next()

		luaScript := `
			if redis.call("get", KEYS[1]) == ARGV[1] then
				return redis.call("del", KEYS[1])
			else
				return 0
			end
		`
		if err := redis.Eval(context.Background(), luaScript, []string{lockKey}, lockValue).Err(); err != nil {
			logger.Logger.Error("[RedisLockMiddleware] failed to release lock", zap.String("lockKey", lockKey), zap.Error(err))
			util.SendHTTPResponse(c, nil, protocol.ErrInternalError)
			c.Abort()
			return
		}
	}
}
