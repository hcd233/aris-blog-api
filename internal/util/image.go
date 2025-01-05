// Package util 工具函数
//
//	update 2024-12-09 17:21:16
package util

import (
	"path/filepath"
	"strings"
)

// IsValidImageFormat 判断文件是否为图片格式
//
//	param filename string
//	return bool
//	author centonhuang
//	update 2024-12-09 17:18:26
func IsValidImageFormat(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".webp":
		return true
	default:
		return false
	}
}

// IsValidImageContentType 判断内容类型是否为图片
//
//	param contentType string
//	return bool
//	author centonhuang
//	update 2024-12-09 17:18:26
func IsValidImageContentType(contentType string) bool {
	switch contentType {
	case "image/jpeg", "image/png", "image/gif", "image/bmp", "image/tiff", "image/webp":
		return true
	default:
		return false
	}
}
