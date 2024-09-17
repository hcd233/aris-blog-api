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

	// CodeUnauthorized ResponseCode 未授权
	//	@update 2024-09-16 05:32:38
	CodeUnauthorized ResponseCode = 1001

	// CodeTokenVerifyError ResponseCode 鉴权错误
	//	@update 2024-09-16 03:49:16
	CodeTokenVerifyError ResponseCode = 1002

	// CodeStateError ResponseCode 状态错误
	//	@update 2024-09-16 03:49:36
	CodeStateError ResponseCode = 1003

	// CodeTokenError ResponseCode 令牌错误
	//	@update 2024-09-16 03:49:57
	CodeTokenError ResponseCode = 1004

	// CodeGetUserError ResponseCode 获取用户错误
	//	@update 2024-09-16 03:52:59
	CodeGetUserError ResponseCode = 1005

	// CodeUserNotFoundError ResponseCode 用户未找到
	//	@update 2024-09-16 06:12:45
	CodeUserNotFoundError ResponseCode = 1006

	// CodeParamError ResponseCode 参数错误
	//	@update 2024-09-16 06:37:45
	CodeParamError ResponseCode = 1007

	// CodeRouterError ResponseCode 路由错误
	//	@update 2024-09-16 06:37:45
	CodeRouterError ResponseCode = 1008

	// CodeQueryUserError ResponseCode 查询用户错误
	//	@update 2024-09-17 08:41:08
	CodeQueryUserError ResponseCode = 1009

	// CodeInvalidQueryError ResponseCode 非法查询参数错误
	//	@update 2024-09-17 08:47:30
	CodeInvalidQueryError ResponseCode = 1010
)

// CodeMessageMapping 响应码消息映射
//
//	@update 2024-09-16 04:12:17
var CodeMessageMapping = map[ResponseCode]string{
	CodeUnauthorized:      "访问鉴权接口未提供鉴权信息",
	CodeTokenVerifyError:  "鉴权信息校验失败",
	CodeStateError:        "登录状态错误",
	CodeTokenError:        "登录令牌错误",
	CodeGetUserError:      "获取用户信息错误",
	CodeUserNotFoundError: "用户未找到",
	CodeParamError:        "参数错误",
	CodeRouterError:       "路由错误",
	CodeQueryUserError:    "查询用户错误",
	CodeInvalidQueryError: "非法查询参数错误",
}
