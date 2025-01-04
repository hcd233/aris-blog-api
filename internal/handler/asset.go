package handler

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/logger"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/database"
	"github.com/hcd233/Aris-blog/internal/resource/database/dao"
	"github.com/hcd233/Aris-blog/internal/resource/database/model"
	obj_dao "github.com/hcd233/Aris-blog/internal/resource/storage/obj_dao"
	"github.com/hcd233/Aris-blog/internal/util"
	"github.com/samber/lo"
	"go.uber.org/zap"
	webp "golang.org/x/image/webp"
	"gorm.io/gorm"
)

// AssetService 资产服务
//
//	@author centonhuang
//	@update 2024-12-08 16:59:38
type AssetService interface {
	ListUserLikeArticlesHandler(c *gin.Context)
	ListUserLikeCommentsHandler(c *gin.Context)
	ListUserLikeTagsHandler(c *gin.Context)
	CreateBucketHandler(c *gin.Context)
	ListImagesHandler(c *gin.Context)
	UploadImageHandler(c *gin.Context)
	GetImageHandler(c *gin.Context)
	DeleteImageHandler(c *gin.Context)
	GetUserViewArticleHandler(c *gin.Context)
	ListUserViewArticlesHandler(c *gin.Context)
	DeleteUserViewHandler(c *gin.Context)
}

type assetService struct {
	db              *gorm.DB
	userDAO         *dao.UserDAO
	tagDAO          *dao.TagDAO
	articleDAO      *dao.ArticleDAO
	commentDAO      *dao.CommentDAO
	userLikeDAO     *dao.UserLikeDAO
	userViewDAO     *dao.UserViewDAO
	imageObjDAO     *obj_dao.BaseMinioObjDAO
	thumbnailObjDAO *obj_dao.BaseMinioObjDAO
}

// NewAssetService 创建资产服务
//
//	@return AssetService
//	@author centonhuang
//	@update 2024-12-08 17:02:21
func NewAssetService() AssetService {
	return &assetService{
		db:              database.GetDBInstance(),
		userDAO:         dao.GetUserDAO(),
		tagDAO:          dao.GetTagDAO(),
		articleDAO:      dao.GetArticleDAO(),
		commentDAO:      dao.GetCommentDAO(),
		userLikeDAO:     dao.GetUserLikeDAO(),
		userViewDAO:     dao.GetUserViewDAO(),
		imageObjDAO:     obj_dao.GetImageObjDAO(),
		thumbnailObjDAO: obj_dao.GetThumbnailObjDAO(),
	}
}

// ListUserLikeArticlesHandler 列出用户喜欢的文章
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-11-03 06:45:42
func (s *assetService) ListUserLikeArticlesHandler(c *gin.Context) {
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.UserURI)
	param := c.MustGet("param").(*protocol.PageParam)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to view other user's like articles",
		})
		return
	}

	user, err := s.userDAO.GetByName(s.db, userName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	userLikes, pageInfo, err := s.userLikeDAO.PaginateByUserIDAndObjectType(s.db, user.ID, model.LikeObjectTypeArticle, []string{"object_id"}, []string{}, param.Page, param.PageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserLikeError,
			Message: err.Error(),
		})
		return
	}

	articleIDs := lo.Map(*userLikes, func(like model.UserLike, idx int) uint {
		return like.ObjectID
	})

	articles, err := s.articleDAO.BatchGetByIDs(s.db, articleIDs, []string{"id", "title", "slug", "published_at", "likes", "user_id"}, []string{"User", "Tags"})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	if len(articleIDs) != len(*articles) {
		_, deletedIDs := lo.Difference(articleIDs, lo.Map(*articles, func(article model.Article, idx int) uint {
			return article.ID
		}))

		logger.Logger.Warn("[List User Like Articles]", zap.Uints("deletedIDs", deletedIDs))

		deleteLikes := lo.Map(deletedIDs, func(id uint, idx int) model.UserLike {
			return model.UserLike{
				ObjectID:   id,
				ObjectType: model.LikeObjectTypeArticle,
			}
		})

		err = s.userLikeDAO.BatchDelete(s.db, &deleteLikes)
		if err != nil {
			logger.Logger.Error("[List User Like Articles]", zap.Error(err))
			err = nil
		}
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"articles": lo.Map(*articles, func(article model.Article, index int) map[string]interface{} {
				return article.GetLikeInfo()
			}),
			"pageInfo": pageInfo,
		},
	})
}

