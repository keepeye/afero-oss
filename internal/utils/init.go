package utils

import "os"

func init() {
	// Ensure OssObjectManager implements ObjectManager interface
	var _ ObjectManager = (*OssObjectManager)(nil)
	var _ os.FileInfo = (*OssObjectMeta)(nil)
}
