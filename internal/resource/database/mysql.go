// Package database provides mysql database connection.
//
//	@update 2024-06-22 09:04:46
package database

import (
	"fmt"
	"time"

	"github.com/hcd233/Aris-AI-go/internal/config"

	"github.com/samber/lo"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB undefined
//
//	@update 2024-06-22 09:07:38
var (
	DB     *gorm.DB
)

// InitDatabase function
//
//	@author centonhuang
//	@update 2024-06-22 09:23:59
func InitDatabase() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.MysqlUser, config.MysqlPassword, config.MysqlHost, config.MysqlPort, config.MysqlDatabase)

	DB = lo.Must(gorm.Open(mysql.New(mysql.Config{
		DSN:               dsn,
		DefaultStringSize: 256,
	}),
		&gorm.Config{
			DryRun:         false, // 只生成SQL不运行
			TranslateError: true,
		}))


	db := lo.Must(DB.DB())

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(5 * time.Hour)
}