// ListUserLikeCommentsHandler 列出用户喜欢的评论
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-11-03 06:47:41
func (s *assetService) ListUserLikeCommentsHandler(c *gin.Context) {
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.UserURI)
	param := c.MustGet("param").(*protocol.PageParam)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to view other user's like comments",
		})
		return
	}

	user, err := s.userDAO.GetByName(s.db, userName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	userLikes, pageInfo, err := s.userLikeDAO.PaginateByUserIDAndObjectType(s.db, user.ID, model.LikeObjectTypeComment, []string{"object_id"}, []string{}, param.Page, param.PageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserLikeError,
			Message: err.Error(),
		})
		return
	}

	commentIDs := lo.Map(*userLikes, func(like model.UserLike, idx int) uint {
		return like.ObjectID
	})

	comments, err := s.commentDAO.BatchGetByIDs(s.db, commentIDs, []string{"id", "user", "created_at", "content", "likes"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	if len(commentIDs) != len(*comments) {
		_, deletedIDs := lo.Difference(commentIDs, lo.Map(*comments, func(comment model.Comment, idx int) uint {
			return comment.ID
		}))

		logger.Logger.Warn("[List User Like Comments]", zap.Uints("deletedIDs", deletedIDs))

		deleteLikes := lo.Map(deletedIDs, func(id uint, idx int) model.UserLike {
			return model.UserLike{
				ObjectID:   id,
				ObjectType: model.LikeObjectTypeComment,
			}
		})

		err = s.userLikeDAO.BatchDelete(s.db, &deleteLikes)
		if err != nil {
			logger.Logger.Error("[List User Like Comments]", zap.Error(err))
			err = nil
		}
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"comments": lo.Map(*comments, func(comment model.Comment, index int) map[string]interface{} {
				return comment.GetLikeInfo()
			}),
			"pageInfo": pageInfo,
		},
	})
}

// ListUserLikeTagsHandler 列出用户喜欢的标签
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-11-03 06:47:43
func (s *assetService) ListUserLikeTagsHandler(c *gin.Context) {
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.UserURI)
	param := c.MustGet("param").(*protocol.PageParam)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to view other user's like comments",
		})
		return
	}

	user, err := s.userDAO.GetByName(s.db, userName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	userLikes, pageInfo, err := s.userLikeDAO.PaginateByUserIDAndObjectType(s.db, user.ID, model.LikeObjectTypeTag, []string{"object_id"}, []string{}, param.Page, param.PageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserLikeError,
			Message: err.Error(),
		})
		return
	}

	tagIDs := lo.Map(*userLikes, func(like model.UserLike, idx int) uint {
		return like.ObjectID
	})

	tags, err := s.tagDAO.BatchGetByIDs(s.db, tagIDs, []string{"id", "created_at", "name", "slug", "likes"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	if len(tagIDs) != len(*tags) {
		_, deletedIDs := lo.Difference(tagIDs, lo.Map(*tags, func(tag model.Tag, idx int) uint {
			return tag.ID
		}))

		logger.Logger.Warn("[List User Like Tags]", zap.Uints("deletedIDs", deletedIDs))

		deleteLikes := lo.Map(deletedIDs, func(id uint, idx int) model.UserLike {
			return model.UserLike{
				ObjectID:   id,
				ObjectType: model.LikeObjectTypeTag,
			}
		})

		err = s.userLikeDAO.BatchDelete(s.db, &deleteLikes)
		if err != nil {
			logger.Logger.Error("[List User Like Tags]", zap.Error(err))
			err = nil
		}
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"tags": lo.Map(*tags, func(tag model.Tag, index int) map[string]interface{} {
				return tag.GetLikeInfo()
			}),
			"pageInfo": pageInfo,
		},
	})
}

