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

// OperationService 用户操作服务
type OperationService interface {
	LikeArticle(ctx context.Context, req *protocol.LikeArticleRequest) (rsp *protocol.LikeArticleResponse, err error)
	LikeComment(ctx context.Context, req *protocol.LikeCommentRequest) (rsp *protocol.LikeCommentResponse, err error)
	LikeTag(ctx context.Context, req *protocol.LikeTagRequest) (rsp *protocol.LikeTagResponse, err error)
	LogArticleView(ctx context.Context, req *protocol.LogArticleViewRequest) (rsp *protocol.LogArticleViewResponse, err error)
}

type operationService struct {
	userDAO     *dao.UserDAO
	tagDAO      *dao.TagDAO
	articleDAO  *dao.ArticleDAO
	commentDAO  *dao.CommentDAO
	userLikeDAO *dao.UserLikeDAO
	userViewDAO *dao.UserViewDAO
}

// NewOperationService 创建用户操作服务
func NewOperationService() OperationService {
	return &operationService{
		userDAO:     dao.GetUserDAO(),
		tagDAO:      dao.GetTagDAO(),
		articleDAO:  dao.GetArticleDAO(),
		commentDAO:  dao.GetCommentDAO(),
		userLikeDAO: dao.GetUserLikeDAO(),
		userViewDAO: dao.GetUserViewDAO(),
	}
}

