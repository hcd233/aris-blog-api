package storage

import (
	"github.com/hcd233/aris-blog-api/internal/config"
)

// Provider 存储提供商
type Provider string

const (
	// ProviderMinio Minio存储
	ProviderMinio Provider = "minio"
	// ProviderCOS 腾讯云COS存储
	ProviderCOS Provider = "cos"
)

var provider Provider

// InitObjectStorage 初始化对象存储
//
//	author centonhuang
//	update 2024-12-09 15:59:06
func InitObjectStorage() {
	provider = GetProvider()

	switch provider {
	case ProviderMinio:
		initMinioClient()
	case ProviderCOS:
		initCosClient()
	}
}

// GetProvider 获取存储提供商
//
//	return Provider
//	author centonhuang
//	update 2025-01-19 14:13:22
func GetProvider() Provider {
	// 优先使用 COS
	if config.CosAppID != "" {
		return ProviderCOS
	}

	if config.MinioEndpoint != "" {
		return ProviderMinio
	}

	panic("no object storage configured")
}
