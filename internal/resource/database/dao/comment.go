package dao

import (
	"fmt"

	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
	"gorm.io/gorm"
)

// CommentDAO 评论DAO
//
//	author centonhuang
//	update 2024-10-23 05:22:38
type CommentDAO struct {
	baseDAO[model.Comment]
}

// PaginateChildren 获取子评论
//
//	receiver dao *CommentDAO
//	param db *gorm.DB
//	param comment *model.Comment
//	param fields []string
//	param page int
//	param pageSize int
//	return children *[]model.Comment
//	return pageInfo *PageInfo
//	return err error
//	author centonhuang
//	update 2024-11-01 07:09:55
func (dao *CommentDAO) PaginateChildren(db *gorm.DB, comment *model.Comment, fields, preloads []string, param *CommonParam) (children *[]model.Comment, pageInfo *PageInfo, err error) {
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

	err = sql.Where(&model.Comment{ParentID: comment.ID}).Limit(limit).Offset(offset).Find(&children).Error
	if err != nil {
		return
	}

	pageInfo = &PageInfo{
		Page:     param.Page,
		PageSize: param.PageSize,
	}

	err = db.Model(&children).Where(&model.Comment{ParentID: comment.ID}).Count(&pageInfo.Total).Error
	return
}

// GetParent 获取父类别
//
//	receiver dao *CommentDAO
//	param db *gorm.DB
//	param comment *model.Comment
//	param fields []string
//	return parent *model.Comment
//	return err error
//	author centonhuang
//	update 2024-10-23 05:22:55
func (dao *CommentDAO) GetParent(db *gorm.DB, comment *model.Comment, fields, preloads []string) (parent *model.Comment, err error) {
	sql := db.Select(fields)
	for _, preload := range preloads {
		sql = sql.Preload(preload)
	}
	err = sql.Where(&model.Comment{ID: comment.ParentID}).First(&parent).Error
	return
}

// PaginateRootsByArticleID 获取文章的根评论
//
//	receiver dao *CommentDAO
//	param db *gorm.DB
//	param articleID uint
//	param fields []string
//	param page int
//	param pageSize int
//	return comments *[]model.Comment
//	return pageInfo *PageInfo
//	return err error
//	author centonhuang
//	update 2024-11-01 07:10:00
func (dao *CommentDAO) PaginateRootsByArticleID(db *gorm.DB, articleID uint, fields, preloads []string, param *CommonParam) (comments *[]model.Comment, pageInfo *PageInfo, err error) {
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

	err = sql.Where(&model.Comment{ArticleID: articleID}).Where("parent_id IS NULL").Limit(limit).Offset(offset).Find(&comments).Error
	if err != nil {
		return
	}

	pageInfo = &PageInfo{
		Page:     param.Page,
		PageSize: param.PageSize,
	}

	err = db.Model(&comments).Where(&model.Comment{ArticleID: articleID}).Where("parent_id IS NULL").Count(&pageInfo.Total).Error
	return
}

// DeleteReclusiveByID 递归删除评论
//
//	receiver dao *CommentDAO
//	param db *gorm.DB
//	param id uint
//	return err error
//	author centonhuang
//	update 2024-10-23 05:23:03
func (dao *CommentDAO) DeleteReclusiveByID(db *gorm.DB, id uint, fields, preloads []string) (err error) {
	categories, err := dao.reclusiveFindChildrenIDsByID(db, id, fields, preloads)
	if err != nil {
		return
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = fmt.Errorf("panic occurred: %v", r)
		} else if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	rootComment, err := dao.GetByID(db, id, fields, preloads)
	if err != nil {
		return
	}

	*categories = append(*categories, *rootComment)
	for _, comment := range *categories {
		err = dao.Delete(tx, &comment)
		if err != nil {
			return
		}
	}
	return
}

func (dao *CommentDAO) reclusiveFindChildrenIDsByID(db *gorm.DB, commentID uint, fields, preloads []string) (categories *[]model.Comment, err error) {
	param := &CommonParam{
		PageParam: &PageParam{
			Page:     2,
			PageSize: -1,
		},
	}
	categories, _, err = dao.PaginateChildren(db, &model.Comment{ID: commentID}, fields, preloads, param)
	if err != nil {
		return
	}

	for _, comment := range *categories {
		childrenCategories, err := dao.reclusiveFindChildrenIDsByID(db, comment.ID, fields, preloads)
		if err != nil {
			return nil, err
		}
		*categories = append(*categories, *childrenCategories...)
	}

	return
}
