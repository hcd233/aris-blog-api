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

// CommentService 评论服务
type CommentService interface {
	CreateArticleComment(ctx context.Context, req *protocol.CreateArticleCommentRequest) (rsp *protocol.CreateArticleCommentResponse, err error)
	DeleteComment(ctx context.Context, req *protocol.DeleteCommentRequest) (rsp *protocol.DeleteCommentResponse, err error)
	ListArticleComments(ctx context.Context, req *protocol.ListArticleCommentsRequest) (rsp *protocol.ListArticleCommentsResponse, err error)
	ListChildrenComments(ctx context.Context, req *protocol.ListChildrenCommentsRequest) (rsp *protocol.ListChildrenCommentsResponse, err error)
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
func (s *commentService) CreateArticleComment(ctx context.Context, req *protocol.CreateArticleCommentRequest) (rsp *protocol.CreateArticleCommentResponse, err error) {
	rsp = &protocol.CreateArticleCommentResponse{}

	logger := logger.LoggerWithContext(ctx)
	db := database.GetDBInstance(ctx)

	article, err := s.articleDAO.GetByIDAndStatus(db, req.ArticleID, model.ArticleStatusPublish, []string{"id"}, []string{})
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

	var parent *model.Comment
	if req.ReplyTo != 0 {
		parent, err = s.commentDAO.GetByID(db, req.ReplyTo, []string{"id", "article_id"}, []string{})
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logger.Error("[CommentService] parent comment not found", zap.Uint("commentID", req.ReplyTo))
				return nil, protocol.ErrDataNotExists
			}
			logger.Error("[CommentService] failed to get parent comment", zap.Uint("commentID", req.ReplyTo), zap.Error(err))
			return nil, protocol.ErrInternalError
		}

		if parent.ArticleID != article.ID {
			logger.Info("[CommentService] parent comment not belong to article",
				zap.Uint("commentID", req.ReplyTo),
				zap.Uint("articleID", article.ID))
			return nil, protocol.ErrBadRequest
		}
	}

	comment := &model.Comment{
		UserID:    req.UserID,
		ArticleID: article.ID,
		Parent:    parent,
		Content:   req.Content,
	}

	if err := s.commentDAO.Create(db, comment); err != nil {
		logger.Error("[CommentService] failed to create comment",

			zap.Uint("articleID", article.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.Comment = &protocol.Comment{
		CommentID: comment.ID,
		Content:   comment.Content,
		UserID:    comment.UserID,
		ReplyTo:   comment.ParentID,
		CreatedAt: comment.CreatedAt.Format(time.DateTime),
	}

	return rsp, nil
}

// DeleteComment 删除评论
func (s *commentService) DeleteComment(ctx context.Context, req *protocol.DeleteCommentRequest) (rsp *protocol.DeleteCommentResponse, err error) {
	rsp = &protocol.DeleteCommentResponse{}

	logger := logger.LoggerWithContext(ctx)
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

	// 只有文章作者和评论作者可以删除评论
	if article.UserID != req.UserID && comment.UserID != req.UserID {
		logger.Error("[CommentService] no permission to delete comment",
			zap.Uint("commentUserID", comment.UserID),
		)
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
func (s *commentService) ListArticleComments(ctx context.Context, req *protocol.ListArticleCommentsRequest) (rsp *protocol.ListArticleCommentsResponse, err error) {
	rsp = &protocol.ListArticleCommentsResponse{}

	logger := logger.LoggerWithContext(ctx)
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

	if article.UserID != req.UserID && article.Status != model.ArticleStatusPublish {
		logger.Error("[CommentService] no permission to list article comments",
			zap.Uint("articleUserID", article.UserID),
		)
		return nil, protocol.ErrNoPermission
	}

	comments, pageInfo, err := s.commentDAO.PaginateRootsByArticleID(db, article.ID, []string{"id", "content", "created_at", "user_id", "likes"}, []string{}, req.PageParam.Page, req.PageParam.PageSize)
	if err != nil {
		logger.Error("[CommentService] failed to paginate article comments",
			zap.Uint("articleID", article.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.Comments = lo.Map(*comments, func(comment model.Comment, _ int) *protocol.Comment {
		return &protocol.Comment{
			CommentID: comment.ID,
			Content:   comment.Content,
			UserID:    comment.UserID,
			ReplyTo:   comment.ParentID,
			CreatedAt: comment.CreatedAt.Format(time.DateTime),
			Likes:     comment.Likes,
		}
	})

	rsp.PageInfo = &protocol.PageInfo{
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
		Total:    pageInfo.Total,
	}

	return rsp, nil
}

// ListChildrenComments 列出子评论
func (s *commentService) ListChildrenComments(ctx context.Context, req *protocol.ListChildrenCommentsRequest) (rsp *protocol.ListChildrenCommentsResponse, err error) {
	rsp = &protocol.ListChildrenCommentsResponse{}

	logger := logger.LoggerWithContext(ctx)
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

	if article.UserID != req.UserID && article.Status != model.ArticleStatusPublish {
		logger.Error("[CommentService] no permission to list children comments",
			zap.Uint("articleUserID", article.UserID),
		)
		return nil, protocol.ErrNoPermission
	}

	comments, pageInfo, err := s.commentDAO.PaginateChildren(db, parentComment,
		[]string{"id", "content", "created_at", "likes", "user_id", "parent_id"},
		[]string{},
		req.PageParam.Page, req.PageParam.PageSize)
	if err != nil {
		logger.Error("[CommentService] failed to paginate children comments",
			zap.Uint("parentID", parentComment.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.Comments = lo.Map(*comments, func(comment model.Comment, _ int) *protocol.Comment {
		return &protocol.Comment{
			CommentID: comment.ID,
			Content:   comment.Content,
			UserID:    comment.UserID,
			ReplyTo:   comment.ParentID,
			CreatedAt: comment.CreatedAt.Format(time.DateTime),
			Likes:     comment.Likes,
		}
	})

	rsp.PageInfo = &protocol.PageInfo{
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
		Total:    pageInfo.Total,
	}

	return rsp, nil
}
