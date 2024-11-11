package util

import (
	"path/filepath"
	"strings"
)

func IsValidImageFormat(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".webp":
		return true
	default:
		return false
	}
}

func IsValidImageContentType(contentType string) bool {
	switch contentType {
	case "image/jpeg", "image/png", "image/gif", "image/bmp", "image/tiff", "image/webp":
		return true
	default:
		return false
	}
}
