package dto

// PageParam 分页参数
//
//	author centonhuang
//	update 2025-10-31 05:30:00
type PageParam struct {
	Page     int `json:"page" doc:"页码，最小为 1" minimum:"1" default:"1"`
	PageSize int `json:"pageSize" doc:"每页条数，范围 1-50" minimum:"1" maximum:"50" default:"10"`
}

// QueryParam 查询参数
//
//	author centonhuang
//	update 2025-10-31 05:30:00
type QueryParam struct {
	Query string `json:"query,omitempty" doc:"模糊搜索关键字"`
}

// PaginateParam 综合分页参数
//
//	author centonhuang
//	update 2025-10-31 05:30:00
type PaginateParam struct {
	*PageParam  `json:"page"`
	*QueryParam `json:"query,omitempty"`
}

// PaginationQuery Huma 查询参数
//
//	author centonhuang
//	update 2025-10-31 05:30:00
type PaginationQuery struct {
	Page     int    `query:"page" doc:"页码，最小为 1" minimum:"1" default:"1"`
	PageSize int    `query:"pageSize" doc:"每页条数，范围 1-50" minimum:"1" maximum:"50" default:"10"`
	Query    string `query:"query" doc:"模糊搜索关键字"`
}

// ToPaginateParam 转换为通用分页参数
func (q PaginationQuery) ToPaginateParam() *PaginateParam {
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

	return &PaginateParam{
		PageParam: &PageParam{
			Page:     page,
			PageSize: pageSize,
		},
		QueryParam: &QueryParam{
			Query: q.Query,
		},
	}
}

// PageInfo 分页信息
//
//	author centonhuang
//	update 2025-10-31 05:30:00
type PageInfo struct {
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
	Total    int64 `json:"total"`
}
