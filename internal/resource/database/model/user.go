// Package model defines the database schema for the model.
//
//	@update 2024-06-22 09:33:43
package model

import (
	"database/sql"
	"errors"
	"time"

	"github.com/hcd233/Aris-blog/internal/resource/database"
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
}

// Create 创建用户
//
//	@receiver u *User
//	@return error
//	@author centonhuang
//	@update 2024-06-22 10:10:07
func (u *User) Create() (err error) {
	err = database.DB.Create(u).Error
	return
}

// Delete 删除用户
//
//	@receiver u *User
//	@return err error
//	@author centonhuang
//	@update 2024-09-22 04:58:37
func (u *User) Delete() (err error) {
	err = database.DB.Delete(u).Error
	return
}

// BindGithubID 绑定Github ID
//
//	@receiver u *User
//	@param githubID string
//	@return error
//	@author centonhuang
//	@update 2024-09-16 11:28:18
func (u *User) BindGithubID(githubID string) (err error) {
	err = database.DB.Model(u).Update("github_bind_id", githubID).Error
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
		"created_at": u.CreatedAt,
		"last_login": u.LastLogin.Time,
		"permission": u.Permission,
	}
}

// UpdateUserInfoByID 使用ID更新用户信息
//
//	@param id uint
//	@param info map[string]interface{}
//	@return user *User
//	@return err error
//	@author centonhuang
//	@update 2024-09-18 04:22:03
func UpdateUserInfoByID(userID uint, info map[string]interface{}) (user *User, err error) {
	info["updated_at"] = time.Now()
	err = database.DB.Model(&User{}).Where(&User{ID: userID}).Updates(info).Error
	if err != nil {
		return nil, err
	}
	err = database.DB.Where(&User{ID: userID}).First(&user).Error
	return
}

// QueryUsers 查询用户
//
//	@param offset int
//	@param limit int
//	@return users []*User
//	@return err error
//	@author centonhuang
//	@update 2024-09-17 08:18:54
func QueryUsers(offset int, limit int) (users *[]User, err error) {
	err = database.DB.Offset(offset).Limit(limit).Find(&users).Error
	return
}

// QueryUserByID 根据用户ID查询用户
//
//	@param userID uint
//	@return user *User
//	@return err error
//	@author centonhuang
//	@update 2024-06-22 10:12:46
func QueryUserByID(userID uint, fields []string) (user *User, err error) {
	err = database.DB.Select(fields).Where(&User{ID: userID}).First(&user).Error
	return
}

// QueryUserByName 根据用户名查询用户
//
//	@param userName string
//	@return user *User
//	@return err error
//	@author centonhuang
//	@update 2024-09-16 06:05:07
func QueryUserByName(userName string, fields []string) (user *User, err error) {
	err = database.DB.Select(fields).Where(&User{Name: userName}).First(&user).Error
	return
}

// QueryUserByEmail 根据邮箱查询用户
//
//	@param email string
//	@return user *User
//	@return err error
//	@author centonhuang
//	@update 2024-09-16 11:21:25
func QueryUserByEmail(email string, fields []string, allowEmpty bool) (user *User, err error) {
	err = database.DB.Select(fields).Where(&User{Email: email}).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) && allowEmpty {
		user, err = nil, nil
	}
	return
}

// QueryUserFieldsByID 查询用户指定字段
//
//	@param userID int
//	@param fields []string
//	@return user *User
//	@return err error
//	@author centonhuang
//	@update 2024-09-21 03:08:02
func QueryUserFieldsByID(userID uint, fields []string, allowEmpty bool) (user *User, err error) {
	err = database.DB.Select(fields).Where(&User{ID: userID}).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) && allowEmpty {
		user, err = nil, nil
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
func CreateUserByBasicInfo(username string, email string, avatar string, permission Permission) (user *User, err error) {
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
