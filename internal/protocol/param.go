package protocol

// PageParam 列表参数
//
//	@author centonhuang
//	@update 2024-09-21 09:00:57
type PageParam struct {
	Page     int `form:"page" binding:"required,gte=1"`
	PageSize int `form:"pageSize" binding:"min=1,max=50"`
}

// QueryParam 查询参数
//
//	@author centonhuang
//	@update 2024-09-18 02:56:39
type QueryParam struct {
	PageParam
	Query  string   `form:"query" binding:"required,min=2"`
	Filter []string `form:"filter"`
}

// GithubCallbackParam Github回调请求参数
//
//	@author centonhuang
//	@update 2024-09-18 03:14:09
type GithubCallbackParam struct {
	Code  string `form:"code" binding:"required"`
	State string `form:"state" binding:"required"`
}
