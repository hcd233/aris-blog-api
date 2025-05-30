package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/logger"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/resource/cache"
	"github.com/hcd233/aris-blog-api/internal/util"
	"github.com/samber/lo"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	redis_store "github.com/ulule/limiter/v3/drivers/store/redis"
	"go.uber.org/zap"
)

// RateLimiterMiddleware 限频中间件
//
//	param serviceName string
//	param key string
//	param period time.Duration
//	param limit int64
//	return gin.HandlerFunc
//	author centonhuang
//	update 2025-01-05 15:06:44
func RateLimiterMiddleware(serviceName, key string, period time.Duration, limit int64) gin.HandlerFunc {
	// 创建限频规则
	rate := limiter.Rate{
		Period: period,
		Limit:  limit,
	}

	redis := cache.GetRedisClient()
	// 使用Redis存储限频数据
	store := lo.Must1(redis_store.NewStoreWithOptions(redis, limiter.StoreOptions{
		Prefix: serviceName,
	}))

	// 创建限频实例
	instance := limiter.New(store, rate)

	// 创建中间件
	middleware := mgin.NewMiddleware(instance,
		mgin.WithLimitReachedHandler(func(c *gin.Context) {
			fields := []zap.Field{zap.String("serviceName", serviceName)}

			if key == "" {
				fields = append(fields, zap.String("key", "ip"), zap.String("value", c.ClientIP()))
			} else {
				fields = append(fields, zap.String("key", key), zap.String("value", c.GetString(key)))
			}

			logger.LoggerWithContext(c).Error("[RateLimiterMiddleware] rate limit reached", fields...)
			util.SendHTTPResponse(c, nil, protocol.ErrTooManyRequests)
		}),
		mgin.WithKeyGetter(func(c *gin.Context) string {
			return c.MustGet(constant.CtxKeyLimiter).(string)
		}),
	)

	return func(c *gin.Context) {
		var key, value string
		if key == "" {
			key = "ip"
			value = c.ClientIP() // 如果没有指定的参数，则使用 IP 地址作为 key
		} else {
			value = c.GetString(key)
		}

		// 设置限频 key
		c.Set(constant.CtxKeyLimiter, fmt.Sprintf("%s:%v", key, value))

		// 应用限频中间件
		middleware(c)

		// 继续处理请求
		c.Next()
	}
}
