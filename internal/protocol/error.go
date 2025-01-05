// Package protocol API协议
//
//	@update 2024-09-18 02:33:08
package protocol

import "errors"

var (

	// ErrInternalError 内部错误
	//
	//	@update 2025-01-04 17:35:44
	ErrInternalError = errors.New("InternalError")

	// ErrUnauthorized 未授权错误
	//
	//	@update 2025-01-04 17:36:00
	ErrUnauthorized = errors.New("Unauthorized")

	// ErrNoPermission 没有权限错误
	//
	//	@update 2025-01-04 17:36:00
	ErrNoPermission = errors.New("NoPermission")

	// ErrDataNotExists 数据不存在错误
	//
	//	@update 2025-01-04 17:36:00
	ErrDataNotExists = errors.New("DataNotExists")

	// ErrDataExists 数据已存在错误
	//
	//	@update 2025-01-04 17:36:00
	ErrDataExists = errors.New("DataExists")

	// ErrTooManyRequests 请求过于频繁错误
	//
	//	@update 2025-01-04 17:36:00
	ErrTooManyRequests = errors.New("TooManyRequests")

	// ErrBadRequest 请求错误
	//
	//	@update 2025-01-04 17:36:00
	ErrBadRequest = errors.New("BadRequest")

	// ErrInsufficientQuota 配额不足错误
	//
	//	@update 2025-01-05 18:41:32
	ErrInsufficientQuota = errors.New("InsufficientQuota")
)
