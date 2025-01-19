// Package objdao 对象存储DAO
package objdao

import (
	"io"
	"net/url"
	"time"
)

// ObjDAO 对象存储DAO接口
//
//	author centonhuang
//	update 2025-01-05 22:45:30
type ObjDAO interface {
	GetBucketName() string
	CreateBucket() (err error)
	CreateDir(userID uint) (objectInfo *ObjectInfo, err error)
	ListObjects(userID uint) (objectInfos []ObjectInfo, err error)
	UploadObject(userID uint, objectName string, size int64, reader io.Reader) (err error)
	DownloadObject(userID uint, objectName string, writer io.Writer) (objectInfo *ObjectInfo, err error)
	PresignObject(userID uint, objectName string) (presignedURL *url.URL, err error)
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
