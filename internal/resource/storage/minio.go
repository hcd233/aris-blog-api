package storage

import (
	"fmt"

	"github.com/hcd233/Aris-blog/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/samber/lo"
)

var minioClient *minio.Client

// InitMinioClient 初始化Minio客户端
func InitMinioClient() {
	minioClient = lo.Must1(minio.New(fmt.Sprintf("%s:%s", config.MinioHost, config.MinioPort), &minio.Options{
		Creds:  credentials.NewStaticV4(config.MinioAccessID, config.MinioAccessKey, ""),
		Secure: true,
	}))
}

// GetObjectStorage 获取对象存储客户端
func GetObjectStorage() *minio.Client {
	return minioClient
}
