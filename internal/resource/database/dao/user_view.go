package dao

import (
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
	"gorm.io/gorm"
)

// UserViewDAO 用户浏览数据访问对象
//
//	author centonhuang
//	update 2024-10-30 03:49:48
type UserViewDAO struct {
	baseDAO[model.UserView]
}

func (dao *UserViewDAO) GetLatestViewByUserIDAndArticleID(db *gorm.DB, userID uint, articleID uint, fields, preloads []string) (userView *model.UserView, err error) {
	sql := db.Select(fields)
	for _, preload := range preloads {
		sql = sql.Preload(preload)
	}

	err = sql.Where(model.UserView{UserID: userID, ArticleID: articleID}).Order("created_at desc").First(&userView).Error
	return
}

func (dao *UserViewDAO) PaginateByUserID(db *gorm.DB, userID uint, fields, preloads []string, param *PaginateParam) (userViews *[]model.UserView, pageInfo *PageInfo, err error) {
	limit, offset := param.PageSize, (param.Page-1)*param.PageSize

	sql := db.Select(fields)
	for _, preload := range preloads {
		sql = sql.Preload(preload)
	}

	if param.Query != "" && len(param.QueryFields) > 0 {
		sql = sql.Where("? LIKE ?", param.QueryFields[0], "%"+param.Query+"%")
		for _, field := range param.QueryFields[1:] {
			sql = sql.Or("? LIKE ?", field, "%"+param.Query+"%")
		}
	}

	err = sql.Where(model.UserView{UserID: userID}).Limit(limit).Offset(offset).Find(&userViews).Error
	if err != nil {
		return
	}

	pageInfo = &PageInfo{
		Page:     param.Page,
		PageSize: param.PageSize,
	}

	err = db.Model(&userViews).Where(model.UserView{UserID: userID}).Count(&pageInfo.Total).Error
	return
}
