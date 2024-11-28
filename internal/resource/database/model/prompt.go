package model

import (
	"gorm.io/gorm"
)

type Task string

type Template struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

const (
	TaskContentCompletion  Task = "contentCompletion"
	TaskArticleSummary     Task = "articleSummary"
	TaskArticleTranslation Task = "articleTranslation"
	TaskArticleQA          Task = "articleQA"
	TaskTermExplaination   Task = "termExplaination"
)

type Prompt struct {
	gorm.Model
	ID        uint       `json:"id" gorm:"column:id;primary_key;auto_increment;comment:'提示词ID'"`
	Task      Task       `json:"task" gorm:"column:task;index:idx_task;not null;comment:'任务类型'"`
	Templates []Template `json:"templates" gorm:"column:templates;type:json;not null;comment:'提示词模板'"`
	Variables []string   `json:"variables" gorm:"column:variables;type:json;not null;comment:'提示词变量'"`
	Version   uint       `json:"version" gorm:"column:version;not null;comment:'版本'"`
}

func (p *Prompt) GetBasicInfo() map[string]interface{} {
	return map[string]interface{}{
		"id":      p.ID,
		"task":    p.Task,
		"version": p.Version,
	}
}

func (p *Prompt) GetDetailedInfo() map[string]interface{} {
	return map[string]interface{}{
		"id":        p.ID,
		"task":      p.Task,
		"templates": p.Templates,
		"variables": p.Variables,
		"version":   p.Version,
	}
}
