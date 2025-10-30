package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
)

func initUserRouter(userGroup *huma.Group) {
	userHandler := handler.NewUserHandler()

	userGroup.UseMiddleware(middleware.JwtMiddlewareForHuma())

	// 获取当前用户信息
	huma.Register(userGroup, huma.Operation{
		OperationID: "getCurrentUserInfo",
		Method:      http.MethodGet,
		Path:        "/current",
		Summary:     "GetCurrentUserInfo",
		Description: "Get the current user's detailed information, including user ID, username, email, avatar, and permission information",
		Tags:        []string{"user"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, userHandler.HandleGetCurUserInfo)

	// 更新用户信息
	huma.Register(userGroup, huma.Operation{
		OperationID: "updateUserInfo",
		Method:      http.MethodPatch,
		Path:        "/",
		Summary:     "UpdateUserInfo",
		Description: "Update the current user's information, including the username and other fields",
		Tags:        []string{"user"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, userHandler.HandleUpdateInfo)

	// 获取指定用户信息
	huma.Register(userGroup, huma.Operation{
		OperationID: "getUserInfo",
		Method:      http.MethodGet,
		Path:        "/{userID}",
		Summary:     "GetUserInfo",
		Description: "Get the public information of the specified user by user ID, including user ID, username, and avatar",
		Tags:        []string{"user"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, userHandler.HandleGetUserInfo)
}
