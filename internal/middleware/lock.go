package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hcd233/aris-blog-api/internal/logger"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/resource/cache"
	"github.com/hcd233/aris-blog-api/internal/util"
	"go.uber.org/zap"
)

// RedisLockMiddleware Redis锁中间件
//
//	param serviceName string
//	param key string
//	param expire time.Duration
//	return fiber.Handler
//	author centonhuang
//	update 2025-01-05 15:06:51
func RedisLockMiddleware(serviceName, key string, expire time.Duration) fiber.Handler {
	redis := cache.GetRedisClient()

	return func(c *fiber.Ctx) error {
		ctx := context.Background()

		value := c.Locals(key)

		lockKey := fmt.Sprintf("%s:%s:%v", serviceName, key, value)
		lockValue := uuid.New().String()

		success, err := redis.SetNX(ctx, lockKey, lockValue, expire).Result()
		if err != nil {
			logger.WithFCtx(c).Error("[RedisLockMiddleware] failed to get lock", zap.Error(err))
			util.SendHTTPResponse(c, nil, protocol.ErrInternalError)
			return c.Status(fiber.StatusInternalServerError).JSON(protocol.HTTPResponse{
				Error: protocol.ErrInternalError.Error(),
			})
		}

		if !success {
			lockValue, err = redis.Get(ctx, lockKey).Result()
			if err != nil {
				logger.WithFCtx(c).Error("[RedisLockMiddleware] failed to get lock info", zap.String("lockKey", lockKey), zap.Error(err))
				util.SendHTTPResponse(c, nil, protocol.ErrInternalError)
				return c.Status(fiber.StatusInternalServerError).JSON(protocol.HTTPResponse{
					Error: protocol.ErrInternalError.Error(),
				})
			}
			logger.WithFCtx(c).Info("[RedisLockMiddleware] resource is locked", zap.String("lockKey", lockKey), zap.String("lockValue", lockValue))
			util.SendHTTPResponse(c, nil, protocol.ErrTooManyRequests)
			return c.Status(fiber.StatusTooManyRequests).JSON(protocol.HTTPResponse{
				Error: protocol.ErrTooManyRequests.Error(),
			})
		}

		err = c.Next()

		luaScript := `
			if redis.call("get", KEYS[1]) == ARGV[1] then
				return redis.call("del", KEYS[1])
			else
				return 0
			end
		`
		if err := redis.Eval(context.Background(), luaScript, []string{lockKey}, lockValue).Err(); err != nil {
			logger.WithFCtx(c).Error("[RedisLockMiddleware] failed to release lock", zap.String("lockKey", lockKey), zap.Error(err))
			util.SendHTTPResponse(c, nil, protocol.ErrInternalError)
			return c.Status(fiber.StatusInternalServerError).JSON(protocol.HTTPResponse{
				Error: protocol.ErrInternalError.Error(),
			})
		}

		return err
	}
}
