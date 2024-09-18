package protocol

// Response 标准响应体
//
//	@author centonhuang
//	@update 2024-09-16 03:41:34
type Response struct {
	Message string                 `json:"message"`
	Code    ResponseCode           `json:"code"`
	Data    map[string]interface{} `json:"data"`
}
