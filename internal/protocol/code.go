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

	// CodeGetTagError ResponseCode 获取标签错误
	//	@update 2024-09-22 03:39:36
	CodeGetTagError ResponseCode = 1018

	// CodeCreateTagError ResponseCode 创建标签错误
	//	@update 2024-09-22 03:39:19
	CodeCreateTagError ResponseCode = 1019

	// CodeUpdateTagError ResponseCode 更新标签错误
	//	@update 2024-09-22 03:39:19
	CodeUpdateTagError ResponseCode = 1020

	// CodeDeleteTagError ResponseCode 删除标签错误
	//	@update 2024-09-22 03:39:19
	CodeDeleteTagError ResponseCode = 1021

	// CodeQueryTagError ResponseCode 查询标签错误
	//	@update 2024-09-22 05:19:14
	CodeQueryTagError ResponseCode = 1022

	// CodeCreateCategoryError ResponseCode 创建分类错误
	//	@update 2024-09-28 07:11:12
	CodeCreateCategoryError ResponseCode = 1023

	// CodeGetCategoryError ResponseCode 获取分类错误
	//	@update 2024-10-01 03:57:27
	CodeGetCategoryError ResponseCode = 1024

	// CodeUpdateCategoryError ResponseCode 更新分类错误
	//	@update 2024-10-02 04:05:46
	CodeUpdateCategoryError ResponseCode = 1025

	// CodeDeleteCategoryError ResponseCode 删除分类错误
	//	@update 2024-10-02 04:25:05
	CodeDeleteCategoryError ResponseCode = 1026

	// CodeGetArticleVersionError ResponseCode 获取文章版本错误
	//	@update 2024-10-16 11:54:22
	CodeGetArticleVersionError ResponseCode = 1027

	// CodeCreateArticleVersionError ResponseCode 创建文章版本错误
	//	@update 2024-10-17 12:46:33
	CodeCreateArticleVersionError ResponseCode = 1028

	// CodeCreateArticleVersionRateLimitError ResponseCode 达到创建文章版本频率限制错误
	//	@update 2024-10-17 01:10:33
	CodeCreateArticleVersionRateLimitError ResponseCode = 1029

	// CodeQueryArticleError ResponseCode 查询文章错误
	//	@update 2024-10-23 12:05:24
	CodeQueryArticleError ResponseCode = 1030

	// CodeGetCommentError ResponseCode 获取评论错误
	//	@update 2024-10-23 07:16:33
	CodeGetCommentError ResponseCode = 1031

	// CodeCreateCommentError ResponseCode 创建评论错误
	//	@update 2024-10-24 04:45:21
	CodeCreateCommentError ResponseCode = 1032

	// CodeDeleteCommentError ResponseCode 删除评论错误
	//	@update 2024-10-24 07:09:39
	CodeDeleteCommentError ResponseCode = 1033

	// CodeCreateCommentRateLimitError ResponseCode 达到创建评论频率限制错误
	//	@update 2024-10-24 05:52:47
	CodeCreateCommentRateLimitError ResponseCode = 1034

	// CodeLikeTagError ResponseCode 点赞/取消点赞标签错误
	//	@update 2024-10-30 04:16:22
	CodeLikeTagError ResponseCode = 1035

	// CodeLikeArticleError ResponseCode 点赞/取消点赞文章错误
	//	@update 2024-10-30 04:16:43
	CodeLikeArticleError ResponseCode = 1036

	// CodeLikeCommentError ResponseCode 点赞/取消点赞评论错误
	//	@update 2024-10-30 04:16:39
	CodeLikeCommentError ResponseCode = 1037

	// CodeLikeArticleRateLimitError ResponseCode 达到创建文章版本频率限制错误
	//	@update 2024-10-30 06:00:03
	CodeLikeArticleRateLimitError ResponseCode = 1038

	// CodeLikeCommentRateLimitError ResponseCode 达到创建评论频率限制错误
	//	@update 2024-10-30 06:00:15
	CodeLikeCommentRateLimitError ResponseCode = 1039

	// CodeLikeTagRateLimitError ResponseCode 达到创建标签频率限制错误
	//	@update 2024-10-30 06:00:21
	CodeLikeTagRateLimitError ResponseCode = 1040

	// CodeUpdateUserError ResponseCode 更新用户错误
	//	@update 2024-10-30 09:54:53
	CodeUpdateUserError ResponseCode = 1041

	// CodeGetUserLikeError ResponseCode 获取用户点赞信息错误
	//	@update 2024-11-03 07:10:30
	CodeGetUserLikeError ResponseCode = 1042

	// CodeUnknownError ResponseCode 未知错误
	//	@update 2024-09-21 08:22:14
	CodeUnknownError ResponseCode = 10000
)

// CodeMessageMapping 响应码消息映射
//
//	@update 2024-09-16 04:12:17
var CodeMessageMapping = map[ResponseCode]string{
	CodeOk:               "请求成功",
	CodeUnauthorized:     "访问鉴权接口未提供鉴权信息",
	CodeTokenVerifyError: "鉴权信息校验失败",
	CodeStateError:       "登录状态错误",
	CodeTokenError:       "登录令牌错误",

	CodeNotPermissionError: "无操作资源权限",

	CodeParamError: "查询参数错误",
	CodeBodyError:  "请求体错误",
	CodeURIError:   "路径参数错误",

	CodeGetUserError:      "获取用户信息错误",
	CodeUpdateUserError:   "更新用户信息错误",
	CodeUserNotFoundError: "用户未找到",
	CodeQueryUserError:    "查询用户错误",

	CodeGetArticleError:    "获取文章错误",
	CodeCreateArticleError: "创建文章错误",
	CodeUpdateArticleError: "更新文章错误",
	CodeDeleteArticleError: "删除文章错误",
	CodeQueryArticleError:  "查询文章错误",

	CodeGetTagError:    "获取标签错误",
	CodeCreateTagError: "创建标签错误",
	CodeUpdateTagError: "更新标签错误",
	CodeDeleteTagError: "删除标签错误",
	CodeQueryTagError:  "查询标签错误",

	CodeCreateCategoryError: "创建分类错误",
	CodeGetCategoryError:    "获取分类错误",
	CodeUpdateCategoryError: "更新分类错误",
	CodeDeleteCategoryError: "删除分类错误",

	CodeGetArticleVersionError:             "获取文章版本错误",
	CodeCreateArticleVersionError:          "创建文章版本错误",
	CodeCreateArticleVersionRateLimitError: "达到创建文章版本频率限制，请稍后再试",

	CodeGetCommentError:             "获取评论错误",
	CodeCreateCommentError:          "创建评论错误",
	CodeDeleteCommentError:          "删除评论错误",
	CodeCreateCommentRateLimitError: "达到创建评论频率限制，请稍后再试",

	CodeLikeTagError:     "点赞/取消点赞标签错误",
	CodeLikeArticleError: "点赞/取消点赞文章错误",
	CodeLikeCommentError: "点赞/取消点赞评论错误",

	CodeLikeArticleRateLimitError: "达到点赞文章频率限制，请稍后再试",
	CodeLikeCommentRateLimitError: "达到点赞评论频率限制，请稍后再试",
	CodeLikeTagRateLimitError:     "达到点赞标签频率限制，请稍后再试",

	CodeGetUserLikeError: "获取用户点赞信息错误",

	CodeUnknownError: "未知错误",
}
