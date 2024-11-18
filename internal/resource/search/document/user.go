// Package document 文档类
//
//	@update 2024-10-17 09:52:05
package document

import "github.com/hcd233/Aris-blog/internal/resource/database/model"

// UserDocument 用户文档
//
//	@author centonhuang
//	@update 2024-10-17 09:54:11
type UserDocument struct {
	ID       uint   `json:"id"`
	UserName string `json:"userName"`
	Avatar   string `json:"avatar"`
}

// TransformUserToDocument 将用户转换为文档
//
//	@param user *model.User
//	@return *UserDocument
//	@author centonhuang
//	@update 2024-10-18 01:35:58
func TransformUserToDocument(user *model.User) *UserDocument {
	return &UserDocument{
		ID:       user.ID,
		UserName: user.Name,
		Avatar:   user.Avatar,
	}
}
