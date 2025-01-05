// Package logger provides a logger that can be used throughout the application.
package logger

import (
	"os"
	"strings"

	"github.com/hcd233/Aris-blog/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger undefined 全局日志
//
//	update 2024-09-16 12:47:59
var Logger *zap.Logger

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
		Filename:   strings.Join([]string{config.LogDirPath, "aris-ai-go.log"}, "/"),
		MaxSize:    100, // MB
		MaxBackups: 3,
		MaxAge:     7, // days
		Compress:   false,
	})

	// error logger
	errFileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   config.LogDirPath + "/aris-ai-go-error.log",
		MaxSize:    500, // MB
		MaxBackups: 3,
		MaxAge:     30, // days
		Compress:   false,
	})

	// panic logger
	panicFileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   config.LogDirPath + "/aris-ai-go-panic.log",
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

	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapLevelMapping["PANIC"]))
}
