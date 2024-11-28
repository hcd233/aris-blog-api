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
	commentDAOSingleton        *CommentDAO
	userLikeDAOSingleton       *UserLikeDAO
	userViewDAOSingleton       *UserViewDAO
	promptDAOSingleton         *PromptDAO

	categoryOnce       sync.Once
	userOnce           sync.Once
	tagOnce            sync.Once
	articleOnce        sync.Once
	articleVersionOnce sync.Once
	commentOnce        sync.Once
	userLikeOnce       sync.Once
	userViewOnce       sync.Once
	promptOnce         sync.Once
)

// GetCategoryDAO 获取类别DAO
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

// GetUserDAO 获取用户DAO
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

// GetTagDAO 获取标签DAO
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

// GetArticleDAO 获取文章DAO
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

// GetArticleVersionDAO 获取文章版本DAO
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

// GetCommentDAO 获取评论DAO
//
//	@return *CommentDAO
//	@author centonhuang
//	@update 2024-10-23 06:01:15
func GetCommentDAO() *CommentDAO {
	commentOnce.Do(func() {
		commentDAOSingleton = &CommentDAO{}
	})
	return commentDAOSingleton
}

// GetUserLikeDAO 获取用户点赞DAO
//
//	@return *UserLikeDAO
//	@author centonhuang
//	@update 2024-10-30 03:50:52
func GetUserLikeDAO() *UserLikeDAO {
	userLikeOnce.Do(func() {
		userLikeDAOSingleton = &UserLikeDAO{}
	})
	return userLikeDAOSingleton
}

// GetUserViewDAO 获取用户浏览DAO
//
//	@return *UserLikeDAO
//	@author centonhuang
//	@update 2024-10-30 03:50:52
func GetUserViewDAO() *UserViewDAO {
	userViewOnce.Do(func() {
		userViewDAOSingleton = &UserViewDAO{}
	})
	return userViewDAOSingleton
}

// GetPromptDAO 获取提示词DAO
//
//	@return *PromptDAO
//	@author centonhuang
//	@update 2024-11-28 04:26:28
func GetPromptDAO() *PromptDAO {
	promptOnce.Do(func() {
		promptDAOSingleton = &PromptDAO{}
	})
	return promptDAOSingleton
}
