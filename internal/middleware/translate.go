package middleware

import (
	"bytes"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/logger"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	return r.body.Write(b)
}

// TranslateMiddleware 翻译状态码中间件
//
//	@return gin.HandlerFunc
//	@author centonhuang
//	@update 2024-09-16 03:27:12
func TranslateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		c.Next()

		response := protocol.Response{}
		if w.body.String() == "" {
			c.JSON(c.Writer.Status(), protocol.Response{
				Code: protocol.CodeUnknownError,
			})
		}
		if w.Header().Get("Content-Type") != "application/json; charset=utf-8" {
			c.Status(c.Writer.Status())
			return
		}
		err := json.Unmarshal(w.body.Bytes(), &response)
		if err != nil {
			logger.Logger.Error("[TranslateMiddleware]", zap.Error(err))
			response = protocol.Response{
				Code: protocol.CodeUnknownError,
			}
		}

		code := response.Code
		message := protocol.CodeMessageMapping[code]
		appendMessage := response.Message
		if appendMessage != "" {
			appendMessage = ": " + appendMessage
		}
		response.Message = message + appendMessage

		translatedResponse := lo.Must1(json.Marshal(response))
		w.ResponseWriter.Write(translatedResponse)
		w.body.Reset()

		c.JSON(c.Writer.Status(), translatedResponse)
	}
}
