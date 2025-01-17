package service

import (
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

	if req.CurUserName != req.UserName {
		logger.Logger.Info("[ArticleVersionService] no permission to create article version",
			zap.String("curUserName", req.CurUserName),
			zap.String("userName", req.UserName))
		return nil, protocol.ErrNoPermission
	}

	user, err := s.userDAO.GetByName(s.db, req.UserName, []string{"id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[ArticleVersionService] user not found", zap.String("userName", req.UserName))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[ArticleVersionService] failed to get user", zap.String("userName", req.UserName), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	article, err := s.articleDAO.GetBySlugAndUserID(s.db, req.ArticleSlug, user.ID, []string{"id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[ArticleVersionService] article not found",
				zap.String("articleSlug", req.ArticleSlug),
				zap.Uint("userID", user.ID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[ArticleVersionService] failed to get article",
			zap.String("articleSlug", req.ArticleSlug),
			zap.Uint("userID", user.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	latestVersion, err := s.articleVersionDAO.GetLatestByArticleID(s.db, article.ID, []string{"version"}, []string{})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Logger.Error("[ArticleVersionService] failed to get latest version",
			zap.Uint("articleID", article.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
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

	return rsp, nil
}

// GetArticleVersionInfo 获取文章版本信息
func (s *articleVersionService) GetArticleVersionInfo(req *protocol.GetArticleVersionInfoRequest) (rsp *protocol.GetArticleVersionInfoResponse, err error) {
	rsp = &protocol.GetArticleVersionInfoResponse{}

	if req.CurUserName != req.UserName {
		logger.Logger.Info("[ArticleVersionService] no permission to get article version",
			zap.String("curUserName", req.CurUserName),
			zap.String("userName", req.UserName))
		return nil, protocol.ErrNoPermission
	}

	user, err := s.userDAO.GetByName(s.db, req.UserName, []string{"id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[ArticleVersionService] user not found", zap.String("userName", req.UserName))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[ArticleVersionService] failed to get user", zap.String("userName", req.UserName), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	article, err := s.articleDAO.GetBySlugAndUserID(s.db, req.ArticleSlug, user.ID, []string{"id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[ArticleVersionService] article not found",
				zap.String("articleSlug", req.ArticleSlug),
				zap.Uint("userID", user.ID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[ArticleVersionService] failed to get article",
			zap.String("articleSlug", req.ArticleSlug),
			zap.Uint("userID", user.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	version, err := s.articleVersionDAO.GetByArticleIDAndVersion(s.db, article.ID, req.Version,
		[]string{"id", "article_id", "version", "content", "created_at", "updated_at"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[ArticleVersionService] version not found",
				zap.Uint("articleID", article.ID),
				zap.Uint("version", req.Version))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[ArticleVersionService] failed to get version",
			zap.Uint("articleID", article.ID),
			zap.Uint("version", req.Version),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.Version = &protocol.ArticleVersion{
		ArticleID:        version.ArticleID,
		ArticleVersionID: version.ID,
		Version:          version.Version,
		Content:          version.Content,
		CreatedAt:        version.CreatedAt.Format(time.DateTime),
		UpdatedAt:        version.UpdatedAt.Format(time.DateTime),
	}

	return rsp, nil
}

// GetLatestArticleVersionInfo 获取最新文章版本信息
func (s *articleVersionService) GetLatestArticleVersionInfo(req *protocol.GetLatestArticleVersionInfoRequest) (rsp *protocol.GetLatestArticleVersionInfoResponse, err error) {
	rsp = &protocol.GetLatestArticleVersionInfoResponse{}

	if req.CurUserName != req.UserName {
		logger.Logger.Info("[ArticleVersionService] no permission to get latest article version",
			zap.String("curUserName", req.CurUserName),
			zap.String("userName", req.UserName))
		return nil, protocol.ErrNoPermission
	}

	user, err := s.userDAO.GetByName(s.db, req.UserName, []string{"id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[ArticleVersionService] user not found", zap.String("userName", req.UserName))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[ArticleVersionService] failed to get user", zap.String("userName", req.UserName), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	article, err := s.articleDAO.GetBySlugAndUserID(s.db, req.ArticleSlug, user.ID, []string{"id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[ArticleVersionService] article not found",
				zap.String("articleSlug", req.ArticleSlug),
				zap.Uint("userID", user.ID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[ArticleVersionService] failed to get article",
			zap.String("articleSlug", req.ArticleSlug),
			zap.Uint("userID", user.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
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
		Version:          version.Version,
		Content:          version.Content,
		CreatedAt:        version.CreatedAt.Format(time.DateTime),
		UpdatedAt:        version.UpdatedAt.Format(time.DateTime),
	}

	return rsp, nil
}

// ListArticleVersions 列出文章版本
func (s *articleVersionService) ListArticleVersions(req *protocol.ListArticleVersionsRequest) (rsp *protocol.ListArticleVersionsResponse, err error) {
	rsp = &protocol.ListArticleVersionsResponse{}

	if req.CurUserName != req.UserName {
		logger.Logger.Info("[ArticleVersionService] no permission to list article versions",
			zap.String("curUserName", req.CurUserName),
			zap.String("userName", req.UserName))
		return nil, protocol.ErrNoPermission
	}

	user, err := s.userDAO.GetByName(s.db, req.UserName, []string{"id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[ArticleVersionService] user not found", zap.String("userName", req.UserName))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[ArticleVersionService] failed to get user", zap.String("userName", req.UserName), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	article, err := s.articleDAO.GetBySlugAndUserID(s.db, req.ArticleSlug, user.ID, []string{"id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[ArticleVersionService] article not found",
				zap.String("articleSlug", req.ArticleSlug),
				zap.Uint("userID", user.ID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[ArticleVersionService] failed to get article",
			zap.String("articleSlug", req.ArticleSlug),
			zap.Uint("userID", user.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
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
			Version:          version.Version,
			Content:          version.Content,
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
