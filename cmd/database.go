package cmd

import (
	"github.com/hcd233/aris-blog-api/internal/resource/database"
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "数据库相关命令组",
	Long:  `提供一组用于管理和操作数据库的命令，包括迁移、备份和恢复等功能。`,
}

var migrateDatabaseCmd = &cobra.Command{
	Use:   "migrate",
	Short: "迁移数据库",
	Long:  `执行数据库迁移操作，将数据库结构更新到最新的模式。`,
	Run: func(cmd *cobra.Command, _ []string) {
		database.InitDatabase()
		db := database.GetDBInstance(cmd.Context())
		lo.Must0(db.AutoMigrate(model.Models...))
		lo.Must0(database.DropAllForeignKeys(cmd.Context()))
	},
}

func init() {
	databaseCmd.AddCommand(migrateDatabaseCmd)
	rootCmd.AddCommand(databaseCmd)
}