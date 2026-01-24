# gz - Automated Compression/Decompression CLI Tool

An intelligent compression/decompression tool implemented in Go that automatically identifies and selects the appropriate compression/decompression method based on file extensions.

## Features

- 🎯 **Automatic Format Detection**: Automatically selects compression/decompression method based on file extension
- 📦 **Multi-Format Support**: Supports 8 common compression formats
- 🗂️ **Smart Directory Management**: Automatically creates directories when extracting, with directory name as the archive name prefix
- 📁 **Scattered Files Handling**: Automatically detects and organizes scattered files into a unified directory

## Supported Formats

- `zip` - ZIP compression format
- `tar` - TAR archive format
- `tar.gz` / `tgz` - TAR+GZIP compression
- `tar.bz2` - TAR+BZIP2 compression
- `tar.xz` - TAR+XZ compression
- `gz` - GZIP single file compression
- `bz2` - BZIP2 single file compression
- `7z` - 7-Zip compression format

## Installation

### Method 1: Download Release (Recommended)

Download the latest release package:
- https://github.com/panyingyun/gz/releases

### Method 2: Using go install (Recommended)

If you have Go 1.16 or higher installed, you can install directly using the `go install` command:

```bash
go install github.com/panyingyun/gz@latest
```

After installation, make sure `$GOPATH/bin` or `$HOME/go/bin` is in your `PATH` environment variable, then you can use the `gz` command directly.

### Method 3: Build from Source (Not recommended, requires build environment)

```bash
# Clone the repository
git clone https://github.com/panyingyun/gz.git
cd gz
# Build
make build
```

## Usage

### Compress Files/Directories

```bash
# Compress to ZIP format
gz zip images.zip folder/

# Compress to TAR.GZ format
gz zip source.tar.gz folder/

# Compress to TAR.BZ2 format
gz zip archive.tar.bz2 folder/

# Compress to 7Z format
gz zip archive.7z folder/
```

### Extract Files

```bash
# Extract ZIP file
gz unzip archive.zip

# Extract TAR.GZ file
gz unzip archive.tar.gz

# Extract TAR.XZ file
gz unzip archive.tar.xz

# Extract 7Z file
gz unzip archive.7z
```

## Feature Details

### Automatic Directory Creation

When extracting, a directory is automatically created with the name being the prefix of the archive filename.

For example:
- Extract `archive.zip` → creates `archive/` directory
- Extract `source.tar.gz` → creates `source/` directory

### Scattered Files Handling

When an archive contains scattered files (not contained in a unified folder), the tool automatically creates an `extracted/` directory and organizes all files into it, keeping the file structure clean.

## Examples

```bash
# Compress the docs folder in the current directory to ZIP format
gz zip docs.zip docs/

# Extract a downloaded archive
gz unzip download.tar.gz
# Automatically creates download/ directory and extracts contents into it

# If the archive contains scattered files, they will be automatically organized into extracted/ directory
gz unzip messy_files.zip
# Scattered files will be organized into download/extracted/ directory
```

## Notes

- `gz` and `bz2` formats only support single file compression, not directories
- Compressing large files may take some time, please be patient
- Ensure you have sufficient disk space for extraction operations

## Error Handling

The tool provides friendly error messages for:
- File not found
- Unsupported format
- Permission errors
- Other errors during compression/decompression

## License

This project is licensed under the GNU General Public License v3.0.

## Support the Author

If you find gz helpful, consider buying the author a cup of coffee ☕

<div style="display: flex; gap: 10px;">
  <img src="docs/alipay.jpg" alt="Alipay" width="200"  height="373"/>
  <img src="docs/wcpay.png" alt="WeChat Pay" width="200" height="373"/>
</div>
