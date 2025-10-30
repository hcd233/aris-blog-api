package router

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-blog-api/internal/handler"
)

func initUserRouter(userGroup *huma.Group) {
	userHandler := handler.NewUserHandler()

	// 获取当前用户信息
	huma.Register(userGroup, huma.Operation{
		OperationID: "getCurrentUserInfo",
		Path:        "/current",
		Summary:     "获取当前用户信息",
		Description: "获取当前登录用户的详细信息，包括用户ID、用户名、邮箱、头像、权限等信息",
		Tags:        []string{"user"},
	}, userHandler.HandleGetCurUserInfo)

	// 更新用户信息
	huma.Register(userGroup, huma.Operation{
		OperationID: "updateUserInfo",
		Path:        "/",
		Summary:     "更新用户信息",
		Description: "更新当前登录用户的信息，包括用户名等字段",
		Tags:        []string{"user"},
	}, userHandler.HandleUpdateInfo)

	// 获取指定用户信息
	huma.Register(userGroup, huma.Operation{
		OperationID: "getUserInfo",
		Path:        "/{userID}",
		Summary:     "获取用户信息",
		Description: "根据用户ID获取指定用户的公开信息，包括用户ID、用户名、头像等",
		Tags:        []string{"user"},
	}, userHandler.HandleGetUserInfo)
}
