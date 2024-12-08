package document

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

// TextDocument 文本文档
//
//	@author centonhuang
//	@update 2024-12-07 14:13:03
type TextDocument struct {
	Document
	ID       string
	Content  string
	Metadata map[string]interface{}
}

// NewTextDocument 创建文本文档
//
//	@param content string
//	@param metadata map[string]interface{}
//	@return *TextDocument
//	@author centonhuang
//	@update 2024-12-08 15:48:37
func NewTextDocument(content string, metadata map[string]interface{}) *TextDocument {
	return &TextDocument{
		ID:       uuid.New().String(),
		Content:  content,
		Metadata: metadata,
	}
}

// Map 转换为map
//
//	@receiver td *TextDocument
//	@return map
//	@author centonhuang
//	@update 2024-12-07 14:25:33
func (td *TextDocument) Map() map[string]interface{} {
	return lo.Assign(map[string]interface{}{"content": td.Content}, td.Metadata)
}

// String 转换为字符串
//
//	@receiver td *TextDocument
//	@return string
//	@author centonhuang
//	@update 2024-12-07 14:26:03
func (td *TextDocument) String() string {
	str := fmt.Sprintf("<content>\n%s\n</content>\n", td.Content)
	if len(td.Metadata) == 0 {
		return str
	}
	str += "<metadata>\n"
	for k, v := range td.Metadata {
		str += fmt.Sprintf("\t<%s>\n\t%v\n\t</%s>\n", k, v, k)
	}
	str += "</metadata>\n"
	return str
}
