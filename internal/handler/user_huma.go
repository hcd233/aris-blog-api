package handler

import (
	"context"

	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service"
)

// GetCurUserInfoHuma 获取当前用户信息（Huma 版本）
func GetCurUserInfoHuma(ctx context.Context, input *struct{}) (*protocol.CurUserOutput, error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	req := &protocol.GetCurUserInfoRequest{
		UserID: userID,
	}

	rsp, err := service.NewUserService().GetCurUserInfo(ctx, req)
	if err != nil {
		return nil, protocol.HTTPError(ctx, err)
	}

	return &protocol.CurUserOutput{
		Body: *rsp.User,
	}, nil
}

// GetUserInfoHuma 获取用户信息（Huma 版本）
func GetUserInfoHuma(ctx context.Context, input *protocol.UserInput) (*protocol.UserOutput, error) {
	req := &protocol.GetUserInfoRequest{
		UserID: input.UserID,
	}

	rsp, err := service.NewUserService().GetUserInfo(ctx, req)
	if err != nil {
		return nil, protocol.HTTPError(ctx, err)
	}

	return &protocol.UserOutput{
		Body: *rsp.User,
	}, nil
}

// UpdateUserInfoHuma 更新用户信息（Huma 版本）
func UpdateUserInfoHuma(ctx context.Context, input *protocol.UpdateUserInput) (*protocol.EmptyResponse, error) {
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	req := &protocol.UpdateUserInfoRequest{
		UserID:          userID,
		UpdatedUserName: input.UserName,
	}

	_, err := service.NewUserService().UpdateUserInfo(ctx, req)
	if err != nil {
		return nil, protocol.HTTPError(ctx, err)
	}

	return &protocol.EmptyResponse{}, nil
}
