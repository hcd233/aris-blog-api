// Package objdao 对象存储DAO
package objdao

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/samber/lo"
)

// ObjDAO 对象存储DAO接口
//
//	author centonhuang
//	update 2025-01-05 22:45:30
type ObjDAO interface {
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

	createBucketTimeout   = 10 * time.Second
	listObjectsTimeout    = 10 * time.Second
	uploadObjectTimeout   = 30 * time.Second
	downloadObjectTimeout = 30 * time.Second
	deleteObjectTimeout   = 10 * time.Second
	presignObjectTimeout  = 10 * time.Second

	presignObjectExpire = 5 * time.Minute
)

// BaseMinioObjDAO 基础Minio对象存储DAO
//
//	author centonhuang
//	update 2025-01-05 22:45:43
type BaseMinioObjDAO struct {
	ObjectType ObjectType
	BucketName string
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

func (dao *BaseMinioObjDAO) composeDirName(userID uint) string {
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
func (dao *BaseMinioObjDAO) CreateBucket() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), createBucketTimeout)
	defer cancel()

	err = dao.client.MakeBucket(ctx, dao.BucketName, minio.MakeBucketOptions{})

	return
}

// CreateDir 创建目录
//
//	receiver dao *BaseMinioObjDAO
//	param userID uint
//	return objectInfo *ObjectInfo
//	return err error
//	author centonhuang
//	update 2025-01-18 17:37:41
func (dao *BaseMinioObjDAO) CreateDir(userID uint) (objectInfo *ObjectInfo, err error) {
	dirName := dao.composeDirName(userID)

	// 创建目录
	ctx, cancel := context.WithTimeout(context.Background(), createBucketTimeout)
	defer cancel()

	// 创建一个空的目录对象
	object, err := dao.client.PutObject(ctx, dao.BucketName, dirName+"/", nil, 0, minio.PutObjectOptions{})
	if err != nil {
		return
	}

	objectInfo = &ObjectInfo{
		ObjectName:   object.Key,
		ContentType:  "",
		Size:         object.Size,
		LastModified: object.LastModified,
		Expires:      time.Time{},
		ETag:         object.ETag,
	}

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
	dirName := dao.composeDirName(userID)
	dirName += "/"

	ctx, cancel := context.WithTimeout(context.Background(), listObjectsTimeout)
	defer cancel()

	objectCh := dao.client.ListObjects(ctx, dao.BucketName, minio.ListObjectsOptions{
		Prefix:     dirName,
		StartAfter: dirName,
	})

	for object := range objectCh {
		if object.Err != nil {
			err = object.Err
			return
		}

		// 跳过目录本身
		if object.Key == dirName {
			continue
		}

		objectInfos = append(objectInfos, ObjectInfo{
			ObjectName:   strings.TrimPrefix(object.Key, dirName),
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
	dirName := dao.composeDirName(userID)
	objectName = path.Join(dirName, objectName)

	ctx, cancel := context.WithTimeout(context.Background(), uploadObjectTimeout)
	defer cancel()

	_, err = dao.client.PutObject(ctx, dao.BucketName, objectName, reader, size, minio.PutObjectOptions{})
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
	dirName := dao.composeDirName(userID)
	objectName = path.Join(dirName, objectName)

	ctx, cancel := context.WithTimeout(context.Background(), downloadObjectTimeout)
	defer cancel()

	object, err := dao.client.GetObject(ctx, dao.BucketName, objectName, minio.GetObjectOptions{})
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
	dirName := dao.composeDirName(userID)
	objectName = path.Join(dirName, objectName)

	ctx, cancel := context.WithTimeout(context.Background(), presignObjectTimeout)
	defer cancel()

	presignedURL, err = dao.client.PresignedGetObject(ctx, dao.BucketName, objectName, presignObjectExpire, nil)
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
	dirName := dao.composeDirName(userID)
	objectName = path.Join(dirName, objectName)

	ctx, cancel := context.WithTimeout(context.Background(), deleteObjectTimeout)
	defer cancel()

	err = dao.client.RemoveObject(ctx, dao.BucketName, objectName, minio.RemoveObjectOptions{})
	return
}
