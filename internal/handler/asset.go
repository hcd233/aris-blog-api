package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/aris-blog-api/internal/logger"
	"github.com/hcd233/aris-blog-api/internal/protocol"
	"github.com/hcd233/aris-blog-api/internal/service"
	"github.com/hcd233/aris-blog-api/internal/util"
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
//	@Summary		列出用户喜欢的文章
//	@Description	列出用户喜欢的文章
//	@Tags			asset
//	@Accept			json
//	@Produce		json
//	@Param			page	query		protocol.PageParam	true	"分页参数"
//	@Security		ApiKeyAuth
//	@Success		200			{object}	protocol.HTTPResponse{data=protocol.ListUserLikeArticlesResponse,error=nil}
//	@Failure		400			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/asset/like/articles [get]
//	param c *gin.Context
//	author centonhuang
//	update 2024-11-03 06:45:42
func (h *assetHandler) HandleListUserLikeArticles(c *gin.Context) {
	userID := c.GetUint("userID")
	param := c.MustGet("param").(*protocol.PageParam)

	req := &protocol.ListUserLikeArticlesRequest{
		UserID:    userID,
		PageParam: param,
	}

	rsp, err := h.svc.ListUserLikeArticles(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleListUserLikeComments 列出用户喜欢的评论
//
//	@Summary		列出用户喜欢的评论
//	@Description	列出用户喜欢的评论
//	@Tags			asset
//	@Accept			json
//	@Produce		json
//	@Param			page	query		protocol.PageParam	true	"分页参数"
//	@Security		ApiKeyAuth
//	@Success		200			{object}	protocol.HTTPResponse{data=protocol.ListUserLikeCommentsResponse,error=nil}
//	@Failure		400			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/asset/like/comments [get]
//	param c *gin.Context
//	author centonhuang
//	update 2024-11-03 06:47:41
func (h *assetHandler) HandleListUserLikeComments(c *gin.Context) {
	userID := c.GetUint("userID")
	param := c.MustGet("param").(*protocol.PageParam)

	req := &protocol.ListUserLikeCommentsRequest{
		UserID:    userID,
		PageParam: param,
	}

	rsp, err := h.svc.ListUserLikeComments(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleListUserLikeTags 列出用户喜欢的标签
//
//	@Summary		列出用户喜欢的标签
//	@Description	列出用户喜欢的标签
//	@Tags			asset
//	@Accept			json
//	@Produce		json
//	@Param			page	query		protocol.PageParam	true	"分页参数"
//	@Security		ApiKeyAuth
//	@Success		200			{object}	protocol.HTTPResponse{data=protocol.ListUserLikeTagsResponse,error=nil}
//	@Failure		400			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/asset/like/tags [get]
//	param c *gin.Context
//	author centonhuang
//	update 2024-11-03 06:47:43
func (h *assetHandler) HandleListUserLikeTags(c *gin.Context) {
	userID := c.GetUint("userID")
	param := c.MustGet("param").(*protocol.PageParam)

	req := &protocol.ListUserLikeTagsRequest{
		UserID:    userID,
		PageParam: param,
	}

	rsp, err := h.svc.ListUserLikeTags(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleListImages 列出图片
//
//	@Summary		列出图片
//	@Description	列出图片
//	@Tags			asset
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200			{object}	protocol.HTTPResponse{data=protocol.ListImagesResponse,error=nil}
//	@Failure		400			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/asset/object/images [get]
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:46:02
func (h *assetHandler) HandleListImages(c *gin.Context) {
	userID := c.GetUint("userID")

	req := &protocol.ListImagesRequest{
		UserID: userID,
	}

	rsp, err := h.svc.ListImages(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleUploadImage 上传图片
//
//	@Summary		上传图片
//	@Description	上传图片
//	@Tags			asset
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			file	formData	file	true	"图片文件"
//	@Security		ApiKeyAuth
//	@Success		200			{object}	protocol.HTTPResponse{data=protocol.UploadImageResponse,error=nil}
//	@Failure		400			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/asset/object/image [post]
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
		UserID:      userID,
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
//	@Summary		获取图片
//	@Description	获取图片
//	@Tags			asset
//	@Accept			json
//	@Produce		json
//	@Param			path		path		protocol.ObjectURI	true	"对象URI"
//	@Param			param	query		protocol.ImageParam	true	"图片参数"
//	@Security		ApiKeyAuth
//	@Success		200			{object}	protocol.HTTPResponse{data=protocol.GetImageResponse,error=nil}
//	@Failure		400			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/asset/object/image/{objectName} [get]
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:46:02
func (h *assetHandler) HandleGetImage(c *gin.Context) {
	userID := c.GetUint("userID")
	uri := c.MustGet("uri").(*protocol.ObjectURI)
	param := c.MustGet("param").(*protocol.ImageParam)

	req := &protocol.GetImageRequest{
		UserID:    userID,
		ImageName: uri.ObjectName,
		Quality:   param.Quality,
	}

	rsp, err := h.svc.GetImage(req)
	if err != nil {
		util.SendHTTPResponse(c, nil, err)
		return
	}

	c.Header("Content-Type", "image/jpeg")
	c.Redirect(http.StatusFound, rsp.PresignedURL)
}

// HandleDeleteImage 删除图片
//
//	@Summary		删除图片
//	@Description	删除图片
//	@Tags			asset
//	@Accept			json
//	@Produce		json
//	@Param			path		path		protocol.ObjectURI	true	"对象URI"
//	@Security		ApiKeyAuth
//	@Success		200			{object}	protocol.HTTPResponse{data=protocol.DeleteImageResponse,error=nil}
//	@Failure		400			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/asset/object/image/{objectName} [delete]
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:46:02
func (h *assetHandler) HandleDeleteImage(c *gin.Context) {
	userID := c.GetUint("userID")
	uri := c.MustGet("uri").(*protocol.ObjectURI)

	req := &protocol.DeleteImageRequest{
		UserID:    userID,
		ImageName: uri.ObjectName,
	}

	rsp, err := h.svc.DeleteImage(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleListUserViewArticles 列出用户浏览的文章
//
//	@Summary		列出用户浏览的文章
//	@Description	列出用户浏览的文章
//	@Tags			asset
//	@Accept			json
//	@Produce		json
//	@Param			page	query		protocol.PageParam	true	"分页参数"
//	@Security		ApiKeyAuth
//	@Success		200			{object}	protocol.HTTPResponse{data=protocol.ListUserViewArticlesResponse,error=nil}
//	@Failure		400			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/asset/view/articles [get]
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:46:35
func (h *assetHandler) HandleListUserViewArticles(c *gin.Context) {
	userID := c.GetUint("userID")
	pageParam := c.MustGet("param").(*protocol.PageParam)

	req := &protocol.ListUserViewArticlesRequest{
		UserID:    userID,
		PageParam: pageParam,
	}

	rsp, err := h.svc.ListUserViewArticles(req)

	util.SendHTTPResponse(c, rsp, err)
}

// HandleDeleteUserView 删除用户的文章浏览记录
//
//	@Summary		删除用户的文章浏览记录
//	@Description	删除用户的文章浏览记录
//	@Tags			asset
//	@Accept			json
//	@Produce		json
//	@Param			path		path		protocol.ViewURI	true	"浏览URI"
//	@Security		ApiKeyAuth
//	@Success		200			{object}	protocol.HTTPResponse{data=protocol.DeleteUserViewResponse,error=nil}
//	@Failure		400			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		401			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		403			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Failure		500			{object}	protocol.HTTPResponse{data=nil,error=string}
//	@Router			/v1/asset/view/article/{viewID} [delete]
//	param c *gin.Context
//	author centonhuang
//	update 2025-01-04 15:46:35
func (h *assetHandler) HandleDeleteUserView(c *gin.Context) {
	userID := c.GetUint("userID")
	uri := c.MustGet("uri").(*protocol.ViewURI)

	req := &protocol.DeleteUserViewRequest{
		UserID: userID,
		ViewID: uri.ViewID,
	}

	rsp, err := h.svc.DeleteUserView(req)

	util.SendHTTPResponse(c, rsp, err)
}
