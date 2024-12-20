// Package model defines the database schema for the model.
//
//	@update 2024-06-22 09:33:43
package model

import (
	"database/sql"

	"gorm.io/gorm"
)

type (
	// Permission string 权限
	//	@update 2024-09-21 01:34:29
	Permission string

	// PermissionLevel int8 权限等级
	//	@update 2024-09-21 01:34:29

	PermissionLevel int8

	Quota int8

	// Platform string 平台
	//	@update 2024-09-21 01:34:12
	Platform string
)

const (

	// PlatformGithub github user
	//	@update 2024-06-22 10:05:13
	PlatformGithub Platform = "github"

	// PermissionReader general permission
	//	@update 2024-06-22 10:05:15
	PermissionReader Permission = "reader"

	// PermissionCreator creator permission
	//	@update 2024-06-22 10:05:17
	PermissionCreator Permission = "creator"

	// PermissionAdmin admin permission
	//	@update 2024-06-22 10:05:17
	PermissionAdmin Permission = "admin"

	QuotaReader Quota = 5

	QuotaCreator Quota = 30

	QuotaAdmin Quota = 120
)

// PermissionLevelMapping 权限等级映射
//
//	@update 2024-09-21 01:34:29
var (
	PermissionLevelMapping = map[Permission]int8{
		PermissionReader:  1,
		PermissionCreator: 2,
		PermissionAdmin:   3,
	}

	PermissionQuotaMapping = map[Permission]Quota{
		PermissionReader:  QuotaReader,
		PermissionCreator: QuotaCreator,
		PermissionAdmin:   QuotaAdmin,
	}
)

// User 用户数据库模型
//
//	@author centonhuang
//	@update 2024-06-22 09:36:22
type User struct {
	gorm.Model
	ID           uint         `json:"id" gorm:"column:id;primary_key;auto_increment;comment:用户ID"`
	Name         string       `json:"name" gorm:"column:name;unique;not null;comment:用户名"`
	Email        string       `json:"email" gorm:"column:email;unique;not null;comment:邮箱"`
	Avatar       string       `json:"avatar" gorm:"column:avatar;not null;comment:头像"`
	Permission   Permission   `json:"permission" gorm:"column:permission;not null;default:'reader';comment:权限"`
	LastLogin    sql.NullTime `json:"last_login" gorm:"column:last_login;comment:最后登录时间"`
	GithubBindID string       `json:"-" gorm:"unique;comment:Github绑定ID"`
	LLMQuota     Quota        `json:"llm_quota" gorm:"column:llm_quota;not null;default:0;comment:LLM配额"`
	Articles     []Article    `json:"articles" gorm:"foreignKey:UserID"`
	Categories   []Category   `json:"categories" gorm:"foreignKey:UserID"`
	Tags         []Tag        `json:"tags" gorm:"foreignKey:CreatedBy"`
}

// GetBasicInfo 获取用户基本信息
//
//	@receiver u *User
//	@return map
//	@author centonhuang
//	@update 2024-09-18 03:47:14
func (u *User) GetBasicInfo() map[string]interface{} {
	return map[string]interface{}{
		"id":     u.ID,
		"name":   u.Name,
		"avatar": u.Avatar,
	}
}

// GetDetailedInfo 获取用户详细信息
//
//	@receiver u *User
//	@return map
//	@author centonhuang
//	@update 2024-09-18 03:50:04
func (u *User) GetDetailedInfo() map[string]interface{} {
	return map[string]interface{}{
		"id":         u.ID,
		"name":       u.Name,
		"email":      u.Email,
		"avatar":     u.Avatar,
		"createdAt":  u.CreatedAt,
		"lastLogin":  u.LastLogin.Time,
		"permission": u.Permission,
	}
}
