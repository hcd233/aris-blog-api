package objdao

import (
	"sync"

	"github.com/hcd233/aris-blog-api/internal/config"
	"github.com/hcd233/aris-blog-api/internal/resource/storage"
)

var (

	// ImageObjDAOSingleton 图片对象DAO单例
	//	update 2025-01-05 22:45:54
	ImageObjDAOSingleton *BaseMinioObjDAO

	// ThumbnailObjDAOSingleton 缩略图对象DAO单例
	//	update 2025-01-05 22:45:54
	ThumbnailObjDAOSingleton *BaseMinioObjDAO

	imageObjOnce     sync.Once
	thumbnailObjOnce sync.Once
)

// GetImageObjDAO 获取图片对象DAO单例
//
//	return *BaseMinioObjDAO
//	author centonhuang
//	update 2024-10-18 01:10:28
func GetImageObjDAO() *BaseMinioObjDAO {
	imageObjOnce.Do(func() {
		ImageObjDAOSingleton = &BaseMinioObjDAO{
			ObjectType: ObjectTypeImage,
			BucketName: config.MinioBucketName,
			client:     storage.GetObjectStorage(),
		}
	})
	return ImageObjDAOSingleton
}

// GetThumbnailObjDAO 获取缩略图对象DAO单例
//
//	return *BaseMinioObjDAO
//	author centonhuang
//	update 2024-10-18 01:09:59
func GetThumbnailObjDAO() *BaseMinioObjDAO {
	thumbnailObjOnce.Do(func() {
		ThumbnailObjDAOSingleton = &BaseMinioObjDAO{
			ObjectType: ObjectTypeThumbnail,
			BucketName: config.MinioBucketName,
			client:     storage.GetObjectStorage(),
		}
	})
	return ThumbnailObjDAOSingleton
}
