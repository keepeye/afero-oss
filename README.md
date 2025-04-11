# afero-oss

[![Go Report Card](https://goreportcard.com/badge/github.com/messikiller/afero-oss)](https://goreportcard.com/report/github.com/messikiller/afero-oss)
[![GoDoc](https://godoc.org/github.com/messikiller/afero-oss?status.svg)](https://godoc.org/github.com/messikiller/afero-oss)
[![License](https://img.shields.io/github/license/messikiller/afero-oss.svg)](https://github.com/messikiller/afero-oss/blob/main/LICENSE)

afero-oss 是一个基于 [spf13/afero](https://github.com/spf13/afero) 的阿里云对象存储服务（OSS）文件系统实现。它允许你像操作本地文件系统一样操作阿里云 OSS，提供了统一的文件系统抽象接口。

## 特性

- 完整实现 `afero.Fs` 接口
- 支持所有标准文件操作（创建、读取、写入、删除等）
- 提供目录操作支持
- 支持文件预加载和缓存
- 并发安全
- 支持上下文控制

## 安装

```bash
go get github.com/messikiller/afero-oss
```

## 快速开始

### 基础用法

```go
package main

import (
    "fmt"
    ossfs "github.com/messikiller/afero-oss"
)

func main() {
    // 创建 OSS 文件系统实例
    fs := ossfs.NewOssFs(
        "your-access-key-id",
        "your-access-key-secret",
        "your-region",
        "your-bucket-name",
    )

    // 创建文件
    f, err := fs.Create("example.txt")
    if err != nil {
        panic(err)
    }
    defer f.Close()

    // 写入内容
    _, err = f.WriteString("Hello, OSS!")
    if err != nil {
        panic(err)
    }

    // 读取文件
    content, err := afero.ReadFile(fs, "example.txt")
    if err != nil {
        panic(err)
    }
    fmt.Println(string(content))
}
```

### 使用上下文

```go
import "context"

ctx := context.Background()
fs := ossfs.NewOssFs("access-key", "secret", "region", "bucket").WithContext(ctx)
```

### 目录操作

```go
// 创建目录
err := fs.MkdirAll("path/to/dir", 0755)

// 读取目录
dir, err := fs.Open("path/to/dir")
files, err := dir.Readdir(-1)
```

### 文件操作

```go
// 打开文件
f, err := fs.OpenFile("example.txt", os.O_RDWR|os.O_CREATE, 0644)

// 写入数据
data := []byte("Hello, World!")
_, err = f.Write(data)

// 读取数据
buf := make([]byte, 1024)
n, err := f.Read(buf)

// 删除文件
err = fs.Remove("example.txt")
```

## API 文档

### 主要接口

- `NewOssFs(accessKeyId, accessKeySecret, region, bucket string) *Fs`
  - 创建新的 OSS 文件系统实例

- `WithContext(ctx context.Context) *Fs`
  - 设置操作上下文

### 文件操作

- `Create(name string) (afero.File, error)`
  - 创建新文件

- `OpenFile(name string, flag int, perm os.FileMode) (afero.File, error)`
  - 打开文件，支持各种模式

- `Remove(name string) error`
  - 删除文件

- `RemoveAll(path string) error`
  - 递归删除目录及其内容

### 目录操作

- `Mkdir(name string, perm os.FileMode) error`
  - 创建目录

- `MkdirAll(path string, perm os.FileMode) error`
  - 递归创建目录

### 其他操作

- `Stat(name string) (os.FileInfo, error)`
  - 获取文件信息

- `Rename(oldname, newname string) error`
  - 重命名/移动文件

## 高级特性

### 文件预加载

afero-oss 支持文件预加载功能，可以提高读取性能：

```go
file, err := fs.OpenFile("large-file.txt", os.O_RDWR, 0644)
if err != nil {
    panic(err)
}

// 文件会在首次读取时自动预加载到内存
data := make([]byte, 1024)
n, err := file.Read(data)
```

### 自动同步

默认情况下，文件修改会自动同步到 OSS：

```go
file, err := fs.OpenFile("example.txt", os.O_RDWR, 0644)
if err != nil {
    panic(err)
}

// 写入会自动同步到 OSS
file.WriteString("Hello, World!")
```

## 性能优化建议

1. 对于大文件操作，建议使用分片读写
2. 频繁访问的小文件会自动缓存
3. 合理使用上下文控制操作超时

## 限制说明

1. 不支持文件权限修改（chmod）
2. 不支持所有者修改（chown）
3. 不支持时间戳修改（chtimes）
4. 目录必须以 '/' 结尾

## 贡献指南

欢迎提交 Issue 和 Pull Request。在提交 PR 之前，请确保：

1. 代码通过所有测试
2. 新功能包含测试用例
3. 更新相关文档

## 许可证

本项目采用 MIT 许可证，详见 [LICENSE](LICENSE) 文件。