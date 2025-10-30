package handler

import (
	"context"

	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/protocol"
)

// PaginationQuery 通用分页查询参数
//
//	@author centonhuang
//	@update 2025-10-31 04:20:00
type PaginationQuery struct {
	Page     int    `query:"page" doc:"页码，最小为 1" minimum:"1" default:"1"`
	PageSize int    `query:"pageSize" doc:"每页条数，范围 1-50" minimum:"1" maximum:"50" default:"10"`
	Query    string `query:"query" doc:"模糊搜索关键字"`
}

// ToPaginateParam 转换为协议层的分页参数
//
//	receiver q PaginationQuery
//	return *protocol.PaginateParam
//	author centonhuang
//	update 2025-10-31 04:20:00
func (q PaginationQuery) ToPaginateParam() *protocol.PaginateParam {
	page := q.Page
	if page <= 0 {
		page = 1
	}

	pageSize := q.PageSize
	switch {
	case pageSize <= 0:
		pageSize = 10
	case pageSize > 50:
		pageSize = 50
	}

	return &protocol.PaginateParam{
		PageParam: &protocol.PageParam{
			Page:     page,
			PageSize: pageSize,
		},
		QueryParam: &protocol.QueryParam{
			Query: q.Query,
		},
	}
}

// UserIDFromCtx 从上下文中解析用户 ID
//
//	param ctx context.Context
//	return uint
//	return bool
//	author centonhuang
//	update 2025-10-31 04:20:00
func UserIDFromCtx(ctx context.Context) (uint, bool) {
	if value := ctx.Value(constant.CtxKeyUserID); value != nil {
		if userID, ok := value.(uint); ok {
			return userID, true
		}
	}
	return 0, false
}
