package main

import (
	"path/filepath"
	"strings"
)

// detectFormat 根据文件扩展名识别压缩格式
func detectFormat(filename string) string {
	filename = strings.ToLower(filename)

	// 按优先级检查复合扩展名（从长到短）
	if strings.HasSuffix(filename, ".tar.xz") {
		return "tar.xz"
	}
	if strings.HasSuffix(filename, ".tar.bz2") {
		return "tar.bz2"
	}
	if strings.HasSuffix(filename, ".tar.gz") {
		return "tar.gz"
	}
	if strings.HasSuffix(filename, ".tar") {
		return "tar"
	}
	if strings.HasSuffix(filename, ".tgz") {
		return "tar.gz"
	}
	if strings.HasSuffix(filename, ".zip") {
		return "zip"
	}
	if strings.HasSuffix(filename, ".gz") {
		return "gz"
	}
	if strings.HasSuffix(filename, ".bz2") {
		return "bz2"
	}
	if strings.HasSuffix(filename, ".7z") {
		return "7z"
	}

	return ""
}

// getArchiveNamePrefix 提取压缩包名称前缀（用于创建解压目录）
// 例如：archive.tar.gz -> archive
func getArchiveNamePrefix(filename string) string {
	base := filepath.Base(filename)

	// 移除所有压缩扩展名
	extensions := []string{".tar.xz", ".tar.bz2", ".tar.gz", ".tgz", ".tar", ".zip", ".gz", ".bz2", ".7z"}

	for _, ext := range extensions {
		if strings.HasSuffix(strings.ToLower(base), ext) {
			return strings.TrimSuffix(base, ext)
		}
	}

	// 如果没有匹配的扩展名，返回不带扩展名的文件名
	ext := filepath.Ext(base)
	if ext != "" {
		return strings.TrimSuffix(base, ext)
	}

	return base
}
