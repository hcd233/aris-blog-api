// Package database 存储中间件
//
//	@update 2024-06-22 09:04:46
package database

import (
	"fmt"
	"time"

	"github.com/hcd233/Aris-blog/internal/config"
	"github.com/hcd233/Aris-blog/internal/logger"
	"go.uber.org/zap"

	"github.com/samber/lo"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
		}))

	sqlDB := lo.Must(db.DB())

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(5 * time.Hour)

	logger.Logger.Info("[Database] Connected to MySQL database", zap.String("host", config.MysqlHost), zap.String("port", config.MysqlPort), zap.String("database", config.MysqlDatabase))
}
