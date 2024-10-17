package docdao

import (
	"sync"

	"github.com/hcd233/Aris-blog/internal/resource/search/document"
)

var (
	tagDocDAOSingleton     *BaseDocDAO[document.TagDocument]
	userDocDAOSingleton    *BaseDocDAO[document.UserDocument]
	articleDocDAOSingleton *BaseDocDAO[document.ArticleDocument]

	tagDocOnce     sync.Once
	userDocOnce    sync.Once
	articleDocOnce sync.Once

	userIndex    = "user"
	tagIndex     = "tag"
	articleIndex = "article"
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
		}
	})
	return articleDocDAOSingleton
}
