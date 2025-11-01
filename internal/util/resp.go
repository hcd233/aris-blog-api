package util

import (
	"net/http"
	"sync"
	"time"

	"github.com/bytedance/sonic"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/sse"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/protocol/dto"
	"github.com/samber/lo"
)

const (
	heartbeatInterval = 1 * time.Second

	heartbeatEvent = "heartbeat"
	streamEvent    = "stream"
	errorEvent     = "error"
	doneEvent      = "done"
)

// SendStreamEventResponses 发送流式事件响应
//
//	param c *fiber.Ctx
//	param streamChan <-chan string
//	param errChan <-chan error
//	return err error
//	author centonhuang
//	update 2024-12-09 17:18:12
func SendStreamEventResponses(sender sse.Sender, streamChan <-chan string, errChan <-chan error) {
	var mu sync.Mutex
	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			mu.Lock()
			sender.Data(lo.Must1(sonic.Marshal(protocol.SSEResponse{
				Delta: "",
				Stop:  false,
				Error: "",
			})))
			mu.Unlock()
		}
	}()

	for {
		select {
		case token, ok := <-streamChan:
			mu.Lock()
			if !ok {
				sender.Data(lo.Must1(sonic.Marshal(protocol.SSEResponse{
					Delta: "",
					Stop:  true,
					Error: "",
				})))
				mu.Unlock()
				return
			}
			sender.Data(lo.Must1(sonic.Marshal(protocol.SSEResponse{
				Delta: token,
				Stop:  false,
				Error: "",
			})))
			mu.Unlock()
		case err := <-errChan:
			mu.Lock()
			if err != nil {
				sender.Data(lo.Must1(sonic.Marshal(protocol.SSEResponse{
					Delta: "",
					Stop:  true,
					Error: err.Error(),
				})))
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
func WrapHTTPResponse[rspT any](rsp rspT, err error) (*protocol.HTTPResponse[rspT], huma.StatusError) {
	if statusErr := transformError(err); statusErr != nil {
		return nil, statusErr
	}
	return &protocol.HTTPResponse[rspT]{
		Body: rsp,
	}, nil
}

func transformError(err error) (statusErr huma.StatusError) {
	switch err {
	case protocol.ErrDataNotExists: // 404
		statusErr = huma.Error404NotFound(err.Error())
	case protocol.ErrDataExists, protocol.ErrBadRequest, protocol.ErrInsufficientQuota: // 400
		statusErr = huma.Error400BadRequest(err.Error())
	case protocol.ErrUnauthorized: // 401
		statusErr = huma.Error401Unauthorized(err.Error())
	case protocol.ErrNoPermission: // 403
		statusErr = huma.Error403Forbidden(err.Error())
	case protocol.ErrTooManyRequests: // 429
		statusErr = huma.Error429TooManyRequests(err.Error())
	case protocol.ErrInternalError: // 500
		statusErr = huma.Error500InternalServerError(err.Error())
	case protocol.ErrNoImplement: // 501
		statusErr = huma.Error501NotImplemented(err.Error())
	case nil:
		statusErr = nil
	default:
		statusErr = huma.Error500InternalServerError("Unknown error: " + err.Error())
	}
	return
}

// RedirectURL 重定向
//
//	@param url string
//	@return *protocol.RedirectResponse
//	@author centonhuang
//	@update 2025-11-02 04:06:07
func RedirectURL(rsp *dto.URLResponse, err error) (*protocol.RedirectResponse, huma.StatusError) {
	if statusErr := transformError(err); statusErr != nil {
		return nil, statusErr
	}
	return &protocol.RedirectResponse{
		Status: http.StatusTemporaryRedirect,
		Url:    rsp.URL,
	}, nil
}
