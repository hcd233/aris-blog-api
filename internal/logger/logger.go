// Package logger provides a logger that can be used throughout the application.
package logger

import (
	"context"
	"os"
	"path"
	"strings"

	"github.com/hcd233/aris-blog-api/internal/config"
	"github.com/hcd233/aris-blog-api/internal/constant"
	"go.uber.org/zap"
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
	if userName := ctx.Value(constant.CtxKeyUserName); userName != nil {
		logger = logger.With(zap.String(constant.CtxKeyUserName, userName.(string)))
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

	// 配置文件输出的JSON结构化日志编码器
	jsonEncoderConfig := zapcore.EncoderConfig{
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

	// 配置控制台输出的彩色日志编码器
	consoleEncoderConfig := zapcore.EncoderConfig{
		TimeKey:          timeKey,
		LevelKey:         levelKey,
		NameKey:          nameKey,
		CallerKey:        callerKey,
		MessageKey:       messageKey,
		StacktraceKey:    stacktraceKey,
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      zapcore.CapitalColorLevelEncoder, // 彩色级别编码
		EncodeTime:       zapcore.RFC3339TimeEncoder,
		EncodeDuration:   zapcore.SecondsDurationEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder,
		ConsoleSeparator: "  ", // 控制台分隔符
	}

	core := zapcore.NewTee(
		// 控制台输出 - 始终使用彩色Console编码器
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(consoleEncoderConfig),
			zapcore.AddSync(os.Stdout),
			logLevel,
		),
		// 文件输出 - 统一使用JSON编码器
		zapcore.NewCore(
			zapcore.NewJSONEncoder(jsonEncoderConfig),
			zapcore.NewMultiWriteSyncer(logFileWriter),
			logLevel,
		),
		// Error log 输出到 err.log
		zapcore.NewCore(
			zapcore.NewJSONEncoder(jsonEncoderConfig),
			zapcore.NewMultiWriteSyncer(errFileWriter),
			zapLevelMapping[logLevelError],
		),
		// Panic log 输出到 panic.log
		zapcore.NewCore(
			zapcore.NewJSONEncoder(jsonEncoderConfig),
			zapcore.NewMultiWriteSyncer(panicFileWriter),
			zapLevelMapping[logLevelPanic],
		),
	)

	defaultLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapLevelMapping[logLevelPanic]))
}
