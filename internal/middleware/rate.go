package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/cache"
	"github.com/samber/lo"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	redis_store "github.com/ulule/limiter/v3/drivers/store/redis"
)

// RateLimiterMiddleware 限频中间件
//
//	@param period time.Duration
//	@param limit int64
//	@param key string
//	@param errCode protocol.ResponseCode
//	@return gin.HandlerFunc
//	@author centonhuang
//	@update 2024-10-22 05:01:51
func RateLimiterMiddleware(period time.Duration, limit int64, prefix, key string, errorCode protocol.ResponseCode) gin.HandlerFunc {
	// 创建限频规则
	rate := limiter.Rate{
		Period: period,
		Limit:  limit,
	}

	redis := cache.GetRedisClient()
	// 使用内存存储限频数据
	store := lo.Must1(redis_store.NewStoreWithOptions(redis, limiter.StoreOptions{
		Prefix: prefix,
	}))

	// 创建限频实例
	instance := limiter.New(store, rate)

	// 创建中间件
	middleware := mgin.NewMiddleware(instance,
		mgin.WithLimitReachedHandler(func(c *gin.Context) {
			c.JSON(http.StatusTooManyRequests, protocol.Response{
				Code: errorCode,
			})
			c.Abort()
		}),
		mgin.WithKeyGetter(func(c *gin.Context) string {
			return c.MustGet("limiter").(string)
		}),
	)

	return func(c *gin.Context) {
		// 获取限频 key
		value := c.MustGet(key)

		if key == "" {
			value = c.ClientIP() // 如果没有指定的参数，则使用 IP 地址作为 key
		}

		// 设置限频 key
		c.Set("limiter", fmt.Sprintf("%s:%v", key, value))

		// 应用限频中间件
		middleware(c)

		// 继续处理请求
		c.Next()
	}
}
