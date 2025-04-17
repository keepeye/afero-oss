/*
@Package: ossfs
@Link:    https://github.com/messikiller/ossfs
@Copyright (c) 2025 messikiller <messikiller@aliyun.com>
@Licence: MIT

Package ossfs provides a filesystem interface to interact with Alibaba Oss Storage.

It's an implementation of afero (https://github.com/spf13/afero) filesystem interface.

It defines a File struct that represents a file or directory in OSS. The File struct
supports methods for reading, writing, preloading, and syncing files to and from OSS.

The main features of this package include:
- File operations: Read, Write, Seek, Close, etc.
- Directory operations: Readdir, Readdirnames.
- File preloading: Automatically preload files to a local filesystem for faster access.
- File syncing: Sync preloaded files back to OSS.

This package is useful in scenarios where OSS is used as a backend storage for a filesystem.
*/
package ossfs
