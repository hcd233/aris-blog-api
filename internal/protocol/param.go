package protocol

// PageParam 列表参数
//
//	@author centonhuang
//	@update 2024-09-21 09:00:57
type PageParam struct {
	Limit  int `form:"limit" binding:"required,min=1,max=50"`
	Offset int `form:"offset" binding:"gte=0"`
}

// QueryParam 查询参数
//
//	@author centonhuang
//	@update 2024-09-18 02:56:39
type QueryParam struct {
	PageParam
	Query string `form:"query" binding:"required"`
}

// GithubCallbackParam Github回调请求参数
//
//	@author centonhuang
//	@update 2024-09-18 03:14:09
type GithubCallbackParam struct {
	Code  string `form:"code" binding:"required"`
	State string `form:"state" binding:"required"`
}