func (s *assetService) CreateBucketHandler(c *gin.Context) {
	userID := c.GetUint("userID")
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.UserURI)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to create other user's bucket",
		})
		return
	}

	var wg sync.WaitGroup
	var imageBucketExist, thumbnailBucketExist bool
	var imageBucketErr, thumbnailBucketErr error

	wg.Add(2)

	go func() {
		defer wg.Done()
		imageBucketExist, imageBucketErr = s.imageObjDAO.CreateBucket(userID)
	}()

	go func() {
		defer wg.Done()
		thumbnailBucketExist, thumbnailBucketErr = s.thumbnailObjDAO.CreateBucket(userID)
	}()

	wg.Wait()

	if imageBucketExist && thumbnailBucketExist {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeBucketExistError,
			Message: "Bucket already exists",
		})
		return
	}

	if imageBucketErr != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeCreateBucketError,
			Message: imageBucketErr.Error(),
		})
		return
	}

	if thumbnailBucketErr != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeCreateBucketError,
			Message: thumbnailBucketErr.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code:    protocol.CodeOk,
		Message: "Bucket created successfully",
	})
}

func (s *assetService) ListImagesHandler(c *gin.Context) {
	userID := c.GetUint("userID")
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.UserURI)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to list other user's images",
		})
		return
	}

	objectInfos, err := s.imageObjDAO.ListObjects(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeListImagesError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"objects": objectInfos,
		},
	})
}

func (s *assetService) UploadImageHandler(c *gin.Context) {
	userID := c.GetUint("userID")
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.UserURI)
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeUploadImageError,
			Message: err.Error(),
		})
		return
	}

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to upload image to other user's bucket",
		})
		return
	}

	if !util.IsValidImageFormat(file.Filename) {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeUploadImageError,
			Message: fmt.Sprintf("Invalid image format: %s", file.Filename),
		})
		return
	}

	if contentType := file.Header.Get("Content-Type"); !util.IsValidImageContentType(contentType) {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeUploadImageError,
			Message: fmt.Sprintf("Invalid image content type: %s", contentType),
		})
		return
	}

	if maxFileSize := 3 * 1024 * 1024; file.Size > int64(maxFileSize) {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeUploadImageError,
			Message: fmt.Sprintf("File size is too large(%d bytes), max file size is %d bytes", file.Size, maxFileSize),
		})
		return
	}

	rawImageReader, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeUploadImageError,
			Message: err.Error(),
		})
		return
	}
	defer rawImageReader.Close()

	var rawImage image.Image
	var imageFormat imaging.Format

	extension := filepath.Ext(file.Filename)
	switch extension {
	case ".webp":
		rawImage = lo.Must1(webp.Decode(rawImageReader))
		imageFormat = imaging.PNG
	case ".png", ".jpg", ".jpeg", ".gif":
		rawImage, _ = lo.Must2(image.Decode(rawImageReader))
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
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeUploadImageError,
			Message: err.Error(),
		})
		return
	}
	rawImageReader.Seek(0, io.SeekStart)

	var wg sync.WaitGroup
	var imageErr, thumbnailErr error

	wg.Add(2)

	go func() {
		defer wg.Done()
		imageErr = s.imageObjDAO.UploadObject(userID, file.Filename, file.Size, rawImageReader)
	}()

	go func() {
		defer wg.Done()
		thumbnailErr = s.thumbnailObjDAO.UploadObject(userID, file.Filename, int64(thumbnailBuffer.Len()), &thumbnailBuffer)
	}()

	wg.Wait()

	if imageErr != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeUploadImageError,
			Message: imageErr.Error(),
		})
		return
	}

	if thumbnailErr != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeUploadImageError,
			Message: thumbnailErr.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code:    protocol.CodeOk,
		Message: "Image uploaded successfully",
	})
}

