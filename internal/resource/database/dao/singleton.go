package dao

import "sync"

var (
	categoryDAOSingleton *CategoryDAO
	once                 sync.Once
)

// GetCategoryDAO 获取类别数据访问对象
//
//	@return *categoryDAO
//	@author centonhuang
//	@update 2024-10-17 03:45:31
func GetCategoryDAO() *CategoryDAO {
	once.Do(func() {
		categoryDAOSingleton = &CategoryDAO{}
	})
	return categoryDAOSingleton
}
