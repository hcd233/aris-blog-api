package handler

import "context"

import "github.com/hcd233/aris-blog-api/internal/constant"

// UserIDFromCtx 从上下文中解析用户 ID
func UserIDFromCtx(ctx context.Context) (uint, bool) {
	if value := ctx.Value(constant.CtxKeyUserID); value != nil {
		if userID, ok := value.(uint); ok {
			return userID, true
		}
	}
	return 0, false
}
