// Package storage Minio对象存储模块
//
//	update 2024-12-09 15:58:58
package storage

import (
	"context"

	"github.com/hcd233/Aris-blog/internal/config"
	"github.com/hcd233/Aris-blog/internal/logger"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

var minioClient *minio.Client

// InitObjectStorage 初始化对象存储
//
//	author centonhuang
//	update 2024-12-09 15:59:06
func InitObjectStorage() {
	minioClient = lo.Must1(minio.New(config.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.MinioAccessID, config.MinioAccessKey, ""),
		Secure: config.MinioTLS,
	}))

	_ = lo.Must1(minioClient.ListBuckets(context.Background()))

	logger.Logger.Info("[Object Storage] Connected to Minio database", zap.String("endpoint", config.MinioEndpoint))
}

// GetObjectStorage 获取对象存储客户端
func GetObjectStorage() *minio.Client {
	return minioClient
}
