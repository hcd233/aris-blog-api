// Package model defines the database schema for the model.
//
//	@update 2024-06-22 09:33:43
package model

import (
	"database/sql"
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
	Name       string       `json:"name" gorm:"column:name;unique;not null"`
	Email      string       `json:"email" gorm:"column:email;unique;not null"`
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
		"last_login": time.Now(),
	})
	return result.Error
}

// QueryUsers 查询用户
//
//	@param offset int
//	@param limit int
//	@return users []*User
//	@return err error
//	@author centonhuang
//	@update 2024-09-17 08:18:54
func QueryUsers(offset int, limit int) (users []*User, err error) {
	result := database.DB.Offset(offset).Limit(limit).Find(&users)
	err = result.Error
	return
}

// QueryUserByID 根据用户ID查询用户
//
//	@param userID uint
//	@return user *User
//	@return err error
//	@author centonhuang
//	@update 2024-06-22 10:12:46
func QueryUserByID(userID uint) (user *User, err error) {
	result := database.DB.Where("id = ?", userID).Where("delete_at = ?", nil).First(&user)
	err = result.Error
	return
}

// QueryUserByName 根据用户名查询用户
//
//	@param userName string
//	@return user *User
//	@return err error
//	@author centonhuang
//	@update 2024-09-16 06:05:07
func QueryUserByName(userName string) (user *User, err error) {
	result := database.DB.Where(&User{Name: userName}).First(&user)
	err = result.Error
	return
}

// QueryUserByEmail 根据邮箱查询用户
//
//	@param email string
//	@return user *User
//	@return err error
//	@author centonhuang
//	@update 2024-09-16 11:21:25
func QueryUserByEmail(email string) (user *User) {
	result := database.DB.Where(User{Email: email}).First(&user)
	if result.Error != nil {
		return nil
	}
	return
}

// CreateUserByBasicInfo 根据基本信息添加用户
//
// @param username string
// @param email string
// @param avatar string
// @param permission string
// @return user *User
// @return err error
// @author centonhuang
// @update 2024-09-16 11:26:37
func CreateUserByBasicInfo(username string, email string, avatar string, permission string) (user *User, err error) {
	user = &User{
		Name:       username,
		Email:      email,
		Avatar:     avatar,
		Permission: permission,
		LastLogin: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}
	err = user.Create()
	return
}

// BindGithubID 绑定Github ID
//
//	@receiver u *User
//	@param githubID string
//	@return error
//	@author centonhuang
//	@update 2024-09-16 11:28:18
func (u *User) BindGithubID(githubID string) error {
	result := database.DB.Model(u).Update("github_bind_id", githubID)
	return result.Error
}

// GetUserDetailedInfo 获取用户详细信息
//
//	@param user *User
//	@return map
//	@author centonhuang
//	@update 2024-09-18 01:12:24
func GetUserDetailedInfo(user *User) map[string]interface{} {
	return map[string]interface{}{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"created_at": user.CreatedAt,
		"last_login": user.LastLogin.Time,
	}
}

// GetUserBasicInfo 获取用户基本信息
//
//	@param user *User
//	@return map
//	@author centonhuang
//	@update 2024-09-18 01:29:20
func GetUserBasicInfo(user *User) map[string]interface{} {
	return map[string]interface{}{
		"id":     user.ID,
		"name":   user.Name,
		"avatar": user.Avatar,
	}
}
