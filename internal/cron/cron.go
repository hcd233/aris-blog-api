package cron

import (
	"fmt"

	"go.uber.org/zap"
)

func InitCronJobs() {
	quotaCron := NewQuotaCron()
	quotaCron.Start()
}

type cronLoggerAdapter struct {
	prefix string
	logger *zap.Logger
}

func newCronLoggerAdapter(prefix string, logger *zap.Logger) cronLoggerAdapter {
	if prefix == "" {
		prefix = "[Cron]"
	}
	return cronLoggerAdapter{prefix: prefix, logger: logger}
}

func (l cronLoggerAdapter) Error(err error, msg string, keysAndValues ...interface{}) {
	zapKeyValues := []zap.Field{zap.Error(err)}
	zapKeyValues = append(zapKeyValues, convertZapKeyValues(keysAndValues...)...)
	l.logger.Error(fmt.Sprintf("[%s] %s", l.prefix, msg), zapKeyValues...)
}

func (l cronLoggerAdapter) Info(msg string, keysAndValues ...interface{}) {
	zapKeyValues := convertZapKeyValues(keysAndValues...)
	l.logger.Info(fmt.Sprintf("[%s] %s", l.prefix, msg), zapKeyValues...)
}

func convertZapKeyValues(keysAndValues ...interface{}) []zap.Field {
	if len(keysAndValues)%2 != 0 {
		panic("keysAndValues must be a slice of key-value pairs")
	}
	len := len(keysAndValues) / 2
	zapKeyValues := make([]zap.Field, 0, len)
	for i := 0; i < len; i++ {
		key, value := keysAndValues[i*2].(string), keysAndValues[i*2+1]
		zapKeyValues = append(zapKeyValues, zap.Any(key, value))
	}
	return zapKeyValues
}
