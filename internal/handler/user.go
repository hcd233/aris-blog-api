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
	HandleGetCurrentUser(ctx context.Context, req *dto.EmptyRequest) (*protocol.HTTPResponse[*dto.GetCurrentUserResponse], error)
	HandleGetUser(ctx context.Context, req *dto.GetUserRequest) (*protocol.HTTPResponse[*dto.GetUserResponse], error)
	HandleUpdateUser(ctx context.Context, req *dto.UpdateUserRequest) (*protocol.HTTPResponse[*dto.EmptyResponse], error)
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

func (h *userHandler) HandleGetCurrentUser(ctx context.Context, req *dto.EmptyRequest) (*protocol.HTTPResponse[*dto.GetCurrentUserResponse], error) {
	return util.WrapHTTPResponse(h.svc.GetCurUserInfo(ctx, req))
}

func (h *userHandler) HandleGetUser(ctx context.Context, req *dto.GetUserRequest) (*protocol.HTTPResponse[*dto.GetUserResponse], error) {
	return util.WrapHTTPResponse(h.svc.GetUserInfo(ctx, req))
}

func (h *userHandler) HandleUpdateUser(ctx context.Context, req *dto.UpdateUserRequest) (*protocol.HTTPResponse[*dto.EmptyResponse], error) {
	return util.WrapHTTPResponse(h.svc.UpdateUserInfo(ctx, req))
}
