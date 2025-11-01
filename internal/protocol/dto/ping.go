package dto

// PingResponse 健康检查响应
//
//	@author centonhuang
//	@update 2025-11-02 05:21:02
type PingResponse struct {
	Status string `json:"status" doc:"Ping status"`
}
