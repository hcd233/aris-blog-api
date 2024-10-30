package dao

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"gorm.io/gorm"
)

// UserLikeDAO 用户点赞数据访问对象
//
//	@author centonhuang
//	@update 2024-10-30 03:49:48
type UserLikeDAO struct {
	baseDAO[model.UserLike]
}

// Delete 删除点赞信息
//
//	@receiver dao *UserLikeDAO
//	@param db *gorm.DB
//	@param userLike *model.UserLike
//	@return err error
//	@author centonhuang
//	@update 2024-10-30 05:21:20
func (dao *UserLikeDAO) Delete(db *gorm.DB, userLike *model.UserLike) (err error) {
	UUID := uuid.New().String()
	err = db.Model(userLike).Updates(map[string]interface{}{"object_type": fmt.Sprintf("%s-%s", userLike.ObjectType, UUID), "deleted_at": time.Now()}).Error
	return
}

// GetByUserIDAndObject 通过用户和点赞对象访问点赞信息
//
//	@receiver dao *UserLikeDAO
//	@param db *gorm.DB
//	@param userID uint
//	@param objectID uint
//	@param objectType model.LikeObjectType
//	@param fields []string
//	@return userLike *model.UserLike
//	@return err error
//	@author centonhuang
//	@update 2024-10-30 04:46:50
func (dao *UserLikeDAO) GetByUserIDAndObject(db *gorm.DB, userID uint, objectID uint, objectType model.LikeObjectType, fields []string) (userLike *model.UserLike, err error) {
	err = db.Select(fields).Where(model.UserLike{UserID: userID, ObjectID: objectID, ObjectType: objectType}).First(&userLike).Error
	return
}
