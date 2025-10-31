package dto

// PageParam 分页参数
//
//	author centonhuang
//	update 2025-10-31 05:30:00
type PageParam struct {
	Page     int `json:"page" query:"page" doc:"Page number, minimum is 1" default:"1"`
	PageSize int `json:"pageSize" query:"pageSize" doc:"Items per page, range 1-50" minimum:"1" maximum:"50" default:"10"`
}

// QueryParam 查询参数
//
//	author centonhuang
//	update 2025-10-31 05:30:00
type QueryParam struct {
	Query string `json:"query,omitempty" query:"query" doc:"Fuzzy search keyword"`
}

// CommonParam 综合分页参数
//
//	author centonhuang
//	update 2025-10-31 05:30:00
type CommonParam struct {
	PageParam
	QueryParam
}

// PageInfo 分页信息
//
//	author centonhuang
//	update 2025-10-31 05:30:00
type PageInfo struct {
	Page     int   `json:"page" doc:"Page number"`
	PageSize int   `json:"pageSize" doc:"Items per page"`
	Total    int64 `json:"total" doc:"Total items"`
}
