package humahandler

import (
	"context"

	humadto "github.com/hcd233/aris-blog-api/internal/huma"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service"
)

type UserHandlers struct{ svc service.UserService }

func NewUserHandlers() *UserHandlers { return &UserHandlers{svc: service.NewUserService()} }

type (
	userPathInput   struct{ humadto.UserPath }
	updateUserInput struct {
		authHeader
		humadto.UpdateUserInput
	}
)

func (h *UserHandlers) HandleGetCurUserInfo(ctx context.Context, a *authHeader) (*humadto.Output[protocol.GetCurUserInfoResponse], error) {
	req := &protocol.GetCurUserInfoRequest{UserID: a.UserID}
	rsp, err := h.svc.GetCurUserInfo(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.GetCurUserInfoResponse]{Body: *rsp}, nil
}

func (h *UserHandlers) HandleGetUserInfo(ctx context.Context, input *userPathInput) (*humadto.Output[protocol.GetUserInfoResponse], error) {
	req := &protocol.GetUserInfoRequest{UserID: input.UserID}
	rsp, err := h.svc.GetUserInfo(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.GetUserInfoResponse]{Body: *rsp}, nil
}

func (h *UserHandlers) HandleUpdateInfo(ctx context.Context, input *updateUserInput) (*humadto.Output[protocol.UpdateUserInfoResponse], error) {
	req := &protocol.UpdateUserInfoRequest{UserID: input.UserID, UpdatedUserName: input.Body.UserName}
	rsp, err := h.svc.UpdateUserInfo(ctx, req)
	if err != nil {
		return nil, err
	}
	return &humadto.Output[protocol.UpdateUserInfoResponse]{Body: *rsp}, nil
}
