// Package callback 日志回调处理器
package callback

import (
	"context"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/schema"
	"github.com/hcd233/aris-blog-api/internal/logger"
	"go.uber.org/zap"
)

type logCallbackHandler struct {
	logger *zap.Logger
}

// OnStart(ctx context.Context, info *RunInfo, input CallbackInput) context.Context
// OnEnd(ctx context.Context, info *RunInfo, output CallbackOutput) context.Context

// OnError(ctx context.Context, info *RunInfo, err error) context.Context

// OnStartWithStreamInput(ctx context.Context, info *RunInfo,
// 	input *schema.StreamReader[CallbackInput]) context.Context
// OnEndWithStreamOutput(ctx context.Context, info *RunInfo,
// 	output *schema.StreamReader[CallbackOutput]) context.Context

// OnStart 开始事件
//
//	@receiver l *LogCallbackHandler
//	@param ctx context.Context
//	@param runInfo *callbacks.RunInfo
//	@param input CallbackInput
//	@return context.Context
//	@author centonhuang
//	@update 2025-02-09 22:23:04
func (l *logCallbackHandler) OnStart(ctx context.Context, runInfo *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
	l.logger.Info("[LogCallbackHandler] OnStart", zap.Any("runInfo", runInfo), zap.Any("input", input))
	return ctx
}

// OnEnd 结束事件
//
//	@receiver l *LogCallbackHandler
//	@param ctx context.Context
//	@param runInfo *callbacks.RunInfo
//	@param output *model.CallbackOutput
//	@return context.Context
//	@author centonhuang
//	@update 2025-02-09 22:23:04
func (l *logCallbackHandler) OnEnd(ctx context.Context, runInfo *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
	l.logger.Info("[LogCallbackHandler] OnEnd", zap.Any("runInfo", runInfo), zap.Any("output", output))
	return ctx
}

// OnStartWithStreamInput 流式开始事件
//
//	@receiver l *LogCallbackHandler
//	@param ctx context.Context
//	@param runInfo *callbacks.RunInfo
//	@param input *schema.StreamReader[callbacks.CallbackInput]
//	@return context.Context
//	@author centonhuang
//	@update 2025-02-09 22:23:04
func (l *logCallbackHandler) OnStartWithStreamInput(ctx context.Context, runInfo *callbacks.RunInfo, input *schema.StreamReader[callbacks.CallbackInput]) context.Context {
	l.logger.Info("[LogCallbackHandler] OnStartWithStreamInput", zap.Any("runInfo", runInfo), zap.Any("input", input))
	return ctx
}

// OnEndWithStreamOutput 流式结束事件
//
//	@receiver l *LogCallbackHandler
//	@param ctx context.Context
//	@param runInfo *callbacks.RunInfo
//	@param output *schema.StreamReader[callbacks.CallbackOutput]
//	@return context.Context
//	@author centonhuang
//	@update 2025-02-09 22:23:04
func (l *logCallbackHandler) OnEndWithStreamOutput(ctx context.Context, runInfo *callbacks.RunInfo, output *schema.StreamReader[callbacks.CallbackOutput]) context.Context {
	l.logger.Info("[LogCallbackHandler] OnEndWithStreamOutput", zap.Any("runInfo", runInfo), zap.Any("output", output))
	return ctx
}

// OnError 错误事件
//
//	@receiver l *LogCallbackHandler
//	@param ctx context.Context
//	@param runInfo *callbacks.RunInfo
//	@param err error
//	@return context.Context
//	@author centonhuang
//	@update 2025-02-09 22:23:04
func (l *logCallbackHandler) OnError(ctx context.Context, runInfo *callbacks.RunInfo, err error) context.Context {
	l.logger.Error("[LogCallbackHandler] OnError", zap.Any("runInfo", runInfo), zap.Error(err))
	return ctx
}

// NewLogCallbackHandler 创建eino日志回调处理器
//
//	@return *callbacksHelper.ModelCallbackHandler
//	@author centonhuang
//	@update 2025-02-09 22:10:02
func NewLogCallbackHandler() callbacks.Handler {
	return &logCallbackHandler{
		logger: logger.Logger(),
	}
}
