package protocol

// GithubCallbackParam Github回调请求参数
//
//	author centonhuang
//	update 2024-09-18 03:14:09
type GithubCallbackParam struct {
	Code  string `form:"code" binding:"required"`
	State string `form:"state" binding:"required"`
}

// QQCallbackParam QQ回调请求参数
type QQCallbackParam struct {
	Code  string `form:"code" binding:"required"`
	State string `form:"state" binding:"required"`
}

// OAuth2CallbackParam 通用OAuth2回调请求参数
type OAuth2CallbackParam struct {
	Code  string `form:"code" binding:"required"`
	State string `form:"state" binding:"required"`
}

// PageParam 列表参数
//
//	author centonhuang
//	update 2024-09-21 09:00:57
type PageParam struct {
	Page     int `form:"page" binding:"required,gte=1"`
	PageSize int `form:"pageSize" binding:"min=1,max=50"`
}

// QueryParam 查询参数
//
//	author centonhuang
//	update 2024-09-18 02:56:39
type QueryParam struct {
	Query string `form:"query"`
}

// PaginateParam 分页查询参数
//
//	@author centonhuang
//	@update 2025-08-25 12:30:17
type PaginateParam struct {
	*PageParam
	*QueryParam
}

// ArticleParam 文章参数
//
//	author centonhuang
//	update 2024-09-21 09:59:55
type ArticleParam struct {
	ArticleSlug string `form:"articleSlug" binding:"required"`
	Author      string `form:"author" binding:"required"`
}

// ImageParam 图片参数
//
//	author centonhuang
//	update 2024-12-08 16:42:00
type ImageParam struct {
	Quality string `form:"quality" binding:"required,oneof=raw thumb"`
}
