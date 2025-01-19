package protocol

// UserURI 用户路径参数
//
//	author centonhuang
//	update 2024-09-18 02:50:19
type UserURI struct {
	UserID uint `uri:"userID" binding:"required"`
}

// ArticleURI 文章别名路径参数
//
//	author centonhuang
//	update 2024-09-21 06:13:15
type ArticleURI struct {
	ArticleID uint `uri:"articleID" binding:"required"`
}

// TagURI 标签路径参数
//
//	author centonhuang
//	update 2024-10-29 07:43:35
type TagURI struct {
	TagID uint `uri:"tagID" binding:"required"`
}

// CategoryURI 分类路径参数
//
//	author centonhuang
//	update 2024-10-01 04:52:37
type CategoryURI struct {
	CategoryID uint `uri:"categoryID" binding:"required"`
}

// ArticleVersionURI 文章版本路径参数
//
//	author centonhuang
//	update 2024-10-18 03:13:26
type ArticleVersionURI struct {
	ArticleURI
	Version uint `uri:"version" binding:"required,min=1"`
}

// CommentURI 评论路径参数
//
//	author centonhuang
//	update 2024-10-24 05:57:22
type CommentURI struct {
	CommentID uint `uri:"commentID" binding:"required,min=1"`
}

// ViewURI 查看路径参数
//
//	author centonhuang
//	update 2024-10-29 07:43:35
type ViewURI struct {
	UserURI
	ViewID uint `uri:"viewID" binding:"required"`
}

// ObjectURI 对象路径参数
//
//	author centonhuang
//	update 2024-10-29 07:43:35
type ObjectURI struct {
	ObjectName string `uri:"objectName" binding:"required"`
}

// TaskURI 任务路径参数
//
//	author centonhuang
//	update 2024-12-08 16:42:27
type TaskURI struct {
	TaskName string `uri:"taskName" binding:"required,oneof=contentCompletion articleSummary articleTranslation articleQA termExplaination"`
}

// PromptVersionURI 提示词版本路径参数
//
//	author centonhuang
//	update 2024-12-08 16:42:31
type PromptVersionURI struct {
	TaskURI
	Version uint `uri:"version" binding:"required,min=1"`
}
