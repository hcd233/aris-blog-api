package docdao

import (
	"sync"

	"github.com/hcd233/Aris-blog/internal/resource/search/document"
)

var (
	userDocDAOSingleton    *BaseDocDAO[document.UserDocument]
	tagDocDAOSingleton     *BaseDocDAO[document.TagDocument]
	articleDocDAOSingleton *BaseDocDAO[document.ArticleDocument]

	userDocOnce    sync.Once
	tagDocOnce     sync.Once
	articleDocOnce sync.Once

	userIndex    = "user"
	tagIndex     = "tag"
	articleIndex = "article"

	userFilters    = []string{}
	tagFilters     = []string{"creator"}
	articleFilters = []string{"author"}
)

// GetUserDocDAO 获取用户文档DTO单例
//
//	@return *BaseDocDAO
//	@author centonhuang
//	@update 2024-10-18 01:10:28
func GetUserDocDAO() *BaseDocDAO[document.UserDocument] {
	userDocOnce.Do(func() {
		userDocDAOSingleton = &BaseDocDAO[document.UserDocument]{
			IndexName: userIndex,
			Filters:   userFilters,
		}
	})
	return userDocDAOSingleton
}

// GetTagDocDAO 获取标签文档DTO单例
//
//	@return *BaseDocDAO
//	@author centonhuang
//	@update 2024-10-18 01:09:59
func GetTagDocDAO() *BaseDocDAO[document.TagDocument] {
	tagDocOnce.Do(func() {
		tagDocDAOSingleton = &BaseDocDAO[document.TagDocument]{
			IndexName: tagIndex,
			Filters:   tagFilters,
		}
	})
	return tagDocDAOSingleton
}

// GetArticleDocDAO 获取文章文档DTO单例
//
//	@return *BaseDocDAO
//	@author centonhuang
//	@update 2024-10-18 01:10:45
func GetArticleDocDAO() *BaseDocDAO[document.ArticleDocument] {
	articleDocOnce.Do(func() {
		articleDocDAOSingleton = &BaseDocDAO[document.ArticleDocument]{
			IndexName: articleIndex,
			Filters:   articleFilters,
		}
	})
	return articleDocDAOSingleton
}
