package service

import (
	"context"
	"errors"
	"time"

	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/logger"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/protocol/dto"
	"github.com/hcd233/aris-blog-api/internal/resource/database"
	"github.com/hcd233/aris-blog-api/internal/resource/database/dao"
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// OperationService 用户操作服务
type OperationService interface {
	LikeArticle(ctx context.Context, req *dto.LikeArticleRequest) (rsp *dto.EmptyResponse, err error)
	LikeComment(ctx context.Context, req *dto.LikeCommentRequest) (rsp *dto.EmptyResponse, err error)
	LikeTag(ctx context.Context, req *dto.LikeTagRequest) (rsp *dto.EmptyResponse, err error)
	LogArticleView(ctx context.Context, req *dto.LogArticleViewRequest) (rsp *dto.EmptyResponse, err error)
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

func (s *operationService) LikeArticle(ctx context.Context, req *dto.LikeArticleRequest) (rsp *dto.EmptyResponse, err error) {
	rsp = &dto.EmptyResponse{}

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	if req == nil || req.Body == nil {
		logger.Error("[OperationService] request body is nil")
		return nil, protocol.ErrBadRequest
	}

	userIDValue := ctx.Value(constant.CtxKeyUserID)
	if userIDValue == nil {
		logger.Error("[OperationService] user id missing in context")
		return nil, protocol.ErrUnauthorized
	}
	userID := userIDValue.(uint)

	article, err := s.articleDAO.GetByID(db, req.Body.ArticleID, []string{"id", "likes", "status", "user_id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[OperationService] article not found", zap.Uint("articleID", req.Body.ArticleID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[OperationService] failed to get article", zap.Uint("articleID", req.Body.ArticleID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if userID != article.UserID && article.Status != model.ArticleStatusPublish {
		logger.Info("[OperationService] no permission to like article", zap.Uint("articleID", article.ID), zap.String("articleStatus", string(article.Status)))
		return nil, protocol.ErrNoPermission
	}

	userLike := &model.UserLike{
		UserID:     userID,
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

	if req.Body.Undo {
		if err = s.transactUndoLikeArticle(tx, article, userLike); err != nil {
			logger.Error("[OperationService] failed to undo like article", zap.Uint("articleID", article.ID), zap.Error(err))
			return nil, protocol.ErrInternalError
		}
	} else {
		if err = s.transactLikeArticle(tx, article, userLike); err != nil {
			logger.Error("[OperationService] failed to like article", zap.Uint("articleID", article.ID), zap.Error(err))
			return nil, protocol.ErrInternalError
		}
	}

	return rsp, nil
}

func (s *operationService) LikeComment(ctx context.Context, req *dto.LikeCommentRequest) (rsp *dto.EmptyResponse, err error) {
	rsp = &dto.EmptyResponse{}

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	if req == nil || req.Body == nil {
		logger.Error("[OperationService] request body is nil")
		return nil, protocol.ErrBadRequest
	}

	userIDValue := ctx.Value(constant.CtxKeyUserID)
	if userIDValue == nil {
		logger.Error("[OperationService] user id missing in context")
		return nil, protocol.ErrUnauthorized
	}
	userID := userIDValue.(uint)

	comment, err := s.commentDAO.GetByID(db, req.Body.CommentID, []string{"id", "likes", "article_id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[OperationService] comment not found", zap.Uint("commentID", req.Body.CommentID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[OperationService] failed to get comment", zap.Uint("commentID", req.Body.CommentID), zap.Error(err))
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

	if userID != article.UserID && article.Status != model.ArticleStatusPublish {
		logger.Info("[OperationService] no permission to like comment", zap.Uint("articleID", article.ID), zap.String("articleStatus", string(article.Status)))
		return nil, protocol.ErrNoPermission
	}

	userLike := &model.UserLike{
		UserID:     userID,
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

	if req.Body.Undo {
		if err = s.transactUndoLikeComment(tx, comment, userLike); err != nil {
			logger.Error("[OperationService] failed to undo like comment", zap.Uint("commentID", comment.ID), zap.Error(err))
			return nil, protocol.ErrInternalError
		}
	} else {
		if err = s.transactLikeComment(tx, comment, userLike); err != nil {
			logger.Error("[OperationService] failed to like comment", zap.Uint("commentID", comment.ID), zap.Error(err))
			return nil, protocol.ErrInternalError
		}
	}

	return rsp, nil
}

func (s *operationService) LikeTag(ctx context.Context, req *dto.LikeTagRequest) (rsp *dto.EmptyResponse, err error) {
	rsp = &dto.EmptyResponse{}

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	if req == nil || req.Body == nil {
		logger.Error("[OperationService] request body is nil")
		return nil, protocol.ErrBadRequest
	}

	userIDValue := ctx.Value(constant.CtxKeyUserID)
	if userIDValue == nil {
		logger.Error("[OperationService] user id missing in context")
		return nil, protocol.ErrUnauthorized
	}
	userID := userIDValue.(uint)

	tag, err := s.tagDAO.GetByID(db, req.Body.TagID, []string{"id", "likes"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[OperationService] tag not found", zap.Uint("tagID", req.Body.TagID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[OperationService] failed to get tag", zap.Uint("tagID", req.Body.TagID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	userLike := &model.UserLike{
		UserID:     userID,
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

	if req.Body.Undo {
		if err = s.transactUndoLikeTag(tx, tag, userLike); err != nil {
			logger.Error("[OperationService] failed to undo like tag", zap.Uint("tagID", tag.ID), zap.Error(err))
			return nil, protocol.ErrInternalError
		}
	} else {
		if err = s.transactLikeTag(tx, tag, userLike); err != nil {
			logger.Error("[OperationService] failed to like tag", zap.Uint("tagID", tag.ID), zap.Error(err))
			return nil, protocol.ErrInternalError
		}
	}

	return rsp, nil
}

func (s *operationService) LogArticleView(ctx context.Context, req *dto.LogArticleViewRequest) (rsp *dto.EmptyResponse, err error) {
	rsp = &dto.EmptyResponse{}

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	if req == nil || req.Body == nil {
		logger.Error("[OperationService] request body is nil")
		return nil, protocol.ErrBadRequest
	}

	userIDValue := ctx.Value(constant.CtxKeyUserID)
	if userIDValue == nil {
		logger.Error("[OperationService] user id missing in context")
		return nil, protocol.ErrUnauthorized
	}
	userID := userIDValue.(uint)

	article, err := s.articleDAO.GetByID(db, req.Body.ArticleID, []string{"id", "status", "user_id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[OperationService] article not found", zap.Uint("articleID", req.Body.ArticleID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[OperationService] failed to get article", zap.Uint("articleID", req.Body.ArticleID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if userID != article.UserID && article.Status != model.ArticleStatusPublish {
		logger.Info("[OperationService] no permission to view article", zap.Uint("articleID", article.ID), zap.String("articleStatus", string(article.Status)))
		return nil, protocol.ErrNoPermission
	}

	userView, err := s.userViewDAO.GetLatestViewByUserIDAndArticleID(db, userID, article.ID, []string{"id", "created_at", "progress"}, []string{})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("[OperationService] failed to get user view", zap.Uint("articleID", article.ID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if userView.ID == 0 || userView.Progress >= 95 {
		if req.Body.Progress != 0 {
			logger.Error("[OperationService] invalid progress", zap.Uint("articleID", article.ID), zap.Int8("progress", req.Body.Progress))
			return nil, protocol.ErrInternalError
		}

		userView = &model.UserView{
			UserID:       userID,
			ArticleID:    article.ID,
			Progress:     req.Body.Progress,
			LastViewedAt: time.Now().UTC(),
		}

		if err = s.userViewDAO.Create(db, userView); err != nil {
			logger.Error("[OperationService] failed to create user view", zap.Uint("articleID", article.ID), zap.Error(err))
			return nil, protocol.ErrInternalError
		}
	} else {
		if req.Body.Progress-userView.Progress < 5 {
			logger.Info("[OperationService] log view too frequently", zap.Uint("articleID", article.ID), zap.Int8("lastProgress", userView.Progress), zap.Int8("progress", req.Body.Progress))
			return nil, protocol.ErrTooManyRequests
		}

		if err = s.userViewDAO.Update(db, userView, map[string]interface{}{
			"progress":       req.Body.Progress,
			"last_viewed_at": time.Now().UTC(),
		}); err != nil {
			logger.Error("[OperationService] failed to update user view", zap.Uint("articleID", article.ID), zap.Error(err))
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
