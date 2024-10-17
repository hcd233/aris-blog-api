package dao

import (
	"sync"
)

var (
	categoryDAOSingleton *CategoryDAO
	userDAOSingleton     *UserDAO
	categoryOnce         sync.Once
	userOnce             sync.Once
)

// GetCategoryDAO 获取类别数据访问对象
//
//	@return *categoryDAO
//	@author centonhuang
//	@update 2024-10-17 03:45:31
func GetCategoryDAO() *CategoryDAO {
	categoryOnce.Do(func() {
		categoryDAOSingleton = &CategoryDAO{}
	})
	return categoryDAOSingleton
}

// GetUserDAO 获取用户数据访问对象
//
//	@return *baseDAO
//	@author centonhuang
//	@update 2024-10-17 04:59:37
func GetUserDAO() *UserDAO {
	userOnce.Do(func() {
		userDAOSingleton = &UserDAO{}
	})
	return userDAOSingleton
}
