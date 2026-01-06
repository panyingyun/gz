package main

import (
	"archive/tar"
	"archive/zip"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver/v3"
)

// extractArchive 根据压缩包后缀自动选择解压方法
func extractArchive(archivePath string) error {
	format := detectFormat(archivePath)
	if format == "" {
		return fmt.Errorf("不支持的文件格式: %s", archivePath)
	}

	// 获取压缩包名称前缀
	prefix := getArchiveNamePrefix(archivePath)
	extractDir := prefix

	// 创建解压目录
	if err := os.MkdirAll(extractDir, 0o755); err != nil {
		return fmt.Errorf("创建解压目录失败: %v", err)
	}

	// 根据格式解压
	var err error
	switch format {
	case "zip":
		err = extractZip(archivePath, extractDir)
	case "tar":
		err = extractTar(archivePath, extractDir)
	case "tar.gz", "tgz":
		err = extractTarGz(archivePath, extractDir)
	case "tar.bz2":
		err = extractTarBz2(archivePath, extractDir)
	case "tar.xz":
		err = extractTarXz(archivePath, extractDir)
	case "gz":
		err = extractGz(archivePath, extractDir)
	case "bz2":
		err = extractBz2(archivePath, extractDir)
	case "7z":
		err = extract7z(archivePath, extractDir)
	default:
		return fmt.Errorf("不支持的解压格式: %s", format)
	}

	if err != nil {
		return err
	}

	// 检查并处理散乱文件
	return handleScatteredFiles(extractDir)
}

// extractZip 解压zip文件
func extractZip(archivePath, extractDir string) error {
	zipReader, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		path := filepath.Join(extractDir, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.FileInfo().Mode())
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.FileInfo().Mode())
		if err != nil {
			fileReader.Close()
			return err
		}

		_, err = io.Copy(targetFile, fileReader)
		fileReader.Close()
		targetFile.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

// extractTar 解压tar文件
func extractTar(archivePath, extractDir string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	tarReader := tar.NewReader(file)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		path := filepath.Join(extractDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				return err
			}

			outFile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		}
	}

	return nil
}

// extractTarGz 解压tar.gz文件
func extractTarGz(archivePath, extractDir string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		path := filepath.Join(extractDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				return err
			}

			outFile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		}
	}

	return nil
}

// extractTarBz2 解压tar.bz2文件
func extractTarBz2(archivePath, extractDir string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	bz2Reader := bzip2.NewReader(file)
	tarReader := tar.NewReader(bz2Reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		path := filepath.Join(extractDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				return err
			}

			outFile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		}
	}

	return nil
}

// extractTarXz 解压tar.xz文件
func extractTarXz(archivePath, extractDir string) error {
	return extractWithArchiver(archivePath, extractDir)
}

// extractGz 解压gz文件（单文件）
func extractGz(archivePath, extractDir string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	// 获取原始文件名（如果gz文件头中有）
	filename := gzReader.Name
	if filename == "" {
		// 如果没有，使用压缩包名称去掉.gz后缀
		filename = getArchiveNamePrefix(archivePath)
	}

	outputPath := filepath.Join(extractDir, filename)
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, gzReader)
	return err
}

// extractBz2 解压bz2文件（单文件）
func extractBz2(archivePath, extractDir string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	bz2Reader := bzip2.NewReader(file)

	// 使用压缩包名称去掉.bz2后缀作为输出文件名
	filename := getArchiveNamePrefix(archivePath)
	outputPath := filepath.Join(extractDir, filename)

	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, bz2Reader)
	return err
}

// extract7z 解压7z文件
func extract7z(archivePath, extractDir string) error {
	return extractWithArchiver(archivePath, extractDir)
}

// extractWithArchiver 使用archiver库进行解压
func extractWithArchiver(archivePath, extractDir string) error {
	// 使用ByExtension获取格式对象
	formatObj, err := archiver.ByExtension(archivePath)
	if err != nil {
		return fmt.Errorf("无法识别格式: %v", err)
	}

	unarchiver, ok := formatObj.(archiver.Unarchiver)
	if !ok {
		return fmt.Errorf("格式不支持解压")
	}

	// archiver v3使用Unarchive函数
	return unarchiver.Unarchive(archivePath, extractDir)
}

// handleScatteredFiles 检查并处理散乱文件
func handleScatteredFiles(extractDir string) error {
	// 读取解压目录中的所有文件和目录
	var allItems []string

	err := filepath.Walk(extractDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(extractDir, path)
		if err != nil {
			return err
		}

		// 跳过根目录本身
		if relPath == "." {
			return nil
		}

		allItems = append(allItems, relPath)
		return nil
	})
	if err != nil {
		return err
	}

	if len(allItems) == 0 {
		return nil
	}

	// 检查所有项的第一级路径
	firstLevelPaths := make(map[string]bool)
	hasRootLevelItems := false

	for _, item := range allItems {
		parts := strings.Split(filepath.ToSlash(item), "/")
		if len(parts) == 1 {
			// 项直接在根级别
			hasRootLevelItems = true
			firstLevelPaths[parts[0]] = true
		} else if len(parts) > 1 {
			// 项在子目录中
			firstLevelPaths[parts[0]] = true
		}
	}

	// 如果第一级路径超过1个，或者有项直接在根级别，判定为散乱文件
	isScattered := len(firstLevelPaths) > 1 || (hasRootLevelItems && len(firstLevelPaths) > 0)

	// 如果是散乱文件，创建包装目录
	if isScattered {
		wrapperDir := filepath.Join(extractDir, "extracted")
		if err := os.MkdirAll(wrapperDir, 0o755); err != nil {
			return err
		}

		// 移动所有文件和目录到包装目录
		// 需要按深度排序，先移动深层文件，避免移动父目录时影响子目录
		for _, item := range allItems {
			oldPath := filepath.Join(extractDir, item)
			newPath := filepath.Join(wrapperDir, item)

			// 检查源路径是否存在（可能在移动过程中已经被移动）
			if _, err := os.Stat(oldPath); os.IsNotExist(err) {
				continue
			}

			if err := os.MkdirAll(filepath.Dir(newPath), 0o755); err != nil {
				return err
			}

			if err := os.Rename(oldPath, newPath); err != nil {
				return err
			}
		}
	}

	return nil
}
