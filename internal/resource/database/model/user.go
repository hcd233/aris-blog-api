// Package model defines the database schema for the model.
//
//	@update 2024-06-22 09:33:43
package model

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type (
	// Permission string 权限
	//	@update 2024-09-21 01:34:29
	Permission string

	// Platform string 平台
	//	@update 2024-09-21 01:34:12
	Platform string
)

const (

	// PlatformGithub github user
	//	@update 2024-06-22 10:05:13
	PlatformGithub Platform = "github"

	// PermissionGeneral general permission
	//	@update 2024-06-22 10:05:15
	PermissionGeneral Permission = "general"
	// PermissionAdmin admin permission
	//	@update 2024-06-22 10:05:17
	PermissionAdmin Permission = "admin"
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
	Permission   Permission   `json:"permission" gorm:"column:permission;not null;default:'general';comment:权限"`
	LastLogin    sql.NullTime `json:"last_login" gorm:"column:last_login;not null;comment:最后登录时间"`
	GithubBindID string       `json:"-" gorm:"unique;comment:Github绑定ID"`
	Articles     []Article    `json:"articles" gorm:"foreignKey:UserID"`
	Categories   []Category   `json:"categories" gorm:"foreignKey:UserID"`
	Tags         []Tag        `json:"tags" gorm:"foreignKey:CreateBy"`
}

// BeforeCreate 创建用户前
//
//	@receiver u *User
//	@param _ *gorm.DB
//	@return err error
//	@update 2024-06-22 10:10:07
func (u *User) BeforeCreate(_ *gorm.DB) (err error) {
	if !u.LastLogin.Valid {
		u.LastLogin = sql.NullTime{Time: time.Now(), Valid: true}
	}
	return
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
