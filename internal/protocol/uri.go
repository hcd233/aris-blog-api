package protocol

// UserURI 用户路径参数
//
//	@author centonhuang
//	@update 2024-09-18 02:50:19
type UserURI struct {
	UserName string `uri:"userName" binding:"required"`
}

// ArticleURI 文章路径参数
//
//	@author centonhuang
//	@update 2024-09-21 06:13:15
type ArticleURI struct {
	UserURI
	ArticleSlug string `uri:"articleSlug" binding:"required"`
}

// TagURI 标签路径参数
//
//	@author centonhuang
//	@update 2024-09-22 03:20:00
type TagURI struct {
	TagSlug string `uri:"tagSlug" binding:"required"`
}

// CategoryURI 分类路径参数
//
//	@author centonhuang
//	@update 2024-10-01 04:52:37
type CategoryURI struct {
	UserURI
	CategoryID uint `uri:"categoryID" binding:"required"`
}

// ArticleVersionURI 文章版本路径参数
//
//	@author centonhuang
//	@update 2024-10-18 03:13:26
type ArticleVersionURI struct {
	ArticleURI
	Version uint `uri:"version" binding:"required,min=1"`
}
