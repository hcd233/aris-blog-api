package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"io"
	"net/url"
	"path/filepath"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"github.com/hcd233/aris-blog-api/internal/constant"
	"github.com/hcd233/aris-blog-api/internal/logger"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/protocol/dto"
	"github.com/hcd233/aris-blog-api/internal/resource/database"
	"github.com/hcd233/aris-blog-api/internal/resource/database/dao"
	"github.com/hcd233/aris-blog-api/internal/resource/database/model"
	objdao "github.com/hcd233/aris-blog-api/internal/resource/storage/obj_dao"
	"github.com/hcd233/aris-blog-api/internal/util"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"golang.org/x/image/webp"
	"gorm.io/gorm"
)

// AssetService 资产服务
//
//	author centonhuang
//	update 2025-01-05 17:57:43
type AssetService interface {
	ListUserLikeArticles(ctx context.Context, req *dto.ListUserLikeArticlesRequest) (rsp *dto.ListUserLikeArticlesResponse, err error)
	ListUserLikeComments(ctx context.Context, req *dto.ListUserLikeCommentsRequest) (rsp *dto.ListUserLikeCommentsResponse, err error)
	ListUserLikeTags(ctx context.Context, req *dto.ListUserLikeTagsRequest) (rsp *dto.ListUserLikeTagsResponse, err error)
	ListImages(ctx context.Context, req *dto.EmptyRequest) (rsp *dto.ListImagesResponse, err error)
	UploadImage(ctx context.Context, req *dto.UploadImageRequest) (rsp *dto.EmptyResponse, err error)
	GetImage(ctx context.Context, req *dto.GetImageRequest) (rsp *dto.URLResponse, err error)
	DeleteImage(ctx context.Context, req *dto.DeleteImageRequest) (rsp *dto.EmptyResponse, err error)
	ListUserViewArticles(ctx context.Context, req *dto.ListUserViewArticlesRequest) (rsp *dto.ListUserViewArticlesResponse, err error)
	DeleteUserView(ctx context.Context, req *dto.DeleteUserViewRequest) (rsp *dto.EmptyResponse, err error)
}

type assetService struct {
	userDAO         *dao.UserDAO
	tagDAO          *dao.TagDAO
	articleDAO      *dao.ArticleDAO
	commentDAO      *dao.CommentDAO
	userLikeDAO     *dao.UserLikeDAO
	userViewDAO     *dao.UserViewDAO
	imageObjDAO     objdao.ObjDAO
	thumbnailObjDAO objdao.ObjDAO
}

// NewAssetService 创建资产服务
//
//	return AssetService
//	author centonhuang
//	update 2025-01-05 16:41:39
func NewAssetService() AssetService {
	return &assetService{
		userDAO:         dao.GetUserDAO(),
		tagDAO:          dao.GetTagDAO(),
		articleDAO:      dao.GetArticleDAO(),
		commentDAO:      dao.GetCommentDAO(),
		userLikeDAO:     dao.GetUserLikeDAO(),
		userViewDAO:     dao.GetUserViewDAO(),
		imageObjDAO:     objdao.GetImageObjDAO(),
		thumbnailObjDAO: objdao.GetThumbnailObjDAO(),
	}
}

