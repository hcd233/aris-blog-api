package database

import (
    "context"
    "fmt"
    "strings"
)

// DropAllForeignKeys 删除当前 schema 下所有表的物理外键约束，实现逻辑外键
//
//	param ctx context.Context
//	return error
//	author centonhuang
//	update 2025-01-19 20:20:00
func DropAllForeignKeys(ctx context.Context) error {
    gdb := db.WithContext(ctx)

    // 查询当前 schema 下的所有外键约束名称与表名
    type row struct {
        ConstraintName string
        TableName      string
        SchemaName     string
    }

    rows, err := gdb.Raw(`
        SELECT con.conname AS constraint_name,
               rel.relname AS table_name,
               nsp.nspname AS schema_name
        FROM pg_constraint con
        JOIN pg_class rel ON rel.oid = con.conrelid
        JOIN pg_namespace nsp ON nsp.oid = con.connamespace
        WHERE con.contype = 'f' AND nsp.nspname = current_schema();
    `).Rows()
    if err != nil {
        return err
    }
    defer rows.Close()

    drops := make([]string, 0, 16)
    for rows.Next() {
        var r row
        if err := rows.Scan(&r.ConstraintName, &r.TableName, &r.SchemaName); err != nil {
            return err
        }
        // 生成 drop 语句，使用双引号保证大小写与特殊字符
        drops = append(drops, fmt.Sprintf(`ALTER TABLE "%s"."%s" DROP CONSTRAINT IF EXISTS "%s";`, r.SchemaName, r.TableName, r.ConstraintName))
    }

    if len(drops) == 0 {
        return nil
    }

    // 合并为一条批量语句以减少往返
    batch := strings.Join(drops, "\n")
    return gdb.Exec(batch).Error
}

package database

import (
    "context"
    "fmt"
    "strings"

    "github.com/samber/lo"
)

// DropAllForeignKeys 删除当前 schema 下所有表的物理外键约束，实现逻辑外键
//
//	param ctx context.Context
//	return error
//	author centonhuang
//	update 2025-01-19 20:20:00
func DropAllForeignKeys(ctx context.Context) error {
    gdb := db.WithContext(ctx)

    // 查询当前 schema 下的所有外键约束名称与表名
    type row struct {
        ConstraintName string
        TableName      string
        SchemaName     string
    }

    rows, err := gdb.Raw(`
        SELECT con.conname AS constraint_name,
               rel.relname AS table_name,
               nsp.nspname AS schema_name
        FROM pg_constraint con
        JOIN pg_class rel ON rel.oid = con.conrelid
        JOIN pg_namespace nsp ON nsp.oid = con.connamespace
        WHERE con.contype = 'f' AND nsp.nspname = current_schema();
    `).Rows()
    if err != nil {
        return err
    }
    defer rows.Close()

    drops := make([]string, 0, 16)
    for rows.Next() {
        var r row
        if err := rows.Scan(&r.ConstraintName, &r.TableName, &r.SchemaName); err != nil {
            return err
        }
        // 生成 drop 语句，使用双引号保证大小写与特殊字符
        drops = append(drops, fmt.Sprintf(`ALTER TABLE "%s"."%s" DROP CONSTRAINT IF EXISTS "%s";`, r.SchemaName, r.TableName, r.ConstraintName))
    }

    if len(drops) == 0 {
        return nil
    }

    // 合并为一条批量语句以减少往返
    batch := strings.Join(drops, "\n")
    return gdb.Exec(batch).Error
}

