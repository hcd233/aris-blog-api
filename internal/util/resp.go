package util

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
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
//	@param c *gin.Context
//	@param data interface{}
//	@param err error
//	@author centonhuang
//	@update 2025-01-04 17:34:06
func SendHTTPResponse(c *gin.Context, data interface{}, err error) {
	switch err {
	case protocol.ErrDataNotExists: // 404
		c.JSON(http.StatusOK, protocol.HTTPResponse{
			Error: err.Error(),
		})
	case protocol.ErrDataExists: // 400
		c.JSON(http.StatusOK, protocol.HTTPResponse{
			Error: err.Error(),
		})
	case protocol.ErrBadRequest: // 400
		c.JSON(http.StatusBadRequest, protocol.HTTPResponse{
			Error: err.Error(),
		})
	case protocol.ErrInsufficientQuota: // 400
		c.JSON(http.StatusBadRequest, protocol.HTTPResponse{
			Error: err.Error(),
		})
	case protocol.ErrUnauthorized: // 401
		c.JSON(http.StatusUnauthorized, protocol.HTTPResponse{
			Error: err.Error(),
		})
	case protocol.ErrNoPermission: // 403
		c.JSON(http.StatusForbidden, protocol.HTTPResponse{
			Error: err.Error(),
		})
	case protocol.ErrTooManyRequests: // 429
		c.JSON(http.StatusTooManyRequests, protocol.HTTPResponse{
			Error: err.Error(),
		})
	case protocol.ErrInternalError: // 500
		c.JSON(http.StatusInternalServerError, protocol.HTTPResponse{
			Error: err.Error(),
		})
	case nil:
		c.JSON(http.StatusOK, protocol.HTTPResponse{
			Data: data,
		})
	}
}

// SendStreamEventResponses 发送流式事件响应
//
//	@param c *gin.Context
//	@param streamChan <-chan string
//	@param errChan <-chan error
//	@return err error
//	@author centonhuang
//	@update 2024-12-09 17:18:12
func SendStreamEventResponses(c *gin.Context, streamChan <-chan string, errChan <-chan error) (err error) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Writer.Flush()

	var mu sync.Mutex
	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			mu.Lock()
			c.SSEvent(heartbeatEvent, string(lo.Must1(json.Marshal(protocol.SSEResponse{
				Delta: "",
				Stop:  false,
				Error: "",
			}))))
			c.Writer.Flush()
			mu.Unlock()
		}
	}()

	for {
		select {
		case token, ok := <-streamChan:
			mu.Lock()
			if !ok {
				c.SSEvent(doneEvent, string(lo.Must1(json.Marshal(protocol.SSEResponse{
					Delta: "",
					Stop:  true,
					Error: "",
				}))))
				c.Writer.Flush()
				mu.Unlock()
				return
			}
			c.SSEvent(streamEvent, string(lo.Must1(json.Marshal(protocol.SSEResponse{
				Delta: token,
				Stop:  false,
				Error: "",
			}))))
			c.Writer.Flush()
			mu.Unlock()
		case err = <-errChan:
			mu.Lock()
			if err != nil {
				c.SSEvent(errorEvent, string(lo.Must1(json.Marshal(protocol.SSEResponse{
					Delta: "",
					Stop:  true,
					Error: err.Error(),
				}))))
				c.Writer.Flush()
				mu.Unlock()
				return
			}
			mu.Unlock()
		}
	}
}
