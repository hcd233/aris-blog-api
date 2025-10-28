package middleware

import (
	"strings"

	"github.com/danielgtaylor/huma/v2"
)

// HumaJWTMiddleware 为 Huma 路由创建 JWT 中间件适配器
func HumaJWTMiddleware(api huma.API) {
	// 添加全局中间件来处理认证
	api.UseMiddleware(func(ctx huma.Context, next func(huma.Context)) {
		// 检查是否需要认证的路由
		path := ctx.Operation().Path
		
		// 定义需要认证的路径模式
		authRequiredPaths := []string{
			"/v1/user/current",
			"/v1/user",
		}
		
		needsAuth := false
		for _, authPath := range authRequiredPaths {
			if path == authPath || strings.HasPrefix(path, "/v1/user/{userID}") {
				needsAuth = true
				break
			}
		}
		
		if needsAuth {
			// 从请求头获取 Authorization
			auth := ctx.Header("Authorization")
			if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
				huma.WriteErr(api, ctx, 401, "缺少有效的认证头")
				return
			}
			
			// 这里应该验证 JWT token
			// 为了演示，我们暂时使用模拟的用户ID
			// 在实际应用中，你需要解析和验证 JWT token
			token := strings.TrimPrefix(auth, "Bearer ")
			if token == "test-token" { // 简单的测试令牌
				mockUserID := uint(1)
				// 将用户ID存储到上下文中，供 handler 使用
				ctx = huma.WithValue(ctx, "userID", mockUserID)
			} else {
				huma.WriteErr(api, ctx, 401, "无效的认证令牌")
				return
			}
		}
		
		next(ctx)
	})
}