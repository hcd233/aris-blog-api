// Package objdao 对象存储DAO
package objdao

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/samber/lo"
)

// ObjDAO 对象存储DAO接口
//
//	author centonhuang
//	update 2025-01-05 22:45:30
type ObjDAO interface {
	composeBucketName(userID uint) (bucketName string)
	CreateBucket(userID uint) (exist bool, err error)
	ListObjects(userID uint) (objectInfos []ObjectInfo, err error)
	UploadObject(userID uint, objectName string, size int64, reader io.Reader) (err error)
	DownloadObject(userID uint, objectName string, writer io.Writer) (objectInfo *ObjectInfo, err error)
	DeleteObject(userID uint, objectName string) (err error)
}

// ObjectType 对象类型
//
//	author centonhuang
//	update 2025-01-05 22:45:37
type ObjectType string

const (

	// ObjectTypeImage ObjectType
	//	update 2025-01-05 17:36:01
	ObjectTypeImage ObjectType = "image"

	// ObjectTypeThumbnail ObjectType
	//	update 2025-01-05 17:36:05
	ObjectTypeThumbnail ObjectType = "thumbnail"

	createBucketTimeout   = 5 * time.Second
	listObjectsTimeout    = 5 * time.Second
	uploadObjectTimeout   = 20 * time.Second
	downloadObjectTimeout = 20 * time.Second
	deleteObjectTimeout   = 5 * time.Second
	presignObjectTimeout  = 5 * time.Second

	presignObjectExpire = 5 * time.Minute
)

// BaseMinioObjDAO 基础Minio对象存储DAO
//
//	author centonhuang
//	update 2025-01-05 22:45:43
type BaseMinioObjDAO struct {
	ObjectType ObjectType
	client     *minio.Client
}

// ObjectInfo 对象信息
//
//	author centonhuang
//	update 2025-01-05 22:45:48
type ObjectInfo struct {
	ObjectName   string    `json:"objectName"`
	ContentType  string    `json:"contentType"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"lastModified"`
	Expires      time.Time `json:"expires"`
	ETag         string    `json:"etag"`
}

func (dao *BaseMinioObjDAO) composeBucketName(userID uint) string {
	return fmt.Sprintf("user-%d-%s", userID, dao.ObjectType)
}

// CreateBucket 创建桶
//
//	receiver dao *BaseMinioObjDAO
//	param userID uint
//	return exist bool
//	return err error
//	author centonhuang
//	update 2025-01-05 17:37:41
func (dao *BaseMinioObjDAO) CreateBucket(userID uint) (exist bool, err error) {
	bucketName := dao.composeBucketName(userID)

	ctx, cancel := context.WithTimeout(context.Background(), createBucketTimeout)
	defer cancel()

	exist, err = dao.client.BucketExists(ctx, bucketName)

	if exist {
		return
	}

	err = dao.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		return
	}

	// 确保不设置任何允许公共访问的策略
	// 这里我们不需要设置任何策略，因为默认情况下，桶和对象是私有的
	return
}

// ListObjects 列出桶中的对象
//
//	receiver dao *BaseMinioObjDAO
//	param userID uint
//	return objectInfos []ObjectInfo
//	return err error
//	author centonhuang
//	update 2025-01-05 17:37:45
func (dao *BaseMinioObjDAO) ListObjects(userID uint) (objectInfos []ObjectInfo, err error) {
	bucketName := dao.composeBucketName(userID)

	ctx, cancel := context.WithTimeout(context.Background(), listObjectsTimeout)
	defer cancel()

	objectCh := dao.client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{})
	for object := range objectCh {
		if object.Err != nil {
			err = object.Err
			return
		}

		objectInfos = append(objectInfos, ObjectInfo{
			ObjectName:   object.Key,
			ContentType:  object.ContentType,
			Size:         object.Size,
			LastModified: object.LastModified,
			Expires:      object.Expires,
			ETag:         object.ETag,
		})
	}
	return
}

// UploadObject 上传对象
//
//	receiver dao *BaseMinioObjDAO
//	param userID uint
//	param objectName string
//	param size int64
//	param reader io.Reader
//	return err error
//	author centonhuang
//	update 2025-01-05 17:37:50
func (dao *BaseMinioObjDAO) UploadObject(userID uint, objectName string, size int64, reader io.Reader) (err error) {
	bucketName := dao.composeBucketName(userID)

	ctx, cancel := context.WithTimeout(context.Background(), uploadObjectTimeout)
	defer cancel()

	_, err = dao.client.PutObject(ctx, bucketName, objectName, reader, size, minio.PutObjectOptions{})
	return
}

// DownloadObject 下载对象
//
//	receiver dao *BaseMinioObjDAO
//	param userID uint
//	param objectName string
//	param writer io.Writer
//	return objectInfo *ObjectInfo
//	return err error
//	author centonhuang
//	update 2025-01-05 17:37:57
func (dao *BaseMinioObjDAO) DownloadObject(userID uint, objectName string, writer io.Writer) (objectInfo *ObjectInfo, err error) {
	bucketName := dao.composeBucketName(userID)

	ctx, cancel := context.WithTimeout(context.Background(), downloadObjectTimeout)
	defer cancel()

	object, err := dao.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return
	}
	defer object.Close()

	stat := lo.Must1(object.Stat())

	objectInfo = &ObjectInfo{
		ObjectName:   stat.Key,
		ContentType:  stat.ContentType,
		Size:         stat.Size,
		LastModified: stat.LastModified,
		Expires:      stat.Expires,
		ETag:         stat.ETag,
	}

	_, err = io.Copy(writer, object)

	return
}

// PresignObject 生成对象的预签名URL
//
//	receiver dao *BaseMinioObjDAO
//	param userID uint
//	param objectName string
//	param writer io.Writer
//	return url *url.URL
//	return err error
//	author centonhuang
//	update 2025-01-05 17:38:03
func (dao *BaseMinioObjDAO) PresignObject(userID uint, objectName string) (presignedURL *url.URL, err error) {
	bucketName := dao.composeBucketName(userID)

	ctx, cancel := context.WithTimeout(context.Background(), presignObjectTimeout)
	defer cancel()

	presignedURL, err = dao.client.PresignedGetObject(ctx, bucketName, objectName, presignObjectExpire, nil)
	return
}

// DeleteObject 删除对象
//
//	receiver dao *BaseMinioObjDAO
//	param userID uint
//	param objectName string
//	return err error
//	author centonhuang
//	update 2025-01-05 17:38:09
func (dao *BaseMinioObjDAO) DeleteObject(userID uint, objectName string) (err error) {
	bucketName := dao.composeBucketName(userID)

	ctx, cancel := context.WithTimeout(context.Background(), deleteObjectTimeout)
	defer cancel()

	err = dao.client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	return
}
