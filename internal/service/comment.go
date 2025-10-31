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

// CommentService 评论服务
type CommentService interface {
	CreateArticleComment(ctx context.Context, req *dto.CreateCommentRequest) (rsp *dto.CreateCommentResponse, err error)
	DeleteComment(ctx context.Context, req *dto.DeleteCommentRequest) (rsp *dto.EmptyResponse, err error)
	ListArticleComments(ctx context.Context, req *dto.ListArticleCommentRequest) (rsp *dto.ListArticleCommentResponse, err error)
	ListChildrenComments(ctx context.Context, req *dto.ListChildrenCommentRequest) (rsp *dto.ListChildrenCommentResponse, err error)
}

type commentService struct {
	userDAO    *dao.UserDAO
	articleDAO *dao.ArticleDAO
	commentDAO *dao.CommentDAO
}

// NewCommentService 创建评论服务
func NewCommentService() CommentService {
	return &commentService{
		userDAO:    dao.GetUserDAO(),
		articleDAO: dao.GetArticleDAO(),
		commentDAO: dao.GetCommentDAO(),
	}
}

// CreateArticleComment 创建文章评论
func (s *commentService) CreateArticleComment(ctx context.Context, req *dto.CreateCommentRequest) (rsp *dto.CreateCommentResponse, err error) {
	logger := logger.WithCtx(ctx)

	if req == nil || req.Body == nil {
		logger.Error("[CommentService] request body is nil")
		return nil, protocol.ErrBadRequest
	}

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	rsp = &dto.CreateCommentResponse{}

	db := database.GetDBInstance(ctx)

	article, err := s.articleDAO.GetByIDAndStatus(db, req.Body.ArticleID, model.ArticleStatusPublish, []string{"id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[CommentService] article not found",
				zap.Uint("articleID", req.Body.ArticleID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[CommentService] failed to get article",
			zap.Uint("articleID", req.Body.ArticleID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	var parent *model.Comment
	if req.Body.ReplyTo != 0 {
		parent, err = s.commentDAO.GetByID(db, req.Body.ReplyTo, []string{"id", "article_id"}, []string{})
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logger.Error("[CommentService] parent comment not found", zap.Uint("commentID", req.Body.ReplyTo))
				return nil, protocol.ErrDataNotExists
			}
			logger.Error("[CommentService] failed to get parent comment", zap.Uint("commentID", req.Body.ReplyTo), zap.Error(err))
			return nil, protocol.ErrInternalError
		}

		if parent.ArticleID != article.ID {
			logger.Info("[CommentService] parent comment not belong to article",
				zap.Uint("commentID", req.Body.ReplyTo),
				zap.Uint("articleID", article.ID))
			return nil, protocol.ErrBadRequest
		}
	}

	comment := &model.Comment{
		UserID:    userID,
		ArticleID: article.ID,
		Parent:    parent,
		Content:   req.Body.Content,
	}

	if err := s.commentDAO.Create(db, comment); err != nil {
		logger.Error("[CommentService] failed to create comment",
			zap.Uint("articleID", article.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.Comment = &dto.Comment{
		CommentID: comment.ID,
		Content:   comment.Content,
		UserID:    comment.UserID,
		ReplyTo:   comment.ParentID,
		CreatedAt: comment.CreatedAt.Format(time.DateTime),
		Likes:     comment.Likes,
	}

	return rsp, nil
}

// DeleteComment 删除评论
func (s *commentService) DeleteComment(ctx context.Context, req *dto.DeleteCommentRequest) (rsp *dto.EmptyResponse, err error) {
	logger := logger.WithCtx(ctx)

	rsp = &dto.EmptyResponse{}

	db := database.GetDBInstance(ctx)

	comment, err := s.commentDAO.GetByID(db, req.CommentID, []string{"id", "user_id", "article_id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[CommentService] comment not found",
				zap.Uint("commentID", req.CommentID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[CommentService] failed to get comment",
			zap.Uint("commentID", req.CommentID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	article, err := s.articleDAO.GetByID(db, comment.ArticleID, []string{"id", "user_id"}, []string{})
	if err != nil {
		logger.Error("[CommentService] failed to get article",
			zap.Uint("articleID", comment.ArticleID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	if article.UserID != userID && comment.UserID != userID {
		logger.Error("[CommentService] no permission to delete comment",
			zap.Uint("commentUserID", comment.UserID))
		return nil, protocol.ErrNoPermission
	}

	if err := s.commentDAO.DeleteReclusiveByID(db, comment.ID, []string{"id"}, []string{}); err != nil {
		logger.Error("[CommentService] failed to delete comment",
			zap.Uint("commentID", comment.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	return rsp, nil
}

// ListArticleComments 列出文章评论
func (s *commentService) ListArticleComments(ctx context.Context, req *dto.ListArticleCommentRequest) (rsp *dto.ListArticleCommentResponse, err error) {
	logger := logger.WithCtx(ctx)

	rsp = &dto.ListArticleCommentResponse{}

	db := database.GetDBInstance(ctx)

	article, err := s.articleDAO.GetByID(db, req.ArticleID, []string{"id", "user_id", "status"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[CommentService] article not found",
				zap.Uint("articleID", req.ArticleID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[CommentService] failed to get article",
			zap.Uint("articleID", req.ArticleID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	if article.UserID != userID && article.Status != model.ArticleStatusPublish {
		logger.Error("[CommentService] no permission to list article comments",
			zap.Uint("articleUserID", article.UserID))
		return nil, protocol.ErrNoPermission
	}

	param := &dao.CommonParam{
		PageParam: &dao.PageParam{
			Page:     req.PageParam.Page,
			PageSize: req.PageParam.PageSize,
		},
		QueryParam: &dao.QueryParam{
			Query:       req.QueryParam.Query,
			QueryFields: []string{"content"},
		},
	}
	comments, pageInfo, err := s.commentDAO.PaginateRootsByArticleID(db, article.ID, []string{"id", "content", "created_at", "user_id", "likes"}, []string{}, param)
	if err != nil {
		logger.Error("[CommentService] failed to paginate article comments",
			zap.Uint("articleID", article.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.Comments = lo.Map(*comments, func(comment model.Comment, _ int) *dto.Comment {
		return &dto.Comment{
			CommentID: comment.ID,
			Content:   comment.Content,
			UserID:    comment.UserID,
			ReplyTo:   comment.ParentID,
			CreatedAt: comment.CreatedAt.Format(time.DateTime),
			Likes:     comment.Likes,
		}
	})

	rsp.PageInfo = &dto.PageInfo{
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
		Total:    pageInfo.Total,
	}

	return rsp, nil
}

// ListChildrenComments 列出子评论
func (s *commentService) ListChildrenComments(ctx context.Context, req *dto.ListChildrenCommentRequest) (rsp *dto.ListChildrenCommentResponse, err error) {
	logger := logger.WithCtx(ctx)

	rsp = &dto.ListChildrenCommentResponse{}

	db := database.GetDBInstance(ctx)

	parentComment, err := s.commentDAO.GetByID(db, req.CommentID, []string{"id", "article_id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[CommentService] parent comment not found", zap.Uint("commentID", req.CommentID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[CommentService] failed to get parent comment", zap.Uint("commentID", req.CommentID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	article, err := s.articleDAO.GetByID(db, parentComment.ArticleID, []string{"id", "user_id", "status"}, []string{})
	if err != nil {
		logger.Error("[CommentService] failed to get article", zap.Uint("articleID", parentComment.ArticleID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	if article.UserID != userID && article.Status != model.ArticleStatusPublish {
		logger.Error("[CommentService] no permission to list children comments",
			zap.Uint("articleUserID", article.UserID))
		return nil, protocol.ErrNoPermission
	}

	param := &dao.CommonParam{
		PageParam: &dao.PageParam{
			Page:     req.PageParam.Page,
			PageSize: req.PageParam.PageSize,
		},
		QueryParam: &dao.QueryParam{
			Query:       req.QueryParam.Query,
			QueryFields: []string{"content"},
		},
	}

	comments, pageInfo, err := s.commentDAO.PaginateChildren(db, parentComment,
		[]string{"id", "content", "created_at", "likes", "user_id", "parent_id"},
		[]string{},
		param)
	if err != nil {
		logger.Error("[CommentService] failed to paginate children comments",
			zap.Uint("parentID", parentComment.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.Comments = lo.Map(*comments, func(comment model.Comment, _ int) *dto.Comment {
		return &dto.Comment{
			CommentID: comment.ID,
			Content:   comment.Content,
			UserID:    comment.UserID,
			ReplyTo:   comment.ParentID,
			CreatedAt: comment.CreatedAt.Format(time.DateTime),
			Likes:     comment.Likes,
		}
	})

	rsp.PageInfo = &dto.PageInfo{
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
		Total:    pageInfo.Total,
	}

	return rsp, nil
}
