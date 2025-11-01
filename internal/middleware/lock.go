package middleware

import (
	"fmt"
	"time"

	"github.com/danielgtaylor/huma/v2"
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
		ctx := c.Context()

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
		if err := redis.Eval(ctx, luaScript, []string{lockKey}, lockValue).Err(); err != nil {
			logger.WithFCtx(c).Error("[RedisLockMiddleware] failed to release lock", zap.String("lockKey", lockKey), zap.Error(err))
			util.SendHTTPResponse(c, nil, protocol.ErrInternalError)
			return c.Status(fiber.StatusInternalServerError).JSON(protocol.HTTPResponse{
				Error: protocol.ErrInternalError.Error(),
			})
		}

		return err
	}
}

// RedisLockMiddlewareForHuma Redis锁中间件 for Huma
//
//	param serviceName string
//	param key string
//	param expire time.Duration
//	return func(ctx huma.Context, next func(huma.Context))
//	author centonhuang
//	update 2025-11-01 18:10:00
func RedisLockMiddlewareForHuma(serviceName, key string, expire time.Duration) func(ctx huma.Context, next func(huma.Context)) {
	redis := cache.GetRedisClient()

	return func(ctx huma.Context, next func(huma.Context)) {
		value := ctx.Context().Value(key)

		lockKey := fmt.Sprintf("%s:%s:%v", serviceName, key, value)
		lockValue := uuid.New().String()

		success, err := redis.SetNX(ctx.Context(), lockKey, lockValue, expire).Result()
		if err != nil {
			logger.WithCtx(ctx.Context()).Error("[RedisLockMiddlewareForHuma] failed to get lock", zap.Error(err))
			ctx.SetStatus(500)
			return
		}

		if !success {
			lockValue, err = redis.Get(ctx.Context(), lockKey).Result()
			if err != nil {
				logger.WithCtx(ctx.Context()).Error("[RedisLockMiddlewareForHuma] failed to get lock info", 
					zap.String("lockKey", lockKey), zap.Error(err))
				ctx.SetStatus(500)
				return
			}
			logger.WithCtx(ctx.Context()).Info("[RedisLockMiddlewareForHuma] resource is locked", 
				zap.String("lockKey", lockKey), zap.String("lockValue", lockValue))
			ctx.SetStatus(429)
			return
		}

		next(ctx)

		luaScript := `
			if redis.call("get", KEYS[1]) == ARGV[1] then
				return redis.call("del", KEYS[1])
			else
				return 0
			end
		`
		if err := redis.Eval(ctx.Context(), luaScript, []string{lockKey}, lockValue).Err(); err != nil {
			logger.WithCtx(ctx.Context()).Error("[RedisLockMiddlewareForHuma] failed to release lock", 
				zap.String("lockKey", lockKey), zap.Error(err))
		}
	}
}
