package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "zip":
		if len(os.Args) < 4 {
			fmt.Fprintf(os.Stderr, "错误: zip命令需要两个参数: 输出文件名 源目录\n")
			printUsage()
			os.Exit(1)
		}
		outputFile := os.Args[2]
		sourcePath := os.Args[3]

		if err := compressArchive(outputFile, sourcePath); err != nil {
			fmt.Fprintf(os.Stderr, "压缩失败: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("成功压缩到: %s\n", outputFile)

	case "unzip":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "错误: unzip命令需要一个参数: 压缩包路径\n")
			printUsage()
			os.Exit(1)
		}
		archivePath := os.Args[2]

		if err := extractArchive(archivePath); err != nil {
			fmt.Fprintf(os.Stderr, "解压失败: %v\n", err)
			os.Exit(1)
		}

		prefix := getArchiveNamePrefix(archivePath)
		fmt.Printf("成功解压到: %s\n", prefix)

	default:
		fmt.Fprintf(os.Stderr, "错误: 未知命令 '%s'\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "用法:\n")
	fmt.Fprintf(os.Stderr, "  gz zip <输出文件> <源目录>    压缩文件或目录\n")
	fmt.Fprintf(os.Stderr, "  gz unzip <压缩包>             解压文件\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "支持的格式:\n")
	fmt.Fprintf(os.Stderr, "  zip, tar, tar.gz, tar.bz2, tar.xz, gz, bz2, 7z\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "示例:\n")
	fmt.Fprintf(os.Stderr, "  gz zip images.zip folder/\n")
	fmt.Fprintf(os.Stderr, "  gz zip source.tar.gz folder/\n")
	fmt.Fprintf(os.Stderr, "  gz unzip archive.zip\n")
	fmt.Fprintf(os.Stderr, "  gz unzip archive.tar.xz\n")
}

