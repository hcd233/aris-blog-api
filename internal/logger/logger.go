// Package logger provides a logger that can be used throughout the application.
package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path"
	"strings"

	"github.com/hcd233/aris-blog-api/internal/config"
	"github.com/hcd233/aris-blog-api/internal/constant"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger undefined 全局日志
//
//	update 2024-09-16 12:47:59
var defaultLogger *zap.Logger

const (
	infoLogFile  = "aris-blog-api.log"
	errLogFile   = "aris-blog-api-error.log"
	panicLogFile = "aris-blog-api-panic.log"

	logLevelDebug  = "DEBUG"
	logLevelInfo   = "INFO"
	logLevelWarn   = "WARN"
	logLevelError  = "ERROR"
	logLevelDPanic = "DPANIC"
	logLevelPanic  = "PANIC"
	logLevelFatal  = "FATAL"

	timeKey       = "timestamp"
	levelKey      = "level"
	nameKey       = "logger"
	callerKey     = "caller"
	messageKey    = "message"
	stacktraceKey = "stacktrace"
)

// indentedJSONEncoder 带缩进的JSON编码器
type indentedJSONEncoder struct {
	zapcore.Encoder
}

// newIndentedJSONEncoder 创建带缩进的JSON编码器
func newIndentedJSONEncoder(cfg zapcore.EncoderConfig) zapcore.Encoder {
	return &indentedJSONEncoder{
		Encoder: zapcore.NewJSONEncoder(cfg),
	}
}

// EncodeEntry 重写编码方法，添加JSON缩进
func (enc *indentedJSONEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	// 使用原始编码器获取JSON
	buf, err := enc.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		return nil, err
	}

	// 获取原始JSON字节
	jsonBytes := buf.Bytes()

	// 移除最后的换行符
	if len(jsonBytes) > 0 && jsonBytes[len(jsonBytes)-1] == '\n' {
		jsonBytes = jsonBytes[:len(jsonBytes)-1]
	}

	// 创建缩进的JSON
	var indentedBuf bytes.Buffer
	if err := json.Indent(&indentedBuf, jsonBytes, "", "  "); err != nil {
		// 如果缩进失败，返回原始内容
		return buf, nil
	}

	// 添加换行符
	indentedBuf.WriteByte('\n')

	// 创建新的buffer并写入缩进后的内容
	newBuf := buffer.NewPool().Get()
	newBuf.Write(indentedBuf.Bytes())

	return newBuf, nil
}

// Clone 克隆编码器
func (enc *indentedJSONEncoder) Clone() zapcore.Encoder {
	return &indentedJSONEncoder{
		Encoder: enc.Encoder.Clone(),
	}
}

func Logger() *zap.Logger {
	return defaultLogger
}

func LoggerWithContext(ctx context.Context) *zap.Logger {
	logger := defaultLogger
	if traceID := ctx.Value(constant.CtxKeyTraceID); traceID != nil {
		logger = logger.With(zap.String(constant.CtxKeyTraceID, traceID.(string)))
	}
	if userID := ctx.Value(constant.CtxKeyUserID); userID != nil {
		logger = logger.With(zap.Uint(constant.CtxKeyUserID, userID.(uint)))
	}
	return logger
}

func init() {
	zapLevelMapping := map[string]zap.AtomicLevel{
		logLevelDebug:  zap.NewAtomicLevelAt(zap.DebugLevel),
		logLevelInfo:   zap.NewAtomicLevelAt(zap.InfoLevel),
		logLevelWarn:   zap.NewAtomicLevelAt(zap.WarnLevel),
		logLevelError:  zap.NewAtomicLevelAt(zap.ErrorLevel),
		logLevelDPanic: zap.NewAtomicLevelAt(zap.DPanicLevel),
		logLevelPanic:  zap.NewAtomicLevelAt(zap.PanicLevel),
		logLevelFatal:  zap.NewAtomicLevelAt(zap.FatalLevel),
	}

	logLevel, ok := zapLevelMapping[strings.ToUpper(config.LogLevel)]
	if !ok {
		logLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	// general logger
	logFileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   path.Join(config.LogDirPath, infoLogFile),
		MaxSize:    100, // MB
		MaxBackups: 3,
		MaxAge:     7, // days
		Compress:   false,
	})

	// error logger
	errFileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   path.Join(config.LogDirPath, errLogFile),
		MaxSize:    500, // MB
		MaxBackups: 3,
		MaxAge:     30, // days
		Compress:   false,
	})

	// panic logger
	panicFileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   path.Join(config.LogDirPath, panicLogFile),
		MaxSize:    500, // MB
		MaxBackups: 3,
		MaxAge:     30, // days
		Compress:   false,
	})

	// 配置结构化日志编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        timeKey,
		LevelKey:       levelKey,
		NameKey:        nameKey,
		CallerKey:      callerKey,
		MessageKey:     messageKey,
		StacktraceKey:  stacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 控制台输出使用彩色编码器（开发模式）
	consoleEncoderConfig := encoderConfig
	if logLevel == zap.NewAtomicLevelAt(zap.DebugLevel) {
		consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		consoleEncoderConfig.ConsoleSeparator = "  "
	}

	core := zapcore.NewTee(
		// 控制台输出 - 根据调试模式选择编码器
		func() zapcore.Core {
			if logLevel == zap.NewAtomicLevelAt(zap.DebugLevel) {
				return zapcore.NewCore(
					zapcore.NewConsoleEncoder(consoleEncoderConfig),
					zapcore.AddSync(os.Stdout),
					logLevel,
				)
			}
			// 生产模式使用带缩进的JSON编码器
			return zapcore.NewCore(
				newIndentedJSONEncoder(encoderConfig),
				zapcore.AddSync(os.Stdout),
				logLevel,
			)
		}(),
		// 文件输出 - 统一使用JSON编码器（不缩进，节省空间）
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.NewMultiWriteSyncer(logFileWriter),
			logLevel,
		),
		// Error log 输出到 err.log
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.NewMultiWriteSyncer(errFileWriter),
			zapLevelMapping[logLevelError],
		),
		// Panic log 输出到 panic.log
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.NewMultiWriteSyncer(panicFileWriter),
			zapLevelMapping[logLevelPanic],
		),
	)

	defaultLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapLevelMapping[logLevelPanic]))
}
