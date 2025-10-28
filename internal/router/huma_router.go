package router

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
	"github.com/hcd233/aris-blog-api/internal/handler"
	"github.com/hcd233/aris-blog-api/internal/middleware"
)

// RegisterHumaRouter 注册 Huma 路由
func RegisterHumaRouter(app *fiber.App) {
	// 创建 Huma 配置
	config := huma.DefaultConfig("Aris Blog API", "1.0.0")
	config.Info.Description = "一个现代化的博客 API 服务"
	config.Info.Contact = &huma.Contact{
		Name:  "centonhuang",
		Email: "example@example.com",
	}
	config.Info.License = &huma.License{
		Name: "MIT",
		URL:  "https://opensource.org/licenses/MIT",
	}
	
	// 添加服务器信息
	config.Servers = []*huma.Server{
		{
			URL:         "http://localhost:8080",
			Description: "开发服务器",
		},
	}

	// 添加安全定义
	config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"bearerAuth": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
			Description:  "JWT Bearer 认证",
		},
	}

	// 创建 Huma API 实例
	api := humafiber.New(app, config)

	// 应用认证中间件
	middleware.HumaJWTMiddleware(api)

	// 注册健康检查路由
	registerPingRoutes(api)

	// 注册用户路由
	registerUserRoutes(api)

	// 注册标签路由 (未来扩展)
	// registerTagRoutes(api)

	// 注册文章路由 (未来扩展)
	// registerArticleRoutes(api)

	// 添加 OpenAPI 文档路由
	// Huma 会自动在以下路径提供 OpenAPI 文档:
	// - GET /openapi.json - OpenAPI 3.0 JSON 文档
	// - GET /openapi.yaml - OpenAPI 3.0 YAML 文档
	// - GET /docs - 交互式文档 (Swagger UI)
}

// registerPingRoutes 注册健康检查路由
func registerPingRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "ping",
		Method:      "GET",
		Path:        "/",
		Summary:     "健康检查",
		Description: "检查服务是否正常运行",
		Tags:        []string{"health"},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "服务正常",
			},
		},
	}, func(ctx context.Context, input *struct{}) (*struct {
		Body struct {
			Status string `json:"status" example:"ok" doc:"服务状态"`
		}
	}, error) {
		return &struct {
			Body struct {
				Status string `json:"status" example:"ok" doc:"服务状态"`
			}
		}{
			Body: struct {
				Status string `json:"status" example:"ok" doc:"服务状态"`
			}{
				Status: "ok",
			},
		}, nil
	})
}

// registerUserRoutes 注册用户路由
func registerUserRoutes(api huma.API) {
	userHandler := handler.NewHumaUserHandler()
	userHandler.RegisterRoutes(api)
}