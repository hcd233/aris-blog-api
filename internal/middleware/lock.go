package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/cache"
)

func RedisLockMiddleware(prefix, key string, expire time.Duration) gin.HandlerFunc {
	redis := cache.GetRedisClient()

	return func(c *gin.Context) {
		value := c.MustGet(key)

		lockKey := fmt.Sprintf("%s:%s:%v", prefix, key, value)
		lockValue := uuid.New().String()

		success, err := redis.SetNX(context.Background(), lockKey, lockValue, expire).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, protocol.Response{
				Code:    protocol.CodeGetLockError,
				Message: "failed to get lock",
			})
			c.Abort()
			return
		}

		if !success {
			c.JSON(http.StatusTooManyRequests, protocol.Response{
				Code:    protocol.CodeGetLockError,
				Message: "resource is locked",
			})
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
			c.JSON(http.StatusInternalServerError, protocol.Response{
				Code:    protocol.CodeReleaseLockError,
				Message: "failed to release lock",
			})
			c.Abort()
		}
	}
}
