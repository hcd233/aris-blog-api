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
	return logger
}

func init() {
	var (
		cfg             zap.Config
		zapLevelMapping = map[string]zap.AtomicLevel{
			"DEBUG":  zap.NewAtomicLevelAt(zap.DebugLevel),
			"INFO":   zap.NewAtomicLevelAt(zap.InfoLevel),
			"WARN":   zap.NewAtomicLevelAt(zap.WarnLevel),
			"ERROR":  zap.NewAtomicLevelAt(zap.ErrorLevel),
			"DPANIC": zap.NewAtomicLevelAt(zap.DPanicLevel),
			"PANIC":  zap.NewAtomicLevelAt(zap.PanicLevel),
			"FATAL":  zap.NewAtomicLevelAt(zap.FatalLevel),
		}
	)

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

	if logLevel == zap.NewAtomicLevelAt(zap.DebugLevel) {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}
	// Set log level
	cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	cfg.EncoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	cfg.EncoderConfig.ConsoleSeparator = "  "
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewConsoleEncoder(cfg.EncoderConfig), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), logFileWriter), logLevel),
		// Error log / Panic log output to err.log
		zapcore.NewCore(zapcore.NewConsoleEncoder(cfg.EncoderConfig), zapcore.NewMultiWriteSyncer(errFileWriter), zapLevelMapping["ERROR"]),
		// PanicLog output to panic.log
		zapcore.NewCore(zapcore.NewConsoleEncoder(cfg.EncoderConfig), zapcore.NewMultiWriteSyncer(panicFileWriter), zapLevelMapping["PANIC"]),
	)

	defaultLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapLevelMapping["PANIC"]))
}
