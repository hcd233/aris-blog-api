package handler

import (
	"context"

	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/protocol/dto"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
)

// UserHandler 用户处理器
//
//	author centonhuang
//	update 2025-01-04 15:56:20
type UserHandler interface {
	HandleGetCurUserInfo(ctx context.Context, req *dto.EmptyRequest) (*protocol.HumaHTTPResponse[*dto.GetCurUserInfoResponse], error)
	HandleGetUserInfo(ctx context.Context, req *dto.GetUserInfoRequest) (*protocol.HumaHTTPResponse[*dto.GetUserResponse], error)
	HandleUpdateInfo(ctx context.Context, req *dto.UpdateUserRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error)
}

type userHandler struct {
	svc service.UserService
}

// NewUserHandler 创建用户处理器
//
//	return UserHandler
//	author centonhuang
//	update 2024-12-08 16:59:38
func NewUserHandler() UserHandler {
	return &userHandler{
		svc: service.NewUserService(),
	}
}

func (h *userHandler) HandleGetCurUserInfo(ctx context.Context, req *dto.EmptyRequest) (*protocol.HumaHTTPResponse[*dto.GetCurUserInfoResponse], error) {
	return util.WrapHTTPResponse(h.svc.GetCurUserInfo(ctx, req))
}

func (h *userHandler) HandleGetUserInfo(ctx context.Context, req *dto.GetUserInfoRequest) (*protocol.HumaHTTPResponse[*dto.GetUserResponse], error) {
	return util.WrapHTTPResponse(h.svc.GetUserInfo(ctx, req))
}

func (h *userHandler) HandleUpdateInfo(ctx context.Context, req *dto.UpdateUserRequest) (*protocol.HumaHTTPResponse[*dto.EmptyResponse], error) {
	return util.WrapHTTPResponse(h.svc.UpdateUserInfo(ctx, req))
}
