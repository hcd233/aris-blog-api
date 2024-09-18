package protocol

// UserURI 用户路径参数
//
//	@author centonhuang
//	@update 2024-09-18 02:50:19
type UserURI struct {
	UserName string `uri:"userName" binding:"required"`
}

// QueryUserParams 查询用户请求参数
//
//	@author centonhuang
//	@update 2024-09-18 02:56:39
type QueryUserParams struct {
	Query  string `form:"query" binding:"required"`
	Limit  int64  `form:"limit" binding:"required,min=1,max=50"`
	Offset int64  `form:"offset" binding:"gte=0"`
}

// UpdateUserBody 更新用户请求体
//
//	@author centonhuang
//	@update 2024-09-18 02:39:31
type UpdateUserBody struct {
	UserName string `json:"userName" binding:"required"`
}
