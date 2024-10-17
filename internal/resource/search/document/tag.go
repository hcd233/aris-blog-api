package document

import "github.com/hcd233/Aris-blog/internal/resource/database/model"

// TagDocument 标签文档
//
//	@author centonhuang
//	@update 2024-10-17 09:55:25
type TagDocument struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

// TransformTagToDocument 将标签转换为文档
//
//	@param tag *model.Tag
//	@return *TagDocument
//	@author centonhuang
//	@update 2024-10-18 01:35:13
func TransformTagToDocument(tag *model.Tag) *TagDocument {
	return &TagDocument{
		ID:          tag.ID,
		Name:        tag.Name,
		Slug:        tag.Slug,
		Description: tag.Description,
	}
}
