# Afero-OSS: Alibaba Cloud OSS Implementation for Afero File System

## Project Introduction

> Based on the project: [spf13/afero](https://github.com/spf13/afero)

> [中文文档](./README_CN.md)

Afero-OSS is an implementation of Alibaba Cloud Object Storage Service (OSS) based on the Afero file system interface. It allows developers to operate Alibaba Cloud OSS as easily as using a local file system, providing a unified and flexible file storage solution.

## Features

- 🔄 Full compatibility with Afero file system interface
- 🚀 Support for Alibaba Cloud OSS v2 SDK
- 📂 Provides standard file and directory operations
- 🔒 Concurrent safe design
- 💾 Default memory preloading mechanism enabled to improve read and write performance

## Installation

Install using Go modules:

```bash
go get github.com/messikiller/afero-oss
```

## Quick Start

```go
package main

import (
    "github.com/spf13/afero"
    "github.com/messikiller/afero-oss"
)

func main() {
    // Create OSS file system instance
    ossFs := ossfs.NewOssFs(
        "your-access-key-id", 
        "your-access-key-secret", 
        "your-region", 
        "your-bucket-name"
    )

    // Operate OSS like a local file system
    file, err := ossFs.Create("example.txt")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    file.WriteString("Hello, Afero-OSS!")
}
```

## Main Functionalities

- File creation, reading, writing, and deletion
- Directory management (creation, listing)
- File metadata retrieval
- File preloading and synchronization

## Limitations

- Does not support file system-specific operations like `Chmod`, `Chown`
- Depends on Alibaba Cloud OSS service
- Default memory preloading for file objects is not recommended for large files. Switch preloading file system objects based on usage:

```go
ossFs := NewOssFs(...)
ossFs.WithPreloadFs(afero.NewBasePathFs(afero.NewOsFs(), "/tmp"))
```

## Contributing

Welcome to submit Issues and Pull Requests!

## License

This project is open-sourced under the MIT License.

## Dependencies

- Afero: File system abstraction interface
- Alibaba Cloud OSS Go SDK v2