package dao

import (
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
)

// ArticleDAO 标签数据访问对象
//
//	@author centonhuang
//	@update 2024-10-17 06:34:00
type ArticleDAO struct {
	baseDAO[model.Article]
}
