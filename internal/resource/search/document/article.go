package document

import (
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"github.com/samber/lo"
)

// ArticleDocument 文章文档
//
//	@author centonhuang
//	@update 2024-10-17 10:05:45
type ArticleDocument struct {
	ID      uint     `json:"id"`
	Title   string   `json:"title,omitempty"`
	Author  string   `json:"author,omitempty"`
	Content string   `json:"content,omitempty"`
	Tags    []string `json:"tags,omitempty"`
}

// TransformArticleToDocument 将文章转换为文档
//
//	@param article *model.Article
//	@param latestVersion *model.ArticleVersion
//	@return *ArticleDocument
//	@author centonhuang
//	@update 2024-10-18 01:34:41
func TransformArticleToDocument(article *model.Article, latestVersion *model.ArticleVersion) *ArticleDocument {
	return &ArticleDocument{
		ID:      article.ID,
		Title:   article.Title,
		Author:  article.User.Name,
		Content: latestVersion.Content,
		Tags:    lo.Map(article.Tags, func(tag model.Tag, idx int) string { return tag.Name }),
	}
}
