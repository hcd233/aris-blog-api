package dao

import (
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
	"gorm.io/gorm"
)

// PromptDAO 提示词DAO
//
//	author centonhuang
//	update 2024-10-23 05:22:38
type PromptDAO struct {
	baseDAO[model.Prompt]
}

func (dao *PromptDAO) GetLatestPromptByTask(db *gorm.DB, task model.Task, fields, preloads []string) (prompt *model.Prompt, err error) {
	sql := db.Select(fields)
	for _, preload := range preloads {
		sql = sql.Preload(preload)
	}
	err = sql.Where(&model.Prompt{Task: task}).Last(&prompt).Error
	return
}

func (dao *PromptDAO) PaginateByTask(db *gorm.DB, task model.Task, fields, preloads []string, param *PaginateParam) (prompts []*model.Prompt, pageInfo *PageInfo, err error) {
	limit, offset := param.PageSize, (param.Page-1)*param.PageSize

	sql := db.Select(fields)
	for _, preload := range preloads {
		sql = sql.Preload(preload)
	}

	// 添加模糊查询支持
	if param.Query != "" && len(param.QueryFields) > 0 {
		sql = sql.Where("? LIKE ?", param.QueryFields[0], "%"+param.Query+"%")
		for _, field := range param.QueryFields[1:] {
			sql = sql.Or("? LIKE ?", field, "%"+param.Query+"%")
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

func (dao *PromptDAO) GetPromptByTaskAndVersion(db *gorm.DB, task model.Task, version uint, fields, preloads []string) (prompt *model.Prompt, err error) {
	sql := db.Select(fields)
	for _, preload := range preloads {
		sql = sql.Preload(preload)
	}

	err = sql.Where(&model.Prompt{Task: task, Version: version}).First(&prompt).Error
	return
}
