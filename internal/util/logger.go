package util

import (
	"sync"

	"go.uber.org/zap"
)

var (
	logger     *zap.Logger
	loggerOnce sync.Once
)

// GetLogger 获取日志记录器
//
//	return *zap.Logger
//	author system
//	update 2025-01-19 12:00:00
func GetLogger() *zap.Logger {
	loggerOnce.Do(func() {
		// 这里应该从配置文件读取日志配置，暂时使用默认配置
		var err error
		logger, err = zap.NewProduction()
		if err != nil {
			// 如果创建生产环境日志器失败，使用开发环境日志器
			logger, _ = zap.NewDevelopment()
		}
	})
	return logger
}