// Package model defines the database schema for the model.
//
//	@update 2024-06-22 09:33:43
package model

import (
	"database/sql"
	"errors"
	"time"

	"github.com/hcd233/Aris-AI-go/internal/resource/database"
	"gorm.io/gorm"
)

const (

	// PlatformGithub github user
	//	@update 2024-06-22 10:05:13
	PlatformGithub = "github"

	// PermissionGeneral general permission
	//	@update 2024-06-22 10:05:15
	PermissionGeneral = "general"
	// PermissionAdmin admin permission
	//	@update 2024-06-22 10:05:17
	PermissionAdmin = "admin"
)

// User Model Schema
//
//	@author centonhuang
//	@update 2024-06-22 09:36:22
type User struct {
	gorm.Model
	Username   string       `json:"username" gorm:"column:username;not null"`
	Platform   string       `json:"platform" gorm:"column:platform;not null"`
	Avatar     string       `json:"avatar" gorm:"column:avatar;not null"`
	Permission string       `json:"permission" gorm:"column:permission;not null"`
	LastLogin  sql.NullTime `json:"last_login" gorm:"column:last_login;not null"`

	GithubBindID string `gorm:"unique" json:"-"`
}

// Create 创建用户
//
//	@receiver u *User
//	@return error
//	@author centonhuang
//	@update 2024-06-22 10:10:07
func (u *User) Create() error {
	result := database.DB.Create(u)
	return result.Error
}

// UpdateUserInfo 更新用户
//
//	@receiver u *User
//	@return error
//	@author centonhuang
//	@update 2024-06-22 10:24:05
func (u *User) UpdateUserInfo() error {
	result := database.DB.Model(u).Updates(map[string]interface{}{
		"username":   u.Username,
		"avatar":     u.Avatar,
		"last_login": time.Now(),
	})
	return result.Error
}

// QueryUserByID 根据用户ID查询用户
//
//	@param userID uint
//	@return user *User
//	@return err error
//	@author centonhuang
//	@update 2024-06-22 10:12:46
func QueryUserByID(userID uint) (user *User, err error) {
	result := database.DB.First(&user, userID)
	err = result.Error
	return
}

// QueryUserByPlatformAndID 根据绑定平台和ID查询用户
//
//	@param platform string
//	@param platformID string
//	@return user * User
//	@return err error
//	@author centonhuang
//	@update 2024-06-22 10:29:32
func QueryUserByPlatformAndID(platform string, platformID string) (user *User, err error) {
	condition := "platform = ?"

	switch platform {
	case PlatformGithub:
		condition += " AND github_bind_id = ?"
	default:
		err = errors.New("Unknown platform: " + platform)
		return
	}

	result := database.DB.First(&user, condition, platform, platformID)

	err = result.Error
	if err != nil && err.Error() == "record not found" {
		user, err = nil, nil
	}
	return
}

// AddUserByBasicInfo 根据基本信息添加用户
//
//	@param username string
//	@param platform string
//	@param avatar string
//	@param permission string
//	@return err error
//	@author centonhuang
//	@update 2024-06-22 10:18:46
func AddUserByBasicInfo(username string, avatar string, permission string, platform string, bindID string) (user *User, err error) {
	user = &User{
		Username:   username,
		Platform:   platform,
		Avatar:     avatar,
		Permission: permission,
		LastLogin: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}

	switch platform {
	case PlatformGithub:
		user.GithubBindID = bindID
	default:
		err = errors.New("Unknown platform: " + platform)
		return
	}

	err = user.Create()
	return
}
