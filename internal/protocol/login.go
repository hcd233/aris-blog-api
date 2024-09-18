package protocol

// GithubCallbackParams Github回调请求参数
//
//	@author centonhuang
//	@update 2024-09-18 03:14:09
type GithubCallbackParams struct {
	Code  string `form:"code" binding:"required"`
	State string `form:"state" binding:"required"`
}
