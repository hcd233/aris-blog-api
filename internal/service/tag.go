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
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// TagService 标签服务
//
//	author centonhuang
//	update 2025-01-04 17:16:27
type TagService interface {
	CreateTag(ctx context.Context, req *protocol.CreateTagRequest) (rsp *protocol.CreateTagResponse, err error)
	GetTagInfo(ctx context.Context, req *protocol.GetTagInfoRequest) (rsp *protocol.GetTagInfoResponse, err error)
	UpdateTag(ctx context.Context, req *protocol.UpdateTagRequest) (rsp *protocol.UpdateTagResponse, err error)
	DeleteTag(ctx context.Context, req *protocol.DeleteTagRequest) (rsp *protocol.DeleteTagResponse, err error)
	ListTags(ctx context.Context, req *protocol.ListTagsRequest) (rsp *protocol.ListTagsResponse, err error)
}

type tagService struct {
	userDAO *dao.UserDAO
	tagDAO  *dao.TagDAO
}

// NewTagService 创建标签服务
//
//	return TagService
//	author centonhuang
//	update 2025-01-05 11:48:36
func NewTagService() TagService {
	return &tagService{
		userDAO: dao.GetUserDAO(),
		tagDAO:  dao.GetTagDAO(),
	}
}

// CreateTag 创建标签
//
//	receiver s *tagService
//	param req *protocol.CreateTagRequest
//	return rsp *protocol.CreateTagResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 11:52:33
func (s *tagService) CreateTag(ctx context.Context, req *protocol.CreateTagRequest) (rsp *protocol.CreateTagResponse, err error) {
	rsp = &protocol.CreateTagResponse{}

	logger := logger.LoggerWithContext(ctx)
	db := database.GetDBInstance(ctx)

	tag := &model.Tag{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		UserID:      req.UserID,
	}

	if err := s.tagDAO.Create(db, tag); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			logger.Error("[TagService] Tag is already exists", zap.String("name", req.Name), zap.String("slug", req.Slug), zap.Error(err))
			return nil, protocol.ErrDataExists
		}
		logger.Error("[TagService] Failed to create tag", zap.String("name", req.Name), zap.String("slug", req.Slug), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.Tag = &protocol.Tag{
		TagID:       tag.ID,
		Name:        tag.Name,
		Slug:        tag.Slug,
		Description: tag.Description,
		UserID:      tag.UserID,
		CreatedAt:   tag.CreatedAt.Format(time.DateTime),
		UpdatedAt:   tag.UpdatedAt.Format(time.DateTime),
		Likes:       tag.Likes,
	}

	return rsp, nil
}

