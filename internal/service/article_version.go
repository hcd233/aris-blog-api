package service

import (
	"errors"
	"time"

	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/logger"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/resource/database"
	"github.com/hcd233/aris-blog-api/internal/resource/database/dao"
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ArticleVersionService 文章版本服务
type ArticleVersionService interface {
	CreateArticleVersion(req *protocol.CreateArticleVersionRequest) (rsp *protocol.CreateArticleVersionResponse, err error)
	GetArticleVersionInfo(req *protocol.GetArticleVersionInfoRequest) (rsp *protocol.GetArticleVersionInfoResponse, err error)
	GetLatestArticleVersionInfo(req *protocol.GetLatestArticleVersionInfoRequest) (rsp *protocol.GetLatestArticleVersionInfoResponse, err error)
	ListArticleVersions(req *protocol.ListArticleVersionsRequest) (rsp *protocol.ListArticleVersionsResponse, err error)
}

type articleVersionService struct {
	db                *gorm.DB
	userDAO           *dao.UserDAO
	articleDAO        *dao.ArticleDAO
	articleVersionDAO *dao.ArticleVersionDAO
}

// NewArticleVersionService 创建文章版本服务
func NewArticleVersionService() ArticleVersionService {
	return &articleVersionService{
		db:                database.GetDBInstance(),
		userDAO:           dao.GetUserDAO(),
		articleDAO:        dao.GetArticleDAO(),
		articleVersionDAO: dao.GetArticleVersionDAO(),
	}
}

// CreateArticleVersion 创建文章版本
func (s *articleVersionService) CreateArticleVersion(req *protocol.CreateArticleVersionRequest) (rsp *protocol.CreateArticleVersionResponse, err error) {
	rsp = &protocol.CreateArticleVersionResponse{}

	article, err := s.articleDAO.GetByID(s.db, req.ArticleID, []string{"id", "user_id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[ArticleVersionService] article not found",
				zap.Uint("articleID", req.ArticleID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[ArticleVersionService] failed to get article",
			zap.Uint("articleID", req.ArticleID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if article.UserID != req.UserID {
		logger.Logger.Error("[ArticleVersionService] no permission to create article version",
			zap.Uint("articleID", req.ArticleID),
			zap.Uint("userID", req.UserID))
		return nil, protocol.ErrNoPermission
	}

	latestVersion, err := s.articleVersionDAO.GetLatestByArticleID(s.db, article.ID, []string{"version", "content"}, []string{})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Logger.Error("[ArticleVersionService] failed to get latest version",
			zap.Uint("articleID", article.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if latestVersion.Content == req.Content {
		logger.Logger.Warn("[ArticleVersionService] content is the same as the latest version",
			zap.Uint("articleID", article.ID),
			zap.Uint("articleVersionID", latestVersion.ID),
			zap.String("content", req.Content))
		return nil, protocol.ErrDataExists
	}

	nextVersion := uint(1)
	if latestVersion != nil {
		nextVersion = latestVersion.Version + 1
	}

	version := &model.ArticleVersion{
		ArticleID: article.ID,
		Version:   nextVersion,
		Content:   req.Content,
	}

	if err := s.articleVersionDAO.Create(s.db, version); err != nil {
		logger.Logger.Error("[ArticleVersionService] failed to create version",
			zap.Uint("articleID", article.ID),
			zap.Uint("version", version.Version),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.ArticleVersion = &protocol.ArticleVersion{
		ArticleID:        version.ArticleID,
		ArticleVersionID: version.ID,
		VersionID:        version.Version,
		Content:          version.Content,
		CreatedAt:        version.CreatedAt.Format(time.DateTime),
		UpdatedAt:        version.UpdatedAt.Format(time.DateTime),
	}

	return rsp, nil
}

// GetArticleVersionInfo 获取文章版本信息
func (s *articleVersionService) GetArticleVersionInfo(req *protocol.GetArticleVersionInfoRequest) (rsp *protocol.GetArticleVersionInfoResponse, err error) {
	rsp = &protocol.GetArticleVersionInfoResponse{}

	article, err := s.articleDAO.GetByID(s.db, req.ArticleID, []string{"id", "user_id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[ArticleVersionService] article not found",
				zap.Uint("articleID", req.ArticleID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[ArticleVersionService] failed to get article",
			zap.Uint("articleID", req.ArticleID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	version, err := s.articleVersionDAO.GetByArticleIDAndVersion(s.db, article.ID, req.VersionID,
		[]string{"id", "article_id", "version", "content", "created_at", "updated_at"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[ArticleVersionService] version not found",
				zap.Uint("articleID", article.ID),
				zap.Uint("versionID", req.VersionID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[ArticleVersionService] failed to get version",
			zap.Uint("articleID", article.ID),
			zap.Uint("versionID", req.VersionID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.Version = &protocol.ArticleVersion{
		ArticleID:        version.ArticleID,
		ArticleVersionID: version.ID,
		VersionID:        version.Version,
		Content:          version.Content,
		CreatedAt:        version.CreatedAt.Format(time.DateTime),
		UpdatedAt:        version.UpdatedAt.Format(time.DateTime),
	}

	return rsp, nil
}

// GetLatestArticleVersionInfo 获取最新文章版本信息
func (s *articleVersionService) GetLatestArticleVersionInfo(req *protocol.GetLatestArticleVersionInfoRequest) (rsp *protocol.GetLatestArticleVersionInfoResponse, err error) {
	rsp = &protocol.GetLatestArticleVersionInfoResponse{}

	article, err := s.articleDAO.GetByID(s.db, req.ArticleID, []string{"id", "user_id", "status"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[ArticleVersionService] article not found",
				zap.Uint("articleID", req.ArticleID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[ArticleVersionService] failed to get article",
			zap.Uint("articleID", req.ArticleID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	// 如果文章不是公开的，则只有作者本人可以查看
	if article.UserID != req.UserID && article.Status != model.ArticleStatusPublish {
		logger.Logger.Error("[ArticleVersionService] no permission to get latest article version",
			zap.Uint("articleID", req.ArticleID),
			zap.Uint("userID", req.UserID))
		return nil, protocol.ErrNoPermission
	}

	version, err := s.articleVersionDAO.GetLatestByArticleID(s.db, article.ID,
		[]string{"id", "article_id", "version", "content", "created_at", "updated_at"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[ArticleVersionService] latest version not found",
				zap.Uint("articleID", article.ID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[ArticleVersionService] failed to get latest version",
			zap.Uint("articleID", article.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.Version = &protocol.ArticleVersion{
		ArticleID:        version.ArticleID,
		ArticleVersionID: version.ID,
		VersionID:        version.Version,
		Content:          version.Content,
		CreatedAt:        version.CreatedAt.Format(time.DateTime),
		UpdatedAt:        version.UpdatedAt.Format(time.DateTime),
	}

	return rsp, nil
}

// ListArticleVersions 列出文章版本
func (s *articleVersionService) ListArticleVersions(req *protocol.ListArticleVersionsRequest) (rsp *protocol.ListArticleVersionsResponse, err error) {
	rsp = &protocol.ListArticleVersionsResponse{}

	article, err := s.articleDAO.GetByID(s.db, req.ArticleID, []string{"id", "user_id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[ArticleVersionService] article not found",
				zap.Uint("articleID", req.ArticleID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[ArticleVersionService] failed to get article",
			zap.Uint("articleID", req.ArticleID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if article.UserID != req.UserID {
		logger.Logger.Error("[ArticleVersionService] no permission to list article versions",
			zap.Uint("articleID", req.ArticleID),
			zap.Uint("userID", req.UserID))
		return nil, protocol.ErrNoPermission
	}

	versions, pageInfo, err := s.articleVersionDAO.PaginateByArticleID(s.db, article.ID,
		[]string{"id", "article_id", "version", "content", "created_at", "updated_at"}, []string{},
		req.PageParam.Page, req.PageParam.PageSize)
	if err != nil {
		logger.Logger.Error("[ArticleVersionService] failed to paginate versions",
			zap.Uint("articleID", article.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.Versions = lo.Map(*versions, func(version model.ArticleVersion, _ int) *protocol.ArticleVersion {
		return &protocol.ArticleVersion{
			ArticleID:        version.ArticleID,
			ArticleVersionID: version.ID,
			VersionID:        version.Version,
			Content:          string([]rune(version.Content)[:constant.ListArticleVersionContentLength]) + "...",
			CreatedAt:        version.CreatedAt.Format(time.DateTime),
			UpdatedAt:        version.UpdatedAt.Format(time.DateTime),
		}
	})

	rsp.PageInfo = &protocol.PageInfo{
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
		Total:    pageInfo.Total,
	}

	return rsp, nil
}
