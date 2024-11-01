package dao

import (
	"fmt"

	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"gorm.io/gorm"
)

// CommentDAO 评论DAO
//
//	@author centonhuang
//	@update 2024-10-23 05:22:38
type CommentDAO struct {
	baseDAO[model.Comment]
}

// PaginateChildren 获取子评论
//
//	@receiver dao *CommentDAO
//	@param db *gorm.DB
//	@param comment *model.Comment
//	@param fields []string
//	@param limit int
//	@param offset int
//	@return children *[]model.Comment
//	@return err error
//	@author centonhuang
//	@update 2024-10-23 05:22:43
func (dao *CommentDAO) PaginateChildren(db *gorm.DB, comment *model.Comment, fields []string, page, pageSize int) (children *[]model.Comment, pageInfo *PageInfo, err error) {
	limit, offset := pageSize, (page-1)*pageSize
	err = db.Select(fields).Limit(limit).Offset(offset).Where(&model.Comment{ParentID: comment.ID}).Find(&children).Error
	if err != nil {
		return
	}

	pageInfo = &PageInfo{
		Page:     page,
		PageSize: pageSize,
	}

	err = db.Model(&model.Comment{}).Where(&model.Comment{ParentID: comment.ID}).Count(&pageInfo.Total).Error
	return
}

// GetParent 获取父类别
//
//	@receiver dao *CommentDAO
//	@param db *gorm.DB
//	@param comment *model.Comment
//	@param fields []string
//	@return parent *model.Comment
//	@return err error
//	@author centonhuang
//	@update 2024-10-23 05:22:55
func (dao *CommentDAO) GetParent(db *gorm.DB, comment *model.Comment, fields []string) (parent *model.Comment, err error) {
	err = db.Select(fields).Where(&model.Comment{ID: comment.ParentID}).First(&parent).Error
	return
}

// PaginateRootsByArticleID 获取文章的根评论
//
//	@receiver dao *CommentDAO
//	@param db *gorm.DB
//	@param articleID uint
//	@param fields []string
//	@param limit int
//	@param offset int
//	@return comments *[]model.Comment
//	@return err error
//	@author centonhuang
//	@update 2024-10-23 07:15:21
func (dao *CommentDAO) PaginateRootsByArticleID(db *gorm.DB, articleID uint, fields []string, page, pageSize int) (comments *[]model.Comment, pageInfo *PageInfo, err error) {
	limit, offset := pageSize, (page-1)*pageSize
	err = db.Select(fields).Limit(limit).Offset(offset).Where(&model.Comment{ArticleID: articleID}).Where("parent_id IS NULL").Find(&comments).Error
	if err != nil {
		return
	}

	pageInfo = &PageInfo{
		Page:     page,
		PageSize: pageSize,
	}

	err = db.Model(&model.Comment{}).Where(&model.Comment{ArticleID: articleID}).Where("parent_id IS NULL").Count(&pageInfo.Total).Error
	return
}

// GetByArticleIDAndID 根据文章ID和评论ID获取评论
//
//	@receiver dao *CommentDAO
//	@param db *gorm.DB
//	@param articleID uint
//	@param id uint
//	@param fields []string
//	@return comment *model.Comment
//	@return err error
//	@author centonhuang
//	@update 2024-10-24 06:01:04
func (dao *CommentDAO) GetByArticleIDAndID(db *gorm.DB, articleID, id uint, fields []string) (comment *model.Comment, err error) {
	err = db.Select(fields).Where(&model.Comment{ArticleID: articleID, ID: id}).First(&comment).Error
	return
}

// DeleteReclusiveByID 递归删除评论
//
//	@receiver dao *CommentDAO
//	@param db *gorm.DB
//	@param id uint
//	@return err error
//	@author centonhuang
//	@update 2024-10-23 05:23:03
func (dao *CommentDAO) DeleteReclusiveByID(db *gorm.DB, id uint) (err error) {
	categories, err := dao.reclusiveFindChildrenIDsByID(db, id)
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

	rootComment, err := dao.GetByID(db, id, []string{"id"})
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

func (dao *CommentDAO) reclusiveFindChildrenIDsByID(db *gorm.DB, commentID uint) (categories *[]model.Comment, err error) {
	categories, _, err = dao.PaginateChildren(db, &model.Comment{ID: commentID}, []string{"id"}, 2, -1)
	if err != nil {
		return
	}

	for _, comment := range *categories {
		childrenCategories, err := dao.reclusiveFindChildrenIDsByID(db, comment.ID)
		if err != nil {
			return nil, err
		}
		*categories = append(*categories, *childrenCategories...)
	}

	return
}
