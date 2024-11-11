package asset

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hcd233/Aris-blog/internal/protocol"
	"github.com/hcd233/Aris-blog/internal/resource/storage"
	"github.com/hcd233/Aris-blog/internal/util"
	"github.com/minio/minio-go/v7"
	"github.com/samber/lo"
)

func MakeBucketHandler(c *gin.Context) {
	userID, userName := c.MustGet("userID").(uint), c.MustGet("userName").(string)
	uri := c.MustGet("uri").(*protocol.UserURI)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to create other user's bucket",
		})
		return
	}

	objectStorage := storage.GetObjectStorage()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	exist, err := objectStorage.BucketExists(ctx, fmt.Sprintf("user-%d", userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeMakeBucketError,
			Message: err.Error(),
		})
		return
	}

	if exist {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeMakeBucketError,
			Message: "Bucket already exists",
		})
		return
	}

	err = objectStorage.MakeBucket(ctx, fmt.Sprintf("user-%d", userID), minio.MakeBucketOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeMakeBucketError,
			Message: err.Error(),
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

	objectStorage := storage.GetObjectStorage()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	objectsCh := objectStorage.ListObjects(ctx, fmt.Sprintf("user-%d", userID), minio.ListObjectsOptions{
		Recursive:    true,
		WithVersions: false,
		WithMetadata: false,
	})

	objects := make([]minio.ObjectInfo, 0)
	for object := range objectsCh {
		if object.Err != nil {
			c.JSON(http.StatusInternalServerError, protocol.Response{
				Code:    protocol.CodeListImagesError,
				Message: object.Err.Error(),
			})
			return
		}

		objects = append(objects, object)
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code: protocol.CodeOk,
		Data: map[string]interface{}{
			"objects": lo.Map(objects, func(object minio.ObjectInfo, idx int) map[string]interface{} {
				return map[string]interface{}{
					"name":    object.Key,
					"size":    object.Size,
					"modTime": object.LastModified,
				}
			}),
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

	objectStorage := storage.GetObjectStorage()

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeUploadImageError,
			Message: err.Error(),
		})
		return
	}
	defer src.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	uploadInfo, err := objectStorage.PutObject(ctx, fmt.Sprintf("user-%d", userID), file.Filename, src, file.Size, minio.PutObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeUploadImageError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, protocol.Response{
		Code:    protocol.CodeOk,
		Message: "Image uploaded successfully",
		Data: map[string]interface{}{
			"object": map[string]interface{}{
				"name":    uploadInfo.Key,
				"size":    uploadInfo.Size,
				"modTime": uploadInfo.LastModified,
			},
		},
	})
}

func GetImageHandler(c *gin.Context) {
	userID, userName := c.MustGet("userID").(uint), c.MustGet("userName").(string)
	uri := c.MustGet("uri").(*protocol.ObjectURI)

	if userName != uri.UserName {
		c.JSON(http.StatusForbidden, protocol.Response{
			Code:    protocol.CodeNotPermissionError,
			Message: "You have no permission to get other user's image",
		})
		return
	}

	objectStorage := storage.GetObjectStorage()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	object, err := objectStorage.GetObject(ctx, fmt.Sprintf("user-%d", userID), uri.ObjectName, minio.GetObjectOptions{})
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetImageError,
			Message: err.Error(),
		})
		return
	}

	objectInfo, err := object.Stat()
	if err != nil {
		c.JSON(http.StatusBadRequest, protocol.Response{
			Code:    protocol.CodeGetImageError,
			Message: err.Error(),
		})
		return
	}
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", objectInfo.Key))
	c.Writer.Header().Set("Content-Type", objectInfo.ContentType)
	c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", objectInfo.Size))
	c.Writer.Header().Set("Last-Modified", objectInfo.LastModified.Format(http.TimeFormat))
	c.Writer.Header().Set("ETag", objectInfo.ETag)
	c.Writer.Header().Set("Cache-Control", "public, max-age=31536000")
	c.Writer.Header().Set("Expires", time.Now().AddDate(1, 0, 0).Format(http.TimeFormat))

	defer object.Close()
	if _, err := io.Copy(c.Writer, object); err != nil {
		c.JSON(http.StatusInternalServerError, protocol.Response{
			Code:    protocol.CodeGetImageError,
			Message: err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}
