package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/logger"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/resource/cache"
	"github.com/hcd233/aris-blog-api/internal/util"
	"github.com/samber/lo"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/redis"
	"go.uber.org/zap"
)

// RateLimiterMiddleware 限频中间件
//
//	param serviceName string
//	param key string
//	param period time.Duration
//	param limit int64
//	return fiber.Handler
//	author centonhuang
//	update 2025-01-05 15:06:44
func RateLimiterMiddleware(serviceName, key string, period time.Duration, limit int64) fiber.Handler {
	// 创建限频规则
	rate := limiter.Rate{
		Period: period,
		Limit:  limit,
	}

	redisClient := cache.GetRedisClient()
	// 使用Redis存储限频数据
	store := lo.Must1(redis.NewStoreWithOptions(redisClient, limiter.StoreOptions{
		Prefix: serviceName,
	}))

	// 创建限频实例
	instance := limiter.New(store, rate)

	return func(c *fiber.Ctx) error {
		var keyValue, value string
		if key == "" {
			keyValue = "ip"
			value = c.IP() // 如果没有指定的参数，则使用 IP 地址作为 key
		} else {
			value = c.Locals(key).(string)
		}

		// 设置限频 key
		limiterKey := fmt.Sprintf("%s:%v", keyValue, value)
		c.Locals(constant.CtxKeyLimiter, limiterKey)

		// 检查限频
		context, err := instance.Get(c.Context(), limiterKey)
		if err != nil {
			logger.LoggerWithFiberContext(c).Error("[RateLimiterMiddleware] failed to get rate limit", zap.Error(err))
			return c.Status(fiber.StatusInternalServerError).JSON(protocol.HTTPResponse{
				Error: protocol.ErrInternalError.Error(),
			})
		}

		if context.Reached {
			fields := []zap.Field{zap.String("serviceName", serviceName)}

			if key == "" {
				fields = append(fields, zap.String("key", "ip"), zap.String("value", c.IP()))
			} else {
				fields = append(fields, zap.String("key", key), zap.String("value", value))
			}

			logger.LoggerWithFiberContext(c).Error("[RateLimiterMiddleware] rate limit reached", fields...)
			util.SendHTTPResponse(c, nil, protocol.ErrTooManyRequests)
			return c.Status(fiber.StatusTooManyRequests).JSON(protocol.HTTPResponse{
				Error: protocol.ErrTooManyRequests.Error(),
			})
		}

		return c.Next()
	}
}
