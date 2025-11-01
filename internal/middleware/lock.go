package middleware

import (
	"fmt"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/hcd233/aris-blog-api/internal/api"
	"github.com/hcd233/aris-blog-api/internal/logger"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/resource/cache"
	"github.com/hcd233/aris-blog-api/internal/util"
	"go.uber.org/zap"
)

// RedisLockMiddleware Redis锁中间件
//	@param serviceName string 
//	@param key string 
//	@param expire time.Duration 
//	@return ctx huma.Context 
//	@return next func(huma.Context) 
//	@return func(ctx huma.Context, next func(huma.Context)) 
//	@author centonhuang 
//	@update 2025-11-02 04:16:37 
func RedisLockMiddleware(serviceName, key string, expire time.Duration) func(ctx huma.Context, next func(huma.Context)) {
	redis := cache.GetRedisClient()

	return func(ctx huma.Context, next func(huma.Context)) {
		value := ctx.Context().Value(key)

		lockKey := fmt.Sprintf("%s:%s:%v", serviceName, key, value)
		lockValue := uuid.New().String()

		success, err := redis.SetNX(ctx.Context(), lockKey, lockValue, expire).Result()
		if err != nil {
			logger.WithCtx(ctx.Context()).Error("[RedisLockMiddleware] failed to get lock", zap.Error(err))
			_, err := util.WrapHTTPResponse[any](nil, protocol.ErrInternalError)
			huma.WriteErr(api.GetHumaAPI(), ctx, err.GetStatus(), err.Error(), err)
			return
		}

		if !success {
			lockValue, err = redis.Get(ctx.Context(), lockKey).Result()
			if err != nil {
				logger.WithCtx(ctx.Context()).Error("[RedisLockMiddleware] failed to get lock info",
					zap.String("lockKey", lockKey), zap.Error(err))
				_, err := util.WrapHTTPResponse[any](nil, protocol.ErrInternalError)
				huma.WriteErr(api.GetHumaAPI(), ctx, err.GetStatus(), err.Error(), err)
				return
			}
			logger.WithCtx(ctx.Context()).Info("[RedisLockMiddleware] resource is locked",
				zap.String("lockKey", lockKey), zap.String("lockValue", lockValue))
			_, err := util.WrapHTTPResponse[any](nil, protocol.ErrTooManyRequests)
			huma.WriteErr(api.GetHumaAPI(), ctx, err.GetStatus(), err.Error(), err)
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
			logger.WithCtx(ctx.Context()).Error("[RedisLockMiddleware] failed to release lock",
				zap.String("lockKey", lockKey), zap.Error(err))
		}
	}
}
