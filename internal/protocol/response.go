package protocol

// SSEResponse SSE响应
//
//	author centonhuang
//	update 2024-12-08 16:42:20
type SSEResponse struct {
	Delta string `json:"delta"`
	Stop  bool   `json:"stop"`
	Error string `json:"error,omitempty"`
}

// HTTPResponse HTTP响应
//
//	author centonhuang
//	update 2025-10-31 01:38:26
type HTTPResponse[BodyT any] struct {
	Body BodyT `json:"data"`
}

// RedirectResponse 重定向响应
//
//	@author centonhuang
//	@update 2025-11-02 04:01:39
type RedirectResponse struct {
	Status int    `json:"status" doc:"Status code"`
	// Cookie string `json:"cookie" header:"Set-Cookie" doc:"Cookie for redirect"`
	URL    string `json:"url" header:"Location" doc:"URL for redirect"`
}
