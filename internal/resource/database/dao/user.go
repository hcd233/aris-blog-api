package dao

import (
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"gorm.io/gorm"
)

// UserDAO 用户数据访问对象
//
//	@author centonhuang
//	@update 2024-10-17 02:30:24
type UserDAO struct {
	baseDAO[model.User]
}

// GetByEmail 通过邮箱获取用户
//
//	@receiver dao *UserDAO
//	@param db *gorm.DB
//	@param email string
//	@param fields []string
//	@return user *model.User
//	@return err error
//	@author centonhuang
//	@update 2024-10-17 05:08:00
func (dao *UserDAO) GetByEmail(db *gorm.DB, email string, fields []string) (user *model.User, err error) {
	err = db.Select(fields).Where(model.User{Email: email}).First(&user).Error
	return
}

// GetByName 通过用户名获取用户
//
//	@receiver dao *UserDAO
//	@param db *gorm.DB
//	@param name string
//	@param fields []string
//	@return user *model.User
//	@return err error
//	@author centonhuang
//	@update 2024-10-17 05:18:46
func (dao *UserDAO) GetByName(db *gorm.DB, name string, fields []string) (user *model.User, err error) {
	err = db.Select(fields).Where(model.User{Name: name}).First(&user).Error
	return
}
