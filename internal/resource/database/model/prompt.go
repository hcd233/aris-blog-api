package model

import (
	"gorm.io/gorm"
)

// Task 任务类型
//
//	@author centonhuang
//	@update 2024-12-09 16:13:30
type Task string

// Template 提示词模板
//
//	@author centonhuang
//	@update 2024-12-09 16:13:30
type Template struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

const (

	// TaskContentCompletion Task 内容补全
	//	@update 2024-12-09 16:13:42
	TaskContentCompletion Task = "contentCompletion"

	// TaskArticleSummary Task 文章摘要
	//	@update 2024-12-09 16:13:42
	TaskArticleSummary Task = "articleSummary"

	// TaskArticleTranslation Task 文章翻译
	//	@update 2024-12-09 16:13:42
	TaskArticleTranslation Task = "articleTranslation"

	// TaskArticleQA Task 文章问答
	//	@update 2024-12-09 16:13:42
	TaskArticleQA Task = "articleQA"

	// TaskTermExplaination Task 术语解释
	//	@update 2024-12-09 16:13:42
	TaskTermExplaination Task = "termExplaination"
)

// Prompt 提示词
//
//	@author centonhuang
//	@update 2024-12-09 16:13:42
type Prompt struct {
	gorm.Model
	ID        uint       `json:"id" gorm:"column:id;primary_key;auto_increment;comment:'提示词ID'"`
	Task      Task       `json:"task" gorm:"column:task;uniqueIndex:idx_task_version;not null;comment:'任务类型'"`
	Templates []Template `json:"templates" gorm:"column:templates;type:json;not null;serializer:json;comment:'提示词模板'"`
	Variables []string   `json:"variables" gorm:"column:variables;type:json;not null;serializer:json;comment:'提示词变量'"`
	Version   uint       `json:"version" gorm:"column:version;uniqueIndex:idx_task_version;not null;comment:'版本'"`
}

// GetBasicInfo 获取提示词基本信息
//
//	@receiver p *Prompt
//	@return map[string]interface{}
//	@author centonhuang
//	@update 2024-12-09 16:13:42
func (p *Prompt) GetBasicInfo() map[string]interface{} {
	return map[string]interface{}{
		"id":        p.ID,
		"createdAt": p.CreatedAt,
		"task":      p.Task,
		"version":   p.Version,
	}
}

// GetDetailedInfo 获取提示词详细信息
//
//	@receiver p *Prompt
//	@return map[string]interface{}
//	@author centonhuang
//	@update 2024-12-09 16:13:42
func (p *Prompt) GetDetailedInfo() map[string]interface{} {
	return map[string]interface{}{
		"id":        p.ID,
		"task":      p.Task,
		"templates": p.Templates,
		"variables": p.Variables,
		"version":   p.Version,
	}
}
