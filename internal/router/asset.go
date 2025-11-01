package router

import (
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
)

func initAssetRouter(assetGroup *huma.Group) {
	assetHandler := handler.NewAssetHandlerForHuma()

	assetGroup.UseMiddleware(middleware.JwtMiddlewareForHuma())

	// Like????
	likeGroup := huma.NewGroup(assetGroup, "/like")

	huma.Register(likeGroup, huma.Operation{
		OperationID: "listUserLikeArticles",
		Method:      http.MethodGet,
		Path:        "/articles",
		Summary:     "ListUserLikeArticles",
		Description: "List articles liked by the current user",
		Tags:        []string{"asset"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, assetHandler.HandleListUserLikeArticles)

	huma.Register(likeGroup, huma.Operation{
		OperationID: "listUserLikeComments",
		Method:      http.MethodGet,
		Path:        "/comments",
		Summary:     "ListUserLikeComments",
		Description: "List comments liked by the current user",
		Tags:        []string{"asset"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, assetHandler.HandleListUserLikeComments)

	huma.Register(likeGroup, huma.Operation{
		OperationID: "listUserLikeTags",
		Method:      http.MethodGet,
		Path:        "/tags",
		Summary:     "ListUserLikeTags",
		Description: "List tags liked by the current user",
		Tags:        []string{"asset"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, assetHandler.HandleListUserLikeTags)

	// View????
	viewGroup := huma.NewGroup(assetGroup, "/view")

	huma.Register(viewGroup, huma.Operation{
		OperationID: "listUserViewArticles",
		Method:      http.MethodGet,
		Path:        "/articles",
		Summary:     "ListUserViewArticles",
		Description: "List articles viewed by the current user",
		Tags:        []string{"asset"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, assetHandler.HandleListUserViewArticles)

	huma.Register(viewGroup, huma.Operation{
		OperationID: "deleteUserView",
		Method:      http.MethodDelete,
		Path:        "/{viewID}",
		Summary:     "DeleteUserView",
		Description: "Delete a view record",
		Tags:        []string{"asset"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, assetHandler.HandleDeleteUserView)

	// Object????
	objectGroup := huma.NewGroup(assetGroup, "/object")
	objectGroup.UseMiddleware(middleware.LimitUserPermissionMiddlewareForHuma("objectService", model.PermissionCreator))

	huma.Register(objectGroup, huma.Operation{
		OperationID: "listImages",
		Method:      http.MethodGet,
		Path:        "/images",
		Summary:     "ListImages",
		Description: "List all images uploaded by the current user",
		Tags:        []string{"asset"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, assetHandler.HandleListImages)

	// Image????
	imageGroup := huma.NewGroup(objectGroup, "/image")

	uploadImageGroup := huma.NewGroup(imageGroup, "")
	uploadImageGroup.UseMiddleware(middleware.RateLimiterMiddlewareForHuma("uploadImage", constant.CtxKeyUserID, 10*time.Second, 1))

	huma.Register(uploadImageGroup, huma.Operation{
		OperationID: "uploadImage",
		Method:      http.MethodPost,
		Path:        "/",
		Summary:     "UploadImage",
		Description: "Upload an image",
		Tags:        []string{"asset"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, assetHandler.HandleUploadImage)

	huma.Register(imageGroup, huma.Operation{
		OperationID: "getImage",
		Method:      http.MethodGet,
		Path:        "/{objectName}",
		Summary:     "GetImage",
		Description: "Get presigned URL for an image",
		Tags:        []string{"asset"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, assetHandler.HandleGetImage)

	huma.Register(imageGroup, huma.Operation{
		OperationID: "deleteImage",
		Method:      http.MethodDelete,
		Path:        "/{objectName}",
		Summary:     "DeleteImage",
		Description: "Delete an image",
		Tags:        []string{"asset"},
		Security:    []map[string][]string{{"jwtAuth": {}}},
	}, assetHandler.HandleDeleteImage)
}
