package service

import (
	"context"
	"errors"
	"time"

	"github.com/hcd233/aris-blog-api/internal/logger"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	dto "github.com/hcd233/aris-blog-api/internal/protocol/dto"
	"github.com/hcd233/aris-blog-api/internal/resource/database"
	"github.com/hcd233/aris-blog-api/internal/resource/database/dao"
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// TagService 标签服务
//
//  author centonhuang
//  update 2025-10-31 05:45:00
type TagService interface {
	CreateTag(ctx context.Context, req *dto.TagCreateRequest) (rsp *dto.TagCreateResponse, err error)
	GetTagInfo(ctx context.Context, req *dto.TagGetRequest) (rsp *dto.TagGetResponse, err error)
	UpdateTag(ctx context.Context, req *dto.TagUpdateRequest) (rsp *dto.TagUpdateResponse, err error)
	DeleteTag(ctx context.Context, req *dto.TagDeleteRequest) (rsp *dto.TagDeleteResponse, err error)
	ListTags(ctx context.Context, req *dto.TagListRequest) (rsp *dto.TagListResponse, err error)
}

type tagService struct {
	userDAO *dao.UserDAO
	tagDAO  *dao.TagDAO
}

// NewTagService 创建标签服务
func NewTagService() TagService {
	return &tagService{
		userDAO: dao.GetUserDAO(),
		tagDAO:  dao.GetTagDAO(),
	}
}

// CreateTag 创建标签
func (s *tagService) CreateTag(ctx context.Context, req *dto.TagCreateRequest) (rsp *dto.TagCreateResponse, err error) {
	if req == nil || req.Body == nil {
		return nil, protocol.ErrBadRequest
	}

	rsp = &dto.TagCreateResponse{}

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	tag := &model.Tag{
		Name:        req.Body.Name,
		Slug:        req.Body.Slug,
		Description: req.Body.Description,
		UserID:      req.UserID,
	}

	if err := s.tagDAO.Create(db, tag); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			logger.Error("[TagService] tag already exists", zap.String("name", req.Body.Name), zap.String("slug", req.Body.Slug), zap.Error(err))
			return nil, protocol.ErrDataExists
		}
		logger.Error("[TagService] failed to create tag", zap.String("name", req.Body.Name), zap.String("slug", req.Body.Slug), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.Tag = &dto.Tag{
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
func (s *tagService) GetTagInfo(ctx context.Context, req *dto.TagGetRequest) (rsp *dto.TagGetResponse, err error) {
	rsp = &dto.TagGetResponse{}

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	tag, err := s.tagDAO.GetByID(db, req.TagID,
		[]string{"id", "name", "slug", "description", "user_id", "created_at", "updated_at", "likes"},
		[]string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[TagService] tag not found", zap.Uint("tagID", req.TagID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[TagService] get tag info failed", zap.Uint("tagID", req.TagID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.Tag = &dto.Tag{
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
func (s *tagService) UpdateTag(ctx context.Context, req *dto.TagUpdateRequest) (rsp *dto.TagUpdateResponse, err error) {
	if req == nil || req.Body == nil {
		return nil, protocol.ErrBadRequest
	}

	rsp = &dto.TagUpdateResponse{}

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	tag, err := s.tagDAO.GetByID(db, req.TagID, []string{"id", "user_id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[TagService] tag not found", zap.Uint("tagID", req.TagID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[TagService] get tag info failed", zap.Uint("tagID", req.TagID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if tag.UserID != req.UserID {
		logger.Error("[TagService] no permission to update tag",
			zap.Uint("tagUserID", tag.UserID))
		return nil, protocol.ErrNoPermission
	}

	updateFields := make(map[string]interface{})
	if req.Body.Name != "" {
		updateFields["name"] = req.Body.Name
	}
	if req.Body.Slug != "" {
		updateFields["slug"] = req.Body.Slug
	}
	if req.Body.Description != "" {
		updateFields["description"] = req.Body.Description
	}

	if len(updateFields) == 0 {
		logger.Warn("[TagService] no fields to update",
			zap.Uint("tagID", req.TagID))
		return rsp, nil
	}

	if err := s.tagDAO.Update(db, tag, updateFields); err != nil {
		logger.Error("[TagService] update tag failed",
			zap.Uint("tagID", req.TagID),
			zap.Any("updateFields", updateFields),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	return rsp, nil
}

// DeleteTag 删除标签
func (s *tagService) DeleteTag(ctx context.Context, req *dto.TagDeleteRequest) (rsp *dto.TagDeleteResponse, err error) {
	rsp = &dto.TagDeleteResponse{}

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	tag, err := s.tagDAO.GetByID(db, req.TagID, []string{"id", "name", "slug", "user_id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[TagService] tag not found", zap.Uint("tagID", req.TagID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[TagService] get tag info failed", zap.Uint("tagID", req.TagID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if tag.UserID != req.UserID {
		logger.Error("[TagService] no permission to delete tag",
			zap.Uint("tagUserID", tag.UserID))
		return nil, protocol.ErrNoPermission
	}

	if err := s.tagDAO.Delete(db, tag); err != nil {
		logger.Error("[TagService] delete tag failed", zap.Uint("tagID", req.TagID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	return rsp, nil
}

// ListTags 列出标签
func (s *tagService) ListTags(ctx context.Context, req *dto.TagListRequest) (rsp *dto.TagListResponse, err error) {
	rsp = &dto.TagListResponse{}

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	paginate := req.PaginationQuery.ToPaginateParam()
	param := &dao.PaginateParam{
		PageParam: &dao.PageParam{
			Page:     paginate.PageParam.Page,
			PageSize: paginate.PageParam.PageSize,
		},
		QueryParam: &dao.QueryParam{
			Query:       paginate.QueryParam.Query,
			QueryFields: []string{"name", "description"},
		},
	}

	tags, pageInfo, err := s.tagDAO.Paginate(db,
		[]string{"id", "slug", "name", "description", "user_id", "created_at", "updated_at", "likes"},
		[]string{},
		param,
	)
	if err != nil {
		logger.Error("[TagService] list tags failed", zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.Tags = lo.Map(*tags, func(tag model.Tag, _ int) *dto.Tag {
		return &dto.Tag{
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

	rsp.PageInfo = &dto.PageInfo{
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
		Total:    pageInfo.Total,
	}

	return rsp, nil
}
