package dao

import (
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"gorm.io/gorm"
)

// UserDAO 用户DAO
//
//	author centonhuang
//	update 2024-10-17 02:30:24
type UserDAO struct {
	baseDAO[model.User]
}

// GetByEmail 通过邮箱获取用户
//
//	receiver dao *UserDAO
//	param db *gorm.DB
//	param email string
//	param fields []string
//	return user *model.User
//	return err error
//	author centonhuang
//	update 2024-10-17 05:08:00
func (dao *UserDAO) GetByEmail(db *gorm.DB, email string, fields, preloads []string) (user *model.User, err error) {
	sql := db.Select(fields)
	for _, preload := range preloads {
		sql = sql.Preload(preload)
	}
	err = sql.Where(model.User{Email: email}).First(&user).Error
	return
}

// GetByName 通过用户名获取用户
//
//	receiver dao *UserDAO
//	param db *gorm.DB
//	param name string
//	param fields []string
//	return user *model.User
//	return err error
//	author centonhuang
//	update 2024-10-17 05:18:46
func (dao *UserDAO) GetByName(db *gorm.DB, name string, fields, preloads []string) (user *model.User, err error) {
	sql := db.Select(fields)
	for _, preload := range preloads {
		sql = sql.Preload(preload)
	}
	err = sql.Where(model.User{Name: name}).First(&user).Error
	return
}

// DeliverLLMQuota 发放LLM配额
//
//	receiver dao *UserDAO
//	param db *gorm.DB
//	param userID uint
//	param quota model.Quota
//	return err error
//	author centonhuang
//	update 2024-11-26 02:48:55
func (dao *UserDAO) DeliverLLMQuota(db *gorm.DB, userID uint, quota model.Quota) error {
	return dao.Update(db, &model.User{ID: userID}, map[string]interface{}{
		"llm_quota": quota,
	})
}

// BatchDeliverLLMQuota 批量发放LLM配额
//
//	receiver dao *UserDAO
//	param db *gorm.DB
//	param userIDs []uint
//	param quota model.Quota
//	return err error
//	author centonhuang
//	update 2024-11-26 02:50:15
func (dao *UserDAO) BatchDeliverLLMQuota(db *gorm.DB, userIDs []uint, quota model.Quota) error {
	return db.Model(&model.User{}).Where("id IN ?", userIDs).Update("llm_quota", quota).Error
}
