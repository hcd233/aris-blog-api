package protocol

// CategoryURI 分类路径参数
//
//	@author centonhuang
//	@update 2024-10-01 04:52:37
type CategoryURI struct {
	UserURI
	CategoryID uint `uri:"categoryID" binding:"required"`
}

// CreateCategoryBody 创建分类请求体
//
//	@author centonhuang
//	@update 2024-09-28 07:02:11
type CreateCategoryBody struct {
	ParentID uint   `json:"parentID" binding:"omitempty"`
	Name     string `json:"name" binding:"required"`
}
