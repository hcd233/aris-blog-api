package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-blog-api/internal/handler"
)

func initOauth2Router(oauth2Group *huma.Group) {
	oauth2Handler := handler.NewOauth2Handler()

	// OAuth2登录
	huma.Register(oauth2Group, huma.Operation{
		OperationID: "oauth2Login",
		Method:      http.MethodGet,
		Path:        "/{provider}/login",
		Summary:     "OAuth2Login",
		Description: "Get OAuth2 authorization URL for the specified provider (github/google/qq)",
		Tags:        []string{"oauth2"},
	}, oauth2Handler.HandleLogin)

	// OAuth2回调
	huma.Register(oauth2Group, huma.Operation{
		OperationID: "oauth2Callback",
		Method:      http.MethodGet,
		Path:        "/{provider}/callback",
		Summary:     "OAuth2Callback",
		Description: "Handle OAuth2 callback with authorization code and state",
		Tags:        []string{"oauth2"},
	}, oauth2Handler.HandleCallback)
}
