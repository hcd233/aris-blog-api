// Package dao DAO
//
//	@update 2024-10-17 02:31:49
package dao

import (
	"time"

	"gorm.io/gorm"
)

// baseDAO 基础DAO
//
//	@author centonhuang
//	@update 2024-10-17 02:32:22
type baseDAO[ModelT interface{}] struct{}

// PageInfo 分页信息
//
//	@author centonhuang
//	@update 2024-11-01 05:17:51
type PageInfo struct {
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
	Total    int64 `json:"total"`
}

// Create 创建数据
//
//	@param dao *BaseDAO[T]
//	@return Create
//	@author centonhuang
//	@update 2024-10-17 02:51:49
func (dao *baseDAO[ModelT]) Create(db *gorm.DB, data *ModelT) (err error) {
	err = db.Create(&data).Error
	return
}

// Update 使用ID更新数据
//
//	@param dao *BaseDAO[T]
//	@return Update
//	@author centonhuang
//	@update 2024-10-17 02:52:18
func (dao *baseDAO[ModelT]) Update(db *gorm.DB, data *ModelT, info map[string]interface{}) (err error) {
	info["updated_at"] = time.Now()
	err = db.Model(&data).Updates(info).Error
	return
}

// Delete 删除
//
//	@param dao *BaseDAO[T]
//	@return Delete
//	@author centonhuang
//	@update 2024-10-17 02:52:33
func (dao *baseDAO[ModelT]) Delete(db *gorm.DB, data *ModelT) (err error) {
	err = db.Delete(&data).Error
	return
}

func (dao *baseDAO[ModelT]) BatchDelete(db *gorm.DB, data *[]ModelT) (err error) {
	err = db.Delete(&data).Error
	return
}

// GetByID 使用ID查询指定数据
//
//	@param dao *BaseDAO[T]
//	@return GetByID
//	@author centonhuang
//	@update 2024-10-17 03:06:57
func (dao *baseDAO[ModelT]) GetByID(db *gorm.DB, id uint, fields []string, preloads []string) (data *ModelT, err error) {
	sql := db.Select(fields)
	for _, preload := range preloads {
		sql = sql.Preload(preload)
	}

	err = sql.Where("id = ?", id).First(&data).Error
	return
}

// BatchGetByIDs 批量使用ID查询指定数据
//
//	@param dao *baseDAO[T]
//	@return BatchGetByIDs
//	@author centonhuang
//	@update 2024-11-03 07:34:47
func (dao *baseDAO[ModelT]) BatchGetByIDs(db *gorm.DB, ids []uint, fields []string, preloads []string) (data *[]ModelT, err error) {
	sql := db.Select(fields)
	for _, preload := range preloads {
		sql = sql.Preload(preload)
	}
	err = sql.Where("id IN ?", ids).Find(&data).Error
	return
}

// Paginate 分页查询
//
//	@param dao *BaseDAO[T]
//	@return Paginate
//	@author centonhuang
//	@update 2024-10-17 03:09:11
func (dao *baseDAO[ModelT]) Paginate(db *gorm.DB, fields []string, preloads []string, page, pageSize int) (data *[]ModelT, pageInfo *PageInfo, err error) {
	limit, offset := pageSize, (page-1)*pageSize

	sql := db.Select(fields)
	for _, preload := range preloads {
		sql = sql.Preload(preload)
	}
	err = sql.Limit(limit).Offset(offset).Find(&data).Error

	if err != nil {
		return
	}

	pageInfo = &PageInfo{
		Page:     page,
		PageSize: pageSize,
	}

	err = db.Model(&data).Count(&pageInfo.Total).Error

	return
}
