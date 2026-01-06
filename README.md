# gz - 自动化压缩解压命令行工具

一个使用Go语言实现的智能压缩解压工具，能够根据文件后缀自动识别并选择合适的压缩/解压方法。

## 功能特性

- 🎯 **自动格式识别**：根据文件扩展名自动选择压缩/解压方法
- 📦 **多格式支持**：支持8种常见压缩格式
- 🗂️ **智能目录管理**：解压时自动创建目录，目录名为压缩包名称前缀
- 📁 **散乱文件处理**：自动检测并整理散乱文件到统一目录

## 支持的格式

- `zip` - ZIP压缩格式
- `tar` - TAR归档格式
- `tar.gz` / `tgz` - TAR+GZIP压缩
- `tar.bz2` - TAR+BZIP2压缩
- `tar.xz` - TAR+XZ压缩
- `gz` - GZIP单文件压缩
- `bz2` - BZIP2单文件压缩
- `7z` - 7-Zip压缩格式

## 安装

```bash
# 克隆仓库
git clone <repository-url>
cd gz

# 构建
go build -o gz.exe

# 或者直接运行
go run .
```

## 使用方法

### 压缩文件/目录

```bash
# 压缩为ZIP格式
gz zip images.zip folder/

# 压缩为TAR.GZ格式
gz zip source.tar.gz folder/

# 压缩为TAR.BZ2格式
gz zip archive.tar.bz2 folder/

# 压缩为7Z格式
gz zip archive.7z folder/
```

### 解压文件

```bash
# 解压ZIP文件
gz unzip archive.zip

# 解压TAR.GZ文件
gz unzip archive.tar.gz

# 解压TAR.XZ文件
gz unzip archive.tar.xz

# 解压7Z文件
gz unzip archive.7z
```

## 特性说明

### 自动目录创建

解压时会自动创建一个目录，目录名为压缩包名称的前缀。

例如：
- 解压 `archive.zip` → 创建 `archive/` 目录
- 解压 `source.tar.gz` → 创建 `source/` 目录

### 散乱文件处理

当压缩包内包含散乱的文件（未包含在一个统一的文件夹中）时，工具会自动创建一个 `extracted/` 目录，将所有文件整理到该目录中，保持文件结构的整洁。

## 示例

```bash
# 压缩当前目录下的docs文件夹为ZIP格式
gz zip docs.zip docs/

# 解压下载的压缩包
gz unzip download.tar.gz
# 会自动创建 download/ 目录，并将内容解压到其中

# 如果压缩包内文件散乱，会自动整理到 extracted/ 目录
gz unzip messy_files.zip
# 散乱文件会被整理到 download/extracted/ 目录中
```

## 注意事项

- `gz` 和 `bz2` 格式只支持单文件压缩，不支持目录
- 压缩大文件时可能需要一些时间，请耐心等待
- 确保有足够的磁盘空间用于解压操作

## 错误处理

工具会提供友好的错误提示：
- 文件不存在
- 格式不支持
- 权限错误
- 压缩/解压过程中的其他错误

## 许可证

本项目采用 GNU General Public License v3.0 许可证。

