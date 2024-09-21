// Package protocol API协议
//
//	@update 2024-09-18 02:33:08
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

	// CodeStateError ResponseCode 登录状态错误
	//	@update 2024-09-16 03:49:36
	CodeStateError ResponseCode = 1003

	// CodeTokenError ResponseCode 登录令牌错误
	//	@update 2024-09-16 03:49:57
	CodeTokenError ResponseCode = 1004

	// CodeGetUserError ResponseCode 获取用户错误
	//	@update 2024-09-16 03:52:59
	CodeGetUserError ResponseCode = 1005

	// CodeUserNotFoundError ResponseCode 用户未找到
	//	@update 2024-09-16 06:12:45
	CodeUserNotFoundError ResponseCode = 1006

	// CodeURIError ResponseCode 路径参数错误
	//	@update 2024-09-16 06:37:45
	CodeURIError ResponseCode = 1007

	// CodeQueryUserError ResponseCode 查询用户错误
	//	@update 2024-09-17 08:41:08
	CodeQueryUserError ResponseCode = 1009

	// CodeParamError ResponseCode 查询参数错误
	//	@update 2024-09-17 08:47:30
	CodeParamError ResponseCode = 1010

	// CodeBodyError ResponseCode 请求体错误
	//	@update 2024-09-18 03:19:10
	CodeBodyError ResponseCode = 1011

	// CodeNotPermissionError ResponseCode 无操作资源权限
	//	@update 2024-09-18 04:02:51
	CodeNotPermissionError ResponseCode = 1012

	// CodeGetArticleError ResponseCode 获取文章错误
	//	@update 2024-09-21 09:11:39
	CodeGetArticleError ResponseCode = 1013

	// CodeCreateArticleError ResponseCode 创建文章错误
	//	@update 2024-09-21 10:23:54
	CodeCreateArticleError ResponseCode = 1015

	// CodeUpdateArticleError ResponseCode 更新文章错误
	//	@update 2024-09-22 04:09:25
	CodeUpdateArticleError ResponseCode = 1016

	// CodeDeleteArticleError ResponseCode 删除文章错误
	//	@update 2024-09-22 04:09:25
	CodeDeleteArticleError ResponseCode = 1017

	// CodeUnknownError ResponseCode 未知错误
	//	@update 2024-09-21 08:22:14
	CodeUnknownError ResponseCode = 10000
)

// CodeMessageMapping 响应码消息映射
//
//	@update 2024-09-16 04:12:17
var CodeMessageMapping = map[ResponseCode]string{
	CodeOk:                 "请求成功",
	CodeUnauthorized:       "访问鉴权接口未提供鉴权信息",
	CodeTokenVerifyError:   "鉴权信息校验失败",
	CodeStateError:         "登录状态错误",
	CodeTokenError:         "登录令牌错误",
	CodeGetUserError:       "获取用户信息错误",
	CodeUserNotFoundError:  "用户未找到",
	CodeURIError:           "路径参数错误",
	CodeQueryUserError:     "查询用户错误",
	CodeParamError:         "查询参数错误",
	CodeBodyError:          "请求体错误",
	CodeNotPermissionError: "无操作资源权限",
	CodeGetArticleError:    "获取文章错误",
	CodeCreateArticleError: "创建文章错误",
	CodeUpdateArticleError: "更新文章错误",
	CodeDeleteArticleError: "删除文章错误",
	CodeUnknownError:       "未知错误",
}
