package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-blog-api/internal/handler"
)

func initTokenRouter(tokenGroup *huma.Group) {
	tokenHandler := handler.NewTokenHandler()

	// 刷新令牌
	huma.Register(tokenGroup, huma.Operation{
		OperationID: "refreshToken",
		Method:      http.MethodPost,
		Path:        "/refresh",
		Summary:     "RefreshToken",
		Description: "Refresh the access token using a refresh token",
		Tags:        []string{"token"},
	}, tokenHandler.HandleRefreshToken)
}
