package dao

import (
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// PromptDAO 提示词DAO
//
//	author centonhuang
//	update 2024-10-23 05:22:38
type PromptDAO struct {
	baseDAO[model.Prompt]
}

// GetLatestPromptByTask 获取最新提示词
//
//	receiver dao *PromptDAO
//	param db
//	param task
//	param fields
//	param preloads
//	return prompt
//	return err
//	author centonhuang
//	update 2025-08-25 14:17:49
func (dao *PromptDAO) GetLatestPromptByTask(db *gorm.DB, task model.Task, fields, preloads []string) (prompt *model.Prompt, err error) {
	sql := db.Select(fields)
	for _, preload := range preloads {
		sql = sql.Preload(preload)
	}
	err = sql.Where(&model.Prompt{Task: task}).Last(&prompt).Error
	return
}

// PaginateByTask 分页查询提示词
//
//	author centonhuang
//	update 2024-10-23 05:22:38
//	receiver dao *PromptDAO
//	param db
//	param task
//	param fields
//	param preloads
//	param param
//	return prompts
//	return pageInfo
//	return err
//	author centonhuang
//	update 2025-08-25 14:17:57
func (dao *PromptDAO) PaginateByTask(db *gorm.DB, task model.Task, fields, preloads []string, param *CommonParam) (prompts []*model.Prompt, pageInfo *PageInfo, err error) {
	limit, offset := param.PageSize, (param.Page-1)*param.PageSize

	sql := db.Select(fields)
	for _, preload := range preloads {
		sql = sql.Preload(preload)
	}

	if param.Query != "" && len(param.QueryFields) > 0 {
		like := "%" + param.Query + "%"
		expressions := make([]clause.Expression, 0, len(param.QueryFields))
		for _, field := range param.QueryFields {
			if field == "" {
				continue
			}
			expressions = append(expressions, clause.Like{Column: clause.Column{Name: field}, Value: like})
		}

		if len(expressions) > 0 {
			sql = sql.Where(expressions[0])
			for _, expr := range expressions[1:] {
				sql = sql.Or(expr)
			}
		}
	}

	err = sql.Where(&model.Prompt{Task: task}).Offset(offset).Limit(limit).Find(&prompts).Error
	if err != nil {
		return
	}

	pageInfo = &PageInfo{
		Page:     param.Page,
		PageSize: param.PageSize,
	}

	err = db.Model(&prompts).Where(&model.Prompt{Task: task}).Count(&pageInfo.Total).Error
	return
}

// GetPromptByTaskAndVersion 获取指定任务和版本的提示词
//
//	author centonhuang
//	update 2024-10-23 05:22:38
//	receiver dao *PromptDAO
//	param db
//	param task
//	param version
//	param fields
//	param preloads
//	return prompt
//	return err
//	author centonhuang
//	update 2025-08-25 14:18:05
func (dao *PromptDAO) GetPromptByTaskAndVersion(db *gorm.DB, task model.Task, version uint, fields, preloads []string) (prompt *model.Prompt, err error) {
	sql := db.Select(fields)
	for _, preload := range preloads {
		sql = sql.Preload(preload)
	}

	err = sql.Where(&model.Prompt{Task: task, Version: version}).First(&prompt).Error
	return
}
