package humadto

// UserPath 用户路径参数（Huma 风格，使用 path 标签）
//
//	author centonhuang
//	update 2025-10-28 00:00:00
type UserPath struct {
	UserID uint `path:"userID"`
}

// ArticlePath 文章路径参数
//
//	author centonhuang
//	update 2025-10-28 00:00:00
type ArticlePath struct {
	ArticleID uint `path:"articleID"`
}

// ArticleSlugPath 文章别名路径参数
//
//	author centonhuang
//	update 2025-10-28 00:00:00
type ArticleSlugPath struct {
	AuthorName  string `path:"authorName"`
	ArticleSlug string `path:"articleSlug"`
}

// TagPath 标签路径参数
//
//	author centonhuang
//	update 2025-10-28 00:00:00
type TagPath struct {
	TagID uint `path:"tagID"`
}

// CategoryPath 分类路径参数
//
//	author centonhuang
//	update 2025-10-28 00:00:00
type CategoryPath struct {
	CategoryID uint `path:"categoryID"`
}

// ArticleVersionPath 文章版本路径参数
//
//	author centonhuang
//	update 2025-10-28 00:00:00
type ArticleVersionPath struct {
	ArticlePath
	Version uint `path:"version"`
}

// CommentPath 评论路径参数
//
//	author centonhuang
//	update 2025-10-28 00:00:00
type CommentPath struct {
	CommentID uint `path:"commentID"`
}

// ViewPath 查看路径参数
//
//	author centonhuang
//	update 2025-10-28 00:00:00
type ViewPath struct {
	UserPath
	ViewID uint `path:"viewID"`
}

// ObjectPath 对象路径参数
//
//	author centonhuang
//	update 2025-10-28 00:00:00
type ObjectPath struct {
	ObjectName string `path:"objectName"`
}

// TaskPath 任务路径参数
//
//	author centonhuang
//	update 2025-10-28 00:00:00
type TaskPath struct {
	TaskName string `path:"taskName"`
}

// PromptVersionPath 提示词版本路径参数
//
//	author centonhuang
//	update 2025-10-28 00:00:00
type PromptVersionPath struct {
	TaskPath
	Version uint `path:"version"`
}

// ProviderPath OAuth2 提供商路径参数
//
//	author centonhuang
//	update 2025-10-28 00:00:00
type ProviderPath struct {
	Provider string `path:"provider"`
}
