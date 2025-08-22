package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-blog-api/internal/logger"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/util"
	"go.uber.org/zap"
)

// ValidateURIMiddleware 验证URI中间件
//
//	param uri interface{}
//	return fiber.Handler
//	author centonhuang
//	update 2024-09-21 07:47:53
func ValidateURIMiddleware(uri interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := c.ParamsParser(uri); err != nil {
			logger.WithFCtx(c).Info("[ValidateURIMiddleware] failed to bind uri", zap.Error(err))
			util.SendHTTPResponse(c, nil, protocol.ErrBadRequest)
			return c.Status(fiber.StatusBadRequest).JSON(protocol.HTTPResponse{
				Error: protocol.ErrBadRequest.Error(),
			})
		}
		c.Locals("uri", uri)
		return c.Next()
	}
}

// ValidateParamMiddleware 验证参数中间件
//
//	param param interface{}
//	return fiber.Handler
//	author centonhuang
//	update 2024-09-21 07:48:40
func ValidateParamMiddleware(param interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := c.QueryParser(param); err != nil {
			logger.WithFCtx(c).Info("[ValidateParamMiddleware] failed to bind param", zap.Error(err))
			util.SendHTTPResponse(c, nil, protocol.ErrBadRequest)
			return c.Status(fiber.StatusBadRequest).JSON(protocol.HTTPResponse{
				Error: protocol.ErrBadRequest.Error(),
			})
		}
		c.Locals("param", param)
		return c.Next()
	}
}

// ValidateBodyMiddleware 验证请求体中间件
//
//	param body interface{}
//	return fiber.Handler
//	author centonhuang
//	update 2024-09-21 08:48:25
func ValidateBodyMiddleware(body interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := c.BodyParser(body); err != nil {
			logger.WithFCtx(c).Info("[ValidateBodyMiddleware] failed to bind body", zap.Error(err))
			util.SendHTTPResponse(c, nil, protocol.ErrBadRequest)
			return c.Status(fiber.StatusBadRequest).JSON(protocol.HTTPResponse{
				Error: protocol.ErrBadRequest.Error(),
			})
		}
		c.Locals("body", body)
		return c.Next()
	}
}
