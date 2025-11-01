package api

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
)

var humaAPI huma.API

func init() {
	humaAPI = humafiber.New(GetFiberApp(), huma.Config{
		OpenAPI: &huma.OpenAPI{
			OpenAPI: "3.1.0",
			Info: &huma.Info{
				Title:   "Aris-blog",
				Version: "1.0",
			},
			Components: &huma.Components{
				Schemas: huma.NewMapRegistry("#/components/schemas/", huma.DefaultSchemaNamer),
				SecuritySchemes: map[string]*huma.SecurityScheme{
					"jwtAuth": {
						Type:        "apiKey",
						Name:        "Authorization",
						In:          "header",
						Description: "JWT Authentication，Please pass the JWT token in the Authorization header.",
					},
				},
			},
		},
		OpenAPIPath:   "/openapi",
		DocsPath:      "/docs",
		SchemasPath:   "/schemas",
		Formats:       huma.DefaultFormats,
		DefaultFormat: "application/json",
	})
}

// GetHumaAPI 获取 Huma API 实例
//
//	@return *huma.API
//	@author centonhuang
//	@update 2025-11-02 02:35:59
func GetHumaAPI() huma.API {
	return humaAPI
}
