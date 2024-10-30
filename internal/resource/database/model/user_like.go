package model

import "gorm.io/gorm"

// LikeObjectType 点赞对象类型
//
//	@author centonhuang
//	@update 2024-10-30 03:46:24
type LikeObjectType string

const (

	// LikeObjectTypeTag LikeObjectType 点赞标签类型
	//	@update 2024-10-30 03:47:06
	LikeObjectTypeTag LikeObjectType = "tag"

	// LikeObjectTypeComment LikeObjectType 点赞评论类型
	//	@update 2024-10-30 03:47:24
	LikeObjectTypeComment LikeObjectType = "comment"

	// LikeObjectTypeArticle LikeObjectType 点赞文章类型
	//	@update 2024-10-30 03:47:28
	LikeObjectTypeArticle LikeObjectType = "article"
)

// UserLike 用户点赞
//
//	@author centonhuang
//	@update 2024-10-29 12:31:01
type UserLike struct {
	gorm.Model
	UserID     uint           `json:"user_id" gorm:"not null;index:user_id_object_type;uniqueIndex:user_object;comment:用户ID"`
	User       *User          `json:"user" gorm:"foreignKey:UserID;references:ID"`
	ObjectID   uint           `json:"object_id" gorm:"not null;index:user_id_object_type;uniqueIndex:user_object;comment:对象ID"`
	ObjectType LikeObjectType `json:"object_type" gorm:"not null;index:user_id_object_type;uniqueIndex:user_object;comment:对象类型"`
}
