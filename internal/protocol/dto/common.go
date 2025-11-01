package dto

// EmptyRequest 空请求
//
//	@author centonhuang
//	@update 2025-10-31 02:32:07
type EmptyRequest struct{}

// EmptyResponse 空响应
//
//	author centonhuang
//	update 2025-01-05 15:33:11
type EmptyResponse struct{}

// URLResponse 链接响应
//
//	@author centonhuang
//	@update 2025-11-02 04:11:50
type URLResponse struct {
	URL string `json:"url" doc:"URL"`
}