// LikeArticle 点赞文章
func (s *operationService) LikeArticle(ctx context.Context, req *protocol.LikeArticleRequest) (rsp *protocol.LikeArticleResponse, err error) {
	rsp = &protocol.LikeArticleResponse{}

	logger := logger.LoggerWithContext(ctx)
	db := database.GetDBInstance(ctx)

	article, err := s.articleDAO.GetByID(db, req.ArticleID, []string{"id", "likes", "status", "user_id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[OperationService] article not found",
				zap.Uint("articleID", req.ArticleID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[OperationService] failed to get article",
			zap.Uint("articleID", req.ArticleID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	// 如果用户不是文章作者，且文章状态不是已发布，则不允许点赞
	if req.UserID != article.UserID && article.Status != model.ArticleStatusPublish {
		logger.Info("[OperationService] no permission to like article",

			zap.String("articleStatus", string(article.Status)))
		return nil, protocol.ErrNoPermission
	}

	userLike := &model.UserLike{
		UserID:     req.UserID,
		ObjectID:   article.ID,
		ObjectType: model.LikeObjectTypeArticle,
	}

	tx := db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if req.Undo {
		if err = s.transactUndoLikeArticle(tx, article, userLike); err != nil {
			logger.Error("[OperationService] failed to undo like article",

				zap.Uint("articleID", article.ID),
				zap.Error(err))
			return nil, protocol.ErrInternalError
		}
	} else {
		if err = s.transactLikeArticle(tx, article, userLike); err != nil {
			logger.Error("[OperationService] failed to like article",

				zap.Uint("articleID", article.ID),
				zap.Error(err))
			return nil, protocol.ErrInternalError
		}
	}

	return rsp, nil
}

// LikeComment 点赞评论
func (s *operationService) LikeComment(ctx context.Context, req *protocol.LikeCommentRequest) (rsp *protocol.LikeCommentResponse, err error) {
	rsp = &protocol.LikeCommentResponse{}

	logger := logger.LoggerWithContext(ctx)
	db := database.GetDBInstance(ctx)

	comment, err := s.commentDAO.GetByID(db, req.CommentID, []string{"id", "likes", "article_id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[OperationService] comment not found", zap.Uint("commentID", req.CommentID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[OperationService] failed to get comment", zap.Uint("commentID", req.CommentID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	article, err := s.articleDAO.GetByID(db, comment.ArticleID, []string{"id", "status", "user_id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[OperationService] article not found", zap.Uint("articleID", comment.ArticleID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[OperationService] failed to get article", zap.Uint("articleID", comment.ArticleID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	// 如果用户不是文章作者，且文章状态不是已发布，则不允许点赞
	if req.UserID != article.UserID && article.Status != model.ArticleStatusPublish {
		logger.Info("[OperationService] no permission to like comment",

			zap.Uint("articleID", article.ID),
			zap.String("articleStatus", string(article.Status)))
		return nil, protocol.ErrNoPermission
	}

	userLike := &model.UserLike{
		UserID:     req.UserID,
		ObjectID:   comment.ID,
		ObjectType: model.LikeObjectTypeComment,
	}

	tx := db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if req.Undo {
		if err = s.transactUndoLikeComment(tx, comment, userLike); err != nil {
			logger.Error("[OperationService] failed to undo like comment",

				zap.Uint("commentID", comment.ID),
				zap.Error(err))
			return nil, protocol.ErrInternalError
		}
	} else {
		if err = s.transactLikeComment(tx, comment, userLike); err != nil {
			logger.Error("[OperationService] failed to like comment",

				zap.Uint("commentID", comment.ID),
				zap.Error(err))
			return nil, protocol.ErrInternalError
		}
	}

	return rsp, nil
}

// LikeTag 点赞标签
func (s *operationService) LikeTag(ctx context.Context, req *protocol.LikeTagRequest) (rsp *protocol.LikeTagResponse, err error) {
	rsp = &protocol.LikeTagResponse{}

	logger := logger.LoggerWithContext(ctx)
	db := database.GetDBInstance(ctx)

	tag, err := s.tagDAO.GetByID(db, req.TagID, []string{"id", "likes"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[OperationService] tag not found", zap.Uint("tagID", req.TagID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[OperationService] failed to get tag", zap.Uint("tagID", req.TagID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	userLike := &model.UserLike{
		UserID:     req.UserID,
		ObjectID:   tag.ID,
		ObjectType: model.LikeObjectTypeTag,
	}

	tx := db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if req.Undo {
		if err = s.transactUndoLikeTag(tx, tag, userLike); err != nil {
			logger.Error("[OperationService] failed to undo like tag",

				zap.Uint("tagID", tag.ID),
				zap.Error(err))
			return nil, protocol.ErrInternalError
		}
	} else {
		if err = s.transactLikeTag(tx, tag, userLike); err != nil {
			logger.Error("[OperationService] failed to like tag",

				zap.Uint("tagID", tag.ID),
				zap.Error(err))
			return nil, protocol.ErrInternalError
		}
	}

	return rsp, nil
}

// LogArticleView 记录文章浏览
func (s *operationService) LogArticleView(ctx context.Context, req *protocol.LogArticleViewRequest) (rsp *protocol.LogArticleViewResponse, err error) {
	rsp = &protocol.LogArticleViewResponse{}

	logger := logger.LoggerWithContext(ctx)
	db := database.GetDBInstance(ctx)

	article, err := s.articleDAO.GetByID(db, req.ArticleID, []string{"id", "status"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[OperationService] article not found",
				zap.Uint("articleID", req.ArticleID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[OperationService] failed to get article",
			zap.Uint("articleID", req.ArticleID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	// 如果用户不是文章作者，且文章状态不是已发布，则不允许浏览
	if req.UserID != article.UserID && article.Status != model.ArticleStatusPublish {
		logger.Info("[OperationService] no permission to view article",

			zap.String("articleStatus", string(article.Status)))
		return nil, protocol.ErrNoPermission
	}

	userView, err := s.userViewDAO.GetLatestViewByUserIDAndArticleID(db, req.UserID, article.ID, []string{"id", "created_at", "progress"}, []string{})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("[OperationService] failed to get user view",
			zap.Uint("articleID", article.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if userView.ID == 0 || userView.Progress >= 95 {
		if req.Progress != 0 {
			logger.Error("[OperationService] invalid progress",

				zap.Uint("articleID", article.ID),
				zap.Int8("progress", req.Progress))
			return nil, protocol.ErrInternalError
		}

		userView = &model.UserView{
			UserID:       req.UserID,
			ArticleID:    article.ID,
			Progress:     req.Progress,
			LastViewedAt: time.Now(),
		}

		if err = s.userViewDAO.Create(db, userView); err != nil {
			logger.Error("[OperationService] failed to create user view",

				zap.Uint("articleID", article.ID),
				zap.Error(err))
			return nil, protocol.ErrInternalError
		}
	} else {
		if req.Progress-userView.Progress < 5 {
			logger.Info("[OperationService] log view too frequently",

				zap.Uint("articleID", article.ID),
				zap.Int8("lastProgress", userView.Progress),
				zap.Int8("progress", req.Progress))
			return nil, protocol.ErrTooManyRequests
		}

		if err = s.userViewDAO.Update(db, userView, map[string]interface{}{
			"progress":       req.Progress,
			"last_viewed_at": time.Now(),
		}); err != nil {
			logger.Error("[OperationService] failed to update user view",

				zap.Uint("articleID", article.ID),
				zap.Error(err))
			return nil, protocol.ErrInternalError
		}
	}

	return rsp, nil
}

// 辅助方法
func (s *operationService) transactLikeArticle(tx *gorm.DB, article *model.Article, userLike *model.UserLike) error {
	if err := s.userLikeDAO.Create(tx, userLike); err != nil {
		return err
	}

	if err := s.articleDAO.Update(tx, article, map[string]interface{}{"likes": article.Likes + 1}); err != nil {
		return err
	}

	return nil
}

func (s *operationService) transactUndoLikeArticle(tx *gorm.DB, article *model.Article, userLike *model.UserLike) error {
	userLikeWithID, err := s.userLikeDAO.GetByUserIDAndObject(tx, userLike.UserID, userLike.ObjectID, userLike.ObjectType, []string{"id"}, []string{})
	if err != nil {
		return err
	}

	userLike.ID = userLikeWithID.ID

	if err := s.userLikeDAO.Delete(tx, userLike); err != nil {
		return err
	}

	if err := s.articleDAO.Update(tx, article, map[string]interface{}{"likes": article.Likes - 1}); err != nil {
		return err
	}

	return nil
}

func (s *operationService) transactLikeComment(tx *gorm.DB, comment *model.Comment, userLike *model.UserLike) error {
	if err := s.userLikeDAO.Create(tx, userLike); err != nil {
		return err
	}

	if err := s.commentDAO.Update(tx, comment, map[string]interface{}{"likes": comment.Likes + 1}); err != nil {
		return err
	}

	return nil
}

func (s *operationService) transactUndoLikeComment(tx *gorm.DB, comment *model.Comment, userLike *model.UserLike) error {
	userLikeWithID, err := s.userLikeDAO.GetByUserIDAndObject(tx, userLike.UserID, userLike.ObjectID, userLike.ObjectType, []string{"id"}, []string{})
	if err != nil {
		return err
	}

	userLike.ID = userLikeWithID.ID

	if err := s.userLikeDAO.Delete(tx, userLike); err != nil {
		return err
	}

	if err := s.commentDAO.Update(tx, comment, map[string]interface{}{"likes": comment.Likes - 1}); err != nil {
		return err
	}

	return nil
}

func (s *operationService) transactLikeTag(tx *gorm.DB, tag *model.Tag, userLike *model.UserLike) error {
	if err := s.userLikeDAO.Create(tx, userLike); err != nil {
		return err
	}

	if err := s.tagDAO.Update(tx, tag, map[string]interface{}{"likes": tag.Likes + 1}); err != nil {
		return err
	}

	return nil
}

func (s *operationService) transactUndoLikeTag(tx *gorm.DB, tag *model.Tag, userLike *model.UserLike) error {
	userLikeWithID, err := s.userLikeDAO.GetByUserIDAndObject(tx, userLike.UserID, userLike.ObjectID, userLike.ObjectType, []string{"id"}, []string{})
	if err != nil {
		return err
	}

	userLike.ID = userLikeWithID.ID

	if err := s.userLikeDAO.Delete(tx, userLike); err != nil {
		return err
	}

	if err := s.tagDAO.Update(tx, tag, map[string]interface{}{"likes": tag.Likes - 1}); err != nil {
		return err
	}

	return nil
}
