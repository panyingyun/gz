# Go压缩解压命令行工具实现计划

## 项目结构

```
gz/
├── main.go          # 主程序入口，命令行参数解析
├── compress.go      # 压缩功能实现
├── extract.go       # 解压功能实现
├── format.go        # 格式检测和工具函数
├── go.mod           # Go模块依赖管理
└── README.md        # 使用说明文档
```

## 核心功能设计

### 1. 命令行接口 (`main.go`)

- 解析命令行参数：`gz zip <输出文件> <源目录>` 和 `gz unzip <压缩包>`
- 根据操作类型（zip/unzip）调用相应处理函数
- 错误处理和用户友好的错误提示

### 2. 格式识别 (`format.go`)

- `detectFormat(filename string) string` - 根据文件扩展名识别格式
- 支持的格式映射：
  - `.zip` → zip
  - `.tar` → tar
  - `.tar.gz`, `.tgz` → tar.gz
  - `.tar.bz2` → tar.bz2
  - `.tar.xz` → tar.xz
  - `.gz` → gz
  - `.bz2` → bz2
  - `.7z` → 7z
- `getArchiveNamePrefix(filename string) string` - 提取压缩包名称前缀（用于创建解压目录）

### 3. 压缩功能 (`compress.go`)

- `compressArchive(outputFile, sourcePath string) error` - 主压缩函数
- 根据输出文件名后缀自动选择压缩方法：
  - `zip` → 使用 `archive/zip`
  - `tar`, `tar.gz`, `tar.bz2`, `tar.xz` → 使用 `archive/tar` + 相应压缩
  - `gz`, `bz2` → 单文件压缩
  - `7z` → 使用 `github.com/mholt/archiver`

### 4. 解压功能 (`extract.go`)

- `extractArchive(archivePath string) error` - 主解压函数
- 根据压缩包后缀自动选择解压方法
- 自动创建目录：使用压缩包名称前缀创建目录
- 智能处理散乱文件：
  - 检查解压后的文件是否都在根目录（没有统一的父目录）
  - 如果是散乱文件，自动创建目录包装它们
  - 使用 `isScatteredFiles()` 函数检测

### 5. 依赖管理

- 使用Go标准库：`archive/tar`, `archive/zip`, `compress/gzip`, `compress/bzip2`
- 第三方库：`github.com/mholt/archiver` (v4) 用于7z支持

## 实现细节

### 压缩流程

1. 解析输出文件名，提取格式
2. 检查源路径是否存在
3. 根据格式调用相应的压缩函数
4. 处理目录和文件的递归打包

### 解压流程

1. 解析压缩包文件名，提取格式和前缀
2. 创建解压目录（使用压缩包名称前缀）
3. 根据格式调用相应的解压函数
4. 检查解压后的文件结构
5. 如果是散乱文件，创建包装目录

### 散乱文件检测逻辑

- 列出解压后的所有文件路径
- 检查是否有统一的根目录
- 如果所有文件都在解压目录的根级别（没有共同的父目录），判定为散乱文件
- 创建新目录，移动所有文件到新目录中

## 错误处理

- 文件不存在检查
- 格式不支持提示
- 权限错误处理
- 压缩/解压过程中的错误捕获和友好提示

## 测试用例覆盖

- 各种格式的压缩和解压
- 散乱文件的自动包装

## 实现状态

- ✅ 初始化Go模块，创建go.mod文件，添加依赖（archiver库）
- ✅ 实现format.go：格式检测函数、文件名前缀提取函数
- ✅ 实现compress.go：支持所有8种格式的压缩功能
- ✅ 实现extract.go：支持所有8种格式的解压功能，包含自动目录创建和散乱文件处理
- ✅ 实现main.go：命令行参数解析，调用压缩/解压功能
- ✅ 创建README.md：使用说明和示例

