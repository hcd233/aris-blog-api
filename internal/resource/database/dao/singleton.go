package dao

import (
	"sync"
)

var (
	categoryDAOSingleton       *CategoryDAO
	userDAOSingleton           *UserDAO
	tagDAOSingleton            *TagDAO
	articleDAOSingleton        *ArticleDAO
	articleVersionDAOSingleton *ArticleVersionDAO
	categoryOnce               sync.Once
	userOnce                   sync.Once
	tagOnce                    sync.Once
	articleOnce                sync.Once
	articleVersionOnce         sync.Once
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

// GetTagDAO 获取标签数据访问对象
//
//	@return *tagDAO
//	@author centonhuang
//	@update 2024-10-17 05:30:24
func GetTagDAO() *TagDAO {
	tagOnce.Do(func() {
		tagDAOSingleton = &TagDAO{}
	})
	return tagDAOSingleton
}

// GetArticleDAO 获取文章数据访问对象
//
//	@return *ArticleDAO
//	@author centonhuang
//	@update 2024-10-17 06:34:28
func GetArticleDAO() *ArticleDAO {
	articleOnce.Do(func() {
		articleDAOSingleton = &ArticleDAO{}
	})
	return articleDAOSingleton
}

// GetArticleVersionDAO 获取文章版本数据访问对象
//
//	@return *ArticleVersionDAO
//	@author centonhuang
//	@update 2024-10-17 08:12:02
func GetArticleVersionDAO() *ArticleVersionDAO {
	articleVersionOnce.Do(func() {
		articleVersionDAOSingleton = &ArticleVersionDAO{}
	})
	return articleVersionDAOSingleton
}
