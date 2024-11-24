package docdao

import (
	"sync"

	"github.com/hcd233/Aris-blog/internal/resource/search"
	"github.com/hcd233/Aris-blog/internal/resource/search/document"
)

type (
	UserDocDAO    = BaseMeiliSearchDocDAO[document.UserDocument]
	TagDocDAO     = BaseMeiliSearchDocDAO[document.TagDocument]
	ArticleDocDAO = BaseMeiliSearchDocDAO[document.ArticleDocument]
)

var (
	userDocDAOSingleton    *UserDocDAO
	tagDocDAOSingleton     *TagDocDAO
	articleDocDAOSingleton *ArticleDocDAO

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
func GetUserDocDAO() *UserDocDAO {
	userDocOnce.Do(func() {
		userDocDAOSingleton = &UserDocDAO{
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
func GetTagDocDAO() *TagDocDAO {
	tagDocOnce.Do(func() {
		tagDocDAOSingleton = &TagDocDAO{
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
func GetArticleDocDAO() *ArticleDocDAO {
	articleDocOnce.Do(func() {
		articleDocDAOSingleton = &ArticleDocDAO{
			IndexName: articleIndex,
			Filters:   articleFilters,
			client:    search.GetSearchEngine(),
		}
	})
	return articleDocDAOSingleton
}
