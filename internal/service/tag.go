package service

import (
	"errors"
	"time"

	"github.com/hcd233/Aris-blog/internal/logger"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// TagService 标签服务
//
//	author centonhuang
//	update 2025-01-04 17:16:27
type TagService interface {
	CreateTag(req *protocol.CreateTagRequest) (rsp *protocol.CreateTagResponse, err error)
	GetTagInfo(req *protocol.GetTagInfoRequest) (rsp *protocol.GetTagInfoResponse, err error)
	UpdateTag(req *protocol.UpdateTagRequest) (rsp *protocol.UpdateTagResponse, err error)
	DeleteTag(req *protocol.DeleteTagRequest) (rsp *protocol.DeleteTagResponse, err error)
	ListTags(req *protocol.ListTagsRequest) (rsp *protocol.ListTagsResponse, err error)
	ListUserTags(req *protocol.ListUserTagsRequest) (rsp *protocol.ListUserTagsResponse, err error)
	QueryTag(req *protocol.QueryTagRequest) (rsp *protocol.QueryTagResponse, err error)
	QueryUserTag(req *protocol.QueryUserTagRequest) (rsp *protocol.QueryUserTagResponse, err error)
}

type tagService struct {
	db      *gorm.DB
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
		db:      database.GetDBInstance(),
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
func (s *tagService) CreateTag(req *protocol.CreateTagRequest) (rsp *protocol.CreateTagResponse, err error) {
	rsp = &protocol.CreateTagResponse{}
	tag := &model.Tag{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		UserID:      req.CurUserID,
	}

	if err := s.tagDAO.Create(s.db, tag); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, protocol.ErrDataExists
		}
		return nil, protocol.ErrInternalError
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
func (s *tagService) GetTagInfo(req *protocol.GetTagInfoRequest) (rsp *protocol.GetTagInfoResponse, err error) {
	rsp = &protocol.GetTagInfoResponse{}

	tag, err := s.tagDAO.GetBySlug(s.db, req.TagSlug,
		[]string{"id", "name", "slug", "description", "user_id", "created_at", "updated_at", "likes"},
		[]string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[TagService] Tag not found", zap.String("tagSlug", req.TagSlug))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[TagService] Get tag info failed", zap.String("tagSlug", req.TagSlug), zap.Error(err))
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
func (s *tagService) UpdateTag(req *protocol.UpdateTagRequest) (rsp *protocol.UpdateTagResponse, err error) {
	rsp = &protocol.UpdateTagResponse{}

	tag, err := s.tagDAO.GetBySlug(s.db, req.TagSlug, []string{"id", "user_id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[TagService] Tag not found", zap.String("tagSlug", req.TagSlug))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[TagService] Get tag info failed", zap.String("tagSlug", req.TagSlug), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if tag.UserID != req.CurUserID {
		logger.Logger.Info("[TagService] no permission to update tag",
			zap.Uint("curUserID", req.CurUserID),
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
		logger.Logger.Warn("[TagService] No fields to update",
			zap.Uint("curUserID", req.CurUserID),
			zap.String("tagSlug", req.TagSlug),
			zap.Any("updateFields", updateFields))
		return rsp, nil
	}

	if err := s.tagDAO.Update(s.db, tag, updateFields); err != nil {
		logger.Logger.Error("[TagService] Update tag failed",
			zap.String("tagSlug", req.TagSlug),
			zap.Any("updateFields", updateFields),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	tag = lo.Must1(s.tagDAO.GetBySlug(s.db, req.TagSlug, []string{"id", "name", "slug", "description", "user_id"}, []string{"User"}))

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
func (s *tagService) DeleteTag(req *protocol.DeleteTagRequest) (rsp *protocol.DeleteTagResponse, err error) {
	rsp = &protocol.DeleteTagResponse{}

	tag, err := s.tagDAO.GetBySlug(s.db, req.TagName, []string{"id", "name", "slug", "user_id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[TagService] Tag not found", zap.String("tagName", req.TagName))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[TagService] Get tag info failed", zap.String("tagName", req.TagName), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if tag.UserID != req.CurUserID {
		logger.Logger.Info("[TagService] no permission to delete tag",
			zap.Uint("curUserID", req.CurUserID),
			zap.Uint("tagUserID", tag.UserID))
		return nil, protocol.ErrNoPermission
	}

	if err := s.tagDAO.Delete(s.db, tag); err != nil {
		logger.Logger.Error("[TagService] Delete tag failed", zap.String("tagName", req.TagName), zap.Error(err))
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
func (s *tagService) ListTags(req *protocol.ListTagsRequest) (rsp *protocol.ListTagsResponse, err error) {
	rsp = &protocol.ListTagsResponse{}

	tags, pageInfo, err := s.tagDAO.Paginate(s.db,
		[]string{"id", "slug", "name", "description", "user_id", "created_at", "updated_at", "likes"},
		[]string{},
		req.PageParam.Page, req.PageParam.PageSize)
	if err != nil {
		logger.Logger.Error("[TagService] List tags failed", zap.Error(err))
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

// ListUserTags 列出用户标签
//
//	receiver s *tagService
//	param req *protocol.ListUserTagsRequest
//	return rsp *protocol.ListUserTagsResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 12:37:54
func (s *tagService) ListUserTags(*protocol.ListUserTagsRequest) (*protocol.ListUserTagsResponse, error) {
	// TODO: 和合并标签合并
	return nil, protocol.ErrNoImplement
}

// QueryTag
//
//	receiver s *tagService
//	param req *protocol.QueryTagRequest
//	return rsp *protocol.QueryTagResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 11:52:53
func (s *tagService) QueryTag(*protocol.QueryTagRequest) (*protocol.QueryTagResponse, error) {
	// TODO: 和合并标签合并
	return nil, protocol.ErrNoImplement
}

// QueryUserTag
//
//	receiver s *tagService
//	param req *protocol.QueryUserTagRequest
//	return rsp *protocol.QueryUserTagResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 11:52:55
func (s *tagService) QueryUserTag(*protocol.QueryUserTagRequest) (*protocol.QueryUserTagResponse, error) {
	// TODO: 和合并标签合并
	return nil, protocol.ErrInternalError
}
