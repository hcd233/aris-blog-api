package service

import (
	"context"
	"errors"
	"time"

	"github.com/hcd233/aris-blog-api/internal/logger"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/resource/database"
	"github.com/hcd233/aris-blog-api/internal/resource/database/dao"
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// UserService 用户服务
//
//	author centonhuang
//	update 2025-01-04 21:04:00
type UserService interface {
	GetCurUserInfo(ctx context.Context, req *protocol.GetCurUserInfoRequest) (rsp *protocol.GetCurUserInfoResponse, err error)
	GetUserInfo(ctx context.Context, req *protocol.GetUserInfoRequest) (rsp *protocol.GetUserInfoResponse, err error)
	UpdateUserInfo(ctx context.Context, req *protocol.UpdateUserInfoRequest) (rsp *protocol.UpdateUserInfoResponse, err error)
}

type userService struct {
	userDAO    *dao.UserDAO
	tagDAO     *dao.TagDAO
	articleDAO *dao.ArticleDAO
}

// NewUserService 创建用户服务
//
//	return UserService
//	author centonhuang
//	update 2025-01-04 21:03:45
func NewUserService() UserService {
	return &userService{
		userDAO:    dao.GetUserDAO(),
		tagDAO:     dao.GetTagDAO(),
		articleDAO: dao.GetArticleDAO(),
	}
}

// GetCurUserInfo 获取当前用户信息
//
//	receiver s *userService
//	param ctx context.Context
//	param req *protocol.GetCurUserInfoRequest
//	return rsp *protocol.GetCurUserInfoResponse
//	return err error
//	author centonhuang
//	update 2025-01-04 21:04:03
func (s *userService) GetCurUserInfo(ctx context.Context, req *protocol.GetCurUserInfoRequest) (rsp *protocol.GetCurUserInfoResponse, err error) {
	rsp = &protocol.GetCurUserInfoResponse{}

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	user, err := s.userDAO.GetByID(db, req.UserID, []string{"id", "name", "email", "avatar", "created_at", "last_login", "permission"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[UserService] user not found")
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[UserService] failed to get user by id", zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.User = &protocol.CurUser{
		User: protocol.User{
			UserID:    user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Avatar:    user.Avatar,
			CreatedAt: user.CreatedAt.Format(time.DateTime),
			LastLogin: user.LastLogin.Format(time.DateTime),
		},
		Permission: string(user.Permission),
	}

	logger.Info("[UserService] get cur user info",
		zap.String("email", user.Email),
		zap.String("avatar", user.Avatar),
		zap.Time("createdAt", user.CreatedAt),
		zap.Time("lastLogin", user.LastLogin),
		zap.String("permission", string(user.Permission)))

	return rsp, nil
}

// GetUserInfo 获取用户信息
//
//	receiver s *userService
//	param ctx context.Context
//	param req *protocol.GetUserInfoRequest
//	return *protocol.GetUserInfoResponse
//	return error
//	author centonhuang
//	update 2025-01-04 21:09:04
func (s *userService) GetUserInfo(ctx context.Context, req *protocol.GetUserInfoRequest) (rsp *protocol.GetUserInfoResponse, err error) {
	logger := logger.WithCtx(ctx)

	rsp = &protocol.GetUserInfoResponse{}
	db := database.GetDBInstance(ctx)

	user, err := s.userDAO.GetByID(db, req.UserID, []string{"id", "name", "email", "avatar", "created_at", "last_login", "permission"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[UserService] user not found")
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[UserService] failed to get user by id", zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	logger.Info("[UserService] get user info",
		zap.String("email", user.Email),
		zap.String("avatar", user.Avatar),
		zap.Time("createdAt", user.CreatedAt),
		zap.Time("lastLogin", user.LastLogin))

	rsp.User = &protocol.User{
		UserID:    user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt.Format(time.DateTime),
		LastLogin: user.LastLogin.Format(time.DateTime),
	}

	return rsp, nil
}

func (s *userService) UpdateUserInfo(ctx context.Context, req *protocol.UpdateUserInfoRequest) (rsp *protocol.UpdateUserInfoResponse, err error) {
	logger := logger.WithCtx(ctx)

	rsp = &protocol.UpdateUserInfoResponse{}
	db := database.GetDBInstance(ctx)

	if err := s.userDAO.Update(db, &model.User{ID: req.UserID}, map[string]interface{}{
		"name": req.UpdatedUserName,
	}); err != nil {
		logger.Error("[UserService] failed to update user", zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	return rsp, nil
}
