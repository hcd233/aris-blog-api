package dao

import (
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UserViewDAO 用户浏览数据访问对象
//
//	author centonhuang
//	update 2024-10-30 03:49:48
type UserViewDAO struct {
	baseDAO[model.UserView]
}

// GetLatestViewByUserIDAndArticleID 获取用户最新浏览记录
//
//	@receiver dao *UserViewDAO
//	@param db
//	@param userID
//	@param articleID
//	@param fields
//	@param preloads
//	@return userView
//	@return err
//	@author centonhuang
//	@update 2025-10-31 18:16:21
func (dao *UserViewDAO) GetLatestViewByUserIDAndArticleID(db *gorm.DB, userID uint, articleID uint, fields, preloads []string) (userView *model.UserView, err error) {
	sql := db.Select(fields)
	for _, preload := range preloads {
		sql = sql.Preload(preload)
	}

	err = sql.Where(model.UserView{UserID: userID, ArticleID: articleID}).Order("created_at desc").First(&userView).Error
	return
}

// PaginateByUserID 分页查询用户浏览记录
//
//	@receiver dao *UserViewDAO
//	@param db
//	@param userID
//	@param fields
//	@param preloads
//	@param param
//	@return userViews
//	@return pageInfo
//	@return err
//	@author centonhuang
//	@update 2025-10-31 18:16:28
func (dao *UserViewDAO) PaginateByUserID(db *gorm.DB, userID uint, fields, preloads []string, param *CommonParam) (userViews *[]model.UserView, pageInfo *PageInfo, err error) {
	limit, offset := param.PageSize, (param.Page-1)*param.PageSize

	sql := db.Select(fields)
	for _, preload := range preloads {
		sql = sql.Preload(preload)
	}

	if param.Query != "" && len(param.QueryFields) > 0 {
		like := "%" + param.Query + "%"
		expressions := make([]clause.Expression, 0, len(param.QueryFields))
		for _, field := range param.QueryFields {
			if field == "" {
				continue
			}
			expressions = append(expressions, clause.Like{Column: clause.Column{Name: field}, Value: like})
		}

		if len(expressions) > 0 {
			sql = sql.Where(expressions[0])
			for _, expr := range expressions[1:] {
				sql = sql.Or(expr)
			}
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
