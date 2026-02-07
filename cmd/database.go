package cmd

import (
	"fmt"
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
	},
}

var dropFKCmd = &cobra.Command{
    Use:   "drop-fk",
    Short: "移除所有物理外键约束（Postgres）",
    Long:  "扫描当前数据库的所有非系统schema，删除所有外键约束，以便完全使用逻辑外键。",
    Run: func(cmd *cobra.Command, _ []string) {
        database.InitDatabase()
        db := database.GetDBInstance(cmd.Context())

        type fkRow struct {
            TableSchema string
            TableName   string
            Constraint  string
        }

        var rows []fkRow
        // 查询所有外键约束
        raw := db.Raw(`
            SELECT tc.table_schema AS table_schema,
                   tc.table_name   AS table_name,
                   tc.constraint_name AS constraint
            FROM information_schema.table_constraints tc
            WHERE tc.constraint_type = 'FOREIGN KEY'
              AND tc.table_schema NOT IN ('pg_catalog','information_schema')
        `)
        if err := raw.Scan(&rows).Error; err != nil {
            lo.Must0(err)
        }

        for _, r := range rows {
            stmt := fmt.Sprintf("ALTER TABLE \"%s\".\"%s\" DROP CONSTRAINT IF EXISTS \"%s\";", r.TableSchema, r.TableName, r.Constraint)
            if err := db.Exec(stmt).Error; err != nil {
                lo.Must0(err)
            }
        }
    },
}

func init() {
	databaseCmd.AddCommand(migrateDatabaseCmd)
	databaseCmd.AddCommand(dropFKCmd)
	rootCmd.AddCommand(databaseCmd)
}
