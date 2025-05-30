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
	"github.com/hcd233/aris-blog-api/internal/logger"
	"github.com/hcd233/aris-blog-api/internal/protocol"
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
	ListUserLikeArticles(ctx context.Context, req *protocol.ListUserLikeArticlesRequest) (rsp *protocol.ListUserLikeArticlesResponse, err error)
	ListUserLikeComments(ctx context.Context, req *protocol.ListUserLikeCommentsRequest) (rsp *protocol.ListUserLikeCommentsResponse, err error)
	ListUserLikeTags(ctx context.Context, req *protocol.ListUserLikeTagsRequest) (rsp *protocol.ListUserLikeTagsResponse, err error)
	ListImages(ctx context.Context, req *protocol.ListImagesRequest) (rsp *protocol.ListImagesResponse, err error)
	UploadImage(ctx context.Context, req *protocol.UploadImageRequest) (rsp *protocol.UploadImageResponse, err error)
	GetImage(ctx context.Context, req *protocol.GetImageRequest) (rsp *protocol.GetImageResponse, err error)
	DeleteImage(ctx context.Context, req *protocol.DeleteImageRequest) (rsp *protocol.DeleteImageResponse, err error)
	ListUserViewArticles(ctx context.Context, req *protocol.ListUserViewArticlesRequest) (rsp *protocol.ListUserViewArticlesResponse, err error)
	DeleteUserView(ctx context.Context, req *protocol.DeleteUserViewRequest) (rsp *protocol.DeleteUserViewResponse, err error)
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
//	param req *protocol.ListUserLikeArticlesRequest
//	return rsp *protocol.ListUserLikeArticlesResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 16:47:40
func (s *assetService) ListUserLikeArticles(ctx context.Context, req *protocol.ListUserLikeArticlesRequest) (rsp *protocol.ListUserLikeArticlesResponse, err error) {
	rsp = &protocol.ListUserLikeArticlesResponse{}

	logger := logger.LoggerWithContext(ctx)
	db := database.GetDBInstance(ctx)

	userLikes, pageInfo, err := s.userLikeDAO.PaginateByUserIDAndObjectType(db, req.UserID, model.LikeObjectTypeArticle, []string{"object_id"}, []string{}, req.PageParam.Page, req.PageParam.PageSize)
	if err != nil {
		logger.Error("[AssetService] failed to get user likes", zap.Uint("userID", req.UserID), zap.Error(err))
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

	rsp.Articles = lo.Map(*articles, func(article model.Article, _ int) *protocol.Article {
		return &protocol.Article{
			ArticleID:   article.ID,
			Title:       article.Title,
			Slug:        article.Slug,
			Status:      string(article.Status),
			UserID:      article.UserID,
			CreatedAt:   article.CreatedAt.Format(time.DateTime),
			UpdatedAt:   article.UpdatedAt.Format(time.DateTime),
			PublishedAt: article.PublishedAt.Format(time.DateTime),
			Likes:       article.Likes,
			Views:       article.Views,
			Tags:        lo.Map(article.Tags, func(tag model.Tag, _ int) string { return tag.Slug }),
			Comments:    len(article.Comments),
		}
	})
	rsp.PageInfo = &protocol.PageInfo{
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
		Total:    pageInfo.Total,
	}

	return rsp, nil
}

func (s *assetService) ListUserLikeComments(ctx context.Context, req *protocol.ListUserLikeCommentsRequest) (rsp *protocol.ListUserLikeCommentsResponse, err error) {
	rsp = &protocol.ListUserLikeCommentsResponse{}

	logger := logger.LoggerWithContext(ctx)
	db := database.GetDBInstance(ctx)

	userLikes, pageInfo, err := s.userLikeDAO.PaginateByUserIDAndObjectType(db, req.UserID, model.LikeObjectTypeComment, []string{"object_id"}, []string{}, req.PageParam.Page, req.PageParam.PageSize)
	if err != nil {
		logger.Error("[AssetService] failed to get user likes", zap.Uint("userID", req.UserID), zap.Error(err))
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

func (s *assetService) ListUserLikeTags(ctx context.Context, req *protocol.ListUserLikeTagsRequest) (rsp *protocol.ListUserLikeTagsResponse, err error) {
	rsp = &protocol.ListUserLikeTagsResponse{}

	logger := logger.LoggerWithContext(ctx)
	db := database.GetDBInstance(ctx)

	userLikes, pageInfo, err := s.userLikeDAO.PaginateByUserIDAndObjectType(db, req.UserID, model.LikeObjectTypeTag, []string{"object_id"}, []string{}, req.PageParam.Page, req.PageParam.PageSize)
	if err != nil {
		logger.Error("[AssetService] failed to get user likes", zap.Uint("userID", req.UserID), zap.Error(err))
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

	rsp.Tags = lo.Map(*tags, func(tag model.Tag, _ int) *protocol.Tag {
		return &protocol.Tag{
			TagID:       tag.ID,
			Name:        tag.Name,
			Slug:        tag.Slug,
			Description: tag.Description,
			UserID:      tag.UserID,
			CreatedAt:   tag.CreatedAt.Format(time.DateTime),
			UpdatedAt:   tag.UpdatedAt.Format(time.DateTime),
			Likes:       tag.Likes,
		}
	})
	rsp.PageInfo = &protocol.PageInfo{
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
		Total:    pageInfo.Total,
	}
	return rsp, nil
}

// ListImages 列出图片
//
//	receiver s *assetService
//	param req *protocol.ListImagesRequest
//	return rsp *protocol.ListImagesResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 17:15:38
func (s *assetService) ListImages(ctx context.Context, req *protocol.ListImagesRequest) (rsp *protocol.ListImagesResponse, err error) {
	rsp = &protocol.ListImagesResponse{}

	logger := logger.LoggerWithContext(ctx)

	objectInfos, err := s.imageObjDAO.ListObjects(req.UserID)
	if err != nil {
		logger.Error("[AssetService] failed to list images", zap.Uint("userID", req.UserID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.Images = lo.Map(objectInfos, func(objectInfo objdao.ObjectInfo, _ int) *protocol.Image {
		return &protocol.Image{
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
//	return rsp *protocol.UploadImageResponse
//	return err error
//	author centonhuang
//	update 2025-01-05 17:20:31
func (s *assetService) UploadImage(ctx context.Context, req *protocol.UploadImageRequest) (rsp *protocol.UploadImageResponse, err error) {
	rsp = &protocol.UploadImageResponse{}

	logger := logger.LoggerWithContext(ctx)

	if !util.IsValidImageFormat(req.FileName) {
		logger.Error("[AssetService] invalid image format", zap.String("fileName", req.FileName))
		return nil, protocol.ErrBadRequest
	}

	if !util.IsValidImageContentType(req.ContentType) {
		logger.Error("[AssetService] invalid image content type", zap.String("contentType", req.ContentType))
		return nil, protocol.ErrBadRequest
	}

	if maxFileSize := 3 * 1024 * 1024; req.Size > int64(maxFileSize) {
		logger.Error("[AssetService] file size is too large", zap.Int64("fileSize", req.Size), zap.Int("maxFileSize", maxFileSize))
		return nil, protocol.ErrBadRequest
	}

	var rawImage image.Image
	var imageFormat imaging.Format

	extension := filepath.Ext(req.FileName)
	switch extension {
	case ".webp":
		rawImage = lo.Must1(webp.Decode(req.ReadSeeker))
		imageFormat = imaging.PNG
	case ".png", ".jpg", ".jpeg", ".gif":
		rawImage, _ = lo.Must2(image.Decode(req.ReadSeeker))
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
	req.ReadSeeker.Seek(0, io.SeekStart)

	var wg sync.WaitGroup
	var imageErr, thumbnailErr error

	wg.Add(2)

	go func() {
		defer wg.Done()
		imageErr = s.imageObjDAO.UploadObject(req.UserID, req.FileName, int64(req.Size), req.ReadSeeker)
	}()

	go func() {
		defer wg.Done()
		thumbnailErr = s.thumbnailObjDAO.UploadObject(req.UserID, req.FileName, int64(thumbnailBuffer.Len()), &thumbnailBuffer)
	}()

	wg.Wait()

	if imageErr != nil {
		logger.Error("[AssetService] failed to upload image", zap.Uint("userID", req.UserID), zap.String("fileName", req.FileName), zap.Error(imageErr))
		return nil, protocol.ErrInternalError
	}

	if thumbnailErr != nil {
		logger.Error("[AssetService] failed to upload thumbnail image", zap.Uint("userID", req.UserID), zap.String("fileName", req.FileName), zap.Error(thumbnailErr))
		return nil, protocol.ErrInternalError
	}

	logger.Info("[AssetService] image uploaded successfully",
		zap.Uint("userID", req.UserID),
		zap.String("fileName", req.FileName),
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
func (s *assetService) GetImage(ctx context.Context, req *protocol.GetImageRequest) (rsp *protocol.GetImageResponse, err error) {
	rsp = &protocol.GetImageResponse{}

	logger := logger.LoggerWithContext(ctx)

	var presignedURL *url.URL
	switch req.Quality {
	case "raw":
		presignedURL, err = s.imageObjDAO.PresignObject(req.UserID, req.ImageName)
	case "thumb":
		presignedURL, err = s.thumbnailObjDAO.PresignObject(req.UserID, req.ImageName)
	}
	if err != nil {
		logger.Error("[AssetService] failed to presign object",
			zap.Uint("userID", req.UserID),
			zap.String("imageName", req.ImageName),
			zap.String("quality", req.Quality), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.PresignedURL = presignedURL.String()
	logger.Info("[AssetService] presigned URL", zap.Uint("userID", req.UserID), zap.String("imageName", req.ImageName), zap.String("quality", req.Quality), zap.String("url", rsp.PresignedURL))

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
func (s *assetService) DeleteImage(ctx context.Context, req *protocol.DeleteImageRequest) (rsp *protocol.DeleteImageResponse, err error) {
	rsp = &protocol.DeleteImageResponse{}

	logger := logger.LoggerWithContext(ctx)

	var wg sync.WaitGroup
	var imageErr, thumbnailErr error

	wg.Add(2)

	go func() {
		defer wg.Done()
		imageErr = s.imageObjDAO.DeleteObject(req.UserID, req.ImageName)
	}()

	go func() {
		defer wg.Done()
		thumbnailErr = s.thumbnailObjDAO.DeleteObject(req.UserID, req.ImageName)
	}()

	wg.Wait()

	if imageErr != nil {
		logger.Error("[AssetService] failed to delete image", zap.Uint("userID", req.UserID), zap.String("imageName", req.ImageName), zap.Error(imageErr))
		return nil, protocol.ErrInternalError
	}

	if thumbnailErr != nil {
		logger.Error("[AssetService] failed to delete thumbnail image", zap.Uint("userID", req.UserID), zap.String("imageName", req.ImageName), zap.Error(thumbnailErr))
		return nil, protocol.ErrInternalError
	}

	logger.Info("[AssetService] image deleted successfully", zap.Uint("userID", req.UserID), zap.String("imageName", req.ImageName))
	return rsp, nil
}

func (s *assetService) ListUserViewArticles(ctx context.Context, req *protocol.ListUserViewArticlesRequest) (rsp *protocol.ListUserViewArticlesResponse, err error) {
	rsp = &protocol.ListUserViewArticlesResponse{}

	logger := logger.LoggerWithContext(ctx)
	db := database.GetDBInstance(ctx)

	userViews, pageInfo, err := s.userViewDAO.PaginateByUserID(db, req.UserID, []string{"id", "progress", "last_viewed_at", "user_id", "article_id"}, []string{"User", "Article", "Article.Tags", "Article.User"}, req.PageParam.Page, req.PageParam.PageSize)
	if err != nil {
		logger.Error("[AssetService] failed to list user view articles", zap.Uint("userID", req.UserID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	rsp.UserViews = lo.Map(*userViews, func(userView model.UserView, _ int) *protocol.UserView {
		return &protocol.UserView{
			ViewID:       userView.ID,
			Progress:     userView.Progress,
			LastViewedAt: userView.LastViewedAt.Format(time.DateTime),
			UserID:       userView.UserID,
			ArticleID:    userView.ArticleID,
		}
	})
	rsp.PageInfo = &protocol.PageInfo{
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
		Total:    pageInfo.Total,
	}
	return rsp, nil
}

func (s *assetService) DeleteUserView(ctx context.Context, req *protocol.DeleteUserViewRequest) (rsp *protocol.DeleteUserViewResponse, err error) {
	rsp = &protocol.DeleteUserViewResponse{}

	logger := logger.LoggerWithContext(ctx)
	db := database.GetDBInstance(ctx)

	userView, err := s.userViewDAO.GetByID(db, req.ViewID, []string{"id", "user_id"}, []string{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("[AssetService] user view not found", zap.Uint("userID", req.UserID), zap.Uint("viewID", req.ViewID))
			return nil, protocol.ErrDataNotExists
		}
		logger.Error("[AssetService] failed to get user view", zap.Uint("userID", req.UserID), zap.Uint("viewID", req.ViewID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	if userView.UserID != req.UserID {
		logger.Error("[AssetService] no permission to delete view", zap.Uint("curUserID", req.UserID), zap.Uint("viewID", req.ViewID))
		return nil, protocol.ErrNoPermission
	}

	err = s.userViewDAO.Delete(db, userView)
	if err != nil {
		logger.Error("[AssetService] failed to delete user view", zap.Uint("userID", req.UserID), zap.Uint("viewID", req.ViewID), zap.Error(err))
		return nil, protocol.ErrInternalError
	}

	logger.Info("[AssetService] user view deleted successfully", zap.Uint("userID", req.UserID), zap.Uint("viewID", req.ViewID))
	return rsp, nil
}
