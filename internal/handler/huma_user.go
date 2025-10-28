package handler

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service"
)

// HumaUserHandler Huma 用户处理器
type HumaUserHandler struct {
	svc service.UserService
}

// NewHumaUserHandler 创建 Huma 用户处理器
func NewHumaUserHandler() *HumaUserHandler {
	return &HumaUserHandler{
		svc: service.NewUserService(),
	}
}

// GetCurrentUserInfo 获取当前用户信息
func (h *HumaUserHandler) GetCurrentUserInfo(ctx context.Context, input *protocol.HumaGetCurUserInfoRequest) (*protocol.HumaGetCurUserInfoResponse, error) {
	// 从上下文获取用户ID（这里需要适配 Fiber 的上下文）
	userID, err := protocol.ResolveUserID(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("无法获取用户ID")
	}

	req := &protocol.GetCurUserInfoRequest{
		UserID: userID,
	}

	rsp, err := h.svc.GetCurUserInfo(ctx, req)
	if err != nil {
		return nil, huma.Error500InternalServerError("获取用户信息失败", err)
	}

	return &protocol.HumaGetCurUserInfoResponse{
		Body: struct {
			User *protocol.CurUser `json:"user" doc:"当前用户信息"`
		}{
			User: rsp.User,
		},
	}, nil
}

// GetUserInfo 获取用户信息
func (h *HumaUserHandler) GetUserInfo(ctx context.Context, input *protocol.HumaGetUserInfoRequest) (*protocol.HumaGetUserInfoResponse, error) {
	req := &protocol.GetUserInfoRequest{
		UserID: input.UserID,
	}

	rsp, err := h.svc.GetUserInfo(ctx, req)
	if err != nil {
		return nil, huma.Error500InternalServerError("获取用户信息失败", err)
	}

	return &protocol.HumaGetUserInfoResponse{
		Body: struct {
			User *protocol.User `json:"user" doc:"用户信息"`
		}{
			User: rsp.User,
		},
	}, nil
}

// UpdateUserInfo 更新用户信息
func (h *HumaUserHandler) UpdateUserInfo(ctx context.Context, input *protocol.HumaUpdateUserInfoRequest) (*protocol.HumaUpdateUserInfoResponse, error) {
	// 从上下文获取用户ID
	userID, err := protocol.ResolveUserID(ctx)
	if err != nil {
		return nil, huma.Error400BadRequest("无法获取用户ID")
	}

	req := &protocol.UpdateUserInfoRequest{
		UserID:          userID,
		UpdatedUserName: input.Body.UserName,
	}

	_, err = h.svc.UpdateUserInfo(ctx, req)
	if err != nil {
		return nil, huma.Error500InternalServerError("更新用户信息失败", err)
	}

	return &protocol.HumaUpdateUserInfoResponse{
		Body: struct {
			Success bool `json:"success" example:"true" doc:"更新是否成功"`
		}{
			Success: true,
		},
	}, nil
}

// RegisterUserRoutes 注册用户相关路由
func (h *HumaUserHandler) RegisterRoutes(api huma.API) {
	// 获取当前用户信息
	huma.Register(api, huma.Operation{
		OperationID:   "getCurrentUserInfo",
		Method:        "GET",
		Path:          "/v1/user/current",
		Summary:       "获取当前用户信息",
		Description:   "获取当前登录用户的详细信息",
		Tags:          []string{"user"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "成功获取用户信息",
			},
			"401": {
				Description: "未授权",
			},
			"500": {
				Description: "内部服务器错误",
			},
		},
	}, h.GetCurrentUserInfo)

	// 获取指定用户信息
	huma.Register(api, huma.Operation{
		OperationID:   "getUserInfo",
		Method:        "GET",
		Path:          "/v1/user/{userID}",
		Summary:       "获取用户信息",
		Description:   "根据用户ID获取用户信息",
		Tags:          []string{"user"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "成功获取用户信息",
			},
			"400": {
				Description: "请求参数错误",
			},
			"404": {
				Description: "用户不存在",
			},
			"500": {
				Description: "内部服务器错误",
			},
		},
	}, h.GetUserInfo)

	// 更新用户信息
	huma.Register(api, huma.Operation{
		OperationID:   "updateUserInfo",
		Method:        "PATCH",
		Path:          "/v1/user",
		Summary:       "更新用户信息",
		Description:   "更新当前用户的信息",
		Tags:          []string{"user"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
		Responses: map[string]*huma.Response{
			"200": {
				Description: "成功更新用户信息",
			},
			"400": {
				Description: "请求参数错误",
			},
			"401": {
				Description: "未授权",
			},
			"500": {
				Description: "内部服务器错误",
			},
		},
	}, h.UpdateUserInfo)
}