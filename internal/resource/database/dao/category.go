package dao

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"gorm.io/gorm"
)

// CategoryDAO 类别DAO
//
//	author centonhuang
//	update 2024-10-17 02:30:24
type CategoryDAO struct {
	baseDAO[model.Category]
}

// Delete 删除类别
//
//	receiver dao *CategoryDAO
//	param db *gorm.DB
//	param category *model.Category
//	return err error
//	author centonhuang
//	update 2024-10-17 02:59:11
func (dao *CategoryDAO) Delete(db *gorm.DB, category *model.Category) (err error) {
	UUID := uuid.New().String()
	err = db.Model(category).Updates(map[string]interface{}{"name": fmt.Sprintf("%s-%s", category.Name, UUID), "deleted_at": time.Now()}).Error
	return
}

// PaginateChildren 获取子类别
//
//	receiver dao *CategoryDAO
//	param db *gorm.DB
//	param category *model.Category
//	param fields []string
//	param page int
//	param pageSize int
//	return children *[]model.Category
//	return pageInfo *PageInfo
//	return err error
//	author centonhuang
//	update 2024-11-01 07:09:50
func (dao *CategoryDAO) PaginateChildren(db *gorm.DB, category *model.Category, fields, preloads []string, page, pageSize int) (children *[]model.Category, pageInfo *PageInfo, err error) {
	limit, offset := pageSize, (page-1)*pageSize

	sql := db.Select(fields)
	for _, preload := range preloads {
		sql = sql.Preload(preload)
	}
	err = sql.Where(&model.Category{ParentID: category.ID}).Limit(limit).Offset(offset).Find(&children).Error
	if err != nil {
		return
	}

	pageInfo = &PageInfo{
		Page:     page,
		PageSize: pageSize,
	}

	err = db.Model(&children).Where(&model.Category{ParentID: category.ID}).Count(&pageInfo.Total).Error
	return
}

// GetParent 获取父类别
//
//	receiver dao *CategoryDAO
//	param db *gorm.DB
//	param category *model.Category
//	return parent *model.Category
//	return err error
//	author centonhuang
//	update 2024-10-17 03:04:41
func (dao *CategoryDAO) GetParent(db *gorm.DB, category *model.Category, fields, preloads []string) (parent *model.Category, err error) {
	sql := db.Select(fields)
	for _, preload := range preloads {
		sql = sql.Preload(preload)
	}
	err = sql.Where(&model.Category{ID: category.ParentID}).First(&parent).Error
	return
}

// GetRootByUserID 获取根类别
//
//	receiver dao *CategoryDAO
//	param db *gorm.DB
//	param userID uint
//	param fields []string
//	return category *model.Category
//	return err error
//	author centonhuang
//	update 2024-10-17 03:15:59
func (dao *CategoryDAO) GetRootByUserID(db *gorm.DB, userID uint, fields, preloads []string) (category *model.Category, err error) {
	sql := db.Select(fields)
	for _, preload := range preloads {
		sql = sql.Preload(preload)
	}
	err = sql.Where(&model.Category{UserID: userID}).Where("parent_id IS NULL").First(&category).Error
	return
}

// DeleteReclusiveByID 递归删除类别
//
//	receiver dao *CategoryDAO
//	param db *gorm.DB
//	param id uint
//	return err error
//	author centonhuang
//	update 2024-10-17 03:36:05
func (dao *CategoryDAO) DeleteReclusiveByID(db *gorm.DB, id uint, fields, preloads []string) (err error) {
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

	rootCategory, err := dao.GetByID(db, id, fields, preloads)
	if err != nil {
		return
	}

	*categories = append(*categories, *rootCategory)
	for _, category := range *categories {
		err = dao.Delete(tx, &category)
		if err != nil {
			return
		}
	}
	return
}

func (dao *CategoryDAO) reclusiveFindChildrenIDsByID(db *gorm.DB, categoryID uint, fields, preloads []string) (categories *[]model.Category, err error) {
	categories, _, err = dao.PaginateChildren(db, &model.Category{ID: categoryID}, fields, preloads, 2, -1)
	if err != nil {
		return
	}

	for _, category := range *categories {
		childrenCategories, err := dao.reclusiveFindChildrenIDsByID(db, category.ID, fields, preloads)
		if err != nil {
			return nil, err
		}
		*categories = append(*categories, *childrenCategories...)
	}

	return
}
