package util

import (
	"net/http"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/samber/lo"
)

const (
	heartbeatInterval = 1 * time.Second

	heartbeatEvent = "heartbeat"
	streamEvent    = "stream"
	errorEvent     = "error"
	doneEvent      = "done"
)

// SendHTTPResponse 发送HTTP响应
//
//	param c *fiber.Ctx
//	param data interface{}
//	param err error
//	author centonhuang
//	update 2025-01-04 17:34:06
func SendHTTPResponse(c *fiber.Ctx, data interface{}, err error) {
	switch err {
	case protocol.ErrDataNotExists: // 404
		c.Status(http.StatusOK).JSON(protocol.HTTPResponse{
			Error: err.Error(),
		})
	case protocol.ErrDataExists: // 400
		c.Status(http.StatusOK).JSON(protocol.HTTPResponse{
			Error: err.Error(),
		})
	case protocol.ErrBadRequest: // 400
		c.Status(http.StatusBadRequest).JSON(protocol.HTTPResponse{
			Error: err.Error(),
		})
	case protocol.ErrInsufficientQuota: // 400
		c.Status(http.StatusBadRequest).JSON(protocol.HTTPResponse{
			Error: err.Error(),
		})
	case protocol.ErrUnauthorized: // 401
		c.Status(http.StatusUnauthorized).JSON(protocol.HTTPResponse{
			Error: err.Error(),
		})
	case protocol.ErrNoPermission: // 403
		c.Status(http.StatusForbidden).JSON(protocol.HTTPResponse{
			Error: err.Error(),
		})
	case protocol.ErrTooManyRequests: // 429
		c.Status(http.StatusTooManyRequests).JSON(protocol.HTTPResponse{
			Error: err.Error(),
		})
	case protocol.ErrInternalError: // 500
		c.Status(http.StatusInternalServerError).JSON(protocol.HTTPResponse{
			Error: err.Error(),
		})
	case protocol.ErrNoImplement: // 501
		c.Status(http.StatusNotImplemented).JSON(protocol.HTTPResponse{
			Error: err.Error(),
		})
	case nil:
		c.Status(http.StatusOK).JSON(protocol.HTTPResponse{
			Data: data,
		})
	}
}

// SendStreamEventResponses 发送流式事件响应
//
//	param c *fiber.Ctx
//	param streamChan <-chan string
//	param errChan <-chan error
//	return err error
//	author centonhuang
//	update 2024-12-09 17:18:12
func SendStreamEventResponses(c *fiber.Ctx, streamChan <-chan string, errChan <-chan error) (err error) {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	var mu sync.Mutex
	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			mu.Lock()
			event := "data: " + string(lo.Must1(sonic.Marshal(protocol.SSEResponse{
				Delta: "",
				Stop:  false,
				Error: "",
			}))) + "\n\n"
			c.Write([]byte(event))
			mu.Unlock()
		}
	}()

	for {
		select {
		case token, ok := <-streamChan:
			mu.Lock()
			if !ok {
				event := "data: " + string(lo.Must1(sonic.Marshal(protocol.SSEResponse{
					Delta: "",
					Stop:  true,
					Error: "",
				}))) + "\n\n"
				c.Write([]byte(event))
				mu.Unlock()
				return
			}
			event := "data: " + string(lo.Must1(sonic.Marshal(protocol.SSEResponse{
				Delta: token,
				Stop:  false,
				Error: "",
			}))) + "\n\n"
			c.Write([]byte(event))
			mu.Unlock()
		case err = <-errChan:
			mu.Lock()
			if err != nil {
				event := "data: " + string(lo.Must1(sonic.Marshal(protocol.SSEResponse{
					Delta: "",
					Stop:  true,
					Error: err.Error(),
				}))) + "\n\n"
				c.Write([]byte(event))
				mu.Unlock()
				return
			}
			mu.Unlock()
		}
	}
}

// WrapHTTPResponse 包装HTTP响应错误
//
//	@param rsp rspT
//	@param err error
//	@return *protocol.HumaHTTPResponse[rspT]
//	@return error
//	@author centonhuang
//	@update 2025-10-31 01:47:14
func WrapHTTPResponse[rspT any](rsp rspT, err error) (*protocol.HumaHTTPResponse[rspT], error) {
	switch err {
	case protocol.ErrDataNotExists: // 404
		return nil, huma.Error404NotFound(err.Error())
	case protocol.ErrDataExists, protocol.ErrBadRequest, protocol.ErrInsufficientQuota: // 400
		return nil, huma.Error400BadRequest(err.Error())
	case protocol.ErrUnauthorized: // 401
		return nil, huma.Error401Unauthorized(err.Error())
	case protocol.ErrNoPermission: // 403
		return nil, huma.Error403Forbidden(err.Error())
	case protocol.ErrTooManyRequests: // 429
		return nil, huma.Error429TooManyRequests(err.Error())
	case protocol.ErrInternalError: // 500
		return nil, huma.Error500InternalServerError(err.Error())
	case protocol.ErrNoImplement: // 501
		return nil, huma.Error501NotImplemented(err.Error())
	case nil:
		return &protocol.HumaHTTPResponse[rspT]{
			Body: rsp,
		}, nil
	default:
		return nil, huma.Error500InternalServerError("Unknown error: " + err.Error())
	}
}
