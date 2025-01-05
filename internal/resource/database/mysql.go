// Package database 存储中间件
//
//	@update 2024-06-22 09:04:46
package database

import (
	"context"
	"fmt"
	"time"

	"github.com/hcd233/Aris-blog/internal/config"
	"github.com/hcd233/Aris-blog/internal/logger"
	"go.uber.org/zap"

	"github.com/samber/lo"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// DB undefined mysql数据库连接
//
//	@update 2024-09-16 01:24:51
var db *gorm.DB

// GetDBInstance 获取数据库实例
//
//	@return *gorm.DB
//	@author centonhuang
//	@update 2024-10-17 08:35:47
func GetDBInstance() *gorm.DB {
	return db
}

// InitDatabase 初始化数据库
//
//	@author centonhuang
//	@update 2024-09-22 10:04:36
func InitDatabase() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.MysqlUser, config.MysqlPassword, config.MysqlHost, config.MysqlPort, config.MysqlDatabase)

	db = lo.Must(gorm.Open(mysql.New(mysql.Config{
		DSN:               dsn,
		DefaultStringSize: 256,
	}),
		&gorm.Config{
			DryRun:         false, // 只生成SQL不运行
			TranslateError: true,
			Logger: &GormLoggerAdapter{
				ZapLogger: logger.Logger,
				LogLevel:  4, // Info级别
			},
		}))

	sqlDB := lo.Must(db.DB())

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(5 * time.Hour)

	logger.Logger.Info("[Database] Connected to MySQL database", zap.String("host", config.MysqlHost), zap.String("port", config.MysqlPort), zap.String("database", config.MysqlDatabase))
}

// GormLoggerAdapter 实现gorm的logger接口,使用zap输出SQL日志
//
//	@author centonhuang
//	@update 2025-01-05 21:10:18
type GormLoggerAdapter struct {
	ZapLogger *zap.Logger
	LogLevel  int
}

// LogMode 设置日志级别
//
//	@receiver l *GormLogger
//	@param level gormlogger.LogLevel
//	@return gormlogger.Interface
//	@author centonhuang
//	@update 2025-01-05 21:10:15
func (l *GormLoggerAdapter) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.LogLevel = int(level)
	return &newLogger
}

// Info 打印info级别的日志
//
//	@receiver l *GormLogger
//	@param _ context.Context
//	@param msg string
//	@param data ...interface{}
//	@author centonhuang
//	@update 2025-01-05 21:11:07
func (l *GormLoggerAdapter) Info(_ context.Context, msg string, data ...interface{}) {
	l.ZapLogger.Info(fmt.Sprintf(msg, data...))
}

// Warn 打印warn级别的日志
//
//	@receiver l *GormLogger
//	@param _ context.Context
//	@param msg string
//	@param data ...interface{}
//	@author centonhuang
//	@update 2025-01-05 21:11:08
func (l *GormLoggerAdapter) Warn(_ context.Context, msg string, data ...interface{}) {
	l.ZapLogger.Warn(fmt.Sprintf(msg, data...))
}

// Error 打印error级别的日志
//
//	@receiver l *GormLogger
//	@param _ context.Context
//	@param msg string
//	@param data ...interface{}
//	@author centonhuang
//	@update 2025-01-05 21:11:10
func (l *GormLoggerAdapter) Error(_ context.Context, msg string, data ...interface{}) {
	l.ZapLogger.Error(fmt.Sprintf(msg, data...))
}

// Trace 打印trace级别的日志
//
//	@receiver l *GormLogger
//	@param _ context.Context
//	@param begin time.Time
//	@param fc func() (string, int64)
//	@param err error
//	@author centonhuang
//	@update 2025-01-05 21:11:11
func (l *GormLoggerAdapter) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := []zap.Field{
		zap.String("sql", sql),
		zap.Int64("rows", rows),
		zap.String("elapsed", elapsed.String()),
	}
	if err != nil {
		fields = append(fields, zap.Error(err))
		l.ZapLogger.Error("[SQL]", fields...)
		return
	}

	l.ZapLogger.Info("[SQL]", fields...)
}
