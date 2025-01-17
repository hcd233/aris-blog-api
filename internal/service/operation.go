package service

import (
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
	LikeArticle(req *protocol.LikeArticleRequest) (rsp *protocol.LikeArticleResponse, err error)
	LikeComment(req *protocol.LikeCommentRequest) (rsp *protocol.LikeCommentResponse, err error)
	LikeTag(req *protocol.LikeTagRequest) (rsp *protocol.LikeTagResponse, err error)
	LogArticleView(req *protocol.LogArticleViewRequest) (rsp *protocol.LogArticleViewResponse, err error)
}

type operationService struct {
	db          *gorm.DB
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
		db:          database.GetDBInstance(),
		userDAO:     dao.GetUserDAO(),
		tagDAO:      dao.GetTagDAO(),
		articleDAO:  dao.GetArticleDAO(),
		commentDAO:  dao.GetCommentDAO(),
		userLikeDAO: dao.GetUserLikeDAO(),
		userViewDAO: dao.GetUserViewDAO(),
	}
}

// LikeArticle 点赞文章
func (s *operationService) LikeArticle(req *protocol.LikeArticleRequest) (rsp *protocol.LikeArticleResponse, err error) {
	rsp = &protocol.LikeArticleResponse{}

	user, err := s.userDAO.GetByName(s.db, req.Author, []string{"id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[OperationService] user not found", zap.String("userName", req.Author))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[OperationService] failed to get user", zap.String("userName", req.Author), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	article, err := s.articleDAO.GetBySlugAndUserID(s.db, req.ArticleSlug, user.ID, []string{"id", "likes", "status", "user_id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[OperationService] article not found",
				zap.String("articleSlug", req.ArticleSlug),
				zap.Uint("userID", user.ID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[OperationService] failed to get article",
			zap.String("articleSlug", req.ArticleSlug),
			zap.Uint("userID", user.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if req.CurUserID != user.ID && article.Status != model.ArticleStatusPublish {
		logger.Logger.Info("[OperationService] no permission to like article",
			zap.Uint("curUserID", req.CurUserID),
			zap.String("author", req.Author))
		return nil, protocol.ErrNoPermission
	}

	userLike := &model.UserLike{
		UserID:     req.CurUserID,
		ObjectID:   article.ID,
		ObjectType: model.LikeObjectTypeArticle,
	}

	tx := s.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if req.Undo {
		if err = s.transactUndoLikeArticle(tx, article, userLike); err != nil {
			logger.Logger.Error("[OperationService] failed to undo like article",
				zap.Uint("userID", req.CurUserID),
				zap.Uint("articleID", article.ID),
				zap.Error(err))
			return nil, protocol.ErrInternalError
		}
	} else {
		if err = s.transactLikeArticle(tx, article, userLike); err != nil {
			logger.Logger.Error("[OperationService] failed to like article",
				zap.Uint("userID", req.CurUserID),
				zap.Uint("articleID", article.ID),
				zap.Error(err))
			return nil, protocol.ErrInternalError
		}
	}

	return rsp, nil
}

// LikeComment 点赞评论
func (s *operationService) LikeComment(req *protocol.LikeCommentRequest) (rsp *protocol.LikeCommentResponse, err error) {
	rsp = &protocol.LikeCommentResponse{}

	comment, err := s.commentDAO.GetByID(s.db, req.CommentID, []string{"id", "likes", "article_id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[OperationService] comment not found", zap.Uint("commentID", req.CommentID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[OperationService] failed to get comment", zap.Uint("commentID", req.CommentID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	article, err := s.articleDAO.GetByID(s.db, comment.ArticleID, []string{"id", "status", "user_id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[OperationService] article not found", zap.Uint("articleID", comment.ArticleID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[OperationService] failed to get article", zap.Uint("articleID", comment.ArticleID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if req.CurUserID != article.UserID && article.Status != model.ArticleStatusPublish {
		logger.Logger.Info("[OperationService] no permission to like comment",
			zap.Uint("curUserID", req.CurUserID),
			zap.Uint("articleID", article.ID),
			zap.String("articleStatus", string(article.Status)))
		return nil, protocol.ErrNoPermission
	}

	userLike := &model.UserLike{
		UserID:     req.CurUserID,
		ObjectID:   comment.ID,
		ObjectType: model.LikeObjectTypeComment,
	}

	tx := s.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if req.Undo {
		if err = s.transactUndoLikeComment(tx, comment, userLike); err != nil {
			logger.Logger.Error("[OperationService] failed to undo like comment",
				zap.Uint("userID", req.CurUserID),
				zap.Uint("commentID", comment.ID),
				zap.Error(err))
			return nil, protocol.ErrInternalError
		}
	} else {
		if err = s.transactLikeComment(tx, comment, userLike); err != nil {
			logger.Logger.Error("[OperationService] failed to like comment",
				zap.Uint("userID", req.CurUserID),
				zap.Uint("commentID", comment.ID),
				zap.Error(err))
			return nil, protocol.ErrInternalError
		}
	}

	return rsp, nil
}

// LikeTag 点赞标签
func (s *operationService) LikeTag(req *protocol.LikeTagRequest) (rsp *protocol.LikeTagResponse, err error) {
	rsp = &protocol.LikeTagResponse{}

	tag, err := s.tagDAO.GetBySlug(s.db, req.TagSlug, []string{"id", "likes"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[OperationService] tag not found", zap.String("tagSlug", req.TagSlug))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[OperationService] failed to get tag", zap.String("tagSlug", req.TagSlug), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	userLike := &model.UserLike{
		UserID:     req.CurUserID,
		ObjectID:   tag.ID,
		ObjectType: model.LikeObjectTypeTag,
	}

	tx := s.db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if req.Undo {
		if err = s.transactUndoLikeTag(tx, tag, userLike); err != nil {
			logger.Logger.Error("[OperationService] failed to undo like tag",
				zap.Uint("userID", req.CurUserID),
				zap.Uint("tagID", tag.ID),
				zap.Error(err))
			return nil, protocol.ErrInternalError
		}
	} else {
		if err = s.transactLikeTag(tx, tag, userLike); err != nil {
			logger.Logger.Error("[OperationService] failed to like tag",
				zap.Uint("userID", req.CurUserID),
				zap.Uint("tagID", tag.ID),
				zap.Error(err))
			return nil, protocol.ErrInternalError
		}
	}

	return rsp, nil
}

// LogArticleView 记录文章浏览
func (s *operationService) LogArticleView(req *protocol.LogArticleViewRequest) (rsp *protocol.LogArticleViewResponse, err error) {
	rsp = &protocol.LogArticleViewResponse{}

	user, err := s.userDAO.GetByName(s.db, req.Author, []string{"id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[OperationService] user not found", zap.String("userName", req.Author))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[OperationService] failed to get user", zap.String("userName", req.Author), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	article, err := s.articleDAO.GetBySlugAndUserID(s.db, req.ArticleSlug, user.ID, []string{"id", "status"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[OperationService] article not found",
				zap.String("articleSlug", req.ArticleSlug),
				zap.Uint("userID", user.ID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[OperationService] failed to get article",
			zap.String("articleSlug", req.ArticleSlug),
			zap.Uint("userID", user.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if req.CurUserID != user.ID && article.Status != model.ArticleStatusPublish {
		logger.Logger.Info("[OperationService] no permission to view article",
			zap.Uint("curUserID", req.CurUserID),
			zap.String("author", req.Author))
		return nil, protocol.ErrNoPermission
	}

	userView, err := s.userViewDAO.GetLatestViewByUserIDAndArticleID(s.db, req.CurUserID, article.ID, []string{"id", "created_at", "progress"}, []string{})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Logger.Error("[OperationService] failed to get user view",
			zap.Uint("userID", req.CurUserID),
			zap.Uint("articleID", article.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if userView.ID == 0 || userView.Progress >= 95 {
		if req.Progress != 0 {
			logger.Logger.Error("[OperationService] invalid progress",
				zap.Uint("userID", req.CurUserID),
				zap.Uint("articleID", article.ID),
				zap.Int8("progress", req.Progress))
			return nil, protocol.ErrInternalError
		}

		userView = &model.UserView{
			UserID:    req.CurUserID,
			ArticleID: article.ID,
			Progress:  req.Progress,
		}

		if err = s.userViewDAO.Create(s.db, userView); err != nil {
			logger.Logger.Error("[OperationService] failed to create user view",
				zap.Uint("userID", req.CurUserID),
				zap.Uint("articleID", article.ID),
				zap.Error(err))
			return nil, protocol.ErrInternalError
		}
	} else {
		if req.Progress-userView.Progress < 5 {
			logger.Logger.Info("[OperationService] log view too frequently",
				zap.Uint("userID", req.CurUserID),
				zap.Uint("articleID", article.ID),
				zap.Int8("lastProgress", userView.Progress),
				zap.Int8("progress", req.Progress))
			return nil, protocol.ErrTooManyRequests
		}

		if err = s.userViewDAO.Update(s.db, userView, map[string]interface{}{
			"progress":       req.Progress,
			"last_viewed_at": time.Now(),
		}); err != nil {
			logger.Logger.Error("[OperationService] failed to update user view",
				zap.Uint("userID", req.CurUserID),
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