func (s *assetService) GetImageHandler(c *gin.Context) {
	userID := c.GetUint("userID")
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.ObjectURI)
	param := c.MustGet("param").(*protocol.ImageParam)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to get other user's image",
		})
		return
	}

	var (
		objectInfo *obj_dao.ObjectInfo
		err        error
	)
	switch param.Quality {
	case "raw":
		objectInfo, err = s.imageObjDAO.DownloadObject(userID, uri.ObjectName, c.Writer)
	case "thumb":
		objectInfo, err = s.thumbnailObjDAO.DownloadObject(userID, uri.ObjectName, c.Writer)
	default:
		panic(fmt.Sprintf("Invalid image quality: %s", param.Quality))
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetImageError,
			Message: err.Error(),
		})
		return
	}

	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", objectInfo.ObjectName))
	c.Writer.Header().Set("Content-Type", objectInfo.ContentType)
	c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", objectInfo.Size))
	c.Writer.Header().Set("Last-Modified", objectInfo.LastModified.Format(http.TimeFormat))
	c.Writer.Header().Set("ETag", objectInfo.ETag)
	c.Writer.Header().Set("Cache-Control", "public, max-age=31536000")
	c.Writer.Header().Set("Expires", time.Now().AddDate(1, 0, 0).Format(http.TimeFormat))

	c.Status(http.StatusOK)
}

func (s *assetService) DeleteImageHandler(c *gin.Context) {
	userID := c.GetUint("userID")
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.ObjectURI)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to delete other user's image",
		})
		return
	}

	var wg sync.WaitGroup
	var imageErr, thumbnailErr error

	wg.Add(2)

	go func() {
		defer wg.Done()
		imageErr = s.imageObjDAO.DeleteObject(userID, uri.ObjectName)
	}()

	go func() {
		defer wg.Done()
		thumbnailErr = s.thumbnailObjDAO.DeleteObject(userID, uri.ObjectName)
	}()

	wg.Wait()

	if imageErr != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeDeleteImageError,
			Message: imageErr.Error(),
		})
		return
	}

	if thumbnailErr != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeDeleteImageError,
			Message: thumbnailErr.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code:    protocol.CodeOk,
		Message: "Image deleted successfully",
	})
}

func (s *assetService) GetUserViewArticleHandler(c *gin.Context) {
	userID := c.GetUint("userID")
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.UserURI)
	param := c.MustGet("param").(*protocol.ArticleParam)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to get other user's view",
		})
		return
	}

	user, err := s.userDAO.GetByName(s.db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	article, err := s.articleDAO.GetBySlugAndUserID(s.db, param.ArticleSlug, user.ID, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetArticleError,
			Message: err.Error(),
		})
		return
	}

	userView, err := s.userViewDAO.GetLatestViewByUserIDAndArticleID(s.db, userID, article.ID, []string{"id", "progress", "last_viewed_at", "user_id", "article_id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserViewError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"userView": userView.GetBasicInfo(),
		},
	})
}

// ListUserViewArticlesHandler 列出用户浏览的文章
//
//	@param c *gin.Context
//	@author centonhuang
//	@update 2024-11-03 06:45:42
func (s *assetService) ListUserViewArticlesHandler(c *gin.Context) {
	userID := c.GetUint("userID")
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.UserURI)
	pageParam := c.MustGet("param").(*protocol.PageParam)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to list other user's view",
		})
		return
	}

	userViews, pageInfo, err := s.userViewDAO.PaginateByUserID(s.db, userID, []string{"id", "progress", "last_viewed_at", "user_id", "article_id"}, []string{"User", "Article", "Article.Tags", "Article.User"}, pageParam.Page, pageParam.PageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserViewError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"userViews": lo.Map(*userViews, func(userView model.UserView, idx int) map[string]interface{} {
				return userView.GetDetailedInfo()
			}),
			"pageInfo": pageInfo,
		},
	})
}

func (s *assetService) DeleteUserViewHandler(c *gin.Context) {
	userID := c.GetUint("userID")
	userName := c.GetString("userName")
	uri := c.MustGet("uri").(*protocol.ViewURI)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to delete other user's view",
		})
		return
	}

	_, err := s.userDAO.GetByName(s.db, uri.UserName, []string{"id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserError,
			Message: err.Error(),
		})
		return
	}

	userView, err := s.userViewDAO.GetByID(s.db, uri.ViewID, []string{"id", "user_id"}, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetUserViewError,
			Message: err.Error(),
		})
		return
	}

	if userView.UserID != userID {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to delete other user's view",
		})
		return
	}

	err = s.userViewDAO.Delete(s.db, userView)
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeDeleteUserViewError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
	})
}
