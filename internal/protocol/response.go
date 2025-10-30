package protocol

// HTTPResponse 标准响应体
//
//	author centonhuang
//	update 2024-09-16 03:41:34
type HTTPResponse struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error,omitempty"`
}

// SSEResponse SSE响应
//
//	author centonhuang
//	update 2024-12-08 16:42:20
type SSEResponse struct {
	Delta string `json:"delta"`
	Stop  bool   `json:"stop"`
	Error string `json:"error,omitempty"`
}

// HumaHTTPResponse HTTP响应
//
//	author centonhuang
//	update 2025-10-31 01:38:26
type HumaHTTPResponse[BodyT any] struct {
	Body BodyT `json:"data"`
}