// ListUserLikeArticles 列出用户喜欢的文章
//
//	receiver s *assetService
//	param req *dto.ListUserLikeArticlesRequest
//	return rsp *dto.ListUserLikeArticlesResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 16:47:40
func (s *assetService) ListUserLikeArticles(ctx context.Context, req *dto.ListUserLikeArticlesRequest) (rsp *dto.ListUserLikeArticlesResponse, err error) {
	rsp = &dto.ListUserLikeArticlesResponse{}

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	param := &dao.CommonParam{
		PageParam: &dao.PageParam{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		QueryParam: &dao.QueryParam{
			Query:       req.Query,
			QueryFields: []string{"object_id"},
		},
	}
	userLikes, pageInfo, err := s.userLikeDAO.PaginateByUserIDAndObjectType(db, userID, model.LikeObjectTypeArticle, []string{"object_id"}, []string{}, param)
	if err != nil {
		logger.Error("[AssetService] failed to get user likes", zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	articleIDs := lo.Map(*userLikes, func(like model.UserLike, _ int) uint {
		return like.ObjectID
	})

	articles, err := s.articleDAO.BatchGetByIDs(db, articleIDs,
		[]string{
			"id", "slug", "title", "status", "user_id",
			"created_at", "updated_at", "published_at",
			"likes", "views",
		},
		[]string{"Comments", "Tags"},
	)
	if err != nil {
		logger.Error("[AssetService] failed to get articles", zap.Uints("articleIDs", articleIDs), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if len(articleIDs) != len(*articles) {
		_, deletedIDs := lo.Difference(articleIDs, lo.Map(*articles, func(article model.Article, _ int) uint {
			return article.ID
		}))

		logger.Warn("[AssetService] found deleted articles", zap.Uints("deletedIDs", deletedIDs))

		deleteLikes := lo.Map(deletedIDs, func(id uint, _ int) model.UserLike {
			return model.UserLike{
				ObjectID:   id,
				ObjectType: model.LikeObjectTypeArticle,
			}
		})

		err = s.userLikeDAO.BatchDelete(db, &deleteLikes)
		if err != nil {
			logger.Error("[AssetService] failed to delete user likes", zap.Uints("deletedIDs", deletedIDs), zap.Error(err))
			err = nil
		}
	}

	rsp.Articles = lo.Map(*articles, func(article model.Article, _ int) *dto.Article {
		return &dto.Article{
			ArticleID: article.ID,
			Title:     article.Title,
			Slug:      article.Slug,
			Status:    string(article.Status),
			User: &dto.User{
				UserID: article.User.ID,
				Name:   article.User.Name,
				Avatar: article.User.Avatar,
			},
			CreatedAt:   article.CreatedAt.Format(time.DateTime),
			UpdatedAt:   article.UpdatedAt.Format(time.DateTime),
			PublishedAt: article.PublishedAt.Format(time.DateTime),
			Likes:       article.Likes,
			Views:       article.Views,
			Tags: lo.Map(article.Tags, func(tag model.Tag, _ int) *dto.Tag {
				return &dto.Tag{
					TagID: tag.ID,
					Name:  tag.Name,
					Slug:  tag.Slug,
				}
			}),
			Comments: len(article.Comments),
		}
	})
	rsp.PageInfo = &dto.PageInfo{
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
		Total:    pageInfo.Total,
	}

	return rsp, nil
}

func (s *assetService) ListUserLikeComments(ctx context.Context, req *dto.ListUserLikeCommentsRequest) (rsp *dto.ListUserLikeCommentsResponse, err error) {
	rsp = &dto.ListUserLikeCommentsResponse{}

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	param := &dao.CommonParam{
		PageParam: &dao.PageParam{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		QueryParam: &dao.QueryParam{
			Query:       req.Query,
			QueryFields: []string{"object_id"},
		},
	}
	userLikes, pageInfo, err := s.userLikeDAO.PaginateByUserIDAndObjectType(db, userID, model.LikeObjectTypeComment, []string{"object_id"}, []string{}, param)
	if err != nil {
		logger.Error("[AssetService] failed to get user likes", zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	commentIDs := lo.Map(*userLikes, func(like model.UserLike, _ int) uint {
		return like.ObjectID
	})

	comments, err := s.commentDAO.BatchGetByIDs(db, commentIDs,
		[]string{"id", "user_id", "parent_id", "created_at", "content", "likes"},
		[]string{},
	)
	if err != nil {
		logger.Error("[AssetService] failed to get comments", zap.Uints("commentIDs", commentIDs), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if len(commentIDs) != len(*comments) {
		_, deletedIDs := lo.Difference(commentIDs, lo.Map(*comments, func(comment model.Comment, _ int) uint {
			return comment.ID
		}))

		logger.Warn("[AssetService] found deleted comments", zap.Uints("deletedIDs", deletedIDs))

		deleteLikes := lo.Map(deletedIDs, func(id uint, _ int) model.UserLike {
			return model.UserLike{
				ObjectID:   id,
				ObjectType: model.LikeObjectTypeComment,
			}
		})

		err = s.userLikeDAO.BatchDelete(db, &deleteLikes)
		if err != nil {
			logger.Error("[AssetService] failed to delete user likes", zap.Uints("deletedIDs", deletedIDs), zap.Error(err))
			err = nil
		}
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

func (s *assetService) ListUserLikeTags(ctx context.Context, req *dto.ListUserLikeTagsRequest) (rsp *dto.ListUserLikeTagsResponse, err error) {
	rsp = &dto.ListUserLikeTagsResponse{}

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	param := &dao.CommonParam{
		PageParam: &dao.PageParam{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		QueryParam: &dao.QueryParam{
			Query:       req.Query,
			QueryFields: []string{"object_id"},
		},
	}
	userLikes, pageInfo, err := s.userLikeDAO.PaginateByUserIDAndObjectType(db, userID, model.LikeObjectTypeTag, []string{"object_id"}, []string{}, param)
	if err != nil {
		logger.Error("[AssetService] failed to get user likes", zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	tagIDs := lo.Map(*userLikes, func(like model.UserLike, _ int) uint {
		return like.ObjectID
	})

	tags, err := s.tagDAO.BatchGetByIDs(db, tagIDs,
		[]string{"id", "slug", "name", "description", "user_id", "created_at", "updated_at", "likes"},
		[]string{})
	if err != nil {
		logger.Error("[AssetService] failed to get tags", zap.Uints("tagIDs", tagIDs), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if len(tagIDs) != len(*tags) {
		_, deletedIDs := lo.Difference(tagIDs, lo.Map(*tags, func(tag model.Tag, _ int) uint {
			return tag.ID
		}))

		logger.Warn("[AssetService] found deleted tags", zap.Uints("deletedIDs", deletedIDs))

		deleteLikes := lo.Map(deletedIDs, func(id uint, _ int) model.UserLike {
			return model.UserLike{
				ObjectID:   id,
				ObjectType: model.LikeObjectTypeTag,
			}
		})

		err = s.userLikeDAO.BatchDelete(db, &deleteLikes)
		if err != nil {
			logger.Error("[AssetService] failed to delete user likes", zap.Uints("deletedIDs", deletedIDs), zap.Error(err))
			err = nil
		}
	}

	rsp.Tags = lo.Map(*tags, func(tag model.Tag, _ int) *dto.Tag {
		return &dto.Tag{
			TagID:       tag.ID,
			Name:        tag.Name,
			Slug:        tag.Slug,
			Description: tag.Description,
			CreatedAt:   tag.CreatedAt.Format(time.DateTime),
			UpdatedAt:   tag.UpdatedAt.Format(time.DateTime),
			Likes:       tag.Likes,
		}
	})
	rsp.PageInfo = &dto.PageInfo{
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
		Total:    pageInfo.Total,
	}
	return rsp, nil
}

// ListImages 列出图片
//
//	receiver s *assetService
//	param req *protocol.EmptyRequest
//	return rsp *protocol.ListImagesResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 17:15:38
func (s *assetService) ListImages(ctx context.Context, _ *dto.EmptyRequest) (rsp *dto.ListImagesResponse, err error) {
	rsp = &dto.ListImagesResponse{}

	logger := logger.WithCtx(ctx)

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	objectInfos, err := s.imageObjDAO.ListObjects(ctx, userID)
	if err != nil {
		logger.Error("[AssetService] failed to list images", zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.Images = lo.Map(objectInfos, func(objectInfo objdao.ObjectInfo, _ int) *dto.Image {
		return &dto.Image{
			Name:      objectInfo.ObjectName,
			Size:      objectInfo.Size,
			CreatedAt: objectInfo.LastModified.Format(time.DateTime),
		}
	})

	return rsp, nil
}

// UploadImage 上传图片
//
//	receiver s *assetService
//	param req *protocol.UploadImageRequest
//	return rsp *protocol.EmptyResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 17:20:31
func (s *assetService) UploadImage(ctx context.Context, req *dto.UploadImageRequest) (rsp *dto.EmptyResponse, err error) {
	rsp = &dto.EmptyResponse{}

	logger := logger.WithCtx(ctx)

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	fileName := req.RawBody.Filename
	fileSize := req.RawBody.Size

	if !util.IsValidImageFormat(fileName) {
		logger.Error("[AssetService] invalid image format", zap.String("fileName", fileName))
		return nil, protocol.ErrBadRequest
	}

	if contentType := req.RawBody.Header.Get("Content-Type"); !util.IsValidImageContentType(contentType) {
		logger.Error("[AssetService] invalid image content type", zap.String("contentType", contentType))
		return nil, protocol.ErrBadRequest
	}

	if maxFileSize := 3 * 1024 * 1024; fileSize > int64(maxFileSize) {
		logger.Error("[AssetService] file size is too large", zap.Int64("fileSize", fileSize), zap.Int("maxFileSize", maxFileSize))
		return nil, protocol.ErrBadRequest
	}

	var rawImage image.Image
	var imageFormat imaging.Format

	extension := filepath.Ext(fileName)
	file, err := req.RawBody.Open()
	if err != nil {
		logger.Error("[AssetService] failed to open file", zap.Error(err))
		return nil, protocol.ErrInternalError
	}
	defer file.Close()

	switch extension {
	case ".webp":
		rawImage = lo.Must1(webp.Decode(file))
		imageFormat = imaging.PNG
	case ".png", ".jpg", ".jpeg", ".gif":
		rawImage, _ = lo.Must2(image.Decode(file))
		imageFormat = lo.Must1(imaging.FormatFromExtension(extension))
	default:
		panic(fmt.Sprintf("Invalid image extension: %s", extension))
	}

	// restrict image into 512*512 max size
	x, y := rawImage.Bounds().Dx(), rawImage.Bounds().Dy()

	maxPixel := 512

	for x > maxPixel || y > maxPixel {
		x, y = x/2, y/2
	}

	thumbnailImage := imaging.Thumbnail(rawImage, x, y, imaging.Lanczos)

	var thumbnailBuffer bytes.Buffer
	err = imaging.Encode(&thumbnailBuffer, thumbnailImage, imageFormat)
	if err != nil {
		logger.Error("[AssetService] failed to encode thumbnail image", zap.Error(err))
		return nil, protocol.ErrInternalError
	}
	file.Seek(0, io.SeekStart)

	var wg sync.WaitGroup
	var imageErr, thumbnailErr error

	wg.Add(2)

	go func() {
		defer wg.Done()
		imageErr = s.imageObjDAO.UploadObject(ctx, userID, fileName, fileSize, file)
	}()

	go func() {
		defer wg.Done()
		thumbnailErr = s.thumbnailObjDAO.UploadObject(ctx, userID, fileName, int64(thumbnailBuffer.Len()), &thumbnailBuffer)
	}()

	wg.Wait()

	if imageErr != nil {
		logger.Error("[AssetService] failed to upload image", zap.String("fileName", fileName), zap.Error(imageErr))
		return nil, protocol.ErrInternalError
	}

	if thumbnailErr != nil {
		logger.Error("[AssetService] failed to upload thumbnail image", zap.String("fileName", fileName), zap.Error(thumbnailErr))
		return nil, protocol.ErrInternalError
	}

	logger.Info("[AssetService] image uploaded successfully",
		zap.String("fileName", fileName),
	)
	return rsp, nil
}

// GetImage 获取图片
//
//	receiver s *assetService
//	param req *protocol.GetImageRequest
//	return rsp *protocol.GetImageResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 17:49:39
func (s *assetService) GetImage(ctx context.Context, req *dto.GetImageRequest) (rsp *dto.URLResponse, err error) {
	logger := logger.WithCtx(ctx)

	rsp = &dto.URLResponse{}

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	var presignedURL *url.URL
	switch req.Quality {
	case "low":
		presignedURL, err = s.thumbnailObjDAO.PresignObject(ctx, userID, req.ObjectName)
	case "high", "medium":
		presignedURL, err = s.imageObjDAO.PresignObject(ctx, userID, req.ObjectName)
	default:
		presignedURL, err = s.imageObjDAO.PresignObject(ctx, userID, req.ObjectName)
	}
	if err != nil {
		logger.Error("[AssetService] failed to presign object",
			zap.String("imageName", req.ObjectName),
			zap.String("quality", req.Quality), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.URL = presignedURL.String()

	logger.Info("[AssetService] presigned URL", zap.String("imageName", req.ObjectName), zap.String("quality", req.Quality), zap.String("url", rsp.URL))

	return rsp, nil
}

// DeleteImage 删除图片
//
//	receiver s *assetService
//	param req *protocol.DeleteImageRequest
//	return rsp *protocol.DeleteImageResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 17:49:33
func (s *assetService) DeleteImage(ctx context.Context, req *dto.DeleteImageRequest) (rsp *dto.EmptyResponse, err error) {
	rsp = &dto.EmptyResponse{}

	logger := logger.WithCtx(ctx)

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	var wg sync.WaitGroup
	var imageErr, thumbnailErr error

	wg.Add(2)

	go func() {
		defer wg.Done()
		imageErr = s.imageObjDAO.DeleteObject(ctx, userID, req.ObjectName)
	}()

	go func() {
		defer wg.Done()
		thumbnailErr = s.thumbnailObjDAO.DeleteObject(ctx, userID, req.ObjectName)
	}()

	wg.Wait()

	if imageErr != nil {
		logger.Error("[AssetService] failed to delete image", zap.String("imageName", req.ObjectName), zap.Error(imageErr))
		return nil, protocol.ErrInternalError
	}

	if thumbnailErr != nil {
		logger.Error("[AssetService] failed to delete thumbnail image", zap.String("imageName", req.ObjectName), zap.Error(thumbnailErr))
		return nil, protocol.ErrInternalError
	}

	logger.Info("[AssetService] image deleted successfully", zap.String("imageName", req.ObjectName))
	return rsp, nil
}

func (s *assetService) ListUserViewArticles(ctx context.Context, req *dto.ListUserViewArticlesRequest) (rsp *dto.ListUserViewArticlesResponse, err error) {
	rsp = &dto.ListUserViewArticlesResponse{}

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	param := &dao.CommonParam{
		PageParam: &dao.PageParam{
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		QueryParam: &dao.QueryParam{
			Query:       req.Query,
			QueryFields: []string{"article_id"},
		},
	}
	userViews, pageInfo, err := s.userViewDAO.PaginateByUserID(db, userID, []string{"id", "progress", "last_viewed_at", "user_id", "article_id"}, []string{"User", "Article", "Article.Tags", "Article.User"}, param)
	if err != nil {
		logger.Error("[AssetService] failed to list user view articles", zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.UserViews = lo.Map(*userViews, func(userView model.UserView, _ int) *dto.UserView {
		return &dto.UserView{
			ViewID:       userView.ID,
			Progress:     userView.Progress,
			LastViewedAt: userView.LastViewedAt.Format(time.DateTime),
			UserID:       userView.UserID,
			ArticleID:    userView.ArticleID,
		}
	})
	rsp.PageInfo = &dto.PageInfo{
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
		Total:    pageInfo.Total,
	}
	return rsp, nil
}

func (s *assetService) DeleteUserView(ctx context.Context, req *dto.DeleteUserViewRequest) (rsp *dto.EmptyResponse, err error) {
	rsp = &dto.EmptyResponse{}

	logger := logger.WithCtx(ctx)
	db := database.GetDBInstance(ctx)

	userID := ctx.Value(constant.CtxKeyUserID).(uint)

	userView, err := s.userViewDAO.GetByID(db, req.ViewID, []string{"id", "user_id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[AssetService] user view not found", zap.Uint("viewID", req.ViewID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[AssetService] failed to get user view", zap.Uint("viewID", req.ViewID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if userView.UserID != userID {
		logger.Error("[AssetService] no permission to delete view", zap.Uint("curUserID", userID), zap.Uint("viewID", req.ViewID))
		return nil, protocol.ErrNoPermission
	}

	err = s.userViewDAO.Delete(db, userView)
	if err != nil {
		logger.Error("[AssetService] failed to delete user view", zap.Uint("viewID", req.ViewID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	logger.Info("[AssetService] user view deleted successfully", zap.Uint("viewID", req.ViewID))
	return rsp, nil
}
