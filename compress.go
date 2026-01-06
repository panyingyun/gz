package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mholt/archiver/v3"
)

// compressArchive 根据输出文件名自动选择压缩方法
func compressArchive(outputFile, sourcePath string) error {
	format := detectFormat(outputFile)
	if format == "" {
		return fmt.Errorf("不支持的文件格式: %s", outputFile)
	}

	// 检查源路径是否存在
	info, err := os.Stat(sourcePath)
	if err != nil {
		return fmt.Errorf("源路径不存在: %v", err)
	}

	switch format {
	case "zip":
		return compressZip(outputFile, sourcePath)
	case "tar":
		return compressTar(outputFile, sourcePath)
	case "tar.gz", "tgz":
		return compressTarGz(outputFile, sourcePath)
	case "tar.bz2":
		return compressTarBz2(outputFile, sourcePath)
	case "tar.xz":
		return compressTarXz(outputFile, sourcePath)
	case "gz":
		if info.IsDir() {
			return fmt.Errorf("gz格式只支持单文件压缩，不支持目录")
		}
		return compressGz(outputFile, sourcePath)
	case "bz2":
		if info.IsDir() {
			return fmt.Errorf("bz2格式只支持单文件压缩，不支持目录")
		}
		return compressBz2(outputFile, sourcePath)
	case "7z":
		return compress7z(outputFile, sourcePath)
	default:
		return fmt.Errorf("不支持的压缩格式: %s", format)
	}
}

// compressZip 压缩为zip格式
func compressZip(outputFile, sourcePath string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	return filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 计算相对路径
		relPath, err := filepath.Rel(sourcePath, path)
		if err != nil {
			return err
		}

		// 如果是源目录本身，跳过
		if relPath == "." {
			return nil
		}

		// 创建zip文件头
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// 使用相对路径作为zip内的路径
		header.Name = filepath.ToSlash(relPath)
		if info.IsDir() {
			header.Name += "/"
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// 如果是目录，不需要写入内容
		if info.IsDir() {
			return nil
		}

		// 写入文件内容
		sourceFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer sourceFile.Close()

		_, err = io.Copy(writer, sourceFile)
		return err
	})
}

// compressTar 压缩为tar格式
func compressTar(outputFile, sourcePath string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	tarWriter := tar.NewWriter(file)
	defer tarWriter.Close()

	return filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(sourcePath, path)
		if err != nil {
			return err
		}

		if relPath == "." {
			return nil
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		header.Name = filepath.ToSlash(relPath)

		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		sourceFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer sourceFile.Close()

		_, err = io.Copy(tarWriter, sourceFile)
		return err
	})
}

// compressTarGz 压缩为tar.gz格式
func compressTarGz(outputFile, sourcePath string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	gzWriter := gzip.NewWriter(file)
	defer gzWriter.Close()

	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	return filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(sourcePath, path)
		if err != nil {
			return err
		}

		if relPath == "." {
			return nil
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		header.Name = filepath.ToSlash(relPath)

		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		sourceFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer sourceFile.Close()

		_, err = io.Copy(tarWriter, sourceFile)
		return err
	})
}

// compressTarBz2 压缩为tar.bz2格式
func compressTarBz2(outputFile, sourcePath string) error {
	// Go标准库的compress/bzip2只支持解压，不支持压缩
	// 使用archiver库
	return compressWithArchiver(outputFile, sourcePath, "tar.bz2")
}

// compressTarXz 压缩为tar.xz格式
func compressTarXz(outputFile, sourcePath string) error {
	// Go标准库不支持xz压缩，使用archiver库
	return compressWithArchiver(outputFile, sourcePath, "tar.xz")
}

// compressGz 压缩单文件为gz格式
func compressGz(outputFile, sourcePath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	outputFileHandle, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer outputFileHandle.Close()

	gzWriter := gzip.NewWriter(outputFileHandle)
	defer gzWriter.Close()

	_, err = io.Copy(gzWriter, sourceFile)
	return err
}

// compressBz2 压缩单文件为bz2格式
func compressBz2(outputFile, sourcePath string) error {
	// Go标准库的compress/bzip2只支持解压，使用archiver库
	return compressWithArchiver(outputFile, sourcePath, "bz2")
}

// compress7z 压缩为7z格式
func compress7z(outputFile, sourcePath string) error {
	return compressWithArchiver(outputFile, sourcePath, "7z")
}

// compressWithArchiver 使用archiver库进行压缩
func compressWithArchiver(outputFile, sourcePath, format string) error {
	// 使用ByExtension获取格式对象
	formatObj, err := archiver.ByExtension(outputFile)
	if err != nil {
		return fmt.Errorf("无法识别格式: %v", err)
	}

	archiverObj, ok := formatObj.(archiver.Archiver)
	if !ok {
		return fmt.Errorf("格式不支持压缩: %s", format)
	}

	// archiver v3使用Archive函数，需要传入源文件列表和目标文件
	return archiverObj.Archive([]string{sourcePath}, outputFile)
}
