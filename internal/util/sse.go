package util

import (
	"encoding/json"
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
			c.SSEvent(heartbeatEvent, string(lo.Must1(json.Marshal(protocol.AIStreamResponse{
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
				c.SSEvent(doneEvent, string(lo.Must1(json.Marshal(protocol.AIStreamResponse{
					Delta: "",
					Stop:  true,
					Error: "",
				}))))
				c.Writer.Flush()
				mu.Unlock()
				return
			}
			c.SSEvent(streamEvent, string(lo.Must1(json.Marshal(protocol.AIStreamResponse{
				Delta: token,
				Stop:  false,
				Error: "",
			}))))
			c.Writer.Flush()
			mu.Unlock()
		case err = <-errChan:
			mu.Lock()
			if err != nil {
				c.SSEvent(errorEvent, string(lo.Must1(json.Marshal(protocol.AIStreamResponse{
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
