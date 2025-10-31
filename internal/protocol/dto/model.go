// Package dto 数据传输对象
package dto

// User represents a user entity
//
//	author centonhuang
//	update 2025-01-05 11:37:01
type User struct {
	UserID     uint   `json:"userID" doc:"Unique identifier for the user"`
	Name       string `json:"name" doc:"Display name of the user"`
	Email      string `json:"email,omitempty" doc:"Email address of the user"`
	Avatar     string `json:"avatar" doc:"URL or path to the user's avatar image"`
	CreatedAt  string `json:"createdAt,omitempty" doc:"Timestamp when the user account was created"`
	LastLogin  string `json:"lastLogin,omitempty" doc:"Timestamp of the user's last login"`
	Permission string `json:"permission,omitempty" doc:"Permission level of the user"`
}

// Tag 标签信息
//
//	author centonhuang
//	update 2025-10-31 05:32:00
type Tag struct {
	TagID       uint   `json:"tagID" doc:"Tag ID"`
	Name        string `json:"name" doc:"Tag name"`
	Slug        string `json:"slug" doc:"Tag slug"`
	Description string `json:"description,omitempty" doc:"Tag description"`
	CreatedAt   string `json:"createdAt,omitempty" doc:"Creation timestamp"`
	UpdatedAt   string `json:"updatedAt,omitempty" doc:"Update timestamp"`
	Likes       uint   `json:"likes,omitempty" doc:"Number of likes"`
}

// Article 文章信息
//
//	author centonhuang
//	update 2025-10-31 05:36:00
type Article struct {
	ArticleID   uint   `json:"articleID" doc:"Article ID"`
	Title       string `json:"title" doc:"Article title"`
	Slug        string `json:"slug" doc:"Article slug"`
	Status      string `json:"status" doc:"Article status"`
	User        *User  `json:"user" doc:"Author information"`
	CreatedAt   string `json:"createdAt" doc:"Creation timestamp"`
	UpdatedAt   string `json:"updatedAt" doc:"Update timestamp"`
	PublishedAt string `json:"publishedAt" doc:"Publication timestamp"`
	Likes       uint   `json:"likes" doc:"Number of likes"`
	Views       uint   `json:"views" doc:"Number of views"`
	Tags        []*Tag `json:"tags" doc:"List of tags"`
	Comments    int    `json:"comments" doc:"Number of comments"`
}

// ArticleVersion 文章版本信息
//
//	author centonhuang
//	update 2025-10-31 05:38:00
type ArticleVersion struct {
	ArticleVersionID uint   `json:"versionID" doc:"Version ID"`
	ArticleID        uint   `json:"articleID" doc:"Article ID"`
	VersionID        uint   `json:"version" doc:"Version number"`
	Content          string `json:"content" doc:"Version content"`
	CreatedAt        string `json:"createdAt" doc:"Creation timestamp"`
	UpdatedAt        string `json:"updatedAt" doc:"Update timestamp"`
}
