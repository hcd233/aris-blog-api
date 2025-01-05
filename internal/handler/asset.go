package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/logger"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/service"
	"github.com/hcd233/Aris-blog/internal/util"
	"go.uber.org/zap"
)

// AssetHandler 资产服务
//
//	author centonhuang
//	update 2024-12-08 16:59:38
type AssetHandler interface {
	HandleListUserLikeArticles(c *gin.Context)
	HandleListUserLikeComments(c *gin.Context)
	HandleListUserLikeTags(c *gin.Context)
	HandleCreateBucket(c *gin.Context)
	HandleListImages(c *gin.Context)
	HandleUploadImage(c *gin.Context)
	HandleGetImage(c *gin.Context)
	HandleDeleteImage(c *gin.Context)
	HandleListUserViewArticles(c *gin.Context)
	HandleDeleteUserView(c *gin.Context)
}

type assetHandler struct {
	svc service.AssetService
}

// NewAssetHandler 创建资产处理器
//
//	return AssetHandler
//	author centonhuang
//	update 2024-12-08 17:02:21
func NewAssetHandler() AssetHandler {
	return &assetHandler{
		svc: service.NewAssetService(),
	}
}

// HandleListUserLikeArticles 列出用户喜欢的文章
//
//	param c *gin.Context
//	author centonhuang
//	update 2024-11-03 06:45:42
func (h *assetHandler) HandleListUserLikeArticles(c *gin.Context) {
	userID := c.GetUint("userID")
	param := c.MustGet("param").(*protocol.PageParam)

	req := &protocol.ListUserLikeArticlesRequest{
		CurUserID: userID,
		PageParam: param,
	}

	rsp, err := h.svc.ListUserLikeArticles(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleListUserLikeComments 列出用户喜欢的评论
//
//	param c *gin.Context
//	author centonhuang
//	update 2024-11-03 06:47:41
func (h *assetHandler) HandleListUserLikeComments(c *gin.Context) {
	userID := c.GetUint("userID")
	param := c.MustGet("param").(*protocol.PageParam)

	req := &protocol.ListUserLikeCommentsRequest{
		CurUserID: userID,
		PageParam: param,
	}

	rsp, err := h.svc.ListUserLikeComments(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleListUserLikeTags 列出用户喜欢的标签
//
//	param c *gin.Context
//	author centonhuang
//	update 2024-11-03 06:47:43
func (h *assetHandler) HandleListUserLikeTags(c *gin.Context) {
	userID := c.GetUint("userID")
	param := c.MustGet("param").(*protocol.PageParam)

	req := &protocol.ListUserLikeTagsRequest{
		CurUserID: userID,
		PageParam: param,
	}

	rsp, err := h.svc.ListUserLikeTags(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleCreateBucket 创建桶
//
//	receiver s *assetHandler
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:46:02
func (h *assetHandler) HandleCreateBucket(c *gin.Context) {
	userID := c.GetUint("userID")

	req := &protocol.CreateBucketRequest{
		CurUserID: userID,
	}

	rsp, err := h.svc.CreateBucket(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleListImages 列出图片
//
//	receiver s *assetHandler
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:46:02
func (h *assetHandler) HandleListImages(c *gin.Context) {
	userID := c.GetUint("userID")

	req := &protocol.ListImagesRequest{
		CurUserID: userID,
	}

	rsp, err := h.svc.ListImages(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleUploadImage 上传图片
//
//	receiver s *assetHandler
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:46:02
func (h *assetHandler) HandleUploadImage(c *gin.Context) {
	userID := c.GetUint("userID")
	file, err := c.FormFile("file")
	if err != nil {
		logger.Logger.Error("[HandleUploadImage] get file error", zap.Error(err))
		util.SendHTTPResponse(c, nil, protocol.ErrInternalError)
		return
	}

	contentType := file.Header.Get("Content-Type")
	size := file.Size
	fileName := file.Filename

	reader, err := file.Open()
	if err != nil {
		logger.Logger.Error("[HandleUploadImage] open file error", zap.Error(err))
		util.SendHTTPResponse(c, nil, protocol.ErrInternalError)
		return
	}
	defer reader.Close()

	req := &protocol.UploadImageRequest{
		CurUserID:   userID,
		FileName:    fileName,
		Size:        size,
		ContentType: contentType,
		ReadSeeker:  reader,
	}

	rsp, err := h.svc.UploadImage(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleGetImage 获取图片
//
//	receiver s *assetHandler
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:46:02
func (h *assetHandler) HandleGetImage(c *gin.Context) {
	userID := c.GetUint("userID")
	uri := c.MustGet("uri").(*protocol.ObjectURI)
	param := c.MustGet("param").(*protocol.ImageParam)

	req := &protocol.GetImageRequest{
		CurUserID: userID,
		ImageName: uri.ObjectName,
		Quality:   param.Quality,
	}

	rsp, err := h.svc.GetImage(req)
	if err != nil {
		util.SendHTTPResponse(c, nil, err)
		return
	}

	c.Redirect(http.StatusFound, rsp.PresignedURL)
}

// HandleDeleteImage 删除图片
//
//	receiver s *assetHandler
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:46:02
func (h *assetHandler) HandleDeleteImage(c *gin.Context) {
	userID := c.GetUint("userID")
	uri := c.MustGet("uri").(*protocol.ObjectURI)

	req := &protocol.DeleteImageRequest{
		CurUserID: userID,
		ImageName: uri.ObjectName,
	}

	rsp, err := h.svc.DeleteImage(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleListUserViewArticles 列出用户浏览的文章
//
//	receiver s *assetHandler
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:46:35
func (h *assetHandler) HandleListUserViewArticles(c *gin.Context) {
	userID := c.GetUint("userID")
	pageParam := c.MustGet("param").(*protocol.PageParam)

	req := &protocol.ListUserViewArticlesRequest{
		CurUserID: userID,
		PageParam: pageParam,
	}

	rsp, err := h.svc.ListUserViewArticles(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleDeleteUserView 删除用户浏览的文章
//
//	receiver s *assetHandler
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:46:35
func (h *assetHandler) HandleDeleteUserView(c *gin.Context) {
	userID := c.GetUint("userID")
	uri := c.MustGet("uri").(*protocol.ViewURI)

	req := &protocol.DeleteUserViewRequest{
		CurUserID: userID,
		ViewID:    uri.ViewID,
	}

	rsp, err := h.svc.DeleteUserView(req)

	util.SendHTTPResponse(c, rsp, err)
}