// GetTagInfo 获取标签信息
//
//	receiver s *tagService
//	param req *protocol.GetTagInfoRequest
//	return rsp *protocol.GetTagInfoResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 11:52:46
func (s *tagService) GetTagInfo(ctx context.Context, req *protocol.GetTagInfoRequest) (rsp *protocol.GetTagInfoResponse, err error) {
	rsp = &protocol.GetTagInfoResponse{}

	logger := logger.LoggerWithContext(ctx)
	db := database.GetDBInstance(ctx)

	tag, err := s.tagDAO.GetByID(db, req.TagID,
		[]string{"id", "name", "slug", "description", "user_id", "created_at", "updated_at", "likes"},
		[]string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[TagService] Tag not found", zap.Uint("tagID", req.TagID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[TagService] Get tag info failed", zap.Uint("tagID", req.TagID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.Tag = &protocol.Tag{
		TagID:       tag.ID,
		Name:        tag.Name,
		Slug:        tag.Slug,
		Description: tag.Description,
		UserID:      tag.UserID,
		CreatedAt:   tag.CreatedAt.Format(time.DateTime),
		UpdatedAt:   tag.UpdatedAt.Format(time.DateTime),
		Likes:       tag.Likes,
	}

	return rsp, nil
}

// UpdateTag 更新标签
//
//	receiver s *tagService
//	param req *protocol.UpdateTagRequest
//	return rsp *protocol.UpdateTagResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 11:52:46
func (s *tagService) UpdateTag(ctx context.Context, req *protocol.UpdateTagRequest) (rsp *protocol.UpdateTagResponse, err error) {
	rsp = &protocol.UpdateTagResponse{}

	logger := logger.LoggerWithContext(ctx)
	db := database.GetDBInstance(ctx)

	tag, err := s.tagDAO.GetByID(db, req.TagID, []string{"id", "user_id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[TagService] Tag not found", zap.Uint("tagID", req.TagID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[TagService] Get tag info failed", zap.Uint("tagID", req.TagID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if tag.UserID != req.UserID {
		logger.Error("[TagService] no permission to update tag",
			zap.Uint("tagUserID", tag.UserID))
		return nil, protocol.ErrNoPermission
	}

	updateFields := make(map[string]interface{})
	if req.Name != "" {
		updateFields["name"] = req.Name
	}
	if req.Slug != "" {
		updateFields["slug"] = req.Slug
	}
	if req.Description != "" {
		updateFields["description"] = req.Description
	}

	if len(updateFields) == 0 {
		logger.Warn("[TagService] No fields to update",
			zap.Uint("tagID", req.TagID),
			zap.Any("updateFields", updateFields))
		return rsp, nil
	}

	if err := s.tagDAO.Update(db, tag, updateFields); err != nil {
		logger.Error("[TagService] Update tag failed",
			zap.Uint("tagID", req.TagID),
			zap.Any("updateFields", updateFields),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	return rsp, nil
}

// DeleteTag
//
//	receiver s *tagService
//	param req *protocol.DeleteTagRequest
//	return rsp *protocol.DeleteTagResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 11:52:48
func (s *tagService) DeleteTag(ctx context.Context, req *protocol.DeleteTagRequest) (rsp *protocol.DeleteTagResponse, err error) {
	rsp = &protocol.DeleteTagResponse{}

	logger := logger.LoggerWithContext(ctx)
	db := database.GetDBInstance(ctx)

	tag, err := s.tagDAO.GetByID(db, req.TagID, []string{"id", "name", "slug", "user_id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[TagService] Tag not found", zap.Uint("tagID", req.TagID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[TagService] Get tag info failed", zap.Uint("tagID", req.TagID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if tag.UserID != req.UserID {
		logger.Error("[TagService] no permission to delete tag",
			zap.Uint("tagUserID", tag.UserID))
		return nil, protocol.ErrNoPermission
	}

	if err := s.tagDAO.Delete(db, tag); err != nil {
		logger.Error("[TagService] Delete tag failed", zap.Uint("tagID", req.TagID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	return rsp, nil
}

// ListTags
//
//	receiver s *tagService
//	param req *protocol.ListTagsRequest
//	return rsp *protocol.ListTagsResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 11:52:50
func (s *tagService) ListTags(ctx context.Context, req *protocol.ListTagsRequest) (rsp *protocol.ListTagsResponse, err error) {
	rsp = &protocol.ListTagsResponse{}

	logger := logger.LoggerWithContext(ctx)
	db := database.GetDBInstance(ctx)

	tags, pageInfo, err := s.tagDAO.Paginate(db,
		[]string{"id", "slug", "name", "description", "user_id", "created_at", "updated_at", "likes"},
		[]string{},
		req.PageParam.Page, req.PageParam.PageSize)
	if err != nil {
		logger.Error("[TagService] List tags failed", zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.Tags = lo.Map(*tags, func(tag model.Tag, _ int) *protocol.Tag {
		return &protocol.Tag{
			TagID:       tag.ID,
			Name:        tag.Name,
			Slug:        tag.Slug,
			Description: tag.Description,
			UserID:      tag.UserID,
			CreatedAt:   tag.CreatedAt.Format(time.DateTime),
			UpdatedAt:   tag.UpdatedAt.Format(time.DateTime),
			Likes:       tag.Likes,
		}
	})

	rsp.PageInfo = &protocol.PageInfo{
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
		Total:    pageInfo.Total,
	}

	return rsp, nil
}
