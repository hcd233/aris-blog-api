package objdao

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/samber/lo"
)

type ObjDAO interface {
	ComposeBucketName(userID uint) (bucketName string)
	CreateBucket(userID uint) (exist bool, err error)
	ListObjects(userID uint) (objectInfos []ObjectInfo, err error)
	UploadObject(userID uint, objectName string, size int64, reader io.Reader) (err error)
	DownloadObject(userID uint, objectName string) (reader io.ReadCloser, objectInfo ObjectInfo, err error)
	DeleteObject(userID uint, objectName string) (err error)
}

type ObjectType string

const (
	ObjectTypeImage     ObjectType = "image"
	ObjectTypeThumbnail ObjectType = "thumbnail"

	CreateBucketTimeout   = 5 * time.Second
	ListObjectsTimeout    = 5 * time.Second
	UploadObjectTimeout   = 20 * time.Second
	DownloadObjectTimeout = 20 * time.Second
	DeleteObjectTimeout   = 5 * time.Second
)

type BaseMinioObjDAO struct {
	ObjectType ObjectType
	client     *minio.Client
}

type ObjectInfo struct {
	ObjectName   string    `json:"objectName"`
	ContentType  string    `json:"contentType"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"lastModified"`
	Expires      time.Time `json:"expires"`
	ETag         string    `json:"etag"`
}

func (dao *BaseMinioObjDAO) ComposeBucketName(userID uint) string {
	return fmt.Sprintf("user-%d-%s", userID, dao.ObjectType)
}

func (dao *BaseMinioObjDAO) CreateBucket(userID uint) (exist bool, err error) {
	bucketName := dao.ComposeBucketName(userID)

	ctx, cancel := context.WithTimeout(context.Background(), CreateBucketTimeout)
	defer cancel()

	exist, err = dao.client.BucketExists(ctx, bucketName)

	if exist {
		return
	}

	err = dao.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	return
}

func (dao *BaseMinioObjDAO) ListObjects(userID uint) (objectInfos []ObjectInfo, err error) {
	bucketName := dao.ComposeBucketName(userID)

	ctx, cancel := context.WithTimeout(context.Background(), ListObjectsTimeout)
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

func (dao *BaseMinioObjDAO) UploadObject(userID uint, objectName string, size int64, reader io.Reader) (err error) {
	bucketName := dao.ComposeBucketName(userID)

	ctx, cancel := context.WithTimeout(context.Background(), UploadObjectTimeout)
	defer cancel()

	_, err = dao.client.PutObject(ctx, bucketName, objectName, reader, size, minio.PutObjectOptions{})
	return
}

func (dao *BaseMinioObjDAO) DownloadObject(userID uint, objectName string) (reader io.ReadCloser, objectInfo *ObjectInfo, err error) {
	bucketName := dao.ComposeBucketName(userID)

	ctx, cancel := context.WithTimeout(context.Background(), DownloadObjectTimeout)
	defer cancel()

	object, err := dao.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})

	if err != nil {
		return
	}

	stat := lo.Must1(object.Stat())

	objectInfo = &ObjectInfo{
		ObjectName:   stat.Key,
		ContentType:  stat.ContentType,
		Size:         stat.Size,
		LastModified: stat.LastModified,
		Expires:      stat.Expires,
		ETag:         stat.ETag,
	}

	return object, objectInfo, nil
}

func (dao *BaseMinioObjDAO) DeleteObject(userID uint, objectName string) (err error) {
	bucketName := dao.ComposeBucketName(userID)

	ctx, cancel := context.WithTimeout(context.Background(), DeleteObjectTimeout)
	defer cancel()

	err = dao.client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	return
}
