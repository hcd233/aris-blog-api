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

// CommentService 评论服务
type CommentService interface {
	CreateArticleComment(req *protocol.CreateArticleCommentRequest) (rsp *protocol.CreateArticleCommentResponse, err error)
	GetCommentInfo(req *protocol.GetCommentInfoRequest) (rsp *protocol.GetCommentInfoResponse, err error)
	DeleteComment(req *protocol.DeleteCommentRequest) (rsp *protocol.DeleteCommentResponse, err error)
	ListArticleComments(req *protocol.ListArticleCommentsRequest) (rsp *protocol.ListArticleCommentsResponse, err error)
	ListChildrenComments(req *protocol.ListChildrenCommentsRequest) (rsp *protocol.ListChildrenCommentsResponse, err error)
}

type commentService struct {
	db         *gorm.DB
	userDAO    *dao.UserDAO
	articleDAO *dao.ArticleDAO
	commentDAO *dao.CommentDAO
}

// NewCommentService 创建评论服务
func NewCommentService() CommentService {
	return &commentService{
		db:         database.GetDBInstance(),
		userDAO:    dao.GetUserDAO(),
		articleDAO: dao.GetArticleDAO(),
		commentDAO: dao.GetCommentDAO(),
	}
}

// CreateArticleComment 创建文章评论
func (s *commentService) CreateArticleComment(req *protocol.CreateArticleCommentRequest) (rsp *protocol.CreateArticleCommentResponse, err error) {
	rsp = &protocol.CreateArticleCommentResponse{}

	author, err := s.userDAO.GetByName(s.db, req.Author, []string{"id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[CommentService] user not found", zap.String("userName", req.Author))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[CommentService] failed to get user", zap.String("userName", req.Author), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	article, err := s.articleDAO.GetBySlugAndUserID(s.db, req.ArticleSlug, author.ID, []string{"id", "status"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[CommentService] article not found",
				zap.String("articleSlug", req.ArticleSlug),
				zap.Uint("userID", author.ID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[CommentService] failed to get article",
			zap.String("articleSlug", req.ArticleSlug),
			zap.Uint("userID", author.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if article.Status == model.ArticleStatusDraft {
		logger.Logger.Info("[CommentService] article is draft",
			zap.String("articleSlug", req.ArticleSlug),
			zap.Uint("userID", author.ID))
		return nil, protocol.ErrNoPermission
	}

	var parent *model.Comment
	if req.ReplyTo != 0 {
		parent, err = s.commentDAO.GetByID(s.db, req.ReplyTo, []string{"id", "article_id"}, []string{})
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logger.Logger.Error("[CommentService] parent comment not found", zap.Uint("commentID", req.ReplyTo))
				return nil, protocol.ErrDataNotExists
			}
			logger.Logger.Error("[CommentService] failed to get parent comment", zap.Uint("commentID", req.ReplyTo), zap.Error(err))
			return nil, protocol.ErrInternalError
		}

		if parent.ArticleID != article.ID {
			logger.Logger.Info("[CommentService] parent comment not belong to article",
				zap.Uint("commentID", req.ReplyTo),
				zap.Uint("articleID", article.ID))
			return nil, protocol.ErrBadRequest
		}
	}

	comment := &model.Comment{
		UserID:    req.CurUserID,
		ArticleID: article.ID,
		Parent:    parent,
		Content:   req.Content,
	}

	if err := s.commentDAO.Create(s.db, comment); err != nil {
		logger.Logger.Error("[CommentService] failed to create comment",
			zap.Uint("userID", req.CurUserID),
			zap.Uint("articleID", article.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	comment = lo.Must1(s.commentDAO.GetByID(s.db, comment.ID, []string{"id", "created_at", "content", "parent_id", "user_id"}, []string{}))

	rsp.Comment = &protocol.Comment{
		CommentID: comment.ID,
		Content:   comment.Content,
		UserID:    comment.UserID,
		ReplyTo:   comment.ParentID,
		CreatedAt: comment.CreatedAt.Format(time.DateTime),
	}

	return rsp, nil
}

// GetCommentInfo 获取评论信息
func (s *commentService) GetCommentInfo(req *protocol.GetCommentInfoRequest) (rsp *protocol.GetCommentInfoResponse, err error) {
	rsp = &protocol.GetCommentInfoResponse{}

	user, err := s.userDAO.GetByName(s.db, req.UserName, []string{"id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[CommentService] user not found", zap.String("userName", req.UserName))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[CommentService] failed to get user", zap.String("userName", req.UserName), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	article, err := s.articleDAO.GetBySlugAndUserID(s.db, req.ArticleSlug, user.ID, []string{"id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[CommentService] article not found",
				zap.String("articleSlug", req.ArticleSlug),
				zap.Uint("userID", user.ID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[CommentService] failed to get article",
			zap.String("articleSlug", req.ArticleSlug),
			zap.Uint("userID", user.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	comment, err := s.commentDAO.GetByArticleIDAndID(s.db, article.ID, req.CommentID, []string{"id", "created_at", "content", "user_id", "parent_id", "article_id", "likes"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[CommentService] comment not found",
				zap.Uint("commentID", req.CommentID),
				zap.Uint("articleID", article.ID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[CommentService] failed to get comment",
			zap.Uint("commentID", req.CommentID),
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
		Likes:     comment.Likes,
	}

	return rsp, nil
}

// DeleteComment 删除评论
func (s *commentService) DeleteComment(req *protocol.DeleteCommentRequest) (rsp *protocol.DeleteCommentResponse, err error) {
	rsp = &protocol.DeleteCommentResponse{}

	if req.CurUserName != req.UserName {
		logger.Logger.Info("[CommentService] no permission to delete comment",
			zap.String("curUserName", req.CurUserName),
			zap.String("userName", req.UserName))
		return nil, protocol.ErrNoPermission
	}

	user, err := s.userDAO.GetByName(s.db, req.UserName, []string{"id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[CommentService] user not found", zap.String("userName", req.UserName))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[CommentService] failed to get user", zap.String("userName", req.UserName), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	article, err := s.articleDAO.GetBySlugAndUserID(s.db, req.ArticleSlug, user.ID, []string{"id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[CommentService] article not found",
				zap.String("articleSlug", req.ArticleSlug),
				zap.Uint("userID", user.ID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[CommentService] failed to get article",
			zap.String("articleSlug", req.ArticleSlug),
			zap.Uint("userID", user.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	comment, err := s.commentDAO.GetByArticleIDAndID(s.db, article.ID, req.CommentID, []string{"id", "user_id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[CommentService] comment not found",
				zap.Uint("commentID", req.CommentID),
				zap.Uint("articleID", article.ID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[CommentService] failed to get comment",
			zap.Uint("commentID", req.CommentID),
			zap.Uint("articleID", article.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if comment.UserID != user.ID {
		logger.Logger.Info("[CommentService] no permission to delete comment",
			zap.Uint("commentUserID", comment.UserID),
			zap.Uint("userID", user.ID))
		return nil, protocol.ErrNoPermission
	}

	if err := s.commentDAO.DeleteReclusiveByID(s.db, comment.ID, []string{"id"}, []string{}); err != nil {
		logger.Logger.Error("[CommentService] failed to delete comment",
			zap.Uint("commentID", comment.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	return rsp, nil
}

// ListArticleComments 列出文章评论
func (s *commentService) ListArticleComments(req *protocol.ListArticleCommentsRequest) (rsp *protocol.ListArticleCommentsResponse, err error) {
	rsp = &protocol.ListArticleCommentsResponse{}

	user, err := s.userDAO.GetByName(s.db, req.UserName, []string{"id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[CommentService] user not found", zap.String("userName", req.UserName))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[CommentService] failed to get user", zap.String("userName", req.UserName), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	article, err := s.articleDAO.GetBySlugAndUserID(s.db, req.ArticleSlug, user.ID, []string{"id", "status"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[CommentService] article not found",
				zap.String("articleSlug", req.ArticleSlug),
				zap.Uint("userID", user.ID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[CommentService] failed to get article",
			zap.String("articleSlug", req.ArticleSlug),
			zap.Uint("userID", user.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if req.CurUserName != req.UserName && article.Status != model.ArticleStatusPublish {
		logger.Logger.Info("[CommentService] no permission to list article comments",
			zap.String("curUserName", req.CurUserName),
			zap.String("userName", req.UserName))
		return nil, protocol.ErrNoPermission
	}

	comments, pageInfo, err := s.commentDAO.PaginateRootsByArticleID(s.db, article.ID, []string{"id", "content", "created_at", "user_id", "likes"}, []string{}, req.PageParam.Page, req.PageParam.PageSize)
	if err != nil {
		logger.Logger.Error("[CommentService] failed to paginate article comments",
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
func (s *commentService) ListChildrenComments(req *protocol.ListChildrenCommentsRequest) (rsp *protocol.ListChildrenCommentsResponse, err error) {
	rsp = &protocol.ListChildrenCommentsResponse{}

	user, err := s.userDAO.GetByName(s.db, req.UserName, []string{"id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[CommentService] user not found", zap.String("userName", req.UserName))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[CommentService] failed to get user", zap.String("userName", req.UserName), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	article, err := s.articleDAO.GetBySlugAndUserID(s.db, req.ArticleSlug, user.ID, []string{"id", "status"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[CommentService] article not found",
				zap.String("articleSlug", req.ArticleSlug),
				zap.Uint("userID", user.ID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[CommentService] failed to get article",
			zap.String("articleSlug", req.ArticleSlug),
			zap.Uint("userID", user.ID),
			zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	parentComment, err := s.commentDAO.GetByID(s.db, req.CommentID, []string{"id", "article_id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Logger.Error("[CommentService] parent comment not found", zap.Uint("commentID", req.CommentID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Logger.Error("[CommentService] failed to get parent comment", zap.Uint("commentID", req.CommentID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if parentComment.ArticleID != article.ID {
		logger.Logger.Info("[CommentService] parent comment not belong to article",
			zap.Uint("commentID", req.CommentID),
			zap.Uint("articleID", article.ID))
		return nil, protocol.ErrBadRequest
	}

	comments, pageInfo, err := s.commentDAO.PaginateChildren(s.db, parentComment,
		[]string{"id", "content", "created_at", "likes", "user_id", "parent_id"},
		[]string{},
		req.PageParam.Page, req.PageParam.PageSize)
	if err != nil {
		logger.Logger.Error("[CommentService] failed to paginate children comments",
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
