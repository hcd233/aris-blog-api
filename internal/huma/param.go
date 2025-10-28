package humadto

// PageParam 分页参数（Huma 风格，使用 query 标签）
//
//	author centonhuang
//	update 2025-10-28 00:00:00
type PageParam struct {
	Page     int `query:"page"`
	PageSize int `query:"pageSize"`
}

// QueryParam 通用查询参数（Huma 风格）
//
//	author centonhuang
//	update 2025-10-28 00:00:00
type QueryParam struct {
	Query string `query:"query"`
}

// PaginateParam 分页查询参数（组合）
//
//	author centonhuang
//	update 2025-10-28 00:00:00
type PaginateParam struct {
	*PageParam
	*QueryParam
}

// ImageParam 图片参数（Huma 风格）
//
//	author centonhuang
//	update 2025-10-28 00:00:00
type ImageParam struct {
	Quality string `query:"quality" enum:"raw,thumb"`
}
