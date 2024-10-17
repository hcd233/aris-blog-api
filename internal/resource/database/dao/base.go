// Package dao 数据访问对象
//
//	@update 2024-10-17 02:31:49
package dao

import (
	"time"

	"gorm.io/gorm"
)

// baseDAO 基础数据访问对象
//
//	@author centonhuang
//	@update 2024-10-17 02:32:22
type baseDAO[T interface{}] struct{}

// Create 创建数据
//
//	@param dao *BaseDAO[T]
//	@return Create
//	@author centonhuang
//	@update 2024-10-17 02:51:49
func (dao *baseDAO[T]) Create(db *gorm.DB, data *T) (err error) {
	err = db.Create(&data).Error
	return
}

// Update 使用ID更新数据
//
//	@param dao *BaseDAO[T]
//	@return Update
//	@author centonhuang
//	@update 2024-10-17 02:52:18
func (dao *baseDAO[T]) Update(db *gorm.DB, data *T, info map[string]interface{}) (err error) {
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
func (dao *baseDAO[T]) Delete(db *gorm.DB, data *T) (err error) {
	err = db.Delete(&data).Error
	return
}

// GetByID 使用ID查询指定数据
//
//	@param dao *BaseDAO[T]
//	@return GetByID
//	@author centonhuang
//	@update 2024-10-17 03:06:57
func (dao *baseDAO[T]) GetByID(db *gorm.DB, id uint, fields []string) (data *T, err error) {
	err = db.Select(fields).Where("id = ?", id).First(&data).Error
	return
}

// Paginate 分页查询
//
//	@param dao *BaseDAO[T]
//	@return Paginate
//	@author centonhuang
//	@update 2024-10-17 03:09:11
func (dao *baseDAO[T]) Paginate(db *gorm.DB, fields []string, limit, offset int) (data *[]T, err error) {
	err = db.Select(fields).Limit(limit).Offset(offset).Find(&data).Error
	return
}
