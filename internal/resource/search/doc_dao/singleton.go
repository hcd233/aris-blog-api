package docdao

import (
	"sync"

	"github.com/hcd233/Aris-blog/internal/resource/search"
	"github.com/hcd233/Aris-blog/internal/resource/search/document"
)

var (
	userDocDAOSingleton    *BaseMeiliSearchDocDAO[document.UserDocument]
	tagDocDAOSingleton     *BaseMeiliSearchDocDAO[document.TagDocument]
	articleDocDAOSingleton *BaseMeiliSearchDocDAO[document.ArticleDocument]

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

// GetUserDocDAO 获取用户文档DAO单例
//
//	@return *BaseDocDAO
//	@author centonhuang
//	@update 2024-10-18 01:10:28
func GetUserDocDAO() *BaseMeiliSearchDocDAO[document.UserDocument] {
	userDocOnce.Do(func() {
		userDocDAOSingleton = &BaseMeiliSearchDocDAO[document.UserDocument]{
			IndexName: userIndex,
			Filters:   userFilters,
			client:    search.GetSearchEngine(),
		}
	})
	return userDocDAOSingleton
}

// GetTagDocDAO 获取标签文档DAO单例
//
//	@return *BaseDocDAO
//	@author centonhuang
//	@update 2024-10-18 01:09:59
func GetTagDocDAO() *BaseMeiliSearchDocDAO[document.TagDocument] {
	tagDocOnce.Do(func() {
		tagDocDAOSingleton = &BaseMeiliSearchDocDAO[document.TagDocument]{
			IndexName: tagIndex,
			Filters:   tagFilters,
			client:    search.GetSearchEngine(),
		}
	})
	return tagDocDAOSingleton
}

// GetArticleDocDAO 获取文章文档DAO单例
//
//	@return *BaseDocDAO
//	@author centonhuang
//	@update 2024-10-18 01:10:45
func GetArticleDocDAO() *BaseMeiliSearchDocDAO[document.ArticleDocument] {
	articleDocOnce.Do(func() {
		articleDocDAOSingleton = &BaseMeiliSearchDocDAO[document.ArticleDocument]{
			IndexName: articleIndex,
			Filters:   articleFilters,
			client:    search.GetSearchEngine(),
		}
	})
	return articleDocDAOSingleton
}
