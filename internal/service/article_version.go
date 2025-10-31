package service

import (
	"context"
	"errors"
	"time"

	"github.com/hcd233/aris-blog-api/internal/constant"
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

// ArticleVersionService 文章版本服务
type ArticleVersionService interface {
	CreateArticleVersion(ctx context.Context, req *dto.ArticleVersionCreateRequest) (rsp *dto.ArticleVersionCreateResponse, err error)
	GetArticleVersionInfo(ctx context.Context, req *dto.ArticleVersionGetRequest) (rsp *dto.ArticleVersionGetResponse, err error)
	GetLatestArticleVersionInfo(ctx context.Context, req *dto.ArticleVersionGetLatestRequest) (rsp *dto.ArticleVersionGetLatestResponse, err error)
	ListArticleVersions(ctx context.Context, req *dto.ArticleVersionListRequest) (rsp *dto.ArticleVersionListResponse, err error)
}

type articleVersionService struct {
	userDAO           *dao.UserDAO
	articleDAO        *dao.ArticleDAO
	articleVersionDAO *dao.ArticleVersionDAO
}

// NewArticleVersionService 创建文章版本服务
func NewArticleVersionService() ArticleVersionService {
	return &articleVersionService{
		userDAO:           dao.GetUserDAO(),
		articleDAO:        dao.GetArticleDAO(),
		articleVersionDAO: dao.GetArticleVersionDAO(),
	}
}

// CreateArticleVersion 创建文章版本
func (s *articleVersionService) CreateArticleVersion(ctx context.Context, req *dto.ArticleVersionCreateRequest) (rsp *dto.ArticleVersionCreateResponse, err error) {
	logger := logger.WithCtx(ctx)

	if req == nil || req.Body == nil {
		logger.Error("[ArticleVersionService] request body is nil")
		return nil, protocol.ErrBadRequest
	}

	rsp = &dto.ArticleVersionCreateResponse{}

	db := database.GetDBInstance(ctx)
	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	article, err := s.articleDAO.GetByID(db, req.ArticleID, []string{"id", "user_id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[ArticleVersionService] article not found",
				zap.Uint("articleID", req.ArticleID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[ArticleVersionService] failed to get article",
			zap.Uint("articleID", req.ArticleID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if article.UserID != userID {
		logger.Error("[ArticleVersionService] no permission to create article version",
			zap.Uint("articleID", req.ArticleID))
		return nil, protocol.ErrNoPermission
	}

	latestVersion, err := s.articleVersionDAO.GetLatestByArticleID(db, article.ID, []string{"version", "content"}, []string{})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("[ArticleVersionService] failed to get latest version",
			zap.Uint("articleID", article.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if latestVersion != nil && latestVersion.Content == req.Body.Content {
		logger.Warn("[ArticleVersionService] content is the same as the latest version",
			zap.Uint("articleID", article.ID),
			zap.Uint("articleVersionID", latestVersion.ID))
		return nil, protocol.ErrDataExists
	}

	nextVersion := uint(1)
	if latestVersion != nil {
		nextVersion = latestVersion.Version + 1
	}

	version := &model.ArticleVersion{
		ArticleID: article.ID,
		Version:   nextVersion,
		Content:   req.Body.Content,
	}

	if err := s.articleVersionDAO.Create(db, version); err != nil {
		logger.Error("[ArticleVersionService] failed to create version",
			zap.Uint("articleID", article.ID),
			zap.Uint("version", version.Version),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.ArticleVersion = &dto.ArticleVersion{
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
func (s *articleVersionService) GetArticleVersionInfo(ctx context.Context, req *dto.ArticleVersionGetRequest) (rsp *dto.ArticleVersionGetResponse, err error) {
	logger := logger.WithCtx(ctx)

	rsp = &dto.ArticleVersionGetResponse{}

	db := database.GetDBInstance(ctx)

	// GetArticleVersionInfo 不需要权限校验，任何人都可以查看已发布的文章版本

	article, err := s.articleDAO.GetByID(db, req.ArticleID, []string{"id", "user_id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[ArticleVersionService] article not found",
				zap.Uint("articleID", req.ArticleID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[ArticleVersionService] failed to get article",
			zap.Uint("articleID", req.ArticleID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	version, err := s.articleVersionDAO.GetByArticleIDAndVersion(db, article.ID, req.Version, []string{"id", "article_id", "version", "content", "created_at", "updated_at"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[ArticleVersionService] version not found",
				zap.Uint("articleID", article.ID),
				zap.Uint("versionID", req.Version))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[ArticleVersionService] failed to get version",
			zap.Uint("articleID", article.ID),
			zap.Uint("versionID", req.Version),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.Version = &dto.ArticleVersion{
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
func (s *articleVersionService) GetLatestArticleVersionInfo(ctx context.Context, req *dto.ArticleVersionGetLatestRequest) (rsp *dto.ArticleVersionGetLatestResponse, err error) {
	logger := logger.WithCtx(ctx)

	rsp = &dto.ArticleVersionGetLatestResponse{}

	db := database.GetDBInstance(ctx)

	article, err := s.articleDAO.GetByID(db, req.ArticleID, []string{"id", "user_id", "status"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[ArticleVersionService] article not found",
				zap.Uint("articleID", req.ArticleID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[ArticleVersionService] failed to get article",
			zap.Uint("articleID", req.ArticleID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	if article.UserID != userID && article.Status != model.ArticleStatusPublish {
		logger.Error("[ArticleVersionService] no permission to get latest article version",
			zap.Uint("articleID", req.ArticleID))
		return nil, protocol.ErrNoPermission
	}

	version, err := s.articleVersionDAO.GetLatestByArticleID(db, article.ID, []string{"id", "article_id", "version", "content", "created_at", "updated_at"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[ArticleVersionService] latest version not found",
				zap.Uint("articleID", article.ID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[ArticleVersionService] failed to get latest version",
			zap.Uint("articleID", article.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.Version = &dto.ArticleVersion{
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
func (s *articleVersionService) ListArticleVersions(ctx context.Context, req *dto.ArticleVersionListRequest) (rsp *dto.ArticleVersionListResponse, err error) {
	logger := logger.WithCtx(ctx)

	rsp = &dto.ArticleVersionListResponse{}

	db := database.GetDBInstance(ctx)

	article, err := s.articleDAO.GetByID(db, req.ArticleID, []string{"id", "user_id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[ArticleVersionService] article not found",
				zap.Uint("articleID", req.ArticleID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[ArticleVersionService] failed to get article",
			zap.Uint("articleID", req.ArticleID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	if article.UserID != userID {
		logger.Error("[ArticleVersionService] no permission to list article versions",
			zap.Uint("articleID", req.ArticleID))
		return nil, protocol.ErrNoPermission
	}

	paginate := req.PaginationQuery.ToPaginateParam()
	param := &dao.PaginateParam{
		PageParam: &dao.PageParam{
			Page:     paginate.PageParam.Page,
			PageSize: paginate.PageParam.PageSize,
		},
		QueryParam: &dao.QueryParam{
			Query:       paginate.QueryParam.Query,
			QueryFields: []string{"version", "content"},
		},
	}

	versions, pageInfo, err := s.articleVersionDAO.PaginateByArticleID(db, article.ID,
		[]string{"id", "article_id", "version", "content", "created_at", "updated_at"}, []string{},
		param)
	if err != nil {
		logger.Error("[ArticleVersionService] failed to paginate versions",
			zap.Uint("articleID", article.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.Versions = lo.Map(*versions, func(version model.ArticleVersion, _ int) *dto.ArticleVersion {
		content := version.Content
		if len([]rune(content)) > constant.ListArticleVersionContentLength {
			content = string([]rune(content)[:constant.ListArticleVersionContentLength]) + "..."
		}
		return &dto.ArticleVersion{
			ArticleID:        version.ArticleID,
			ArticleVersionID: version.ID,
			VersionID:        version.Version,
			Content:          content,
			CreatedAt:        version.CreatedAt.Format(time.DateTime),
			UpdatedAt:        version.UpdatedAt.Format(time.DateTime),
		}
	})

	rsp.PageInfo = &dto.PageInfo{
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
		Total:    pageInfo.Total,
	}

	return rsp, nil
}
