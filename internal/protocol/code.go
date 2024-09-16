package protocol

// ResponseCode 响应码
//
//	@author centonhuang
//	@update 2024-09-16 03:45:50
type ResponseCode int

const (

	// CodeOk responseCode 响应成功
	//	@update 2024-09-16 03:41:50
	CodeOk ResponseCode = iota // 0

	// CodeStateError ResponseCode 状态错误
	//	@update 2024-09-16 03:49:36
	CodeStateError ResponseCode = 1001 // 1001

	// CodeTokenError ResponseCode 令牌错误
	//	@update 2024-09-16 03:49:57
	CodeTokenError ResponseCode = 1002 // 1002

	// CodeGetUserError ResponseCode 获取用户错误
	//	@update 2024-09-16 03:52:59
	CodeGetUserError ResponseCode = 1003 // 1003

)

// CodeMessageMapping 响应码消息映射
//
//	@update 2024-09-16 04:12:17
var CodeMessageMapping = map[ResponseCode]string{
	CodeStateError:   "登录状态错误",
	CodeTokenError:   "登录令牌错误",
	CodeGetUserError: "获取用户信息错误",
}
