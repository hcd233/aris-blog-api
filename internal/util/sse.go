package util

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/samber/lo"
)

func SendStreamEventResponses(c *gin.Context, streamChan <-chan string, errChan <-chan error) (err error) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	for {
		select {
		case token, ok := <-streamChan:
			if !ok {
				c.SSEvent("done", string(lo.Must1(json.Marshal(protocol.AIStreamResponse{
					Delta: "",
					Stop:  true,
					Error: "",
				}))))
				c.Writer.Flush()

				return
			}
			c.SSEvent("stream", string(lo.Must1(json.Marshal(protocol.AIStreamResponse{
				Delta: token,
				Stop:  false,
				Error: "",
			}))))
			c.Writer.Flush()
		case err = <-errChan:
			if err != nil {
				c.SSEvent("error", string(lo.Must1(json.Marshal(protocol.AIStreamResponse{
					Delta: "",
					Stop:  true,
					Error: err.Error(),
				}))))
				c.Writer.Flush()
				return
			}
		}
	}
}
