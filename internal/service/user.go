package service

import (
	"errors"
	"time"

	"github.com/hcd233/Aris-blog/internal/logger"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// UserService 用户服务
//
//	author centonhuang
//	update 2025-01-04 21:04:00
type UserService interface {
	GetCurUserInfo(req *protocol.GetCurUserInfoRequest) (rsp *protocol.GetCurUserInfoResponse, err error)
	GetUserInfo(req *protocol.GetUserInfoRequest) (rsp *protocol.GetUserInfoResponse, err error)
	UpdateUserInfo(req *protocol.UpdateUserInfoRequest) (rsp *protocol.UpdateUserInfoResponse, err error)
	QueryUser(req *protocol.QueryUserRequest) (rsp *protocol.QueryUserResponse, err error)
}

type userService struct {
	db         *gorm.DB
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
		db:         database.GetDBInstance(),
		userDAO:    dao.GetUserDAO(),
		tagDAO:     dao.GetTagDAO(),
		articleDAO: dao.GetArticleDAO(),
	}
}

// GetCurUserInfo 获取当前用户信息
//
//	receiver s *userService
//	param req *protocol.GetCurUserInfoRequest
//	return rsp *protocol.GetCurUserInfoResponse
//	return err error
//	author centonhuang
//	update 2025-01-04 21:04:03
func (s *userService) GetCurUserInfo(req *protocol.GetCurUserInfoRequest) (rsp *protocol.GetCurUserInfoResponse, err error) {
	rsp = &protocol.GetCurUserInfoResponse{}

	user, err := s.userDAO.GetByID(s.db, req.CurUserID, []string{"id", "name", "email", "avatar", "created_at", "last_login", "permission"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[UserService] user not found", zap.Uint("userID", req.CurUserID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[UserService] failed to get user by id", zap.Uint("userID", req.CurUserID), zap.Error(err))
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

	logger.Logger.Info("[UserService] get cur user info",
		zap.Uint("userID", user.ID),
		zap.String("name", user.Name),
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
//	param req *protocol.GetUserInfoRequest
//	return *protocol.GetUserInfoResponse
//	return error
//	author centonhuang
//	update 2025-01-04 21:09:04
func (s *userService) GetUserInfo(req *protocol.GetUserInfoRequest) (rsp *protocol.GetUserInfoResponse, err error) {
	rsp = &protocol.GetUserInfoResponse{}

	user, err := s.userDAO.GetByName(s.db, req.UserName, []string{"id", "name", "email", "avatar", "created_at", "last_login", "permission"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[UserService] user not found", zap.String("userName", req.UserName))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[UserService] failed to get user by name", zap.String("userName", req.UserName), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	logger.Logger.Info("[UserService] get user info",
		zap.Uint("userID", user.ID),
		zap.String("name", user.Name),
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

func (s *userService) UpdateUserInfo(req *protocol.UpdateUserInfoRequest) (rsp *protocol.UpdateUserInfoResponse, err error) {
	if req.CurUserName != req.UserName {
		logger.Logger.Info("[UserService] no permission to update user info",
			zap.String("curUserName", req.CurUserName),
			zap.String("userName", req.UserName))
		return nil, protocol.ErrNoPermission
	}

	rsp = &protocol.UpdateUserInfoResponse{}

	user, err := s.userDAO.GetByName(s.db, req.UserName, []string{"id", "name"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[UserService] user not found", zap.String("userName", req.UserName))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[UserService] failed to get user by name", zap.String("userName", req.UserName), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if err := s.userDAO.Update(s.db, &model.User{ID: user.ID}, map[string]interface{}{
		"name": req.UpdatedUserName,
	}); err != nil {
		logger.Logger.Error("[UserService] failed to update user", zap.Uint("userID", user.ID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	user, err = s.userDAO.GetByID(s.db, user.ID, []string{"id", "name", "avatar"}, []string{})
	if err != nil {
		logger.Logger.Error("[UserService] failed to get user after update", zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	return rsp, nil
}

func (s *userService) QueryUser(*protocol.QueryUserRequest) (*protocol.QueryUserResponse, error) {
	// users, queryInfo, err := s.userDocDAO.QueryDocument(req.Query, req.Filter, req.Page, req.PageSize)
	// if err != nil {
	// 	logger.Logger.Error("[UserService] failed to query users", zap.Error(err))
	// 	return nil, protocol.ErrInternalError
	// }

	// TODO: 查询用户和列出用户合并
	return nil, protocol.ErrInternalError
}
