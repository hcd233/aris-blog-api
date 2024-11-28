package dao

import (
	"strings"

	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"github.com/hcd233/Aris-blog/internal/util"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

// PromptDAO 提示词DAO
//
//	@author centonhuang
//	@update 2024-10-23 05:22:38
type PromptDAO struct {
	baseDAO[model.Prompt]
}

func (dao *PromptDAO) Create(db *gorm.DB, prompt *model.Prompt) error {
	contents := lo.Map(prompt.Templates, func(tmplate model.Template, idx int) string {
		return tmplate.Content
	})

	content := strings.Join(contents, "\n")

	prompt.Variables = util.ExtractVariablesFromContent(content)

	return db.Create(prompt).Error
}

func (dao *PromptDAO) GetLatestPromptByTask(db *gorm.DB, task model.Task, fields, preloads []string) (prompt *model.Prompt, err error) {
	sql := db.Select(fields)
	for _, preload := range preloads {
		sql = sql.Preload(preload)
	}
	err = sql.Where(&model.Prompt{Task: task}).Last(&prompt).Error
	return
}

func (dao *PromptDAO) PaginateByTask(db *gorm.DB, task model.Task, fields, preloads []string, page, pageSize int) (prompts []*model.Prompt, pageInfo *PageInfo, err error) {
	limit, offset := pageSize, (page-1)*pageSize
	sql := db.Select(fields)
	for _, preload := range preloads {
		sql = sql.Preload(preload)
	}
	err = sql.Where(&model.Prompt{Task: task}).Offset(offset).Limit(limit).Find(&prompts).Error

	if err != nil {
		return
	}

	pageInfo = &PageInfo{
		Page:     page,
		PageSize: pageSize,
	}
	err = db.Model(&prompts).Where(&model.Prompt{Task: task}).Count(&pageInfo.Total).Error
	return
}
