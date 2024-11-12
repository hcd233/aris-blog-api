package asset

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	obj_dao "github.com/hcd233/Aris-blog/internal/resource/storage/obj_dao"
	"github.com/hcd233/Aris-blog/internal/util"
	"github.com/samber/lo"
)

func CreateBucketHandler(c *gin.Context) {
	userID, userName := c.MustGet("userID").(uint), c.MustGet("userName").(string)
	uri := c.MustGet("uri").(*protocol.UserURI)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to create other user's bucket",
		})
		return
	}

	imageObjDAO, thumbnailObjDAO := obj_dao.GetImageObjDAO(), obj_dao.GetThumbnailObjDAO()

	var wg sync.WaitGroup
	var imageBucketExist, thumbnailBucketExist bool
	var imageBucketErr, thumbnailBucketErr error

	wg.Add(2)

	go func() {
		defer wg.Done()
		imageBucketExist, imageBucketErr = imageObjDAO.CreateBucket(userID)
	}()

	go func() {
		defer wg.Done()
		thumbnailBucketExist, thumbnailBucketErr = thumbnailObjDAO.CreateBucket(userID)
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

func ListImagesHandler(c *gin.Context) {
	userID, userName := c.MustGet("userID").(uint), c.MustGet("userName").(string)
	uri := c.MustGet("uri").(*protocol.UserURI)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to list other user's images",
		})
		return
	}

	imageObjDAO := obj_dao.GetImageObjDAO()

	objectInfos, err := imageObjDAO.ListObjects(userID)
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

func UploadImageHandler(c *gin.Context) {
	userID, userName := c.MustGet("userID").(uint), c.MustGet("userName").(string)
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

	rawImage, _ := lo.Must2(image.Decode(rawImageReader))
	// restrict image into 512*512 max size
	x, y := rawImage.Bounds().Dx(), rawImage.Bounds().Dy()

	maxPixel := 512

	for ; x > maxPixel || y > maxPixel; x, y = x/2, y/2 {
	}

	thumbnailImage := imaging.Thumbnal(rawImage, x, y, imaging.Lanczos)

	var thumbnailBuffer bytes.Buffer
	err = imaging.Encode(&thumbnailBuffer, thumbnailImage, lo.Must1(imaging.FormatFromFilename(file.Filename)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeUploadImageError,
			Message: err.Error(),
		})
		return
	}
	rawImageReader.Seek(0, io.SeekStart)

	imageObjDAO, thumbnailObjDAO := obj_dao.GetImageObjDAO(), obj_dao.GetThumbnailObjDAO()

	var wg sync.WaitGroup
	var imageErr, thumbnailErr error

	wg.Add(2)

	go func() {
		defer wg.Done()
		imageErr = imageObjDAO.UploadObject(userID, file.Filename, file.Size, rawImageReader)
	}()

	go func() {
		defer wg.Done()
		thumbnailErr = thumbnailObjDAO.UploadObject(userID, file.Filename, int64(thumbnailBuffer.Len()), &thumbnailBuffer)
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

func GetImageHandler(c *gin.Context) {
	userID, userName := c.MustGet("userID").(uint), c.MustGet("userName").(string)
	uri := c.MustGet("uri").(*protocol.ObjectURI)
	param := c.MustGet("param").(*protocol.ImageParam)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to get other user's image",
		})
		return
	}

	imageObjDAO, thumbnailObjDAO := obj_dao.GetImageObjDAO(), obj_dao.GetThumbnailObjDAO()

	var (
		object     io.ReadCloser
		objectInfo *obj_dao.ObjectInfo
		err        error
	)
	switch param.Quality {
	case "raw":
		object, objectInfo, err = imageObjDAO.DownloadObject(userID, uri.ObjectName)
	case "thumb":
		object, objectInfo, err = thumbnailObjDAO.DownloadObject(userID, uri.ObjectName)
	default:
		panic(fmt.Sprintf("Invalid image quality: %s", param.Quality))
	}
	defer object.Close()

	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetImageError,
			Message: err.Error(),
		})
		return
	}

	_ = lo.Must1(io.Copy(c.Writer, object))

	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", objectInfo.ObjectName))
	c.Writer.Header().Set("Content-Type", objectInfo.ContentType)
	c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", objectInfo.Size))
	c.Writer.Header().Set("Last-Modified", objectInfo.LastModified.Format(http.TimeFormat))
	c.Writer.Header().Set("ETag", objectInfo.ETag)
	c.Writer.Header().Set("Cache-Control", "public, max-age=31536000")
	c.Writer.Header().Set("Expires", time.Now().AddDate(1, 0, 0).Format(http.TimeFormat))

	c.Status(http.StatusOK)
}

func DeleteImageHandler(c *gin.Context) {
	userID, userName := c.MustGet("userID").(uint), c.MustGet("userName").(string)
	uri := c.MustGet("uri").(*protocol.ObjectURI)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to delete other user's image",
		})
		return
	}

	imageObjDAO, thumbnailObjDAO := obj_dao.GetImageObjDAO(), obj_dao.GetThumbnailObjDAO()

	var wg sync.WaitGroup
	var imageErr, thumbnailErr error

	wg.Add(2)

	go func() {
		defer wg.Done()
		imageErr = imageObjDAO.DeleteObject(userID, uri.ObjectName)
	}()

	go func() {
		defer wg.Done()
		thumbnailErr = thumbnailObjDAO.DeleteObject(userID, uri.ObjectName)
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
