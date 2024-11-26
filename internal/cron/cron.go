package cron

import "go.uber.org/zap"

func InitCronJobs() {
	quotaCron := NewQuotaCron()
	quotaCron.Start()
}

type cronLoggerAdapter struct {
	logger *zap.Logger
}

func (l cronLoggerAdapter) Error(err error, msg string, keysAndValues ...interface{}) {
	l.logger.Error(msg, zap.Error(err))
}

func (l cronLoggerAdapter) Info(msg string, keysAndValues ...interface{}) {
	l.logger.Info(msg)
}
