package dao

import (
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"gorm.io/gorm"
)

// UserViewDAO 用户浏览数据访问对象
//
//	@author centonhuang
//	@update 2024-10-30 03:49:48
type UserViewDAO struct {
	baseDAO[model.UserView]
}

func (dao *UserViewDAO) GetLatestViewByUserIDAndArticleID(db *gorm.DB, userID uint, articleID uint, fields []string) (userView *model.UserView, err error) {
	err = db.Select(fields).Where(model.UserView{UserID: userID, ArticleID: articleID}).Order("created_at desc").First(&userView).Error
	return
}

func (dao *UserViewDAO) PaginateWithPreloadsByUserID(db *gorm.DB, userID uint, page int, pageSize int) (userViews *[]model.UserView, pageInfo *PageInfo, err error) {
	limit, offset := pageSize, (page-1)*pageSize

	err = db.Preload("User").Preload("Article").Preload("Article.Tags").Preload("Article.User").Where(model.UserView{UserID: userID}).Limit(limit).Offset(offset).Find(&userViews).Error
	if err != nil {
		return
	}

	pageInfo = &PageInfo{
		Page:     page,
		PageSize: pageSize,
	}

	err = db.Where(model.UserView{UserID: userID}).Count(&pageInfo.Total).Error
	return
}
